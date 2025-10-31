package contracts

import (
    "bytes"
    "context"
    "encoding/hex"
    "fmt"
    "log"
    "math/big"
    "os"
    "path/filepath"
    "time"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    goeth "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/ethclient"
)

type PrizeDraw struct {
    address common.Address
    abi     abi.ABI
    client  *ethclient.Client
    binder  *bind.BoundContract
    auth    *bind.TransactOpts
}

func NewPrizeDraw(client *ethclient.Client, auth *bind.TransactOpts, addressHex string) (*PrizeDraw, error) {
    addr := common.HexToAddress(addressHex)
    abiData, err := os.ReadFile(filepath.Join("build", "PrizeDraw.abi"))
    if err != nil {
        return nil, fmt.Errorf("read abi: %w", err)
    }
    abiParsed, err := abi.JSON(bytes.NewReader(abiData))
    if err != nil {
        return nil, fmt.Errorf("parse abi: %w", err)
    }
    binder := bind.NewBoundContract(addr, abiParsed, client, client, client)
    return &PrizeDraw{address: addr, abi: abiParsed, client: client, binder: binder, auth: auth}, nil
}

func (p *PrizeDraw) Start(ctx context.Context) (*types.Receipt, error) {
    tx, err := p.binder.Transact(p.auth, "start")
    if err != nil {
        return nil, err
    }
    return bind.WaitMined(ctx, p.client, tx)
}

// StartPrizeResultListener subscribes to PrizeResult events from latest block
func (p *PrizeDraw) StartPrizeResultListener(ctx context.Context) error {
    // There is no generated Filter method; use raw log polling with topic
    // Instead, we poll using FilterLogs in a simple loop from latest
    // Here we poll every 3s
    last := uint64(0)
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            iterStart := last + 1
            // topic = keccak256("PrizeResult(uint256,bool)")
            topicHex := "0xfa0974f074f651f1ec33d39cf7d160a25b998d5aee785ba1fc843b33125c4dde"
            logs, err := p.client.FilterLogs(ctx, goeth.FilterQuery{
                FromBlock: big.NewInt(int64(iterStart)),
                ToBlock:   nil, // latest
                Addresses: []common.Address{p.address},
                Topics:    [][]common.Hash{{common.HexToHash(topicHex)}},
            })
            if err != nil {
                continue
            }
            for _, l := range logs {
                // decode event
                var event struct {
                    PrizeNumber *big.Int
                    IsWin       bool
                }
                if len(l.Data) > 0 {
                    if err := p.abi.UnpackIntoInterface(&event, "PrizeResult", l.Data); err == nil {
                        log.Printf("PrizeResult prizeNumber=%s isWin=%v", event.PrizeNumber.String(), event.IsWin)
                    }
                }
                if l.BlockNumber > last {
                    last = l.BlockNumber
                }
            }
        }
    }
}

func DeployPrizeDraw(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts) (common.Address, error) {
    abiData, err := os.ReadFile(filepath.Join("build", "PrizeDraw.abi"))
    if err != nil {
        return common.Address{}, fmt.Errorf("read abi: %w", err)
    }
    binData, err := os.ReadFile(filepath.Join("build", "PrizeDraw.bin"))
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
    if _, err := bind.WaitMined(ctx, client, tx); err != nil {
        return common.Address{}, err
    }
    return addr, nil
}


