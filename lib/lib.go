package lib

import (
	"github.com/hundredwz/BibTools/cmn"
)

type Lib interface {
	Search(query string) ([]cmn.Record, error)
	BibTex(id string) (string, error)
	SetProxy(proxyUrl string)
}

