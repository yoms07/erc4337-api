// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accountfactory

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// AccountFactoryMetaData contains all meta data concerning the AccountFactory contract.
var AccountFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"name\":\"createAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// AccountFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use AccountFactoryMetaData.ABI instead.
var AccountFactoryABI = AccountFactoryMetaData.ABI

// AccountFactory is an auto generated Go binding around an Ethereum contract.
type AccountFactory struct {
	AccountFactoryCaller     // Read-only binding to the contract
	AccountFactoryTransactor // Write-only binding to the contract
	AccountFactoryFilterer   // Log filterer for contract events
}

// AccountFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccountFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccountFactorySession struct {
	Contract     *AccountFactory   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccountFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccountFactoryCallerSession struct {
	Contract *AccountFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AccountFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccountFactoryTransactorSession struct {
	Contract     *AccountFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AccountFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccountFactoryRaw struct {
	Contract *AccountFactory // Generic contract binding to access the raw methods on
}

// AccountFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccountFactoryCallerRaw struct {
	Contract *AccountFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// AccountFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccountFactoryTransactorRaw struct {
	Contract *AccountFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccountFactory creates a new instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactory(address common.Address, backend bind.ContractBackend) (*AccountFactory, error) {
	contract, err := bindAccountFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccountFactory{AccountFactoryCaller: AccountFactoryCaller{contract: contract}, AccountFactoryTransactor: AccountFactoryTransactor{contract: contract}, AccountFactoryFilterer: AccountFactoryFilterer{contract: contract}}, nil
}

// NewAccountFactoryCaller creates a new read-only instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactoryCaller(address common.Address, caller bind.ContractCaller) (*AccountFactoryCaller, error) {
	contract, err := bindAccountFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccountFactoryCaller{contract: contract}, nil
}

// NewAccountFactoryTransactor creates a new write-only instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*AccountFactoryTransactor, error) {
	contract, err := bindAccountFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccountFactoryTransactor{contract: contract}, nil
}

// NewAccountFactoryFilterer creates a new log filterer instance of AccountFactory, bound to a specific deployed contract.
func NewAccountFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*AccountFactoryFilterer, error) {
	contract, err := bindAccountFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccountFactoryFilterer{contract: contract}, nil
}

// bindAccountFactory binds a generic wrapper to an already deployed contract.
func bindAccountFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AccountFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountFactory *AccountFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccountFactory.Contract.AccountFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountFactory *AccountFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountFactory.Contract.AccountFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountFactory *AccountFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountFactory.Contract.AccountFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountFactory *AccountFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccountFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountFactory *AccountFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountFactory *AccountFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountFactory.Contract.contract.Transact(opts, method, params...)
}

// CreateAccount is a paid mutator transaction binding the contract method 0xde781f55.
//
// Solidity: function createAccount(address owner, bytes32 salt, address entryPoint) returns(address)
func (_AccountFactory *AccountFactoryTransactor) CreateAccount(opts *bind.TransactOpts, owner common.Address, salt [32]byte, entryPoint common.Address) (*types.Transaction, error) {
	return _AccountFactory.contract.Transact(opts, "createAccount", owner, salt, entryPoint)
}

// CreateAccount is a paid mutator transaction binding the contract method 0xde781f55.
//
// Solidity: function createAccount(address owner, bytes32 salt, address entryPoint) returns(address)
func (_AccountFactory *AccountFactorySession) CreateAccount(owner common.Address, salt [32]byte, entryPoint common.Address) (*types.Transaction, error) {
	return _AccountFactory.Contract.CreateAccount(&_AccountFactory.TransactOpts, owner, salt, entryPoint)
}

// CreateAccount is a paid mutator transaction binding the contract method 0xde781f55.
//
// Solidity: function createAccount(address owner, bytes32 salt, address entryPoint) returns(address)
func (_AccountFactory *AccountFactoryTransactorSession) CreateAccount(owner common.Address, salt [32]byte, entryPoint common.Address) (*types.Transaction, error) {
	return _AccountFactory.Contract.CreateAccount(&_AccountFactory.TransactOpts, owner, salt, entryPoint)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AccountFactory *AccountFactoryTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountFactory.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AccountFactory *AccountFactorySession) Receive() (*types.Transaction, error) {
	return _AccountFactory.Contract.Receive(&_AccountFactory.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AccountFactory *AccountFactoryTransactorSession) Receive() (*types.Transaction, error) {
	return _AccountFactory.Contract.Receive(&_AccountFactory.TransactOpts)
}
