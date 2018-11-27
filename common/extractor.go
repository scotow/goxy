package common

import (
	"io"
)

func newExtractor(reader io.Reader, offset int) *Extractor {
	return &Extractor{
		skip:   offset,
		reader: reader,
	}
}

type Extractor struct {
	skip   int
	reader io.Reader
}

func (e *Extractor) Read(p []byte) (int, error) {
	if e.skip > 0 {
		buffer := make([]byte, len(p))
		read, err := e.reader.Read(buffer)

		if read > e.skip {
			copied := copy(p, buffer[e.skip:read])
			e.skip = 0
			return copied, err
		} else {
			e.skip -= read
			return 0, err
		}
	} else {
		return e.reader.Read(p)
	}
}
