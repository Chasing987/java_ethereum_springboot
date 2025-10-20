package com.example.ethereum.service;

import com.example.ethereum.config.EthereumConfig;
import com.example.ethereum.contract.generated.PrizeDraw;
import io.reactivex.disposables.Disposable;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.tx.gas.ContractGasProvider;

import javax.annotation.PostConstruct;

@Service
@RequiredArgsConstructor
@Slf4j
public class PrizeDrawListener {

    private final Web3j web3j;
    private final Credentials credentials;
    private final ContractGasProvider gasProvider;
    private final EthereumConfig ethereumConfig;

    private PrizeDraw contract;

    private Disposable eventSubscription;

    @PostConstruct
    public void init() {
        contract = PrizeDraw.load(
                ethereumConfig.getContractAddress(),
                web3j,
                credentials,
                gasProvider
        );

        eventSubscription = contract.prizeResultEventFlowable(DefaultBlockParameter.valueOf("latest"),
                DefaultBlockParameter.valueOf("latest")).subscribe(event -> {
           // 兑奖 统计开奖信息

           log.info("prizeNumber：{}", event.prizeNumber);
           log.info("isWin:{}", event.isWin);
        });
    }
}
