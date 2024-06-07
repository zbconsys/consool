package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Eth struct {
	ctx context.Context
	cl  *ethclient.Client
}

func NewEth(ctx context.Context, jsonRPCEndpoint string) (*Eth, error) {
	cl, err := ethclient.Dial(jsonRPCEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not connect to eth client: %w", err)
	}

	return &Eth{
		ctx: ctx,
		cl:  cl,
	}, nil
}

func (e *Eth) GetAccountBalance(account common.Address) (*big.Int, error) {
	bal, err := e.cl.BalanceAt(e.ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get balance of account %s: %w", account.Hex(), err)
	}

	return bal, nil
}
