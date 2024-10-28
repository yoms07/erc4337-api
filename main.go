package main

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"web3-account-abstraction-api/generated/abi/accountfactory"
	"web3-account-abstraction-api/generated/abi/entrypoint"
	"web3-account-abstraction-api/generated/abi/paymaster"
	"web3-account-abstraction-api/internal/api"
	"web3-account-abstraction-api/internal/bundler"
	contract "web3-account-abstraction-api/internal/contracts"
	sqlite_store "web3-account-abstraction-api/internal/store/sqlite"
	"web3-account-abstraction-api/internal/usecase"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Config struct {
	RPCURL                    string `mapstructure:"RPC_URL"`
	APIKey                    string `mapstructure:"API_KEY"`
	UserPrivateKey            string `mapstructure:"PRIVATE_KEY"`
	PaymasterSignerPrivateKey string `mapstructure:"PAYMASTER_SIGNER_PRIVATE_KEY"`
	EntryPointAddress         string `mapstructure:"ENTRYPOINT_ADDRESS"`
	AccountFactoryAddress     string `mapstructure:"ACCOUNT_FACTORY_ADDRESS"`
	PaymasterAddress          string `mapstructure:"PAYMASTER_ADDRESS"`
	DBPath                    string `mapstructure:"SQLITE_DB_PATH"`
}

func LoadConfig(path string, env string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(env)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	return config, err
}

func parsePrivateKey(key string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, nil, fmt.Errorf("parsePrivateKey: %+w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("error casting public key to ECDSA")
	}

	return privateKey, publicKeyECDSA, nil
}

func main() {
	config, err := LoadConfig(".", "testnet")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, publicKey, err := parsePrivateKey(config.UserPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	paymasterPrivateKey, paymasterPublicKey, err := parsePrivateKey(config.PaymasterSignerPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	afAddress := common.HexToAddress(config.AccountFactoryAddress)
	epAddress := common.HexToAddress(config.EntryPointAddress)
	pmAddress := common.HexToAddress(config.PaymasterAddress)

	af, err := accountfactory.NewAccountFactory(common.HexToAddress(config.AccountFactoryAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	ep, err := entrypoint.NewEntryPoint(common.HexToAddress(config.EntryPointAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	pm, err := paymaster.NewPaymaster(common.HexToAddress(config.PaymasterAddress), client)
	if err != nil {
		log.Fatal(err)
	}

	contracts := contract.Contracts{
		AccountFactory: af,
		EntryPoint:     ep,
		Paymaster:      pm,

		EntryPointAddress:     epAddress,
		PaymasterAddress:      pmAddress,
		AccountFactoryAddress: afAddress,
	}

	bundler := bundler.NewV07Bundler(client, epAddress)

	contracts.SetChainId(chainId)
	contracts.SetPublicAndPrivateKey(publicKey, privateKey)
	contracts.SetRPCClient(client)
	contracts.SetPaymasterSignerPublicAndPrivateKey(paymasterPublicKey, paymasterPrivateKey)
	contracts.SetPaymasterOwnerPublicAndPrivateKey(publicKey, privateKey)

	u := usecase.NewUseCase(contracts, bundler, client)

	e := echo.New()

	api.SetupAPI(e, sqlite_store.NewStore(db), u, contracts)

	e.Logger.Fatal(e.Start(":8080"))
}
