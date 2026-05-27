// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package priceaggregator

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

// PriceAggregatorMetaData contains all meta data concerning the PriceAggregator contract.
var PriceAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"contractIReporterSet\",\"name\":\"reporterSet_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"assetId_\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description_\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"version_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestFee_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"sent\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"required\",\"type\":\"uint256\"}],\"name\":\"InsufficientFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoRoundData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RefundFailed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"}],\"name\":\"ReqIdAlreadyFulfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"submittedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"latestStartedAt\",\"type\":\"uint256\"}],\"name\":\"StaleTimestamp\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"submittedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentMaxAge\",\"type\":\"uint256\"}],\"name\":\"SubmissionTooOld\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroReporterSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMaxAge\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMaxAge\",\"type\":\"uint256\"}],\"name\":\"MaxAgeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"PriceFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"PriceRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldReporterSet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newReporterSet\",\"type\":\"address\"}],\"name\":\"ReporterSetChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"RequestFeeChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"assetId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"fulfillPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"fulfilled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId_\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundId\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextReqId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reporterSet\",\"outputs\":[{\"internalType\":\"contractIReporterSet\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMaxAge\",\"type\":\"uint256\"}],\"name\":\"setMaxAge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIReporterSet\",\"name\":\"newReporterSet\",\"type\":\"address\"}],\"name\":\"setReporterSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"setRequestFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PriceAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use PriceAggregatorMetaData.ABI instead.
var PriceAggregatorABI = PriceAggregatorMetaData.ABI

// PriceAggregator is an auto generated Go binding around an Ethereum contract.
type PriceAggregator struct {
	PriceAggregatorCaller     // Read-only binding to the contract
	PriceAggregatorTransactor // Write-only binding to the contract
	PriceAggregatorFilterer   // Log filterer for contract events
}

// PriceAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type PriceAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PriceAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PriceAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PriceAggregatorSession struct {
	Contract     *PriceAggregator  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PriceAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PriceAggregatorCallerSession struct {
	Contract *PriceAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// PriceAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PriceAggregatorTransactorSession struct {
	Contract     *PriceAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// PriceAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type PriceAggregatorRaw struct {
	Contract *PriceAggregator // Generic contract binding to access the raw methods on
}

// PriceAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PriceAggregatorCallerRaw struct {
	Contract *PriceAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// PriceAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PriceAggregatorTransactorRaw struct {
	Contract *PriceAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPriceAggregator creates a new instance of PriceAggregator, bound to a specific deployed contract.
func NewPriceAggregator(address common.Address, backend bind.ContractBackend) (*PriceAggregator, error) {
	contract, err := bindPriceAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PriceAggregator{PriceAggregatorCaller: PriceAggregatorCaller{contract: contract}, PriceAggregatorTransactor: PriceAggregatorTransactor{contract: contract}, PriceAggregatorFilterer: PriceAggregatorFilterer{contract: contract}}, nil
}

// NewPriceAggregatorCaller creates a new read-only instance of PriceAggregator, bound to a specific deployed contract.
func NewPriceAggregatorCaller(address common.Address, caller bind.ContractCaller) (*PriceAggregatorCaller, error) {
	contract, err := bindPriceAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorCaller{contract: contract}, nil
}

// NewPriceAggregatorTransactor creates a new write-only instance of PriceAggregator, bound to a specific deployed contract.
func NewPriceAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*PriceAggregatorTransactor, error) {
	contract, err := bindPriceAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorTransactor{contract: contract}, nil
}

// NewPriceAggregatorFilterer creates a new log filterer instance of PriceAggregator, bound to a specific deployed contract.
func NewPriceAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*PriceAggregatorFilterer, error) {
	contract, err := bindPriceAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorFilterer{contract: contract}, nil
}

// bindPriceAggregator binds a generic wrapper to an already deployed contract.
func bindPriceAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PriceAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PriceAggregator *PriceAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceAggregator.Contract.PriceAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PriceAggregator *PriceAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceAggregator.Contract.PriceAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PriceAggregator *PriceAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceAggregator.Contract.PriceAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PriceAggregator *PriceAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PriceAggregator *PriceAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PriceAggregator *PriceAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceAggregator.Contract.contract.Transact(opts, method, params...)
}

// AssetId is a free data retrieval call binding the contract method 0x44de240a.
//
// Solidity: function assetId() view returns(bytes32)
func (_PriceAggregator *PriceAggregatorCaller) AssetId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "assetId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AssetId is a free data retrieval call binding the contract method 0x44de240a.
//
// Solidity: function assetId() view returns(bytes32)
func (_PriceAggregator *PriceAggregatorSession) AssetId() ([32]byte, error) {
	return _PriceAggregator.Contract.AssetId(&_PriceAggregator.CallOpts)
}

// AssetId is a free data retrieval call binding the contract method 0x44de240a.
//
// Solidity: function assetId() view returns(bytes32)
func (_PriceAggregator *PriceAggregatorCallerSession) AssetId() ([32]byte, error) {
	return _PriceAggregator.Contract.AssetId(&_PriceAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PriceAggregator *PriceAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PriceAggregator *PriceAggregatorSession) Decimals() (uint8, error) {
	return _PriceAggregator.Contract.Decimals(&_PriceAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PriceAggregator *PriceAggregatorCallerSession) Decimals() (uint8, error) {
	return _PriceAggregator.Contract.Decimals(&_PriceAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_PriceAggregator *PriceAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_PriceAggregator *PriceAggregatorSession) Description() (string, error) {
	return _PriceAggregator.Contract.Description(&_PriceAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_PriceAggregator *PriceAggregatorCallerSession) Description() (string, error) {
	return _PriceAggregator.Contract.Description(&_PriceAggregator.CallOpts)
}

// Fulfilled is a free data retrieval call binding the contract method 0x1bc404d6.
//
// Solidity: function fulfilled(uint256 ) view returns(bool)
func (_PriceAggregator *PriceAggregatorCaller) Fulfilled(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "fulfilled", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Fulfilled is a free data retrieval call binding the contract method 0x1bc404d6.
//
// Solidity: function fulfilled(uint256 ) view returns(bool)
func (_PriceAggregator *PriceAggregatorSession) Fulfilled(arg0 *big.Int) (bool, error) {
	return _PriceAggregator.Contract.Fulfilled(&_PriceAggregator.CallOpts, arg0)
}

// Fulfilled is a free data retrieval call binding the contract method 0x1bc404d6.
//
// Solidity: function fulfilled(uint256 ) view returns(bool)
func (_PriceAggregator *PriceAggregatorCallerSession) Fulfilled(arg0 *big.Int) (bool, error) {
	return _PriceAggregator.Contract.Fulfilled(&_PriceAggregator.CallOpts, arg0)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_PriceAggregator *PriceAggregatorCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "getRoundData", roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_PriceAggregator *PriceAggregatorSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _PriceAggregator.Contract.GetRoundData(&_PriceAggregator.CallOpts, roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_PriceAggregator *PriceAggregatorCallerSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _PriceAggregator.Contract.GetRoundData(&_PriceAggregator.CallOpts, roundId)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_PriceAggregator *PriceAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_PriceAggregator *PriceAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _PriceAggregator.Contract.LatestRoundData(&_PriceAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_PriceAggregator *PriceAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _PriceAggregator.Contract.LatestRoundData(&_PriceAggregator.CallOpts)
}

// LatestRoundId is a free data retrieval call binding the contract method 0x11a8f413.
//
// Solidity: function latestRoundId() view returns(uint80)
func (_PriceAggregator *PriceAggregatorCaller) LatestRoundId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "latestRoundId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRoundId is a free data retrieval call binding the contract method 0x11a8f413.
//
// Solidity: function latestRoundId() view returns(uint80)
func (_PriceAggregator *PriceAggregatorSession) LatestRoundId() (*big.Int, error) {
	return _PriceAggregator.Contract.LatestRoundId(&_PriceAggregator.CallOpts)
}

// LatestRoundId is a free data retrieval call binding the contract method 0x11a8f413.
//
// Solidity: function latestRoundId() view returns(uint80)
func (_PriceAggregator *PriceAggregatorCallerSession) LatestRoundId() (*big.Int, error) {
	return _PriceAggregator.Contract.LatestRoundId(&_PriceAggregator.CallOpts)
}

// MaxAge is a free data retrieval call binding the contract method 0x687043c5.
//
// Solidity: function maxAge() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCaller) MaxAge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "maxAge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAge is a free data retrieval call binding the contract method 0x687043c5.
//
// Solidity: function maxAge() view returns(uint256)
func (_PriceAggregator *PriceAggregatorSession) MaxAge() (*big.Int, error) {
	return _PriceAggregator.Contract.MaxAge(&_PriceAggregator.CallOpts)
}

// MaxAge is a free data retrieval call binding the contract method 0x687043c5.
//
// Solidity: function maxAge() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCallerSession) MaxAge() (*big.Int, error) {
	return _PriceAggregator.Contract.MaxAge(&_PriceAggregator.CallOpts)
}

// NextReqId is a free data retrieval call binding the contract method 0xef8e7b60.
//
// Solidity: function nextReqId() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCaller) NextReqId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "nextReqId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextReqId is a free data retrieval call binding the contract method 0xef8e7b60.
//
// Solidity: function nextReqId() view returns(uint256)
func (_PriceAggregator *PriceAggregatorSession) NextReqId() (*big.Int, error) {
	return _PriceAggregator.Contract.NextReqId(&_PriceAggregator.CallOpts)
}

// NextReqId is a free data retrieval call binding the contract method 0xef8e7b60.
//
// Solidity: function nextReqId() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCallerSession) NextReqId() (*big.Int, error) {
	return _PriceAggregator.Contract.NextReqId(&_PriceAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PriceAggregator *PriceAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PriceAggregator *PriceAggregatorSession) Owner() (common.Address, error) {
	return _PriceAggregator.Contract.Owner(&_PriceAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PriceAggregator *PriceAggregatorCallerSession) Owner() (common.Address, error) {
	return _PriceAggregator.Contract.Owner(&_PriceAggregator.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_PriceAggregator *PriceAggregatorCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_PriceAggregator *PriceAggregatorSession) PendingOwner() (common.Address, error) {
	return _PriceAggregator.Contract.PendingOwner(&_PriceAggregator.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_PriceAggregator *PriceAggregatorCallerSession) PendingOwner() (common.Address, error) {
	return _PriceAggregator.Contract.PendingOwner(&_PriceAggregator.CallOpts)
}

// ReporterSet is a free data retrieval call binding the contract method 0x4b8b7eaf.
//
// Solidity: function reporterSet() view returns(address)
func (_PriceAggregator *PriceAggregatorCaller) ReporterSet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "reporterSet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReporterSet is a free data retrieval call binding the contract method 0x4b8b7eaf.
//
// Solidity: function reporterSet() view returns(address)
func (_PriceAggregator *PriceAggregatorSession) ReporterSet() (common.Address, error) {
	return _PriceAggregator.Contract.ReporterSet(&_PriceAggregator.CallOpts)
}

// ReporterSet is a free data retrieval call binding the contract method 0x4b8b7eaf.
//
// Solidity: function reporterSet() view returns(address)
func (_PriceAggregator *PriceAggregatorCallerSession) ReporterSet() (common.Address, error) {
	return _PriceAggregator.Contract.ReporterSet(&_PriceAggregator.CallOpts)
}

// RequestFee is a free data retrieval call binding the contract method 0xeb2e578b.
//
// Solidity: function requestFee() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCaller) RequestFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "requestFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestFee is a free data retrieval call binding the contract method 0xeb2e578b.
//
// Solidity: function requestFee() view returns(uint256)
func (_PriceAggregator *PriceAggregatorSession) RequestFee() (*big.Int, error) {
	return _PriceAggregator.Contract.RequestFee(&_PriceAggregator.CallOpts)
}

// RequestFee is a free data retrieval call binding the contract method 0xeb2e578b.
//
// Solidity: function requestFee() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCallerSession) RequestFee() (*big.Int, error) {
	return _PriceAggregator.Contract.RequestFee(&_PriceAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PriceAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_PriceAggregator *PriceAggregatorSession) Version() (*big.Int, error) {
	return _PriceAggregator.Contract.Version(&_PriceAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_PriceAggregator *PriceAggregatorCallerSession) Version() (*big.Int, error) {
	return _PriceAggregator.Contract.Version(&_PriceAggregator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_PriceAggregator *PriceAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_PriceAggregator *PriceAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PriceAggregator.Contract.AcceptOwnership(&_PriceAggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_PriceAggregator *PriceAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PriceAggregator.Contract.AcceptOwnership(&_PriceAggregator.TransactOpts)
}

// FulfillPrice is a paid mutator transaction binding the contract method 0xa844ad4d.
//
// Solidity: function fulfillPrice(uint256 reqId, int256 price, uint256 timestamp, bytes[] signatures) returns()
func (_PriceAggregator *PriceAggregatorTransactor) FulfillPrice(opts *bind.TransactOpts, reqId *big.Int, price *big.Int, timestamp *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "fulfillPrice", reqId, price, timestamp, signatures)
}

// FulfillPrice is a paid mutator transaction binding the contract method 0xa844ad4d.
//
// Solidity: function fulfillPrice(uint256 reqId, int256 price, uint256 timestamp, bytes[] signatures) returns()
func (_PriceAggregator *PriceAggregatorSession) FulfillPrice(reqId *big.Int, price *big.Int, timestamp *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _PriceAggregator.Contract.FulfillPrice(&_PriceAggregator.TransactOpts, reqId, price, timestamp, signatures)
}

// FulfillPrice is a paid mutator transaction binding the contract method 0xa844ad4d.
//
// Solidity: function fulfillPrice(uint256 reqId, int256 price, uint256 timestamp, bytes[] signatures) returns()
func (_PriceAggregator *PriceAggregatorTransactorSession) FulfillPrice(reqId *big.Int, price *big.Int, timestamp *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _PriceAggregator.Contract.FulfillPrice(&_PriceAggregator.TransactOpts, reqId, price, timestamp, signatures)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PriceAggregator *PriceAggregatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PriceAggregator *PriceAggregatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PriceAggregator.Contract.RenounceOwnership(&_PriceAggregator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PriceAggregator *PriceAggregatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PriceAggregator.Contract.RenounceOwnership(&_PriceAggregator.TransactOpts)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x1604f9ea.
//
// Solidity: function requestPrice() payable returns(uint256 reqId)
func (_PriceAggregator *PriceAggregatorTransactor) RequestPrice(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "requestPrice")
}

// RequestPrice is a paid mutator transaction binding the contract method 0x1604f9ea.
//
// Solidity: function requestPrice() payable returns(uint256 reqId)
func (_PriceAggregator *PriceAggregatorSession) RequestPrice() (*types.Transaction, error) {
	return _PriceAggregator.Contract.RequestPrice(&_PriceAggregator.TransactOpts)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x1604f9ea.
//
// Solidity: function requestPrice() payable returns(uint256 reqId)
func (_PriceAggregator *PriceAggregatorTransactorSession) RequestPrice() (*types.Transaction, error) {
	return _PriceAggregator.Contract.RequestPrice(&_PriceAggregator.TransactOpts)
}

// SetMaxAge is a paid mutator transaction binding the contract method 0x5ae28fc9.
//
// Solidity: function setMaxAge(uint256 newMaxAge) returns()
func (_PriceAggregator *PriceAggregatorTransactor) SetMaxAge(opts *bind.TransactOpts, newMaxAge *big.Int) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "setMaxAge", newMaxAge)
}

// SetMaxAge is a paid mutator transaction binding the contract method 0x5ae28fc9.
//
// Solidity: function setMaxAge(uint256 newMaxAge) returns()
func (_PriceAggregator *PriceAggregatorSession) SetMaxAge(newMaxAge *big.Int) (*types.Transaction, error) {
	return _PriceAggregator.Contract.SetMaxAge(&_PriceAggregator.TransactOpts, newMaxAge)
}

// SetMaxAge is a paid mutator transaction binding the contract method 0x5ae28fc9.
//
// Solidity: function setMaxAge(uint256 newMaxAge) returns()
func (_PriceAggregator *PriceAggregatorTransactorSession) SetMaxAge(newMaxAge *big.Int) (*types.Transaction, error) {
	return _PriceAggregator.Contract.SetMaxAge(&_PriceAggregator.TransactOpts, newMaxAge)
}

// SetReporterSet is a paid mutator transaction binding the contract method 0x080721d4.
//
// Solidity: function setReporterSet(address newReporterSet) returns()
func (_PriceAggregator *PriceAggregatorTransactor) SetReporterSet(opts *bind.TransactOpts, newReporterSet common.Address) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "setReporterSet", newReporterSet)
}

// SetReporterSet is a paid mutator transaction binding the contract method 0x080721d4.
//
// Solidity: function setReporterSet(address newReporterSet) returns()
func (_PriceAggregator *PriceAggregatorSession) SetReporterSet(newReporterSet common.Address) (*types.Transaction, error) {
	return _PriceAggregator.Contract.SetReporterSet(&_PriceAggregator.TransactOpts, newReporterSet)
}

// SetReporterSet is a paid mutator transaction binding the contract method 0x080721d4.
//
// Solidity: function setReporterSet(address newReporterSet) returns()
func (_PriceAggregator *PriceAggregatorTransactorSession) SetReporterSet(newReporterSet common.Address) (*types.Transaction, error) {
	return _PriceAggregator.Contract.SetReporterSet(&_PriceAggregator.TransactOpts, newReporterSet)
}

// SetRequestFee is a paid mutator transaction binding the contract method 0xffb9c43f.
//
// Solidity: function setRequestFee(uint256 newFee) returns()
func (_PriceAggregator *PriceAggregatorTransactor) SetRequestFee(opts *bind.TransactOpts, newFee *big.Int) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "setRequestFee", newFee)
}

// SetRequestFee is a paid mutator transaction binding the contract method 0xffb9c43f.
//
// Solidity: function setRequestFee(uint256 newFee) returns()
func (_PriceAggregator *PriceAggregatorSession) SetRequestFee(newFee *big.Int) (*types.Transaction, error) {
	return _PriceAggregator.Contract.SetRequestFee(&_PriceAggregator.TransactOpts, newFee)
}

// SetRequestFee is a paid mutator transaction binding the contract method 0xffb9c43f.
//
// Solidity: function setRequestFee(uint256 newFee) returns()
func (_PriceAggregator *PriceAggregatorTransactorSession) SetRequestFee(newFee *big.Int) (*types.Transaction, error) {
	return _PriceAggregator.Contract.SetRequestFee(&_PriceAggregator.TransactOpts, newFee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PriceAggregator *PriceAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PriceAggregator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PriceAggregator *PriceAggregatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PriceAggregator.Contract.TransferOwnership(&_PriceAggregator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PriceAggregator *PriceAggregatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PriceAggregator.Contract.TransferOwnership(&_PriceAggregator.TransactOpts, newOwner)
}

// PriceAggregatorMaxAgeChangedIterator is returned from FilterMaxAgeChanged and is used to iterate over the raw logs and unpacked data for MaxAgeChanged events raised by the PriceAggregator contract.
type PriceAggregatorMaxAgeChangedIterator struct {
	Event *PriceAggregatorMaxAgeChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PriceAggregatorMaxAgeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceAggregatorMaxAgeChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PriceAggregatorMaxAgeChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PriceAggregatorMaxAgeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceAggregatorMaxAgeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceAggregatorMaxAgeChanged represents a MaxAgeChanged event raised by the PriceAggregator contract.
type PriceAggregatorMaxAgeChanged struct {
	OldMaxAge *big.Int
	NewMaxAge *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMaxAgeChanged is a free log retrieval operation binding the contract event 0xd0c16066bffbe2c853e82f18bad7cfa67f94dcc5f5754ce9a6350d584ff36791.
//
// Solidity: event MaxAgeChanged(uint256 oldMaxAge, uint256 newMaxAge)
func (_PriceAggregator *PriceAggregatorFilterer) FilterMaxAgeChanged(opts *bind.FilterOpts) (*PriceAggregatorMaxAgeChangedIterator, error) {

	logs, sub, err := _PriceAggregator.contract.FilterLogs(opts, "MaxAgeChanged")
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorMaxAgeChangedIterator{contract: _PriceAggregator.contract, event: "MaxAgeChanged", logs: logs, sub: sub}, nil
}

// WatchMaxAgeChanged is a free log subscription operation binding the contract event 0xd0c16066bffbe2c853e82f18bad7cfa67f94dcc5f5754ce9a6350d584ff36791.
//
// Solidity: event MaxAgeChanged(uint256 oldMaxAge, uint256 newMaxAge)
func (_PriceAggregator *PriceAggregatorFilterer) WatchMaxAgeChanged(opts *bind.WatchOpts, sink chan<- *PriceAggregatorMaxAgeChanged) (event.Subscription, error) {

	logs, sub, err := _PriceAggregator.contract.WatchLogs(opts, "MaxAgeChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceAggregatorMaxAgeChanged)
				if err := _PriceAggregator.contract.UnpackLog(event, "MaxAgeChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMaxAgeChanged is a log parse operation binding the contract event 0xd0c16066bffbe2c853e82f18bad7cfa67f94dcc5f5754ce9a6350d584ff36791.
//
// Solidity: event MaxAgeChanged(uint256 oldMaxAge, uint256 newMaxAge)
func (_PriceAggregator *PriceAggregatorFilterer) ParseMaxAgeChanged(log types.Log) (*PriceAggregatorMaxAgeChanged, error) {
	event := new(PriceAggregatorMaxAgeChanged)
	if err := _PriceAggregator.contract.UnpackLog(event, "MaxAgeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceAggregatorOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the PriceAggregator contract.
type PriceAggregatorOwnershipTransferStartedIterator struct {
	Event *PriceAggregatorOwnershipTransferStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PriceAggregatorOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceAggregatorOwnershipTransferStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PriceAggregatorOwnershipTransferStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PriceAggregatorOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceAggregatorOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceAggregatorOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the PriceAggregator contract.
type PriceAggregatorOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_PriceAggregator *PriceAggregatorFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PriceAggregatorOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PriceAggregator.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorOwnershipTransferStartedIterator{contract: _PriceAggregator.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_PriceAggregator *PriceAggregatorFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *PriceAggregatorOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PriceAggregator.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceAggregatorOwnershipTransferStarted)
				if err := _PriceAggregator.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_PriceAggregator *PriceAggregatorFilterer) ParseOwnershipTransferStarted(log types.Log) (*PriceAggregatorOwnershipTransferStarted, error) {
	event := new(PriceAggregatorOwnershipTransferStarted)
	if err := _PriceAggregator.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceAggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PriceAggregator contract.
type PriceAggregatorOwnershipTransferredIterator struct {
	Event *PriceAggregatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PriceAggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceAggregatorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PriceAggregatorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PriceAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceAggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the PriceAggregator contract.
type PriceAggregatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PriceAggregator *PriceAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PriceAggregatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PriceAggregator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorOwnershipTransferredIterator{contract: _PriceAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PriceAggregator *PriceAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PriceAggregatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PriceAggregator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceAggregatorOwnershipTransferred)
				if err := _PriceAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PriceAggregator *PriceAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*PriceAggregatorOwnershipTransferred, error) {
	event := new(PriceAggregatorOwnershipTransferred)
	if err := _PriceAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceAggregatorPriceFulfilledIterator is returned from FilterPriceFulfilled and is used to iterate over the raw logs and unpacked data for PriceFulfilled events raised by the PriceAggregator contract.
type PriceAggregatorPriceFulfilledIterator struct {
	Event *PriceAggregatorPriceFulfilled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PriceAggregatorPriceFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceAggregatorPriceFulfilled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PriceAggregatorPriceFulfilled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PriceAggregatorPriceFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceAggregatorPriceFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceAggregatorPriceFulfilled represents a PriceFulfilled event raised by the PriceAggregator contract.
type PriceAggregatorPriceFulfilled struct {
	ReqId     *big.Int
	Price     *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPriceFulfilled is a free log retrieval operation binding the contract event 0x82c08aa7285d5667568ba3b8821c82fa50deef99dc0c6b75d46fb5c7455ec22a.
//
// Solidity: event PriceFulfilled(uint256 indexed reqId, int256 price, uint256 timestamp)
func (_PriceAggregator *PriceAggregatorFilterer) FilterPriceFulfilled(opts *bind.FilterOpts, reqId []*big.Int) (*PriceAggregatorPriceFulfilledIterator, error) {

	var reqIdRule []interface{}
	for _, reqIdItem := range reqId {
		reqIdRule = append(reqIdRule, reqIdItem)
	}

	logs, sub, err := _PriceAggregator.contract.FilterLogs(opts, "PriceFulfilled", reqIdRule)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorPriceFulfilledIterator{contract: _PriceAggregator.contract, event: "PriceFulfilled", logs: logs, sub: sub}, nil
}

// WatchPriceFulfilled is a free log subscription operation binding the contract event 0x82c08aa7285d5667568ba3b8821c82fa50deef99dc0c6b75d46fb5c7455ec22a.
//
// Solidity: event PriceFulfilled(uint256 indexed reqId, int256 price, uint256 timestamp)
func (_PriceAggregator *PriceAggregatorFilterer) WatchPriceFulfilled(opts *bind.WatchOpts, sink chan<- *PriceAggregatorPriceFulfilled, reqId []*big.Int) (event.Subscription, error) {

	var reqIdRule []interface{}
	for _, reqIdItem := range reqId {
		reqIdRule = append(reqIdRule, reqIdItem)
	}

	logs, sub, err := _PriceAggregator.contract.WatchLogs(opts, "PriceFulfilled", reqIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceAggregatorPriceFulfilled)
				if err := _PriceAggregator.contract.UnpackLog(event, "PriceFulfilled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePriceFulfilled is a log parse operation binding the contract event 0x82c08aa7285d5667568ba3b8821c82fa50deef99dc0c6b75d46fb5c7455ec22a.
//
// Solidity: event PriceFulfilled(uint256 indexed reqId, int256 price, uint256 timestamp)
func (_PriceAggregator *PriceAggregatorFilterer) ParsePriceFulfilled(log types.Log) (*PriceAggregatorPriceFulfilled, error) {
	event := new(PriceAggregatorPriceFulfilled)
	if err := _PriceAggregator.contract.UnpackLog(event, "PriceFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceAggregatorPriceRequestedIterator is returned from FilterPriceRequested and is used to iterate over the raw logs and unpacked data for PriceRequested events raised by the PriceAggregator contract.
type PriceAggregatorPriceRequestedIterator struct {
	Event *PriceAggregatorPriceRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PriceAggregatorPriceRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceAggregatorPriceRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PriceAggregatorPriceRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PriceAggregatorPriceRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceAggregatorPriceRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceAggregatorPriceRequested represents a PriceRequested event raised by the PriceAggregator contract.
type PriceAggregatorPriceRequested struct {
	ReqId     *big.Int
	Requester common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPriceRequested is a free log retrieval operation binding the contract event 0x33b362cfc336cd3829a3a7896832b11a70bdbc0194a75fb7f1003919b23c8600.
//
// Solidity: event PriceRequested(uint256 indexed reqId, address indexed requester)
func (_PriceAggregator *PriceAggregatorFilterer) FilterPriceRequested(opts *bind.FilterOpts, reqId []*big.Int, requester []common.Address) (*PriceAggregatorPriceRequestedIterator, error) {

	var reqIdRule []interface{}
	for _, reqIdItem := range reqId {
		reqIdRule = append(reqIdRule, reqIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _PriceAggregator.contract.FilterLogs(opts, "PriceRequested", reqIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorPriceRequestedIterator{contract: _PriceAggregator.contract, event: "PriceRequested", logs: logs, sub: sub}, nil
}

// WatchPriceRequested is a free log subscription operation binding the contract event 0x33b362cfc336cd3829a3a7896832b11a70bdbc0194a75fb7f1003919b23c8600.
//
// Solidity: event PriceRequested(uint256 indexed reqId, address indexed requester)
func (_PriceAggregator *PriceAggregatorFilterer) WatchPriceRequested(opts *bind.WatchOpts, sink chan<- *PriceAggregatorPriceRequested, reqId []*big.Int, requester []common.Address) (event.Subscription, error) {

	var reqIdRule []interface{}
	for _, reqIdItem := range reqId {
		reqIdRule = append(reqIdRule, reqIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _PriceAggregator.contract.WatchLogs(opts, "PriceRequested", reqIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceAggregatorPriceRequested)
				if err := _PriceAggregator.contract.UnpackLog(event, "PriceRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePriceRequested is a log parse operation binding the contract event 0x33b362cfc336cd3829a3a7896832b11a70bdbc0194a75fb7f1003919b23c8600.
//
// Solidity: event PriceRequested(uint256 indexed reqId, address indexed requester)
func (_PriceAggregator *PriceAggregatorFilterer) ParsePriceRequested(log types.Log) (*PriceAggregatorPriceRequested, error) {
	event := new(PriceAggregatorPriceRequested)
	if err := _PriceAggregator.contract.UnpackLog(event, "PriceRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceAggregatorReporterSetChangedIterator is returned from FilterReporterSetChanged and is used to iterate over the raw logs and unpacked data for ReporterSetChanged events raised by the PriceAggregator contract.
type PriceAggregatorReporterSetChangedIterator struct {
	Event *PriceAggregatorReporterSetChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PriceAggregatorReporterSetChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceAggregatorReporterSetChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PriceAggregatorReporterSetChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PriceAggregatorReporterSetChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceAggregatorReporterSetChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceAggregatorReporterSetChanged represents a ReporterSetChanged event raised by the PriceAggregator contract.
type PriceAggregatorReporterSetChanged struct {
	OldReporterSet common.Address
	NewReporterSet common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterReporterSetChanged is a free log retrieval operation binding the contract event 0xec0431b3bf6c80f50e6df6cb617877002a1c17eda905f696ce48d6207e62b4bd.
//
// Solidity: event ReporterSetChanged(address indexed oldReporterSet, address indexed newReporterSet)
func (_PriceAggregator *PriceAggregatorFilterer) FilterReporterSetChanged(opts *bind.FilterOpts, oldReporterSet []common.Address, newReporterSet []common.Address) (*PriceAggregatorReporterSetChangedIterator, error) {

	var oldReporterSetRule []interface{}
	for _, oldReporterSetItem := range oldReporterSet {
		oldReporterSetRule = append(oldReporterSetRule, oldReporterSetItem)
	}
	var newReporterSetRule []interface{}
	for _, newReporterSetItem := range newReporterSet {
		newReporterSetRule = append(newReporterSetRule, newReporterSetItem)
	}

	logs, sub, err := _PriceAggregator.contract.FilterLogs(opts, "ReporterSetChanged", oldReporterSetRule, newReporterSetRule)
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorReporterSetChangedIterator{contract: _PriceAggregator.contract, event: "ReporterSetChanged", logs: logs, sub: sub}, nil
}

// WatchReporterSetChanged is a free log subscription operation binding the contract event 0xec0431b3bf6c80f50e6df6cb617877002a1c17eda905f696ce48d6207e62b4bd.
//
// Solidity: event ReporterSetChanged(address indexed oldReporterSet, address indexed newReporterSet)
func (_PriceAggregator *PriceAggregatorFilterer) WatchReporterSetChanged(opts *bind.WatchOpts, sink chan<- *PriceAggregatorReporterSetChanged, oldReporterSet []common.Address, newReporterSet []common.Address) (event.Subscription, error) {

	var oldReporterSetRule []interface{}
	for _, oldReporterSetItem := range oldReporterSet {
		oldReporterSetRule = append(oldReporterSetRule, oldReporterSetItem)
	}
	var newReporterSetRule []interface{}
	for _, newReporterSetItem := range newReporterSet {
		newReporterSetRule = append(newReporterSetRule, newReporterSetItem)
	}

	logs, sub, err := _PriceAggregator.contract.WatchLogs(opts, "ReporterSetChanged", oldReporterSetRule, newReporterSetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceAggregatorReporterSetChanged)
				if err := _PriceAggregator.contract.UnpackLog(event, "ReporterSetChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReporterSetChanged is a log parse operation binding the contract event 0xec0431b3bf6c80f50e6df6cb617877002a1c17eda905f696ce48d6207e62b4bd.
//
// Solidity: event ReporterSetChanged(address indexed oldReporterSet, address indexed newReporterSet)
func (_PriceAggregator *PriceAggregatorFilterer) ParseReporterSetChanged(log types.Log) (*PriceAggregatorReporterSetChanged, error) {
	event := new(PriceAggregatorReporterSetChanged)
	if err := _PriceAggregator.contract.UnpackLog(event, "ReporterSetChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceAggregatorRequestFeeChangedIterator is returned from FilterRequestFeeChanged and is used to iterate over the raw logs and unpacked data for RequestFeeChanged events raised by the PriceAggregator contract.
type PriceAggregatorRequestFeeChangedIterator struct {
	Event *PriceAggregatorRequestFeeChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PriceAggregatorRequestFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceAggregatorRequestFeeChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PriceAggregatorRequestFeeChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PriceAggregatorRequestFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceAggregatorRequestFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceAggregatorRequestFeeChanged represents a RequestFeeChanged event raised by the PriceAggregator contract.
type PriceAggregatorRequestFeeChanged struct {
	OldFee *big.Int
	NewFee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRequestFeeChanged is a free log retrieval operation binding the contract event 0xe6d124068d802a9f67830cf20bd0e97e49bb3c83ba28dff35c786a933dc114c4.
//
// Solidity: event RequestFeeChanged(uint256 oldFee, uint256 newFee)
func (_PriceAggregator *PriceAggregatorFilterer) FilterRequestFeeChanged(opts *bind.FilterOpts) (*PriceAggregatorRequestFeeChangedIterator, error) {

	logs, sub, err := _PriceAggregator.contract.FilterLogs(opts, "RequestFeeChanged")
	if err != nil {
		return nil, err
	}
	return &PriceAggregatorRequestFeeChangedIterator{contract: _PriceAggregator.contract, event: "RequestFeeChanged", logs: logs, sub: sub}, nil
}

// WatchRequestFeeChanged is a free log subscription operation binding the contract event 0xe6d124068d802a9f67830cf20bd0e97e49bb3c83ba28dff35c786a933dc114c4.
//
// Solidity: event RequestFeeChanged(uint256 oldFee, uint256 newFee)
func (_PriceAggregator *PriceAggregatorFilterer) WatchRequestFeeChanged(opts *bind.WatchOpts, sink chan<- *PriceAggregatorRequestFeeChanged) (event.Subscription, error) {

	logs, sub, err := _PriceAggregator.contract.WatchLogs(opts, "RequestFeeChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceAggregatorRequestFeeChanged)
				if err := _PriceAggregator.contract.UnpackLog(event, "RequestFeeChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestFeeChanged is a log parse operation binding the contract event 0xe6d124068d802a9f67830cf20bd0e97e49bb3c83ba28dff35c786a933dc114c4.
//
// Solidity: event RequestFeeChanged(uint256 oldFee, uint256 newFee)
func (_PriceAggregator *PriceAggregatorFilterer) ParseRequestFeeChanged(log types.Log) (*PriceAggregatorRequestFeeChanged, error) {
	event := new(PriceAggregatorRequestFeeChanged)
	if err := _PriceAggregator.contract.UnpackLog(event, "RequestFeeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
