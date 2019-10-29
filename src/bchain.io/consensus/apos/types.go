package apos

import "bchain.io/common/types"

//go:generate msgp

type CommonCoinMinHash struct {
	H types.Hash
	J int
}

type BlockCertificate []*CredentialSign
