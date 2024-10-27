package model

import (
	"math/big"
	"web3-account-abstraction-api/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Address = common.Address

type UserOperation struct {
	Sender                        Address
	Nonce                         *big.Int
	Factory                       *Address
	FactoryData                   []byte
	CallData                      []byte
	CallGasLimit                  *big.Int
	VerificationGasLimit          *big.Int
	PreVerificationGas            *big.Int
	MaxFeePerGas                  *big.Int
	MaxPriorityFeePerGas          *big.Int
	Paymaster                     *Address
	PaymasterVerificationGasLimit *big.Int
	PaymasterPostOpGasLimit       *big.Int
	PaymasterData                 []byte
	Signature                     []byte
}

func (u *UserOperation) Pack() PackedUserOperation {
	accountGasLimit := u.packAccountGasLimit(u.VerificationGasLimit, u.CallGasLimit)
	gasFee := u.packAccountGasLimit(u.MaxPriorityFeePerGas, u.MaxFeePerGas)
	paymasterAndData := u.packPaymasterAndPaymasterData()

	return PackedUserOperation{
		Sender:             u.Sender,
		Nonce:              u.Nonce,
		InitCode:           append(u.Factory[:], u.FactoryData...),
		CallData:           u.CallData,
		AccountGasLimits:   accountGasLimit,
		PreVerificationGas: u.PreVerificationGas,
		GasFees:            gasFee,
		PaymasterAndData:   paymasterAndData,
		Signature:          u.Signature,
	}
}

func (u *UserOperation) packAccountGasLimit(
	verificationGasLimit *big.Int,
	callGasLimit *big.Int) [32]byte {

	left := common.LeftPadBytes(verificationGasLimit.Bytes(), 16)
	right := common.LeftPadBytes(callGasLimit.Bytes(), 16)

	return [32]byte(append(left, right...))
}

func (u *UserOperation) packPaymasterAndPaymasterData() []byte {
	if u.Paymaster == nil || utils.IsZeroAddress(u.Paymaster) {
		return []byte{}
	}

	packedPaymasterData := append(
		u.Paymaster.Bytes(),
		append(
			common.LeftPadBytes(u.PaymasterVerificationGasLimit.Bytes(), 16),
			common.LeftPadBytes(u.PaymasterPostOpGasLimit.Bytes(), 16)...,
		)...,
	)

	if len(u.PaymasterData) > 0 {
		packedPaymasterData = append(packedPaymasterData, u.PaymasterData...)
	}

	return packedPaymasterData
}

func (u *UserOperation) EncodePaymasterData(validUntil int, validAfter int, signature []byte) error {
	uint48Type, _ := abi.NewType("uint48", "", nil)
	arguments := abi.Arguments{
		abi.Argument{Type: uint48Type},
		abi.Argument{Type: uint48Type},
	}

	packedOutOfTime, err := arguments.Pack(big.NewInt(int64(validUntil)), big.NewInt(int64(validAfter)))
	if err != nil {
		return err
	}

	u.PaymasterData = append(packedOutOfTime, signature...)
	return nil
}

type PackedUserOperation struct {
	Sender             Address
	Nonce              *big.Int
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas *big.Int
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}
