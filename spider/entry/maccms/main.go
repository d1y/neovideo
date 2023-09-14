package main

import (
	"flag"
	"fmt"
	"sync"

	"d1y.io/neovideo/spider/implement/maccms"
)

var cms = maccms.NewMacCMS(maccms.MacCMSReponseTypeXML, "https://www.hanjuzy.com/inc/api.php")

var maccmsType = flag.String("type", "xml", "接口类型")

var wg sync.WaitGroup

func xmlDemo() {
	wg.Add(4)
	go func() {
		defer wg.Done()
		_, data, _ := cms.XMLGetDetail(5292)
		fmt.Println(data)
	}()
	go func() {
		defer wg.Done()
		data, err := cms.XMLGetSearch("真的出现了", 1)
		fmt.Println(data, err)
	}()
	go func() {
		defer wg.Done()
		var data, _ = cms.XMLGetHome()
		fmt.Printf("%v", data)
	}()
	go func() {
		defer wg.Done()
		var category, _ = cms.XMLGetCategory()
		fmt.Printf("%v", category)
	}()
	wg.Wait()
}

func jsonDemo() {
	// TODO: impl this
	fmt.Println("demo")
}

func main() {
	flag.Parse()
	if *maccmsType == "xml" {
		xmlDemo()
	} else if *maccmsType == "json" {
		jsonDemo()
	}
}
