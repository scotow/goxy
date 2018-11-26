package common

import (
	"github.com/pkg/errors"
	"math/rand"
	"strings"
)

const (
	idLength = 3
)

var (
	ErrInvalidLetter      = errors.New("invalid letter from token")
	ErrInvalidTokenLength = errors.New("invalid token length")
)

var (
	words = [][]string{
		{"load"},                      // A
		{"redirect", "dict"},          // C
		{"video", "view", "twitter"},  // E
		{"info", "buff"},              // F
		{"image", "homepage", "page"}, //G
		{"grid", "static", "historic", "pubic", "login"}, // I
		{"print", "loading"},                             // N
		{"game", "home"},                                 // M
		{"information", "amazon", "facebook"},            // O
		{"card", "return", "picture", "forward"},         // R
		{"reverse", "adjust"},                            // S
		{"data", "date", "site", "state", "path"},        // T
		{"catalogue", "previous"},                        // U
	}
)

func TokenFromPath(path string) string {
	parts := strings.Split(path, "/")
	letters := make([]string, len(parts))
	for index, part := range parts {
		letters[index] = string(part[len(part)-2])
	}
	return strings.Join(letters, "")
}

type Id []int

func NewRandomId() *Id {
	r := make([]int, idLength)
	for i := 0; i < idLength; i++ {
		r[i] = rand.Intn(len(words))
	}
	id := Id(r)
	return &id
}

func NewIdFromToken(token string) (*Id, error) {
	if len(token) != idLength {
		return nil, ErrInvalidTokenLength
	}

	r := make([]int, idLength)
	for index, letter := range []byte(token) {
		r[index] = indexFromLetter(letter)
		if r[index] == -1 {
			return nil, ErrInvalidLetter
		}
	}

	id := Id(r)
	return &id, nil
}

func indexFromLetter(letter byte) int {
	for index, wordsOfSameLetter := range words {
		firstWord := wordsOfSameLetter[0]
		if letter == firstWord[len(firstWord)-2] {
			return index
		}
	}

	return -1
}

func (i *Id) RandomPath() string {
	s := make([]string, len(*i))
	for index, part := range *i {
		s[index] = words[part][rand.Intn(len(words[part]))]
	}
	return strings.Join(s, "/")
}

func (i *Id) Token() string {
	letters := make([]string, len(*i))
	for index, part := range *i {
		word := words[part][0]
		letters[index] = string(word[len(word)-2])
	}
	return strings.Join(letters, "")
}
