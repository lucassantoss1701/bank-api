package entity

import (
	"bytes"
	"regexp"
	"strconv"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUUID() string {
	return uuid.NewString()
}

func hash(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
}

func hashIsValid(hashedPassowrd string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassowrd), []byte(password))
	return err == nil
}

var REGEXCPF = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)

// could be VO
func isCPF(value string) bool {
	const (
		size = 9
		pos  = 10
	)

	return validateDocument(value, REGEXCPF, size, pos)
}

func cleanNonDigits(doc *string) {

	buf := bytes.NewBufferString("")
	for _, r := range *doc {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	*doc = buf.String()
}

func validateDocument(documentValue string, pattern *regexp.Regexp, size int, position int) bool {

	if !pattern.MatchString(documentValue) {
		return false
	}

	cleanNonDigits(&documentValue)

	if valuesAllEquals(documentValue) {
		return false
	}

	d := documentValue[:size]
	digit := calculateDigit(d, position)

	d = d + digit
	digit = calculateDigit(d, position+1)

	return documentValue == d+digit
}

func valuesAllEquals(doc string) bool {

	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}

	return true
}

func convertToInt(r rune) int {
	return int(r - '0')
}

func calculateDigit(doc string, position int) string {

	var sum int
	for _, r := range doc {

		sum += convertToInt(r) * position
		position--

		if position < 2 {
			position = 9
		}
	}

	sum %= 11
	if sum < 2 {
		return "0"
	}

	return strconv.Itoa(11 - sum)
}
