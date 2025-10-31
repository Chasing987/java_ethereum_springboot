package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ethereum-go/internal/config"
	"ethereum-go/internal/contracts"
	"ethereum-go/internal/eth"
)

func main() {
	// Load config
	rootDir, _ := os.Getwd()
	configPath := filepath.Join(rootDir, "config", "application.properties")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Ethereum client
	ethClient, auth, err := eth.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to init ethereum client: %v", err)
	}

	// Contracts
	ss, err := contracts.NewSimpleStorage(ethClient.RPC, auth, cfg.ContractAddress)
	if err != nil {
		log.Printf("warning: could not load SimpleStorage at %s: %v", cfg.ContractAddress, err)
	}
	prize, err := contracts.NewPrizeDraw(ethClient.RPC, auth, cfg.ContractAddress)
	if err != nil {
		log.Printf("warning: could not load PrizeDraw at %s: %v", cfg.ContractAddress, err)
	}

	// PrizeDraw event listener (non-blocking)
	ctx, cancel := context.WithCancel(context.Background())
	if prize != nil {
		go func() {
			if err := prize.StartPrizeResultListener(ctx); err != nil {
				log.Printf("prize listener stopped: %v", err)
			}
		}()
	}

	// HTTP Handlers
	http.HandleFunc("/api/contract/value", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Optional override via query param ?address=0x...
		addr := strings.TrimSpace(r.URL.Query().Get("address"))
		var useSS *contracts.SimpleStorage
		if addr != "" {
			var err error
			useSS, err = contracts.NewSimpleStorage(ethClient.RPC, auth, addr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid address or load failed"})
				return
			}
		} else {
			useSS = ss
		}

		switch r.Method {
		case http.MethodGet:
			if useSS == nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "contract not loaded"})
				return
			}
			val, err := useSS.GetValue(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"value":         val.String(),
				"valueAsString": val.String(),
			})
        case http.MethodPost:
			if useSS == nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "contract not loaded"})
				return
			}
            // Robust parse: accept {"value":"123"} or {"value":123} or query ?value=123
            var payload struct{ Value string `json:"value"` }
            dec := json.NewDecoder(r.Body)
            if err := dec.Decode(&payload); err != nil || payload.Value == "" {
                // Try numeric value
                r.Body.Close()
                // Re-read body is not trivial; fallback to query param
                qv := strings.TrimSpace(r.URL.Query().Get("value"))
                if qv == "" {
                    w.WriteHeader(http.StatusBadRequest)
                    _ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
                    return
                }
                payload.Value = qv
            }
            receipt, err := useSS.SetValue(r.Context(), payload.Value)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"transactionHash": receipt.TxHash.Hex(),
				"blockNumber":     receipt.BlockNumber.String(),
				"gasUsed":         receipt.GasUsed,
				"status":          receipt.Status,
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/contract/deploy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		address, err := contracts.DeploySimpleStorage(r.Context(), ethClient.RPC, auth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"contractAddress": address.Hex(),
			"message":         "合约部署成功",
		})
	})

	http.HandleFunc("/api/contract/deploy-prizedraw", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		address, err := contracts.DeployPrizeDraw(r.Context(), ethClient.RPC, auth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"contractAddress": address.Hex(),
			"message":         "合约部署成功",
		})
	})

	http.HandleFunc("/api/contract/start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Optional override via query param ?address=0x...
		addr := strings.TrimSpace(r.URL.Query().Get("address"))
		var usePrize *contracts.PrizeDraw
		if addr != "" {
			var err error
			usePrize, err = contracts.NewPrizeDraw(ethClient.RPC, auth, addr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid address or load failed"})
				return
			}
		} else {
			usePrize = prize
		}

		if usePrize == nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "contract not loaded"})
			return
		}
		receipt, err := usePrize.Start(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"transactionHash": receipt.TxHash.Hex(),
			"blockNumber":     receipt.BlockNumber.String(),
			"gasUsed":         receipt.GasUsed,
			"status":          receipt.Status,
		})
	})

	server := &http.Server{Addr: ":8080", ReadHeaderTimeout: 10 * time.Second}
	go func() {
		log.Printf("HTTP server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	_ = server.Shutdown(ctxShutdown)
}


