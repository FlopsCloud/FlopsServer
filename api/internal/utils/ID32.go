package utils

import (
	"fmt"
	"strings"
)

// 定义32进制字符集，排除 O、I、L
var base32Chars = []rune("0123456789ABCDEFGHJKMNPQRSTUVWXYZ")[:32]

// toBase32 将整数转换为32进制字符串
func toBase32(num int) string {
	if num == 0 {
		return string(base32Chars[0])
	}

	var result strings.Builder
	base := len(base32Chars) // 32

	for num > 0 {
		remainder := num % base
		result.WriteString(string(base32Chars[remainder]))
		num = num / base
	}

	// 反转字符串
	reversed := []rune(result.String())
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return string(reversed)
}

// GenerateFormattedID 生成格式化的ID
func GenerateFormattedID(number int) string {
	// if number < 0 {
	// 	return "", fmt.Errorf("输入数字必须大于等于0")
	// }

	// 转换为32进制
	base32Str := toBase32(number)

	// 使用前导字符填充到预期长度
	paddedStr := strings.Repeat(string(base32Chars[0]), 10-len(base32Str)) + base32Str

	// 如果结果长度超过10位，返回错误
	// if len(paddedStr) > 10 {
	// 	return "", fmt.Errorf("输入数字过大，生成的ID超过10位")
	// }

	return paddedStr
}

func main_test() {
	// 测试示例
	testNumbers := []int{73, 123, 45678, 1000000}

	for _, num := range testNumbers {
		id := GenerateFormattedID(num)
		// if err != nil {
		// 	fmt.Printf("输入 %d 出错: %v\n", num, err)
		// } else {
		//
		// }
		fmt.Printf("输入: %d -> 生成的ID: %s\n", num, id)
	}

	// 打印可用字符集
	fmt.Printf("\n使用的32进制字符集（共%d个字符）: %s\n",
		len(base32Chars), string(base32Chars))

	// 打印最大可表示的数
	maxValue := 1
	for i := 0; i < 10; i++ {
		maxValue *= 32
	}
	fmt.Printf("10位32进制可表示的最大数值: %d\n", maxValue-1)
}
