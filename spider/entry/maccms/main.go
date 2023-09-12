package main

import (
	"fmt"
	"sync"

	"d1y.io/neovideo/spider/implement/maccms"
)

var cms = maccms.NewMacCMS(maccms.MacCMSReponseTypeXML, "https://www.hanjuzy.com/inc/api.php")

var wg sync.WaitGroup

func main() {
	wg.Add(2)
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
