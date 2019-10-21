package cmn

const (
	IEEE          = `IEEE`
	Google        = `Google`
	TiTleCutLen   = 75
	ContentCutLen = 75
)


type Record struct {
	ID       string
	Title    string
	Authors  string
	Abstract string
}