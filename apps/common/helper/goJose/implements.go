package goJose

import (
	"WenBeego/apps/common/helper"

	"github.com/go-jose/go-jose/v4"
)

var _ jose.NonceSource = (*NonceSourceStruct)(nil)

type NonceSourceStruct struct {
}

func (n *NonceSourceStruct) Nonce() (string, error) {
	return helper.GetUuid()
}
