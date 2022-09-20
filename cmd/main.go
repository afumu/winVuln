package main

import (
	"fmt"
	"github.com/zouchangfu/winVuln/pkg/db"
	"github.com/zouchangfu/winVuln/pkg/detector"
	"github.com/zouchangfu/winVuln/pkg/systeminfo"
)

func main() {
	version, build, win, arch, product, kbResults := systemInfo.DetermineProduct()
	fmt.Println(product, win, build, version, arch, kbResults)
	var kbmap = make(map[string]string)
	for _, v := range kbResults {
		kbmap[v] = ""
	}
	var newKbResults []string
	for key := range kbmap {
		newKbResults = append(newKbResults, key)
	}
	allCve := db.GetAllCve()
	filtered, found := detector.Detect(product, allCve, newKbResults)
	fmt.Println(filtered, found)
}
