package util

import (
	"fmt"
	"github.com/rotisserie/eris"
	"math/rand"
	"streamtg/go-log"
	"strings"
	"time"
)

func IntInSlice(slice []int64, aI interface{}) bool {
	var a int64

	switch aI.(type) {
	case int:
		a = int64(aI.(int))
	case int64:
		a = aI.(int64)
	default:
		return false
	}

	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

func StringInSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ErrCheckFatal(err error, strs ...string) {
	if err != nil {
		if len(strs) != 0 {
			for _, str := range strs {
				err = eris.Wrap(err, str)
			}
		}
		log.Fatal(err)
	}
}

func FloatCut(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.8f", f), "0"), ".")
}

func CreateKey(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteRune(chars[r.Intn(len(chars))])
	}
	return b.String()
}

func TextCut(s, addStr string, n int) string {
	if len([]rune(s)) >= n {
		return string([]rune(s)[:n]) + addStr
	}
	return s
}
