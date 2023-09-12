package main

import (
	"fmt"
	"sync"

	"d1y.io/neovideo/spider/implement/maccms"
)

var cms = maccms.NewMacCMS(maccms.MacCMSReponseTypeXML, "https://www.hanjuzy.com/inc/api.php")

var wg sync.WaitGroup

func main() {
	wg.Add(4)
	go func() {
		defer wg.Done()
		_, data, _ := cms.GetDetail(5292)
		fmt.Println(data)
	}()
	go func() {
		defer wg.Done()
		data, err := cms.GetSearch("真的出现了", 1)
		fmt.Println(data, err)
	}()
	go func() {
		defer wg.Done()
		var data, _ = cms.GetHome()
		fmt.Printf("%v", data)
	}()
	go func() {
		defer wg.Done()
		var category, _ = cms.GetCategory()
		fmt.Printf("%v", category)
	}()
	wg.Wait()
}
