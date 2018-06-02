package serialize

import "io"

// Serializable serialize to or deserialize from byte array
type Serializable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}
