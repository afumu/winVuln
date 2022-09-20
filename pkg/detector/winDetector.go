package detector

import "strings"

type WinVuln struct {
	Id                string `column:"id" json:"id"`
	DatePosted        string `column:"date_posted" json:"datePosted"`
	Cve               string `column:"cve" json:"cve"`
	BulletinKb        string `column:"bulletin_kb" json:"bulletinKb"`
	Title             string `column:"title" json:"title"`
	AffectedProduct   string `column:"affected_product" json:"affectedProduct"`
	AffectedComponent string `column:"affected_component" json:"affectedComponent"`
	Severity          string `column:"severity" json:"severity"`
	Impact            string `column:"impact" json:"impact"`
	Supersedes        string `column:"supersedes" json:"supersedes"`
	Exploits          string `column:"exploits" json:"exploits"`
	Relevant          bool
}

func Detect(product string, list []*WinVuln, kbResults []string) ([]*WinVuln, []*WinVuln) {
	var filtered []*WinVuln

	if strings.Contains(product, "Service Pack") {
		for _, cve := range list {
			if !strings.Contains(cve.AffectedProduct, product) {
				continue
			}
			cve.Relevant = true
			filtered = append(filtered, cve)
			if cve.Supersedes != "" {
				kbResults = append(kbResults, cve.Supersedes)
			}
		}
	} else {
		productSp := product + " Service Pack"
		for _, cve := range list {
			// 判断当前cve漏洞是否影响当前操作系统
			if !strings.Contains(cve.AffectedProduct, product) || strings.Contains(cve.AffectedProduct, productSp) {
				continue
			}

			// 把当前漏洞软件标记为相关
			cve.Relevant = true
			filtered = append(filtered, cve)

			// 为什么需要把这个Supersedes添加到kbResults补丁中呢？
			// Supersedes的意思代表是当前软件包会包含的补丁号，如果当前操作系统有这个软件包，说明当前操作系统也就存在这个补丁号
			// Supersedes 指的是当前包软件存在的补丁包的编号
			if cve.Supersedes != "" {
				kbResults = append(kbResults, cve.Supersedes)
			}
		}
	}

	// 合并补丁包
	join := strings.Join(kbResults, ";")
	markSuperseededHotfix(filtered, join)

	// 获取存在漏洞的软件
	var check []*WinVuln
	for _, v := range filtered {
		if v.Relevant {
			check = append(check, v)
		}
	}

	// 这里其实是做兜底的，应该不会起作用
	var supersedes = make(map[string]*WinVuln)
	for _, v := range check {
		supersedes[v.Supersedes] = v
	}
	var checked []*WinVuln
	for _, v := range check {
		if supersedes[v.BulletinKb] != nil {
			v.Relevant = false
			checked = append(checked, v)
		}
	}

	var found []*WinVuln
	for _, v := range filtered {
		if v.Relevant {
			found = append(found, v)
		}
	}
	return filtered, found
}

func markSuperseededHotfix(filterd []*WinVuln, superseeded string) {
	hotfixes := strings.Split(superseeded, ";")
	// 遍历所有的系统补丁
	for _, kb := range hotfixes {
		// 获取到操作系统已经存在补丁的漏洞软件包
		// 然后把他们设置为不相关
		for _, cve := range filterd {
			if cve.Relevant && cve.BulletinKb == kb {
				cve.Relevant = false
			}
		}
	}
}
func markContain(marked []string, str string) bool {
	for _, v := range marked {
		if str == v {
			return true
		}
	}
	return false
}
