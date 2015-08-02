package utils

import (
	"math/rand"
	"time"
)

func Random() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //以当前纳秒数作为随机数种子
	return r.Int63()
}

func XsrfToken() string {
	//return MD5(strconv.Itoa((int)Random()))
	return ""
}
