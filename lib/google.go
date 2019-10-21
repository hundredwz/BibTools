package lib

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/httplib"
	"github.com/hundredwz/BibTools/cmn"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Google struct {
	ProxyUrl string
}

func NewGoogle(proxyUrl string)*Google{
	return &Google{
		ProxyUrl:proxyUrl,
	}
}
const GoogleSearchUrl = `https://scholar.google.com/scholar`

func (g *Google)SetProxy(proxyUrl string){
	g.ProxyUrl=proxyUrl
}

func (g *Google) Search(query string) ([]cmn.Record, error) {
	req := httplib.Get(GoogleSearchUrl)
	if g.ProxyUrl != "" {
		req.SetProxy(func(request *http.Request) (*url.URL, error) {
			return url.Parse(g.ProxyUrl)
		})
	}
	req.Param("hl", "zh-CN")
	req.Param("q", query)
	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	records := make([]cmn.Record, 0)
	doc.Find(".gs_r.gs_or.gs_scl").Each(func(i int, s *goquery.Selection) {
		id, exist := s.Find(".gs_rt > a").Attr("id")
		if !exist {
			return
		}
		authors := s.Find(".gs_a").Text()
		title := s.Find(".gs_rt > a").Text()
		abstract := s.Find(".gs_rs").Text()
		record := cmn.Record{
			ID:       id,
			Title:    title,
			Authors:  authors,
			Abstract: abstract,
		}
		records = append(records, record)
	})
	return records, nil
}

func (g *Google) BibTex(id string) (string, error) {
	req := httplib.Get(GoogleSearchUrl)
	if g.ProxyUrl != "" {
		req.SetProxy(func(request *http.Request) (*url.URL, error) {
			return url.Parse(g.ProxyUrl)
		})
	}
	req.Param("output", "cite")
	req.Param("q", "info:"+id+":scholar.google.com/")
	resp, err := req.Response()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	bibUrl, exist := doc.Find("a.gs_citi:contains(BibTeX)").Attr("href")
	if !exist {
		return "", errors.New("bib link not exist")
	}
	refReq := httplib.Get(bibUrl)
	if g.ProxyUrl != "" {
		refReq.SetProxy(func(request *http.Request) (*url.URL, error) {
			return url.Parse(g.ProxyUrl)
		})
	}
	refResp, err := refReq.Response()
	if err != nil {
		return "", err
	}
	defer refResp.Body.Close()
	bibText, err := ioutil.ReadAll(refResp.Body)
	if err != nil {
		return "", err
	}
	return string(bibText), nil

}
