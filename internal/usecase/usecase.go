package usecase

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"web3-account-abstraction-api/generated/abi/entrypoint"
	"web3-account-abstraction-api/internal/bundler"
	contract "web3-account-abstraction-api/internal/contracts"
	"web3-account-abstraction-api/internal/model"
	"web3-account-abstraction-api/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vipnode/ether"
)

type SimpleUserOperation struct {
	WalletSalt    []byte
	CallData      []byte
	Paymaster     *common.Address
	PaymasterData []byte
	Sender        *common.Address
}

type Usecase struct {
	contracts contract.Contracts
	bundler   bundler.Bundler
	client    *ethclient.Client

	initialETH *big.Int
}

func NewUseCase(contracts contract.Contracts, bundler bundler.Bundler, client *ethclient.Client) Usecase {
	initialETH, _ := ether.Parse("1 ether")
	return Usecase{
		contracts:  contracts,
		bundler:    bundler,
		client:     client,
		initialETH: initialETH,
	}
}

func markUpGas(gas *big.Int, percentage float64) *big.Int {
	fl, _ := gas.Float64()
	gasAfterMarkup := fl * percentage
	return big.NewInt(int64(math.Floor(gasAfterMarkup)))
}

func (u *Usecase) SendUserOperation(simpleOp SimpleUserOperation) (string, error) {
	var sender common.Address
	var err error
	if simpleOp.Sender != nil {
		sender = *simpleOp.Sender
	} else {
		sender, err = u.contracts.GetSenderAddres([32]byte(simpleOp.WalletSalt))
		if err != nil {
			return "", err
		}
	}

	factory := common.HexToAddress("0x")
	factoryData := []byte{}

	contractCode, err := u.client.CodeAt(context.Background(), sender, nil)
	if err != nil {
		return "", err
	}

	fmt.Printf("sender: %s\n", sender.Hex())
	if len(contractCode) == 0 {
		factory = u.contracts.AccountFactoryAddress
		factoryData, err = u.contracts.GetAccountFactoryCallData(
			u.contracts.OwnerAddress(),
			[32]byte(simpleOp.WalletSalt),
			u.contracts.EntryPointAddress)
		if err != nil {
			return "", err
		}
		err = u.contracts.MintToken(sender, u.initialETH)
		if err != nil {
			return "", err
		}
	}

	nonce, err := u.contracts.EntryPoint.GetNonce(&bind.CallOpts{
		Pending: false,
	}, sender, big.NewInt(0))
	if err != nil {
		return "", err
	}
	maxFeePerGas, _ := ether.Parse("10 gwei")
	maxPriorityFeePerGas, _ := ether.Parse("5 gwei")

	userOp := model.UserOperation{
		Sender:        sender,
		Nonce:         nonce,
		Factory:       &factory,
		FactoryData:   factoryData,
		CallData:      simpleOp.CallData,
		Paymaster:     simpleOp.Paymaster,
		PaymasterData: simpleOp.PaymasterData,
		Signature:     []byte{},

		// Dummy value
		CallGasLimit:                  big.NewInt(400_000),
		VerificationGasLimit:          big.NewInt(400_000),
		PaymasterVerificationGasLimit: big.NewInt(1_000_000),
		PreVerificationGas:            big.NewInt(200_000),
		MaxFeePerGas:                  maxFeePerGas,
		MaxPriorityFeePerGas:          maxPriorityFeePerGas,
		PaymasterPostOpGasLimit:       big.NewInt(21_000),
	}

	estimateGasResult, err := u.bundler.EstimateUserOpGas(userOp)
	if err != nil {
		return "", err
	}

	userOp.PreVerificationGas = markUpGas(estimateGasResult.PreVerificationGas, 1.1)
	userOp.CallGasLimit = markUpGas(estimateGasResult.CallGasLimit, 1.1)
	userOp.VerificationGasLimit = markUpGas(estimateGasResult.VerificationGasLimit, 1.1)
	maxPriorityGasFee, err := u.bundler.GetMaxPriorityFeePerGas()
	if err != nil {
		return "", err
	}
	userOp.MaxPriorityFeePerGas = maxPriorityGasFee

	baseFee, err := u.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	userOp.MaxFeePerGas = baseFee.Add(baseFee, maxPriorityGasFee)

	if simpleOp.Paymaster != nil && !utils.IsZeroAddress(simpleOp.Paymaster) {
		pmSignature, validAfter, validUntil, err := u.contracts.GetPaymasterSignature(userOp)
		if err != nil {
			return "", err
		}
		userOp.EncodePaymasterData(int(validUntil), int(validAfter), pmSignature)
	}

	userOpHash, err := u.contracts.EntryPoint.GetUserOpHash(&bind.CallOpts{Pending: false}, entrypoint.PackedUserOperation(userOp.Pack()))
	if err != nil {
		return "", err
	}

	signature, err := u.contracts.Sign(userOpHash[:])
	if err != nil {
		return "", err
	}
	userOp.Signature = signature

	result, err := u.bundler.SendUserOperation(userOp)
	if err != nil {
		return "", err
	}

	return result.TxHash, nil
}

// only return success/failed
func (u *Usecase) GetUserOperationStatus(hash string) (bool, error) {
	receipt, err := u.bundler.GetUserOperationReceipt(hash)
	if err != nil {
		return false, err
	}

	return receipt.Success, nil
}
