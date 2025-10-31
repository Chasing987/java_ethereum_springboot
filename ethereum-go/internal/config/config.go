package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	RPCURL          string
	ChainID         int64
	PrivateKeyHex   string
	GasPriceWei     uint64
	GasLimit        uint64
	ContractAddress string
}

func Load(path string) (*Config, error) {
	props, err := readProperties(path)
	if err != nil {
		return nil, err
	}
	chainID, _ := strconv.ParseInt(props["ethereum.network.chainId"], 10, 64)
	gasPrice, _ := strconv.ParseUint(props["ethereum.gas-price"], 10, 64)
	gasLimit, _ := strconv.ParseUint(props["ethereum.gas-limit"], 10, 64)
	return &Config{
		RPCURL:          props["ethereum.network.url"],
		ChainID:         chainID,
		PrivateKeyHex:   props["ethereum.account.privateKey"],
		GasPriceWei:     gasPrice,
		GasLimit:        gasLimit,
		ContractAddress: props["ethereum.contract.address"],
	}, nil
}

func readProperties(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open properties: %w", err)
	}
	defer f.Close()
	props := make(map[string]string)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") { // skip comments/empty
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		props[k] = v
	}
	return props, nil
}


