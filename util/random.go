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

// 生成[min,max]之间的随机数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// 生成随机的string
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// 生成随机的owner
func RandomOwner() string {
	return RandomString(6)
}

// 生成随机的存款(0-1000)
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// 生成随机的货币
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
