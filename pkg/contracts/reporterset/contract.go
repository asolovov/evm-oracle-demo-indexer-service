// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package reporterset

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

// ReporterSetMetaData contains all meta data concerning the ReporterSet contract.
var ReporterSetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"initialReporters\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"initialThreshold\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"}],\"name\":\"ReporterAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"}],\"name\":\"ReporterNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reporterCount\",\"type\":\"uint256\"}],\"name\":\"ThresholdExceedsReporterCount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroReporter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroThreshold\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"}],\"name\":\"ReporterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"}],\"name\":\"ReporterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"ThresholdChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"}],\"name\":\"addReporter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReporters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isReporter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"}],\"name\":\"removeReporter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"setThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ReporterSetABI is the input ABI used to generate the binding from.
// Deprecated: Use ReporterSetMetaData.ABI instead.
var ReporterSetABI = ReporterSetMetaData.ABI

// ReporterSet is an auto generated Go binding around an Ethereum contract.
type ReporterSet struct {
	ReporterSetCaller     // Read-only binding to the contract
	ReporterSetTransactor // Write-only binding to the contract
	ReporterSetFilterer   // Log filterer for contract events
}

// ReporterSetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReporterSetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReporterSetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReporterSetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReporterSetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReporterSetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReporterSetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReporterSetSession struct {
	Contract     *ReporterSet      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReporterSetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReporterSetCallerSession struct {
	Contract *ReporterSetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ReporterSetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReporterSetTransactorSession struct {
	Contract     *ReporterSetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ReporterSetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReporterSetRaw struct {
	Contract *ReporterSet // Generic contract binding to access the raw methods on
}

// ReporterSetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReporterSetCallerRaw struct {
	Contract *ReporterSetCaller // Generic read-only contract binding to access the raw methods on
}

// ReporterSetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReporterSetTransactorRaw struct {
	Contract *ReporterSetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReporterSet creates a new instance of ReporterSet, bound to a specific deployed contract.
func NewReporterSet(address common.Address, backend bind.ContractBackend) (*ReporterSet, error) {
	contract, err := bindReporterSet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReporterSet{ReporterSetCaller: ReporterSetCaller{contract: contract}, ReporterSetTransactor: ReporterSetTransactor{contract: contract}, ReporterSetFilterer: ReporterSetFilterer{contract: contract}}, nil
}

// NewReporterSetCaller creates a new read-only instance of ReporterSet, bound to a specific deployed contract.
func NewReporterSetCaller(address common.Address, caller bind.ContractCaller) (*ReporterSetCaller, error) {
	contract, err := bindReporterSet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReporterSetCaller{contract: contract}, nil
}

// NewReporterSetTransactor creates a new write-only instance of ReporterSet, bound to a specific deployed contract.
func NewReporterSetTransactor(address common.Address, transactor bind.ContractTransactor) (*ReporterSetTransactor, error) {
	contract, err := bindReporterSet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReporterSetTransactor{contract: contract}, nil
}

// NewReporterSetFilterer creates a new log filterer instance of ReporterSet, bound to a specific deployed contract.
func NewReporterSetFilterer(address common.Address, filterer bind.ContractFilterer) (*ReporterSetFilterer, error) {
	contract, err := bindReporterSet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReporterSetFilterer{contract: contract}, nil
}

// bindReporterSet binds a generic wrapper to an already deployed contract.
func bindReporterSet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReporterSetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReporterSet *ReporterSetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReporterSet.Contract.ReporterSetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReporterSet *ReporterSetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReporterSet.Contract.ReporterSetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReporterSet *ReporterSetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReporterSet.Contract.ReporterSetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReporterSet *ReporterSetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReporterSet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReporterSet *ReporterSetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReporterSet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReporterSet *ReporterSetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReporterSet.Contract.contract.Transact(opts, method, params...)
}

// GetReporters is a free data retrieval call binding the contract method 0x70f46202.
//
// Solidity: function getReporters() view returns(address[])
func (_ReporterSet *ReporterSetCaller) GetReporters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ReporterSet.contract.Call(opts, &out, "getReporters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetReporters is a free data retrieval call binding the contract method 0x70f46202.
//
// Solidity: function getReporters() view returns(address[])
func (_ReporterSet *ReporterSetSession) GetReporters() ([]common.Address, error) {
	return _ReporterSet.Contract.GetReporters(&_ReporterSet.CallOpts)
}

// GetReporters is a free data retrieval call binding the contract method 0x70f46202.
//
// Solidity: function getReporters() view returns(address[])
func (_ReporterSet *ReporterSetCallerSession) GetReporters() ([]common.Address, error) {
	return _ReporterSet.Contract.GetReporters(&_ReporterSet.CallOpts)
}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_ReporterSet *ReporterSetCaller) GetThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReporterSet.contract.Call(opts, &out, "getThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_ReporterSet *ReporterSetSession) GetThreshold() (*big.Int, error) {
	return _ReporterSet.Contract.GetThreshold(&_ReporterSet.CallOpts)
}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_ReporterSet *ReporterSetCallerSession) GetThreshold() (*big.Int, error) {
	return _ReporterSet.Contract.GetThreshold(&_ReporterSet.CallOpts)
}

// IsReporter is a free data retrieval call binding the contract method 0x044ad7be.
//
// Solidity: function isReporter(address account) view returns(bool)
func (_ReporterSet *ReporterSetCaller) IsReporter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _ReporterSet.contract.Call(opts, &out, "isReporter", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsReporter is a free data retrieval call binding the contract method 0x044ad7be.
//
// Solidity: function isReporter(address account) view returns(bool)
func (_ReporterSet *ReporterSetSession) IsReporter(account common.Address) (bool, error) {
	return _ReporterSet.Contract.IsReporter(&_ReporterSet.CallOpts, account)
}

// IsReporter is a free data retrieval call binding the contract method 0x044ad7be.
//
// Solidity: function isReporter(address account) view returns(bool)
func (_ReporterSet *ReporterSetCallerSession) IsReporter(account common.Address) (bool, error) {
	return _ReporterSet.Contract.IsReporter(&_ReporterSet.CallOpts, account)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReporterSet *ReporterSetCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReporterSet.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReporterSet *ReporterSetSession) Owner() (common.Address, error) {
	return _ReporterSet.Contract.Owner(&_ReporterSet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReporterSet *ReporterSetCallerSession) Owner() (common.Address, error) {
	return _ReporterSet.Contract.Owner(&_ReporterSet.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ReporterSet *ReporterSetCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReporterSet.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ReporterSet *ReporterSetSession) PendingOwner() (common.Address, error) {
	return _ReporterSet.Contract.PendingOwner(&_ReporterSet.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ReporterSet *ReporterSetCallerSession) PendingOwner() (common.Address, error) {
	return _ReporterSet.Contract.PendingOwner(&_ReporterSet.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ReporterSet *ReporterSetTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReporterSet.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ReporterSet *ReporterSetSession) AcceptOwnership() (*types.Transaction, error) {
	return _ReporterSet.Contract.AcceptOwnership(&_ReporterSet.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ReporterSet *ReporterSetTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ReporterSet.Contract.AcceptOwnership(&_ReporterSet.TransactOpts)
}

// AddReporter is a paid mutator transaction binding the contract method 0xdd8755f2.
//
// Solidity: function addReporter(address reporter) returns()
func (_ReporterSet *ReporterSetTransactor) AddReporter(opts *bind.TransactOpts, reporter common.Address) (*types.Transaction, error) {
	return _ReporterSet.contract.Transact(opts, "addReporter", reporter)
}

// AddReporter is a paid mutator transaction binding the contract method 0xdd8755f2.
//
// Solidity: function addReporter(address reporter) returns()
func (_ReporterSet *ReporterSetSession) AddReporter(reporter common.Address) (*types.Transaction, error) {
	return _ReporterSet.Contract.AddReporter(&_ReporterSet.TransactOpts, reporter)
}

// AddReporter is a paid mutator transaction binding the contract method 0xdd8755f2.
//
// Solidity: function addReporter(address reporter) returns()
func (_ReporterSet *ReporterSetTransactorSession) AddReporter(reporter common.Address) (*types.Transaction, error) {
	return _ReporterSet.Contract.AddReporter(&_ReporterSet.TransactOpts, reporter)
}

// RemoveReporter is a paid mutator transaction binding the contract method 0x5de5c212.
//
// Solidity: function removeReporter(address reporter) returns()
func (_ReporterSet *ReporterSetTransactor) RemoveReporter(opts *bind.TransactOpts, reporter common.Address) (*types.Transaction, error) {
	return _ReporterSet.contract.Transact(opts, "removeReporter", reporter)
}

// RemoveReporter is a paid mutator transaction binding the contract method 0x5de5c212.
//
// Solidity: function removeReporter(address reporter) returns()
func (_ReporterSet *ReporterSetSession) RemoveReporter(reporter common.Address) (*types.Transaction, error) {
	return _ReporterSet.Contract.RemoveReporter(&_ReporterSet.TransactOpts, reporter)
}

// RemoveReporter is a paid mutator transaction binding the contract method 0x5de5c212.
//
// Solidity: function removeReporter(address reporter) returns()
func (_ReporterSet *ReporterSetTransactorSession) RemoveReporter(reporter common.Address) (*types.Transaction, error) {
	return _ReporterSet.Contract.RemoveReporter(&_ReporterSet.TransactOpts, reporter)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReporterSet *ReporterSetTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReporterSet.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReporterSet *ReporterSetSession) RenounceOwnership() (*types.Transaction, error) {
	return _ReporterSet.Contract.RenounceOwnership(&_ReporterSet.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReporterSet *ReporterSetTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ReporterSet.Contract.RenounceOwnership(&_ReporterSet.TransactOpts)
}

// SetThreshold is a paid mutator transaction binding the contract method 0x960bfe04.
//
// Solidity: function setThreshold(uint256 newThreshold) returns()
func (_ReporterSet *ReporterSetTransactor) SetThreshold(opts *bind.TransactOpts, newThreshold *big.Int) (*types.Transaction, error) {
	return _ReporterSet.contract.Transact(opts, "setThreshold", newThreshold)
}

// SetThreshold is a paid mutator transaction binding the contract method 0x960bfe04.
//
// Solidity: function setThreshold(uint256 newThreshold) returns()
func (_ReporterSet *ReporterSetSession) SetThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _ReporterSet.Contract.SetThreshold(&_ReporterSet.TransactOpts, newThreshold)
}

// SetThreshold is a paid mutator transaction binding the contract method 0x960bfe04.
//
// Solidity: function setThreshold(uint256 newThreshold) returns()
func (_ReporterSet *ReporterSetTransactorSession) SetThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _ReporterSet.Contract.SetThreshold(&_ReporterSet.TransactOpts, newThreshold)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReporterSet *ReporterSetTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ReporterSet.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReporterSet *ReporterSetSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReporterSet.Contract.TransferOwnership(&_ReporterSet.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReporterSet *ReporterSetTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReporterSet.Contract.TransferOwnership(&_ReporterSet.TransactOpts, newOwner)
}

// ReporterSetOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the ReporterSet contract.
type ReporterSetOwnershipTransferStartedIterator struct {
	Event *ReporterSetOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *ReporterSetOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReporterSetOwnershipTransferStarted)
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
		it.Event = new(ReporterSetOwnershipTransferStarted)
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
func (it *ReporterSetOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReporterSetOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReporterSetOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the ReporterSet contract.
type ReporterSetOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ReporterSet *ReporterSetFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ReporterSetOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReporterSet.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ReporterSetOwnershipTransferStartedIterator{contract: _ReporterSet.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ReporterSet *ReporterSetFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *ReporterSetOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReporterSet.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReporterSetOwnershipTransferStarted)
				if err := _ReporterSet.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_ReporterSet *ReporterSetFilterer) ParseOwnershipTransferStarted(log types.Log) (*ReporterSetOwnershipTransferStarted, error) {
	event := new(ReporterSetOwnershipTransferStarted)
	if err := _ReporterSet.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReporterSetOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ReporterSet contract.
type ReporterSetOwnershipTransferredIterator struct {
	Event *ReporterSetOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ReporterSetOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReporterSetOwnershipTransferred)
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
		it.Event = new(ReporterSetOwnershipTransferred)
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
func (it *ReporterSetOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReporterSetOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReporterSetOwnershipTransferred represents a OwnershipTransferred event raised by the ReporterSet contract.
type ReporterSetOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ReporterSet *ReporterSetFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ReporterSetOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReporterSet.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ReporterSetOwnershipTransferredIterator{contract: _ReporterSet.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ReporterSet *ReporterSetFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ReporterSetOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReporterSet.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReporterSetOwnershipTransferred)
				if err := _ReporterSet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ReporterSet *ReporterSetFilterer) ParseOwnershipTransferred(log types.Log) (*ReporterSetOwnershipTransferred, error) {
	event := new(ReporterSetOwnershipTransferred)
	if err := _ReporterSet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReporterSetReporterAddedIterator is returned from FilterReporterAdded and is used to iterate over the raw logs and unpacked data for ReporterAdded events raised by the ReporterSet contract.
type ReporterSetReporterAddedIterator struct {
	Event *ReporterSetReporterAdded // Event containing the contract specifics and raw log

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
func (it *ReporterSetReporterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReporterSetReporterAdded)
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
		it.Event = new(ReporterSetReporterAdded)
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
func (it *ReporterSetReporterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReporterSetReporterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReporterSetReporterAdded represents a ReporterAdded event raised by the ReporterSet contract.
type ReporterSetReporterAdded struct {
	Reporter common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReporterAdded is a free log retrieval operation binding the contract event 0x52824c4a4d431c73d6d24a91f775d722860111c6d30da30b0b21e062fb4c7365.
//
// Solidity: event ReporterAdded(address indexed reporter)
func (_ReporterSet *ReporterSetFilterer) FilterReporterAdded(opts *bind.FilterOpts, reporter []common.Address) (*ReporterSetReporterAddedIterator, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _ReporterSet.contract.FilterLogs(opts, "ReporterAdded", reporterRule)
	if err != nil {
		return nil, err
	}
	return &ReporterSetReporterAddedIterator{contract: _ReporterSet.contract, event: "ReporterAdded", logs: logs, sub: sub}, nil
}

// WatchReporterAdded is a free log subscription operation binding the contract event 0x52824c4a4d431c73d6d24a91f775d722860111c6d30da30b0b21e062fb4c7365.
//
// Solidity: event ReporterAdded(address indexed reporter)
func (_ReporterSet *ReporterSetFilterer) WatchReporterAdded(opts *bind.WatchOpts, sink chan<- *ReporterSetReporterAdded, reporter []common.Address) (event.Subscription, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _ReporterSet.contract.WatchLogs(opts, "ReporterAdded", reporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReporterSetReporterAdded)
				if err := _ReporterSet.contract.UnpackLog(event, "ReporterAdded", log); err != nil {
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

// ParseReporterAdded is a log parse operation binding the contract event 0x52824c4a4d431c73d6d24a91f775d722860111c6d30da30b0b21e062fb4c7365.
//
// Solidity: event ReporterAdded(address indexed reporter)
func (_ReporterSet *ReporterSetFilterer) ParseReporterAdded(log types.Log) (*ReporterSetReporterAdded, error) {
	event := new(ReporterSetReporterAdded)
	if err := _ReporterSet.contract.UnpackLog(event, "ReporterAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReporterSetReporterRemovedIterator is returned from FilterReporterRemoved and is used to iterate over the raw logs and unpacked data for ReporterRemoved events raised by the ReporterSet contract.
type ReporterSetReporterRemovedIterator struct {
	Event *ReporterSetReporterRemoved // Event containing the contract specifics and raw log

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
func (it *ReporterSetReporterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReporterSetReporterRemoved)
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
		it.Event = new(ReporterSetReporterRemoved)
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
func (it *ReporterSetReporterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReporterSetReporterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReporterSetReporterRemoved represents a ReporterRemoved event raised by the ReporterSet contract.
type ReporterSetReporterRemoved struct {
	Reporter common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReporterRemoved is a free log retrieval operation binding the contract event 0xf125aa5e6cb22c17ea0de4e1f86d2ff289ddd6a37e791b9ad29d0091cba06ac8.
//
// Solidity: event ReporterRemoved(address indexed reporter)
func (_ReporterSet *ReporterSetFilterer) FilterReporterRemoved(opts *bind.FilterOpts, reporter []common.Address) (*ReporterSetReporterRemovedIterator, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _ReporterSet.contract.FilterLogs(opts, "ReporterRemoved", reporterRule)
	if err != nil {
		return nil, err
	}
	return &ReporterSetReporterRemovedIterator{contract: _ReporterSet.contract, event: "ReporterRemoved", logs: logs, sub: sub}, nil
}

// WatchReporterRemoved is a free log subscription operation binding the contract event 0xf125aa5e6cb22c17ea0de4e1f86d2ff289ddd6a37e791b9ad29d0091cba06ac8.
//
// Solidity: event ReporterRemoved(address indexed reporter)
func (_ReporterSet *ReporterSetFilterer) WatchReporterRemoved(opts *bind.WatchOpts, sink chan<- *ReporterSetReporterRemoved, reporter []common.Address) (event.Subscription, error) {

	var reporterRule []interface{}
	for _, reporterItem := range reporter {
		reporterRule = append(reporterRule, reporterItem)
	}

	logs, sub, err := _ReporterSet.contract.WatchLogs(opts, "ReporterRemoved", reporterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReporterSetReporterRemoved)
				if err := _ReporterSet.contract.UnpackLog(event, "ReporterRemoved", log); err != nil {
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

// ParseReporterRemoved is a log parse operation binding the contract event 0xf125aa5e6cb22c17ea0de4e1f86d2ff289ddd6a37e791b9ad29d0091cba06ac8.
//
// Solidity: event ReporterRemoved(address indexed reporter)
func (_ReporterSet *ReporterSetFilterer) ParseReporterRemoved(log types.Log) (*ReporterSetReporterRemoved, error) {
	event := new(ReporterSetReporterRemoved)
	if err := _ReporterSet.contract.UnpackLog(event, "ReporterRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReporterSetThresholdChangedIterator is returned from FilterThresholdChanged and is used to iterate over the raw logs and unpacked data for ThresholdChanged events raised by the ReporterSet contract.
type ReporterSetThresholdChangedIterator struct {
	Event *ReporterSetThresholdChanged // Event containing the contract specifics and raw log

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
func (it *ReporterSetThresholdChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReporterSetThresholdChanged)
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
		it.Event = new(ReporterSetThresholdChanged)
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
func (it *ReporterSetThresholdChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReporterSetThresholdChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReporterSetThresholdChanged represents a ThresholdChanged event raised by the ReporterSet contract.
type ReporterSetThresholdChanged struct {
	OldThreshold *big.Int
	NewThreshold *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterThresholdChanged is a free log retrieval operation binding the contract event 0x3164947cf0f49f08dd0cd80e671535b1e11590d347c55dcaa97ba3c24a96b33a.
//
// Solidity: event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold)
func (_ReporterSet *ReporterSetFilterer) FilterThresholdChanged(opts *bind.FilterOpts) (*ReporterSetThresholdChangedIterator, error) {

	logs, sub, err := _ReporterSet.contract.FilterLogs(opts, "ThresholdChanged")
	if err != nil {
		return nil, err
	}
	return &ReporterSetThresholdChangedIterator{contract: _ReporterSet.contract, event: "ThresholdChanged", logs: logs, sub: sub}, nil
}

// WatchThresholdChanged is a free log subscription operation binding the contract event 0x3164947cf0f49f08dd0cd80e671535b1e11590d347c55dcaa97ba3c24a96b33a.
//
// Solidity: event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold)
func (_ReporterSet *ReporterSetFilterer) WatchThresholdChanged(opts *bind.WatchOpts, sink chan<- *ReporterSetThresholdChanged) (event.Subscription, error) {

	logs, sub, err := _ReporterSet.contract.WatchLogs(opts, "ThresholdChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReporterSetThresholdChanged)
				if err := _ReporterSet.contract.UnpackLog(event, "ThresholdChanged", log); err != nil {
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

// ParseThresholdChanged is a log parse operation binding the contract event 0x3164947cf0f49f08dd0cd80e671535b1e11590d347c55dcaa97ba3c24a96b33a.
//
// Solidity: event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold)
func (_ReporterSet *ReporterSetFilterer) ParseThresholdChanged(log types.Log) (*ReporterSetThresholdChanged, error) {
	event := new(ReporterSetThresholdChanged)
	if err := _ReporterSet.contract.UnpackLog(event, "ThresholdChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
