package transaction

import (
	"bytes"
	"testing"
)

func TestTxAttrSerialze(t *testing.T) {
	var url = "http://foo.com/foo/foo?len=x"
	txAttr := NewTxAttribute(DescriptionURL, []byte(url))
	buf := new(bytes.Buffer)
	if err := txAttr.Serialize(buf); err != nil {
		t.Errorf("tx attribute serialize: %s", err)
	}
	var txAttr2 TxAttribute
	if err := txAttr2.Deserialize(buf); err != nil {
		t.Errorf("tx attribute deserialize: %s", err)
	}
}
