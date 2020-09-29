package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt：生成介于最小值和最大值之间的随机整数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成长度为n的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner 生成一个随机名称
func RandomOwner() string {
	return RandomString(5)
}

// RandomMoney 产生随机的金额
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency 生成随机货币代码
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "RMB"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
