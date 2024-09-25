package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

//RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1) 
}

//RandomStr generates a random string of length n
func RandomStr(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomStr(6)
}

//RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//RandomCur generates a random currency code
func RandomCur() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

//RandomEmail generates a random email address
func RandomEmail(n int) string {
	mails := []string{"gmail.com", "yahoo.com"}
	k := len(mails)

	str := RandomStr(n)

	return fmt.Sprintf("%v@%v",str,mails[rand.Intn(k)])
}