package payload

import (
	"io"

	"github.com/mileschao/echain/smartcontract/types"
)

//InvokeCode invoke code
type InvokeCode struct {
	Code types.VMCode
}

// Serialize implement Payload interface
func (ic *InvokeCode) Serialize(w io.Writer) error {
	return ic.Code.Serialize(w)
}

// Deserialize implement Payload interface
func (ic *InvokeCode) Deserialize(r io.Reader) error {
	return ic.Code.Deserialize(r)
}
