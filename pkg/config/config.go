package config

import (
	"errors"
	"flag"
)

var (
	ErrJsonRPCUndefined = errors.New("json-rpc url undefined")
	ErrModeUndefined    = errors.New("mode of operation undefined")
)

type Config struct {
	JsonRpcURL    string `yaml:"json_rpc_url"`
	Web3SignerURL string `yaml:"web3_signer_url"`
	Mode          string `yaml:"mode"`
	PublicKey     string `yaml:"public_key"`
	SendToAddress string `yaml:"send_to_address"`
}

type flagsRawData struct {
	jsonRpcURL    string
	web3SignerURL string
	mode          string
	publicKey     string
	sendToAddress string
}

func NewConfig() (*Config, error) {
	c := new(Config)
	if err := processFlags(c); err != nil {
		return nil, err
	}

	return c, nil
}

func processFlags(conf *Config) error {
	raw := new(flagsRawData)
	flag.StringVar(&raw.jsonRpcURL, "json-rpc", "", "JSON-RPC URL")
	flag.StringVar(&raw.web3SignerURL, "web3-signer", "", "Web3 signer URL")
	flag.StringVar(&raw.mode, "mode", "", "mode of operation")
	flag.StringVar(&raw.publicKey, "public-key", "", "public key")
	flag.StringVar(&raw.sendToAddress, "send-to-address", "", "send to address")
	flag.Parse()

	if raw.mode == "" {
		return ErrModeUndefined
	}

	if raw.mode == "json-rpc" && raw.jsonRpcURL == "" {
		return ErrJsonRPCUndefined
	}

	// TODO: add more checks - implement cobra

	conf.JsonRpcURL = raw.jsonRpcURL
	conf.Web3SignerURL = raw.web3SignerURL
	conf.Mode = raw.mode
	conf.SendToAddress = raw.sendToAddress
	conf.PublicKey = raw.publicKey

	return nil
}
