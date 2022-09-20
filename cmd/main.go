package main

import (
	"fmt"
	"github.com/zouchangfu/winVuln/pkg/db"
	"github.com/zouchangfu/winVuln/pkg/detector"
	"github.com/zouchangfu/winVuln/pkg/systeminfo"
)

func main() {
	// 获取操作系统名称和操作系统安装的补丁
	product, kbs := systemInfo.GetProductAndKbs()
	var kbmap = make(map[string]string)
	for _, v := range kbs {
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
