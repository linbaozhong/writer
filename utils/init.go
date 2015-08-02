package utils

import (
	"os"
	"regexp"
	"strings"
	"time"
)

var ()

func init() {

}

func Str2Time(s string) time.Time {
	s = strings.Replace(s, "/", "-", 0)
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

//将标准时间转为自1970-1-1开始的毫秒数
func Millisecond(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

//将自1970-1-1开始的毫秒数转为时间
func Msec2Time(msec int64) time.Time {
	if msec == 0 {
		return time.Unix(0, 0)
	}
	t := msec / 1000
	return time.Unix(t, msec*1000000%t)
}

//缩进
func Indent(s string, i int8) string {
	if i > 0 {
		return strings.Repeat("　", int(i)*2) + s
		//return fmt.Sprintf("<span style=\"margin-left:%dem;display:inline-block;\"></span>%s", int(i)*2, s)
	}
	return s
}

//文件夹是否存在
func DirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

//
func GetIp(addr string) string {
	if strings.Index(addr, "[") > -1 {
		reg := regexp.MustCompile(`\[(\S+)\]`)
		return string(reg.FindSubmatch([]byte(addr))[1])
	}
	pos := strings.Index(addr, ":")
	if pos > -1 {
		return addr[:pos]
	}
	return addr
}
