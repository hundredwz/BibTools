package lib

import (
	"fmt"
	"testing"
)

func TestIEEE_Search(t *testing.T) {
	ieee:=new(IEEE)
	result,err:=ieee.Search("blockchain")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v",result)
}

func TestIEEE_BibTex(t *testing.T) {
	ieee:=new(IEEE)
	result,err:=ieee.BibTex("8712639")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}