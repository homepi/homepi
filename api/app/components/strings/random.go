package strings

import (
	"math/rand"
	"time"
)

var numbers = []rune("0123456789")
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Random index generator
func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

// Random digit generator
func RandomDigits(lenght int) (digits string) {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= lenght; i++ {
		if digit := numbers[random(0, 9)]; digit != 0 {
			digits += string(digit)
		}
		continue
	}
	return
}

// Random string generator
func RandomLetters(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}