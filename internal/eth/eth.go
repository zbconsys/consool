package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func (e *Eth) GetAccountNonce(account common.Address) (uint64, error) {
	nonce, err := e.cl.PendingNonceAt(e.ctx, account)
	if err != nil {
		return 0, fmt.Errorf("could not get pending nonce of account %s: %w", account.Hex(), err)
	}

	return nonce, nil
}

func (e *Eth) SuggestGasPrice() (*big.Int, error) {
	return e.cl.SuggestGasPrice(e.ctx)
}

func (e *Eth) ChainID() (*big.Int, error) {
	return e.cl.ChainID(e.ctx)
}

func (e *Eth) SendRawTransaction(transaction *types.Transaction) error {
	return e.cl.SendTransaction(e.ctx, transaction)
}
