package payload

import (
	"bytes"
	"testing"

	"github.com/mileschao/echain/smartcontract/types"
)

func TestDeployCodeSerialize(t *testing.T) {
	var vmcode = types.VMCode{
		VMType: types.Native,
		Code:   []byte{0xFF},
	}
	var dc = &DeployCode{
		Code:        vmcode,
		NeedStorage: false,
		Name:        "name",
		Version:     "version",
		Author:      "author",
		Email:       "email",
		Description: "description",
	}
	buf := new(bytes.Buffer)
	if err := dc.Serialize(buf); err != nil {
		t.Errorf("deploy code serialize")
	}
	var dc2 DeployCode
	if err := dc2.Deserialize(buf); err != nil {
		t.Errorf("deploy code deserialize")
	}
	if dc.Email != dc2.Email {
		t.Errorf("deploy code mail")
	}
}
