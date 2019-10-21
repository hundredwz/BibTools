package lib

import (
	"fmt"
	"testing"
)

func TestGoogle_Search(t *testing.T) {
	google:=new(Google)
	google.ProxyUrl="socks5://127.0.0.1:1080"
	_,err:=google.Search("blockchain")
	fmt.Println(err)
}

func TestGoogle_BibTex(t *testing.T) {
	google:=new(Google)
	google.ProxyUrl="socks5://127.0.0.1:1080"
	result,err:=google.BibTex("2d3nC4UH9DQJ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}