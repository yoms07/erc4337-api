package api

import (
	"fmt"
	"math/big"
	"net/http"
	"web3-account-abstraction-api/generated/abi/account"
	contract "web3-account-abstraction-api/internal/contracts"
	"web3-account-abstraction-api/internal/model"
	"web3-account-abstraction-api/internal/store"
	"web3-account-abstraction-api/internal/usecase"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
)

func handleError(c echo.Context, err error) error {
	return c.String(http.StatusBadRequest, err.Error())
}

func SetupAPI(e *echo.Echo, store store.Store, u usecase.Usecase, contracts contract.Contracts) error {
	e.GET("/wallet", func(c echo.Context) error {
		wallets, err := store.GetAllWallet()
		if err != nil {
			return handleError(c, err)
		}

		return c.JSON(http.StatusOK, wallets)
	})
	e.POST("/wallet", func(c echo.Context) error {
		walletCount, err := store.CountWallet()
		if err != nil {
			return handleError(c, err)
		}

		nextSalt := common.LeftPadBytes(common.FromHex(fmt.Sprintf("%x", walletCount+1)), 32)
		// TODO: adjust can use paymaster or not
		simpleOp := usecase.SimpleUserOperation{
			WalletSalt:    nextSalt,
			CallData:      common.FromHex("0x"),
			Paymaster:     &contracts.PaymasterAddress,
			PaymasterData: common.FromHex("0x"),
		}

		_, err = u.SendUserOperation(simpleOp)
		if err != nil {
			return handleError(c, err)
		}
		addr, _ := contracts.GetSenderAddres([32]byte(nextSalt))

		wallet := model.UserWallet{
			Sender: addr.String(),
		}

		err = store.CreateWallet(wallet)
		if err != nil {
			return handleError(c, err)
		}

		return c.JSON(http.StatusOK, wallet)
	})

	type SendPayload struct {
		CallData string `json:"callData"`
	}
	e.POST("/wallet/:address/send", func(c echo.Context) error {
		s := SendPayload{}
		err := c.Bind(&s)
		if err != nil {
			return handleError(c, err)
		}
		addr := c.Param("address")

		sender := common.HexToAddress(addr)

		simpleOp := usecase.SimpleUserOperation{
			WalletSalt:    nil,
			CallData:      common.FromHex("0x"),
			Paymaster:     &contracts.PaymasterAddress,
			PaymasterData: common.FromHex("0x"),
			Sender:        &sender,
		}

		hash, err := u.SendUserOperation(simpleOp)
		if err != nil {
			return handleError(c, err)
		}

		return c.String(http.StatusOK, hash)
	})
	e.GET("/wallet/tx/:hash/status", func(c echo.Context) error {
		hash := c.Param("hash")
		status, err := u.GetUserOperationStatus(hash)
		if err != nil {
			return handleError(c, err)
		}
		return c.String(http.StatusOK, fmt.Sprintf("%v", status))
	})
	e.GET("/wallet/:wallet/eth/balance", func(c echo.Context) error {
		walletAddress := c.Param("wallet")
		_, err := store.GetWallet(walletAddress)
		if err != nil {
			return handleError(c, err)
		}
		balance, err := contracts.GetETHBalance(common.HexToAddress(walletAddress))
		if err != nil {
			return handleError(c, err)
		}
		return c.String(http.StatusOK, balance.String())
	})
	e.GET("/wallet/:wallet/:tokenAddress/balance", func(c echo.Context) error {
		walletAddress := c.Param("wallet")
		tokenAddress := c.Param("tokenAddress")
		_, err := store.GetWallet(walletAddress)
		if err != nil {
			return handleError(c, err)
		}

		balance, err := contracts.GetERC20Balance(common.HexToAddress(walletAddress), common.HexToAddress(tokenAddress))
		if err != nil {
			handleError(c, err)
		}
		return c.String(http.StatusOK, balance.String())
	})

	type TransferPayload struct {
		Amount int    `json:"amount"`
		To     string `json:"to"`
	}
	e.POST("/wallet/:wallet/eth/transfer", func(c echo.Context) error {
		walletAddress := c.Param("wallet")
		var payload TransferPayload
		err := c.Bind(&payload)
		if err != nil {
			return handleError(c, err)
		}
		_, err = store.GetWallet(walletAddress)
		if err != nil {
			return handleError(c, err)
		}

		abi, _ := account.AccountMetaData.GetAbi()
		callData, err := abi.Pack("withdrawETH", common.HexToAddress(payload.To), big.NewInt(int64(payload.Amount)))
		if err != nil {
			return handleError(c, err)
		}

		simpleOp := usecase.SimpleUserOperation{
			Sender:        (*common.Address)(common.FromHex(walletAddress)),
			CallData:      callData,
			Paymaster:     &contracts.PaymasterAddress,
			PaymasterData: common.FromHex("0x"),
		}

		hash, err := u.SendUserOperation(simpleOp)
		if err != nil {
			return handleError(c, err)
		}
		return c.String(http.StatusOK, hash)
	})

	e.POST("/wallet/:wallet/:tokenAddress/transfer", func(c echo.Context) error {
		walletAddress := c.Param("wallet")
		tokenAddress := c.Param("tokenAddress")
		var payload TransferPayload
		err := c.Bind(&payload)
		if err != nil {
			return handleError(c, err)
		}
		_, err = store.GetWallet(walletAddress)
		if err != nil {
			return handleError(c, err)
		}

		abi, err := account.AccountMetaData.GetAbi()
		if err != nil {
			return handleError(c, err)
		}
		callData, err := abi.Pack("withdrawERC20", common.HexToAddress(tokenAddress), common.HexToAddress(payload.To), big.NewInt(int64(payload.Amount)))
		if err != nil {
			return handleError(c, err)
		}

		simpleOp := usecase.SimpleUserOperation{
			Sender:        (*common.Address)(common.FromHex(walletAddress)),
			CallData:      callData,
			Paymaster:     &contracts.PaymasterAddress,
			PaymasterData: common.FromHex("0x"),
		}

		hash, err := u.SendUserOperation(simpleOp)
		if err != nil {
			return handleError(c, err)
		}
		return c.String(http.StatusOK, hash)
	})

	return nil
}
