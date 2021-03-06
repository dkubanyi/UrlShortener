package encoding

import (
	"fmt"
	"strings"
)

// All characters
const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length   = int64(len(alphabet))
)

type Base62Encoder struct{}

func (e *Base62Encoder) Encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	s := ""
	for ; n > 0; n = n / length {
		s = string(alphabet[n%length]) + s
	}
	return s
}

func (e *Base62Encoder) Decode(s string) (int64, error) {
	var n int64
	for _, c := range []byte(s) {
		i := strings.IndexByte(alphabet, c)
		if i < 0 {
			return 0, fmt.Errorf("unexpected character %c in base62 literal", c)
		}
		n = length*n + int64(i)
	}
	return n, nil
}
