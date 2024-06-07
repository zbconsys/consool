package tools

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"net/http"
)

func GetAddressHexFromPublicKey(publicKeyHex string) (common.Address, error) {
	publicKeyBytes, err := hexutil.Decode(publicKeyHex)
	if err != nil {
		return common.Address{}, fmt.Errorf("invalid public key: %w", err)
	}

	pubKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to unmarshal public key: %w", err)
	}

	return crypto.PubkeyToAddress(*pubKey), err
}

func CreateNewKeystoreFile(privateKeyHex, password, keystoreDir string) error {
	privateKeyBytes, err := hexutil.Decode("0x" + privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	privateKeyECDSA, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("error generating ECDSA key from private key bytes: %w", err)
	}

	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err = ks.ImportECDSA(privateKeyECDSA, password)
	if err != nil {
		return fmt.Errorf("error importing account to keystore: %w", err)
	}

	return nil
}

func FetchPublicKeysFromWeb3Signer(web3SignerURL string) ([]string, error) {
	var publicKeys []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/eth1/publicKeys", web3SignerURL))
	if err != nil {
		return nil, fmt.Errorf("web3signer http request failed: %w", err)
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(bodyByte, &publicKeys); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return publicKeys, nil
}
