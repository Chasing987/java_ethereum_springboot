package com.example.ethereum.contract.generated;

import io.reactivex.Flowable;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.RemoteFunctionCall;
import org.web3j.protocol.core.methods.request.EthFilter;
import org.web3j.protocol.core.methods.response.BaseEventResponse;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/LFDT-web3j/web3j/tree/main/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 1.7.0.
 */
@SuppressWarnings("rawtypes")
public class PrizeDraw extends Contract {
    public static final String BINARY = "608060405234801561001057600080fd5b5061023c806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063be9a655514610030575b600080fd5b61003861003a565b005b6000600a4243604051602001610051929190610127565b6040516020818303038152906040528051906020012060001c6100749190610182565b905060058111156100be577ffa0974f074f651f1ec33d39cf7d160a25b998d5aee785ba1fc843b33125c4dde8160016040516100b19291906101dd565b60405180910390a16100f9565b7ffa0974f074f651f1ec33d39cf7d160a25b998d5aee785ba1fc843b33125c4dde8160006040516100f09291906101dd565b60405180910390a15b50565b6000819050919050565b6000819050919050565b61012161011c826100fc565b610106565b82525050565b60006101338285610110565b6020820191506101438284610110565b6020820191508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061018d826100fc565b9150610198836100fc565b9250826101a8576101a7610153565b5b828206905092915050565b6101bc816100fc565b82525050565b60008115159050919050565b6101d7816101c2565b82525050565b60006040820190506101f260008301856101b3565b6101ff60208301846101ce565b939250505056fea264697066735822122059b67fc96af86ca9d07664e7232db170bd9b548e769fe53df5483be27941261a64736f6c63430008130033";

    private static String librariesLinkedBinary;

    public static final String FUNC_START = "start";

    public static final Event PRIZERESULT_EVENT = new Event("PrizeResult", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Bool>() {}));
    ;

    @Deprecated
    protected PrizeDraw(String contractAddress, Web3j web3j, Credentials credentials,
            BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected PrizeDraw(String contractAddress, Web3j web3j, Credentials credentials,
            ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected PrizeDraw(String contractAddress, Web3j web3j, TransactionManager transactionManager,
            BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected PrizeDraw(String contractAddress, Web3j web3j, TransactionManager transactionManager,
            ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static List<PrizeResultEventResponse> getPrizeResultEvents(
            TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = staticExtractEventParametersWithLog(PRIZERESULT_EVENT, transactionReceipt);
        ArrayList<PrizeResultEventResponse> responses = new ArrayList<PrizeResultEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            PrizeResultEventResponse typedResponse = new PrizeResultEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.prizeNumber = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse.isWin = (Boolean) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public static PrizeResultEventResponse getPrizeResultEventFromLog(Log log) {
        Contract.EventValuesWithLog eventValues = staticExtractEventParametersWithLog(PRIZERESULT_EVENT, log);
        PrizeResultEventResponse typedResponse = new PrizeResultEventResponse();
        typedResponse.log = log;
        typedResponse.prizeNumber = (BigInteger) eventValues.getNonIndexedValues().get(0).getValue();
        typedResponse.isWin = (Boolean) eventValues.getNonIndexedValues().get(1).getValue();
        return typedResponse;
    }

    public Flowable<PrizeResultEventResponse> prizeResultEventFlowable(EthFilter filter) {
        return web3j.ethLogFlowable(filter).map(log -> getPrizeResultEventFromLog(log));
    }

    public Flowable<PrizeResultEventResponse> prizeResultEventFlowable(
            DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        EthFilter filter = new EthFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(PRIZERESULT_EVENT));
        return prizeResultEventFlowable(filter);
    }

    public RemoteFunctionCall<TransactionReceipt> start() {
        final Function function = new Function(
                FUNC_START, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    @Deprecated
    public static PrizeDraw load(String contractAddress, Web3j web3j, Credentials credentials,
            BigInteger gasPrice, BigInteger gasLimit) {
        return new PrizeDraw(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static PrizeDraw load(String contractAddress, Web3j web3j,
            TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new PrizeDraw(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static PrizeDraw load(String contractAddress, Web3j web3j, Credentials credentials,
            ContractGasProvider contractGasProvider) {
        return new PrizeDraw(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static PrizeDraw load(String contractAddress, Web3j web3j,
            TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        return new PrizeDraw(contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static RemoteCall<PrizeDraw> deploy(Web3j web3j, Credentials credentials,
            ContractGasProvider contractGasProvider) {
        return deployRemoteCall(PrizeDraw.class, web3j, credentials, contractGasProvider, getDeploymentBinary(), "");
    }

    public static RemoteCall<PrizeDraw> deploy(Web3j web3j, TransactionManager transactionManager,
            ContractGasProvider contractGasProvider) {
        return deployRemoteCall(PrizeDraw.class, web3j, transactionManager, contractGasProvider, getDeploymentBinary(), "");
    }

    @Deprecated
    public static RemoteCall<PrizeDraw> deploy(Web3j web3j, Credentials credentials,
            BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(PrizeDraw.class, web3j, credentials, gasPrice, gasLimit, getDeploymentBinary(), "");
    }

    @Deprecated
    public static RemoteCall<PrizeDraw> deploy(Web3j web3j, TransactionManager transactionManager,
            BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(PrizeDraw.class, web3j, transactionManager, gasPrice, gasLimit, getDeploymentBinary(), "");
    }


    private static String getDeploymentBinary() {
        if (librariesLinkedBinary != null) {
            return librariesLinkedBinary;
        } else {
            return BINARY;
        }
    }

    public static class PrizeResultEventResponse extends BaseEventResponse {
        public BigInteger prizeNumber;

        public Boolean isWin;
    }
}
