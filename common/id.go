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
		/* A */ {"general", "program", "great", "international", "local", "national", "special", "today", "american", "total", "legal", "areas", "format"},
		/* B */ {"subscribe", "maybe", "clubs", "adobe", "nearby", "describe", "doubt", "toshiba", "unsubscribe", "globe"},
		/* C */ {"search", "which", "contact", "click", "service", "price", "product", "policy", "privacy", "research"},
		/* D */ {"guide", "include", "provide", "canada", "study", "needs", "cards", "friends", "trade", "already"},
		/* E */ {"other", "services", "system", "number", "after", "video", "review", "order", "under", "games"},
		/* F */ {"staff", "microsoft", "stuff", "identify", "draft", "aircraft", "wildlife", "notify", "modify", "shift"},
		/* G */ {"message", "through", "design", "technology", "change", "image", "college", "large", "things", "language", "homepage"},
		/* H */ {"copyright", "right", "night", "light", "might", "months", "weight", "thought", "photography", "brought"},
		/* I */ {"their", "email", "music", "public", "within", "media", "visit", "credit", "movie", "again", "login", "static", "historic"},
		/* K */ {"books", "links", "works", "thanks", "weeks", "networks", "trademarks", "looks", "alaska", "banks"},
		/* L */ {"would", "people", "world", "should", "available", "could", "details", "hotels", "family", "while"},
		/* M */ {"items", "terms", "systems", "forums", "programs", "problems", "become", "welcome", "income", "volume"},
		/* N */ {"online", "company", "management", "development", "using", "phone", "shipping", "being", "found", "following", "print", "loading"},
		/* O */ {"information", "school", "education", "version", "section", "control", "location", "description", "author", "photos", "amazon"},
		/* P */ {"groups", "europe", "happy", "ships", "except", "accept", "perhaps", "therapy", "steps", "shops"},
		/* R */ {"there", "support", "software", "where", "years", "january", "store", "report", "before", "members", "return", "picture", "forward"},
		/* S */ {"business", "first", "these", "please", "because", "those", "address", "house", "access", "class", "reverse", "adjust"},
		/* T */ {"state", "health", "products", "rights", "university", "comments", "results", "community", "website", "south", "accessibility"},
		/* U */ {"about", "group", "forum", "without", "previous", "value", "status", "issue", "various", "linux", "catalogue"},
		/* V */ {"above", "archive", "drive", "receive", "active", "effective", "believe", "executive", "leave", "improve"},
		/* W */ {"reviews", "windows", "known", "shows", "views", "shown", "brown", "allows", "unknown", "follows"},
		/* Y */ {"always", "holidays", "plays", "displays", "surveys", "attorneys", "vinyl", "essays", "tokyo", "kenya"},
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
