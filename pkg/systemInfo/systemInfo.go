package systemInfo

import (
	"bufio"
	"fmt"
	"github.com/zouchangfu/winVuln/pkg/utils"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var buildnumberMap = make(map[int]string)

func init() {
	buildnumberMap[10240] = "1507"
	buildnumberMap[10586] = "1511"
	buildnumberMap[14393] = "1607"
	buildnumberMap[15063] = "1703"
	buildnumberMap[16299] = "1709"
	buildnumberMap[17134] = "1803"
	buildnumberMap[17763] = "1809"
	buildnumberMap[18362] = "1903"
	buildnumberMap[18363] = "1909"
	buildnumberMap[19041] = "2004"
	buildnumberMap[19042] = "20H2"
	buildnumberMap[19043] = "21H1"
	buildnumberMap[19044] = "21H2"
	buildnumberMap[20348] = "21H2"
	buildnumberMap[22000] = "21H2"
}

// GetSystemInfo 获取window的系统信息
func GetSystemInfo() string {
	cmd := exec.Command("systeminfo")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return ""
	}

	in := bufio.NewScanner(stdout)
	builder := strings.Builder{}
	for in.Scan() {
		cmdRe := utils.ConvertByte2String(in.Bytes(), "GB18030")
		builder.Write([]byte(cmdRe + "\n"))
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(err)
		return ""
	}
	systemInfo := builder.String()
	return systemInfo
}

func DetermineProduct() (string, int, string, string, string, []string) {
	info := GetSystemInfo()

	// 系统版本
	versionMatches := utils.GetValueStringByRegex(info, ".*?((\\d+\\.?){3}) ((Service Pack (\\d)|N\\/\\w|.+) )?[ -\\xa5]+ (\\d+).*")
	servicePack := versionMatches[5]
	osBuild := versionMatches[6]
	fmt.Println("servicePack:", servicePack)
	fmt.Println("osbuild:", osBuild)

	var version string
	build, _ := strconv.Atoi(osBuild)
	for key, value := range buildnumberMap {
		if build == key {
			version = value
			break
		}
		if build > key {
			version = value
		} else {
			break
		}
	}
	fmt.Println("version:", version)

	// 系统名称版本
	winMatches := utils.GetValueStringByRegex(info, ".*?Microsoft[\\(R\\)]{0,3} Windows[\\(R\\)?]{0,3} ?(Serverr? )?(\\d+\\.?\\d?( R2)?|XP|VistaT).*")
	win := winMatches[2]
	fmt.Println("win:", win)

	// 系统架构
	osArchs := utils.GetValueStringByRegex(info, ".*?([\\w\\d]+?)-based PC.*")
	fmt.Println("osArchs:", osArchs[1])
	arch := osArchs[1]

	if !isProducts(win) {
		if arch == "X86" {
			arch = "32-bit"
		} else if arch == "x64" {
			arch = "x64-based"
		}
	}

	var product string
	if win == "XP" {
		product = "Microsoft Windows XP"
		if arch != "X86" {
			product += fmt.Sprintf(" Professional %s Edition", arch)
		}
		if servicePack != "" {
			product += fmt.Sprintf(" Service Pack %s", servicePack)
		}
	} else if win == "VistaT" {
		product = "Windows Vista"
		if arch != "X86" {
			product += fmt.Sprintf(" %s Edition", arch)
		}
		if servicePack != "" {
			product += fmt.Sprintf(" Service Pack %s", servicePack)
		}
	} else if win == "7" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
		if servicePack != "" {
			product += fmt.Sprintf(" Service Pack %s", servicePack)
		}
	} else if win == "8" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
	} else if win == "8.1" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
	} else if win == "10" {
		product = fmt.Sprintf("Windows %s Version %s for %s Systems", win, version, arch)
	} else if win == "11" {
		product = fmt.Sprintf("Windows %s for %s Systems", win, arch)
	} else if win == "2003" {
		if arch == "X86" {
			arch = ""
		} else if arch == "x64" {
			arch = " x64 Edition"
		}
		var pversion = " "
		if version != "" {
			pversion += version
		}
		product = fmt.Sprintf("Microsoft Windows Server %s%s%s", win, arch, pversion)
	} else if win == "2008" {
		var pversion = " "
		if version != "" {
			pversion += version
		}
		product = fmt.Sprintf("Windows Server %s for %s Systems%s", win, arch, pversion)
	} else if win == "2008 R2" {
		var pversion = " "
		if version != "" {
			pversion += version
		}
		product = fmt.Sprintf("Windows Server %s for %s Systems%s", win, arch, pversion)
	} else if win == "2012" || win == "2012 R2" || win == "2016" || win == "2019" || win == "2022" {
		product = fmt.Sprintf("Windows Server %s", win)
	}

	// 补丁
	compile := regexp.MustCompile(".*KB\\d+.*")
	all := compile.FindAll([]byte(info), -1)
	var kbResults []string
	for _, v := range all {
		kbs := utils.GetValueStringByRegex(string(v), ".*KB(\\d+).*")
		kbResults = append(kbResults, kbs[1])
	}
	fmt.Printf("kbResults%+v\n:", kbResults)
	return version, build, win, arch, product, kbResults
}

func isProducts(win string) bool {
	var products = []string{"XP", "VistaT", "2003", "2003 R2"}
	for _, v := range products {
		if win == v {
			return true
		}
	}
	return false
}
