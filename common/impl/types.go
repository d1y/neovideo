package impl

type JiexiParse struct {
	Name string
	URL  string
}

type MacCMSParse struct {
	Api        string
	Name       string
	R18        bool
	RespType   string
	JiexiParse bool
	JiexiURL   string
}
