package modes

import (
	"context"
)

type Modes struct {
	ctx      context.Context
	w3Signer *Web3Signer
}

func NewModes(ctx context.Context) (*Modes, error) {
	return &Modes{
		ctx:      ctx,
		w3Signer: &Web3Signer{ctx: ctx},
	}, nil
}

func (m *Modes) W3Signer() *Web3Signer {
	return m.w3Signer
}
