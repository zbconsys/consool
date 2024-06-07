package app

import (
	"context"
	"errors"
	"github.com/zbconsys/consool/internal/eth"
	"github.com/zbconsys/consool/internal/modes"
	"github.com/zbconsys/consool/pkg/config"
)

var (
	ErrModeNotSupported = errors.New("mode not supported")
)

type ModeType string
type ModeRunner func() error

const (
	w3sAddresses ModeType = "w3s-addresses"
	w3sSendTx    ModeType = "w3s-send"
)

type App struct {
	ctx            context.Context
	eth            *eth.Eth
	conf           *config.Config
	modes          *modes.Modes
	availableModes map[ModeType]ModeRunner
}

func New() (*App, error) {
	ctx := context.Background()
	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	ether, err := eth.NewEth(ctx, conf.JsonRpcURL)
	if err != nil {
		return nil, err
	}

	mds, err := modes.NewModes(ctx, conf, ether)
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:            ctx,
		eth:            ether,
		conf:           conf,
		modes:          mds,
		availableModes: make(map[ModeType]ModeRunner),
	}, nil
}

func (a *App) Run() error {
	if err := a.initModes(); err != nil {
		return err
	}

	modeHandler, ok := a.availableModes[ModeType(a.conf.Mode)]
	if !ok {
		return ErrModeNotSupported
	}

	return modeHandler()
}

func (a *App) initModes() error {
	a.availableModes = make(map[ModeType]ModeRunner)
	// TODO: fix mode initialization override
	a.availableModes[w3sAddresses] = a.modes.FetchWeb3AddressesFunds()
	a.availableModes[w3sSendTx] = a.modes.SendFundsFromAccountInWeb3Signer(a.conf.PublicKey, a.conf.SendToAddress)

	return nil
}
