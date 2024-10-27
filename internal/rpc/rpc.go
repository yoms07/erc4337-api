package rpc

import (
	"context"
	"web3-account-abstraction-api/internal/model"
)

type EstimateUserOpGasResult struct {
}

type RPC interface {
	SendUserOperation(ctx context.Context, userOp model.UserOperation) error
	EstimateUserOpGas(ctx context.Context, userOp model.UserOperation) (EstimateUserOpGasResult, error)
}
