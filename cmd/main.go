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

	// 从db中获取所有的漏洞
	allCve := db.GetAllCve()

	// 开始检测漏洞
	result := detector.Detect(product, allCve, kbs)
	for _, v := range result {
		fmt.Printf("%v\n", v)
	}
}
