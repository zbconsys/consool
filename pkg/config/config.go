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
}

type flagsRawData struct {
	jsonRpcURL    string
	web3SignerURL string
	mode          string
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
	flag.Parse()

	if raw.mode == "" {
		return ErrModeUndefined
	}

	if raw.mode == "json-rpc" && raw.jsonRpcURL == "" {
		return ErrJsonRPCUndefined
	}

	conf.JsonRpcURL = raw.jsonRpcURL
	conf.Web3SignerURL = raw.web3SignerURL
	conf.Mode = raw.mode

	return nil
}
