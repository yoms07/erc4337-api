package bundler

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"web3-account-abstraction-api/internal/model"
	"web3-account-abstraction-api/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	epv07DummySignature = common.FromHex("0xfffffffffffffffffffffffffffffff0000000000000000000000000000000007aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1c")
)

type Bundler struct {
	version string

	RPCUrl    *url.URL
	client    *ethclient.Client
	epAddress common.Address

	dummySignature []byte
}

type EstimateUserOpResult struct {
	PreVerificationGas            *big.Int `json:"preVerificationGas"`
	CallGasLimit                  *big.Int `json:"callGasLimit"`
	VerificationGasLimit          *big.Int `json:"verificationGasLimit"`
	PaymasterVerificationGasLimit *big.Int `json:"paymasterVerificationGasLimit"`
}

type SendUserOperationResult struct {
	TxHash string
}

func NewV07Bundler(
	client *ethclient.Client,
	epAddress common.Address,
) Bundler {
	return Bundler{
		version:        "v0.7",
		client:         client,
		epAddress:      epAddress,
		dummySignature: epv07DummySignature,
	}
}

func (b *Bundler) SendUserOperation(userOp model.UserOperation) (SendUserOperationResult, error) {
	requestBody := map[string]interface{}{
		"sender":                        userOp.Sender.Hex(),
		"nonce":                         fmt.Sprintf("0x%s", userOp.Nonce.Text(16)),
		"callData":                      fmt.Sprintf("0x%x", userOp.CallData),
		"callGasLimit":                  fmt.Sprintf("0x%x", userOp.CallGasLimit),
		"verificationGasLimit":          fmt.Sprintf("0x%x", userOp.VerificationGasLimit),
		"preVerificationGas":            fmt.Sprintf("0x%x", userOp.PreVerificationGas),
		"maxFeePerGas":                  fmt.Sprintf("0x%x", userOp.MaxFeePerGas),
		"maxPriorityFeePerGas":          fmt.Sprintf("0x%x", userOp.MaxPriorityFeePerGas),
		"paymasterVerificationGasLimit": fmt.Sprintf("0x%x", userOp.PaymasterVerificationGasLimit),
		"paymasterPostOpGasLimit":       fmt.Sprintf("0x%x", userOp.PaymasterPostOpGasLimit),
		"signature":                     fmt.Sprintf("0x%x", userOp.Signature),
		"paymaster":                     userOp.Paymaster.Hex(),
		"paymasterData":                 fmt.Sprintf("0x%x", userOp.PaymasterData),
	}

	if userOp.Factory != nil && !utils.IsZeroAddress(userOp.Factory) {
		requestBody["factory"] = userOp.Factory.Hex()
		requestBody["factoryData"] = fmt.Sprintf("0x%x", userOp.FactoryData)
	}
	var txHash string
	err := b.client.
		Client().
		Call(&txHash, "eth_sendUserOperation", requestBody, b.epAddress.Hex())

	return SendUserOperationResult{
		TxHash: txHash,
	}, err
}

type UserOperationReceipt struct {
	Success bool    `json:"success"`
	Reason  *string `json:"reason"`
}

func (b *Bundler) GetUserOperationReceipt(opHash string) (UserOperationReceipt, error) {
	var result UserOperationReceipt
	err := b.client.
		Client().
		Call(&result, "eth_getUserOperationReceipt", opHash)

	return result, err
}

type EstimateRequest struct {
	Sender      string  `json:"sender"`
	Nonce       string  `json:"nonce"`
	CallData    string  `json:"callData"`
	Signature   string  `json:"signature"`
	Factory     *string `json:"factory,omitempty"`
	FactoryData *string `json:"factoryData,omitempty"`
}

func newString(s string) *string {
	return &s
}

func (b *Bundler) EstimateUserOpGas(userOp model.UserOperation) (EstimateUserOpResult, error) {
	requestBody := EstimateRequest{
		Sender:    userOp.Sender.Hex(),
		Nonce:     fmt.Sprintf("0x%s", userOp.Nonce.Text(16)),
		CallData:  fmt.Sprintf("0x%x", userOp.CallData),
		Signature: fmt.Sprintf("0x%x", b.dummySignature),
	}
	b.client.SyncProgress(context.Background())

	if userOp.Factory != nil && !utils.IsZeroAddress(userOp.Factory) {
		requestBody.Factory = newString(userOp.Factory.Hex())
		requestBody.FactoryData = newString(fmt.Sprintf("0x%x", userOp.FactoryData))
	}

	result := map[string]string{}
	err := b.client.
		Client().
		CallContext(context.Background(),
			&result,
			"eth_estimateUserOperationGas", requestBody, b.epAddress.Hex())

	return EstimateUserOpResult{
		PreVerificationGas:            big.NewInt(0).SetBytes(common.FromHex(result["preVerificationGas"])),
		CallGasLimit:                  big.NewInt(0).SetBytes(common.FromHex(result["callGasLimit"])),
		VerificationGasLimit:          big.NewInt(0).SetBytes(common.FromHex(result["verificationGasLimit"])),
		PaymasterVerificationGasLimit: nil,
	}, err
}

func (b *Bundler) GetMaxPriorityFeePerGas() (*big.Int, error) {
	var result string
	err := b.client.
		Client().
		Call(&result, "rundler_maxPriorityFeePerGas")
	return big.NewInt(0).SetBytes(common.FromHex(result)), err
}
