package modes

import (
	"context"
	"fmt"
)

type Web3Signer struct {
	ctx context.Context
}

func (s *Web3Signer) Run() error {
	fmt.Println("running w3signer ...")
	return nil
}
