package contracts

import (
    "bytes"
    "context"
    "encoding/hex"
    "fmt"
    "math/big"
    "os"
    "path/filepath"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type SimpleStorage struct {
	address common.Address
	abi     abi.ABI
    client  *ethclient.Client
	binder  *bind.BoundContract
	auth    *bind.TransactOpts
}

func NewSimpleStorage(client *ethclient.Client, auth *bind.TransactOpts, addressHex string) (*SimpleStorage, error) {
	addr := common.HexToAddress(addressHex)
	abiData, err := os.ReadFile(filepath.Join("build", "SimpleStorage.abi"))
	if err != nil {
		return nil, fmt.Errorf("read abi: %w", err)
	}
    abiParsed, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("parse abi: %w", err)
	}
    binder := bind.NewBoundContract(addr, abiParsed, client, client, client)
	return &SimpleStorage{address: addr, abi: abiParsed, client: client, binder: binder, auth: auth}, nil
}

func DeploySimpleStorage(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts) (common.Address, error) {
	abiData, err := os.ReadFile(filepath.Join("build", "SimpleStorage.abi"))
	if err != nil {
		return common.Address{}, fmt.Errorf("read abi: %w", err)
	}
	binData, err := os.ReadFile(filepath.Join("build", "SimpleStorage.bin"))
	if err != nil {
		return common.Address{}, fmt.Errorf("read bin: %w", err)
	}
    abiParsed, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return common.Address{}, fmt.Errorf("parse abi: %w", err)
	}
	bytecode, err := hex.DecodeString(string(bytesTrimSpace(binData)))
	if err != nil {
		return common.Address{}, fmt.Errorf("decode bin: %w", err)
	}
	addr, tx, _, err := bind.DeployContract(auth, abiParsed, bytecode, client)
	if err != nil {
		return common.Address{}, fmt.Errorf("deploy: %w", err)
	}
	// wait for receipt
    _, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("wait mined: %w", err)
	}
	return addr, nil
}

func (s *SimpleStorage) GetValue(ctx context.Context) (*big.Int, error) {
	var out []interface{}
	if err := s.binder.Call(&bind.CallOpts{Context: ctx}, &out, "getValue"); err != nil {
		return nil, err
	}
	if len(out) != 1 {
		return nil, fmt.Errorf("unexpected outputs")
	}
	val, ok := out[0].(*big.Int)
	if !ok {
		// web3 abi may return big.Int by value
		if v, ok2 := out[0].(big.Int); ok2 {
			return &v, nil
		}
		return nil, fmt.Errorf("type assertion failed")
	}
	return val, nil
}

func (s *SimpleStorage) SetValue(ctx context.Context, valueDecimal string) (*types.Receipt, error) {
	val, ok := new(big.Int).SetString(valueDecimal, 10)
	if !ok {
		return nil, fmt.Errorf("invalid decimal value")
	}
	tx, err := s.binder.Transact(s.auth, "setValue", val)
	if err != nil {
		return nil, err
	}
    return bind.WaitMined(ctx, s.client, tx)
}

func bytesTrimSpace(b []byte) []byte {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\n' || b[i] == '\r' || b[i] == '\t') {
		i++
	}
	j := len(b)
	for j > i && (b[j-1] == ' ' || b[j-1] == '\n' || b[j-1] == '\r' || b[j-1] == '\t') {
		j--
	}
	return b[i:j]
}


