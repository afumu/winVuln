package constant

var (
	VERSION_REG   = ".*?((\\d+\\.?){3}) ((Service Pack (\\d)|N\\/\\w|.+) )?[ -\\xa5]+ (\\d+).*"
	NAME_REG      = ".*?Microsoft[\\(R\\)]{0,3} Windows[\\(R\\)?]{0,3} ?(Serverr? )?(\\d+\\.?\\d?( R2)?|XP|VistaT).*"
	ARCH_REG      = ".*?([\\w\\d]+?)-based PC.*"
	KBS_REG       = ".*KB(\\d+).*"
	MATCH_KBS_REG = ".*KB\\d+.*"
)
