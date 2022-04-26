package dragon

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// 将十进制数字转化为二进制字符串
func convertToBin(num int) string {
	s := ""

	if num == 0 {
		return "0"
	}

	// num /= 2 每次循环的时候 都将num除以2  再把结果赋值给num
	for ; num > 0; num /= 2 {
		lsb := num % 2
		// 将数字强制性转化为字符串
		s = strconv.Itoa(lsb) + s
	}
	return s
}

// 字符串通过缓冲的方式进行拼接
func stringJoin(args ...string) string {
	var build strings.Builder

	for _, str := range args {
		build.WriteString(str)
	}

	return build.String()
}

// RegexExtract 正则表达式提取
func regexExtract(s string) (string, string, bool) {
	length := len(s)
	if s[length-1] != '}' {
		return "", "", false
	}
	regex := -1
	for i, v := range s {
		if v == '{' {
			regex = i + 1
		}
	}
	if regex == -1 {
		return "", "", false
	}
	return s[:regex-1], s[regex : length-1], true
}

// 正则表达式验证
func regexVerify(pattern, s string) bool {
	is, _ := regexp.MatchString(pattern, s)
	return is
}

const (
	tiocgwinsz    = 0x5413
	tiocgwinszOsx = 1074295912
)

type window struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// TerminalWidth 得到命令行大小
func terminalWidth() (int, error) {
	w := new(window)
	tio := syscall.TIOCGWINSZ
	if runtime.GOOS == "darwin" {
		tio = tiocgwinszOsx
	}
	res, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tio),
		uintptr(unsafe.Pointer(w)),
	)
	if int(res) == -1 {
		return 0, err
	}
	return int(w.Col), nil
}
