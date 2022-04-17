package fileutils

import (
	"regexp"
	"strings"
)

var (
	// Windows file name: not `?“”/\\<>*|` and len less then 255
	windowsFileReg = regexp.MustCompile(`[\s/\\:*?"<>|]+`)
)

func WindowsName(str string) string {
	if len(str) > 255 {
		str = str[:255]
	}
	str = windowsFileReg.ReplaceAllString(str, " ")
	return strings.TrimSpace(str)
}

// HTMLDecode HTML escape. HTML的&lt;&gt;&amp;&quot;&copy; 分别是<>&"©的转义字符
func HTMLDecode(str string) string {
	str = strings.ReplaceAll(str, "&lt;", "<")
	str = strings.ReplaceAll(str, "&gt;", ">")
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&quot;", `"`)
	str = strings.ReplaceAll(str, "&copy;", "©")
	return str
}

// XMLDecode XML escape. &lt; &gt; &amp; &quot; &apos 分别是<>&"'的转义字符
func XMLDecode(str string) string {
	str = strings.ReplaceAll(str, "&lt;", "<")
	str = strings.ReplaceAll(str, "&gt;", ">")
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&quot;", `"`)
	str = strings.ReplaceAll(str, "&apos;", "'")
	return str
}
