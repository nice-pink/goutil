package random

import (
	"math/rand"
)

// letters

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// alpha numberical

var letterRunesAlphaNum = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringAlphaNum(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunesAlphaNum[rand.Intn(len(letterRunesAlphaNum))]
	}
	return string(b)
}

// alpha numberical extended

var letterRunesAlphaNumExt = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")

func RandStringAlphaNumExt(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunesAlphaNumExt[rand.Intn(len(letterRunesAlphaNumExt))]
	}
	return string(b)
}
