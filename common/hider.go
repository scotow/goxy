package common

import (
	"errors"
	"io"
	"math/rand"
	"strings"
)

var (
	hiders = []Hider{
		{"image/png", "png", []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A")},
		{"image/jpeg", "jpg", []byte("\xFF\xD8\xFF")},
		{"font/ttf", "ttf", []byte("\x00\x01\x00\x00")},
		{"application/zip", "zip", []byte("\x50\x4B\x03\x04")},
		{"text/html", "html", []byte(`<!DOCTYPEhtml><htmllang="en"dir="ltr"><head><metacharset="utf-8"><title>Worldline</title></head><body><divclass="azeazeaz">zaeazeaz</div><divclass="azeaeaze">azaeae</div></body></html>`)},
	}
)

var (
	ErrNoHiderAvailable = errors.New("no hider available")
	ErrInvalidExtension = errors.New("invalid path extension")
)

func RandomHider() (*Hider, error) {
	if len(hiders) == 0 {
		return nil, ErrNoHiderAvailable
	}

	return &hiders[rand.Intn(len(hiders))], nil
}

func HiderFromPath(path string) (*Hider, error) {
	for _, hider := range hiders {
		if strings.HasSuffix(path, hider.Extension) {
			return &hider, nil
		}
	}

	return nil, ErrInvalidExtension
}

type Hider struct {
	Mime      string
	Extension string
	prefix    []byte
}

func (h *Hider) HideData(input []byte) []byte {
	return append(h.prefix, input...)
}

func (h *Hider) GetExtractor(reader io.Reader) *Extractor {
	return newExtractor(reader, len(h.prefix))
}
