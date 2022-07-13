package encoding

type Encoder interface {
	Encode(n int64) string
	Decode(s string) (int64, error)
}
