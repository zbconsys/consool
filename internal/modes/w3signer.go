package modes

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zbconsys/consool/pkg/tools"
)

type Web3Signer struct {
	ctx           context.Context
	web3signerURL string
}

func (s *Web3Signer) GetAddressesAndPublicKeys() (map[string]common.Address, error) {
	resp := map[string]common.Address{}

	pubKeys, err := tools.FetchPublicKeysFromWeb3Signer(s.web3signerURL)
	if err != nil {
		return nil, fmt.Errorf("could not fetch public keys from web3signer: %w", err)
	}

	for _, pubKey := range pubKeys {
		addr, err := tools.GetAddressFromPublicKey(pubKey)
		if err != nil {
			return nil, err
		}

		resp[pubKey] = addr
	}

	return resp, nil
}
