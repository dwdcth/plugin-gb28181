package utils

import (
	"encoding/json"
	"math/rand"
)

func ToJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func ToPrettyString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "    ")
	return string(b)
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString returns a random string with a fixed length
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
