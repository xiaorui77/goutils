package math

const BaseNumber = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_+=!@#$%^&*()"

// Base 将10进制转为其他进制, 最大支持64进制
func Base(num uint64, base int) string {
	if base < 2 || base > len(BaseNumber) {
		return ""
	}
	str := ""
	for num != 0 {
		str = string(BaseNumber[num%uint64(base)]) + str
		num = num / uint64(base)
	}
	return str
}
