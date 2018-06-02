package serialize

import (
	"bytes"
	"testing"
)

func TestVarBytesSerialize(t *testing.T) {
	b := new(bytes.Buffer)
	bts := make([]byte, 4)
	bts[0] = 't'
	bts[1] = 'e'
	bts[2] = 's'
	bts[3] = 't'
	var varbt = VarBytes{Len: uint64(len(bts)), Bytes: bts}
	if err := varbt.Serialize(b); err != nil {
		t.Errorf("varbytes: %s", err)
	}
	bs := b.Bytes()
	if uint8(bs[0]) != 0x04 ||
		bs[1] != 't' ||
		bs[2] != 'e' ||
		bs[3] != 's' ||
		bs[4] != 't' {
		t.Errorf("varbytes: 0-%X, 1-%c, 2-%c, 3-%c, 4-%c", uint8(bs[0]), bs[1], bs[2], bs[3], bs[4])
	}
}

func TestVarBytesDeserialize(t *testing.T) {
	b := new(bytes.Buffer)
	bts := make([]byte, 4)
	bts[0] = 't'
	bts[1] = 'e'
	bts[2] = 's'
	bts[3] = 't'
	var varbt = VarBytes{Len: uint64(len(bts)), Bytes: bts}
	if err := varbt.Serialize(b); err != nil {
		t.Errorf("varbytes: %s", err)
	}

	var vardbt VarBytes
	if err := vardbt.Deserialize(b); err != nil {
		t.Errorf("varbytes; %s", err)
	}
	if vardbt.Len != 0x04 ||
		vardbt.Bytes[0] != 't' ||
		vardbt.Bytes[1] != 'e' ||
		vardbt.Bytes[2] != 's' ||
		vardbt.Bytes[3] != 't' {
		t.Errorf("varbytes: 0-%X, 1-%c, 2-%c, 3-%c, 4-%c", vardbt.Len, vardbt.Bytes[0], vardbt.Bytes[1], vardbt.Bytes[2], vardbt.Bytes[3])
	}

}
