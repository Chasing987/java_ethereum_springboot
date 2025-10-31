package eth

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"ethereum-go/internal/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	RPC *ethclient.Client
}

func NewClient(cfg *config.Config) (*Client, *bind.TransactOpts, error) {
	rpc, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, nil, fmt.Errorf("dial rpc: %w", err)
	}
	privHex := strings.TrimPrefix(cfg.PrivateKeyHex, "0x")
	priv, err := crypto.HexToECDSA(privHex)
	if err != nil {
		return nil, nil, fmt.Errorf("private key: %w", err)
	}
	chainID := big.NewInt(cfg.ChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainID)
	if err != nil {
		return nil, nil, fmt.Errorf("transactor: %w", err)
	}
	auth.GasPrice = big.NewInt(int64(cfg.GasPriceWei))
	auth.GasLimit = uint64(cfg.GasLimit)
	return &Client{RPC: rpc}, auth, nil
}

func AddressFromHex(addr string) common.Address {
	return common.HexToAddress(addr)
}

func PublicKeyFromPrivate(priv *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(priv.PublicKey)
}


