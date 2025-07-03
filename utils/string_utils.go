package utils

import (
	"unicode/utf8"
)

// CheckStringUnicodeLength Unicode编码字符长度检查
func CheckStringUnicodeLength(s string) int {
	// 使用utf8.RuneCountInString检查字符长度
	return utf8.RuneCountInString(s)
}

// TruncateStringByUnicodeLength 将字符串裁剪成指定 Unicode编码长度
func TruncateStringByUnicodeLength(s string, length int) string {
	if length <= 0 {
		return ""
	}

	// 将字符串转换为rune数组
	runes := []rune(s)

	// 确保不截取超过字符串长度的字符
	if length > len(runes) {
		length = len(runes)
	}

	// 截取前length个字符
	truncatedRunes := runes[:length]

	// 将截取后的字符数组转换回字符串
	return string(truncatedRunes)
}
