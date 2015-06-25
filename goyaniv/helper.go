package goyaniv

import (
	"math/rand"
	"net/http"
	"time"
)

func CreateCookie(key string, value string) *http.Cookie {
	return &http.Cookie{
		Name:  key,
		Value: value,
	}
}

func GenerateUnique() string {
	var r = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0")
	uniq := make([]rune, 20)
	for i := range uniq {
		rand.Seed(time.Now().UTC().UnixNano())
		uniq[i] = r[rand.Intn(len(r))]
	}
	return string(uniq)
}
