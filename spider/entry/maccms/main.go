package main

import (
	"fmt"

	"d1y.io/neovideo/spider/implement/maccms"
)

var cms = maccms.NewMacCMS(maccms.MacCMSReponseTypeXML, "https://www.hanjuzy.com/inc/api.php")

func main() {
	var data, _ = cms.GetHome()
	fmt.Printf("%v", data)
}
