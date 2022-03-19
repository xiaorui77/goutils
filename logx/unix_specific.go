//go:build !windows

package logx

import "strconv"

func coloring(str string, color int, enable bool) string {
	if !enable {
		return str
	}
	return "\033[" + strconv.Itoa(color) + "m" + str + "\033[0m"
}
