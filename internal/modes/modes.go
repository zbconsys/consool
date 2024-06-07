package modes

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zbconsys/consool/internal/eth"
	"github.com/zbconsys/consool/pkg/config"
	"github.com/zbconsys/consool/pkg/tools"
	"io"
	"math/big"
	"net/http"
)

type Modes struct {
	ctx      context.Context
	ethCl    *eth.Eth
	w3Signer *Web3Signer

	runnerFunc func() error
}

func NewModes(ctx context.Context, conf *config.Config, ethCl *eth.Eth) (*Modes, error) {
	return &Modes{
		ctx:   ctx,
		ethCl: ethCl,
		w3Signer: &Web3Signer{
			ctx:           ctx,
			web3signerURL: conf.Web3SignerURL,
		},
	}, nil
}

func (m *Modes) Run() error {
	return m.runnerFunc()
}

func (m *Modes) FetchWeb3AddressesFunds() func() error {
	return func() error {
		addrs, err := m.w3Signer.GetAddressesAndPublicKeys()
		if err != nil {
			return err
		}

		for pubKey, addr := range addrs {
			bal, err := m.ethCl.GetAccountBalance(addr)
			if err != nil {
				return err
			}

			nonce, err := m.ethCl.GetAccountNonce(addr)
			if err != nil {
				return err
			}

			//TODO: add central output handler
			fmt.Printf("ACC: %s BAL: %s NONCE: %d PUBLIC_KEY: %s\n", addr, bal, nonce, pubKey)
		}

		return nil
	}
}

func (m *Modes) SendFundsFromAccountInWeb3Signer(publicKey string, toAccount string) func() error {
	return func() error {
		addr, err := tools.GetAddressFromPublicKey(publicKey)
		if err != nil {
			return err
		}

		nonce, err := m.ethCl.GetAccountNonce(addr)
		if err != nil {
			return err
		}

		gas, err := m.ethCl.SuggestGasPrice()
		if err != nil {
			return err
		}

		chainID, err := m.ethCl.ChainID()
		if err != nil {
			return err
		}

		toAddr := common.HexToAddress(toAccount)

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gas,
			Gas:      uint64(21000),
			To:       &toAddr,
			Value:    big.NewInt(5000000000000000000), // 5 ETH
			Data:     []byte{},
		})

		var (
			signer       = types.NewEIP155Signer(chainID)
			dataToEncode = []interface{}{
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				chainID, uint(0), uint(0),
			}
		)

		rawTx, err := rlp.EncodeToBytes(dataToEncode)
		if err != nil {
			return fmt.Errorf("could not rlp encode tx to bytes: %w", err)
		}

		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/eth1/sign/%s", m.w3Signer.web3signerURL, publicKey),
			"application/json",
			bytes.NewBuffer([]byte(fmt.Sprintf("{\"data\":\"%s\"}", hexutil.Encode(rawTx)))),
		)
		if err != nil {
			return fmt.Errorf("web3signer error: %w", err)
		}

		defer resp.Body.Close()

		rawSig, _ := io.ReadAll(resp.Body)

		if bytes.Contains(rawSig, []byte("error")) || bytes.Contains(rawSig, []byte("Error")) || bytes.Contains(rawSig, []byte("Resource not found")) {
			return fmt.Errorf("%s", string(rawSig))
		}

		sig, err := hexutil.Decode(string(rawSig))
		if err != nil {
			return fmt.Errorf("could not decode raw web3signer signature: %w", err)
		}

		// fix legacy V
		if sig[64] == 28 || sig[64] == 27 {
			sig[64] -= 27
		}

		signedTx, err := tx.WithSignature(signer, sig)
		_, err = signer.Sender(signedTx)
		if err != nil {
			return err
		}

		//TODO: proper output
		fmt.Println(signedTx.Hash())

		return m.ethCl.SendRawTransaction(signedTx)

	}
}
