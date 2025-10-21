package util

import (
	"math/rand"
	"strings"
)
const alphabet = "qwertyuiopasdfghjklzxcvbnm"

//Создание генератора чисел для unit тестов с min max
func RandomInt(min, max int64) int64 {
	return rand.Int63n(max - min + 1)
}
//Создание генератора рандомной строки
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i:=0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

return  sb.String()
}

//Создание генератора для колонки owner 

func RandomOwner() string {
	return RandomString(6)
}

//Создание рандомного генератора для денег
func RandomMoney() int64 {
	return  RandomInt(0, 1000)
}

//Генератор случайной валюты
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	l := len(currencies)
	return  currencies[rand.Intn(l)]
}