package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	known = [][]string{
		{"load", "areas", "program", "local", "total", "legal", "format"}, // A
		{"redirect", "dict"},                                       // C
		{"video", "view", "twitter"},                               // E
		{"info", "buff"},                                           // F
		{"image", "homepage", "page"},                              // G
		{"grid", "static", "historic", "pubic", "login"},           // I
		{"print", "loading"},                                       // N
		{"game", "home"},                                           // M
		{"information", "amazon", "facebook"},                      // O
		{"card", "return", "picture", "forward"},                   // R
		{"reverse", "adjust"},                                      // S
		{"data", "date", "site", "state", "path", "accessibility"}, // T
		{"catalogue", "previous"},                                  // U
	}
)

func main() {
	resp, _ := http.Get("https://raw.githubusercontent.com/first20hours/google-10000-english/master/google-10000-english-no-swears.txt")
	data, _ := ioutil.ReadAll(resp.Body)
	words := strings.Split(string(data), "\n")

	for l, f := 'a', 0; l < 'z'; l, f = l+1, 0 {
		matching := make([]string, 0)
		for _, word := range words {
			if len(word) >= 5 && word[len(word)-2] == byte(l) {
				if f < 10 || contains(word) {
					matching = append(matching, word)
					f++
				}
			}
		}

		fmt.Printf("/* %s */ {\"%s\"},\n", string(l-32), strings.Join(matching, "\", \""))
	}
}

func contains(word string) bool {
	for _, group := range known {
		if group[0][len(group[0])-2] != word[len(word)-2] {
			continue
		}

		for _, w := range group {
			if word == w {
				return true
			}
		}

		return false
	}

	return false
}
