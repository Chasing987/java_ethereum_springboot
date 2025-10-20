package com.example.ethereum.service;

import com.example.ethereum.config.EthereumConfig;
import com.example.ethereum.contract.generated.PrizeDraw;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.gas.ContractGasProvider;

import javax.annotation.PostConstruct;
import java.util.concurrent.CompletableFuture;

@Service
@RequiredArgsConstructor
@Slf4j
public class PrizeDrawService {

    private final Web3j web3j;
    private final Credentials credentials;
    private final ContractGasProvider gasProvider;
    private final EthereumConfig ethereumConfig;

    private PrizeDraw contract;

    @PostConstruct
    private void init() {
        contract = PrizeDraw.load(
                ethereumConfig.getContractAddress(),
                web3j,
                credentials,
                gasProvider

        );
    }

    public CompletableFuture<TransactionReceipt> start() {
        return contract.start().sendAsync().thenApply(
                receipt -> {
                    return receipt;
                }).exceptionally(e -> {
            throw new RuntimeException("error", e);
        });
    }
}
