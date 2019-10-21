package lib

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"github.com/hundredwz/BibTools/cmn"
	"io/ioutil"
	"strings"
)


const IEEESearchUrl = `https://ieeexplore.ieee.org/rest/search`
const IEEERefUrl = `https://ieeexplore.ieee.org/xpl/downloadCitations`



type IEEE struct {
	ProxyUrl string
}

type IEEESearchResult struct {
	Records []struct {
		Authors []struct {
			PreferredName           string
			NormalizedName          string
			FirstName               string
			LastName                string
			SearchablePreferredName string
			ID                      int
		}
		ArticleTitle  string
		ArticleNumber string
		Abstract      string
	}
}

func NewIEEE(proxyUrl string)*IEEE{
	return &IEEE{
		ProxyUrl:proxyUrl,
	}
}

func (i *IEEE)SetProxy(proxyUrl string){
	i.ProxyUrl=proxyUrl
}

func (i *IEEE) Search(query string) ([]cmn.Record, error) {
	req := httplib.Post(IEEESearchUrl)
	req.Header("Content-Type", "application/json")
	req.Header("Origin", "https://ieeexplore.ieee.org")
	params := map[string]interface{}{
		"newsearch": true,
		"queryText": query,
	}
	_, err := req.JSONBody(params)
	if err != nil {
		return nil, err
	}
	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var searchResult IEEESearchResult
	err = json.Unmarshal(content, &searchResult)
	if err != nil {
		return nil, err
	}
	records := make([]cmn.Record, 0)
	for _, v := range searchResult.Records {
		var buf bytes.Buffer
		for _, v := range v.Authors {
			buf.WriteString(v.NormalizedName+",")
		}
		record := cmn.Record{
			ID:       v.ArticleNumber,
			Title:    v.ArticleTitle,
			Authors:  buf.String(),
			Abstract: v.Abstract,
		}
		records = append(records, record)
	}
	return records, nil
}

func (i *IEEE) BibTex(id string) (string, error) {
	req := httplib.Post(IEEERefUrl)
	req.Header("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header("Origin", "https://ieeexplore.ieee.org")
	req.Param("recordIds", id)
	req.Param("download-format", "download-bibtex")
	req.Param("citations-format", "citation-only")
	resp, err := req.Response()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ref := strings.Replace(string(content), "<br>", "", -1)
	return ref, nil
}
