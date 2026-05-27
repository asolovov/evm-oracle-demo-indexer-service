// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package oracleregistry

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

// OracleRegistryMetaData contains all meta data concerning the OracleRegistry contract.
var OracleRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAggregator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAssetId\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"assetId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"AssetRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"assetId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAggregator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAggregator\",\"type\":\"address\"}],\"name\":\"AssetUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assetId\",\"type\":\"bytes32\"}],\"name\":\"getAggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"listAssets\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"assetId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"registerAsset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OracleRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleRegistryMetaData.ABI instead.
var OracleRegistryABI = OracleRegistryMetaData.ABI

// OracleRegistry is an auto generated Go binding around an Ethereum contract.
type OracleRegistry struct {
	OracleRegistryCaller     // Read-only binding to the contract
	OracleRegistryTransactor // Write-only binding to the contract
	OracleRegistryFilterer   // Log filterer for contract events
}

// OracleRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleRegistrySession struct {
	Contract     *OracleRegistry   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleRegistryCallerSession struct {
	Contract *OracleRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OracleRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleRegistryTransactorSession struct {
	Contract     *OracleRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OracleRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRegistryRaw struct {
	Contract *OracleRegistry // Generic contract binding to access the raw methods on
}

// OracleRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleRegistryCallerRaw struct {
	Contract *OracleRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// OracleRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleRegistryTransactorRaw struct {
	Contract *OracleRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleRegistry creates a new instance of OracleRegistry, bound to a specific deployed contract.
func NewOracleRegistry(address common.Address, backend bind.ContractBackend) (*OracleRegistry, error) {
	contract, err := bindOracleRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleRegistry{OracleRegistryCaller: OracleRegistryCaller{contract: contract}, OracleRegistryTransactor: OracleRegistryTransactor{contract: contract}, OracleRegistryFilterer: OracleRegistryFilterer{contract: contract}}, nil
}

// NewOracleRegistryCaller creates a new read-only instance of OracleRegistry, bound to a specific deployed contract.
func NewOracleRegistryCaller(address common.Address, caller bind.ContractCaller) (*OracleRegistryCaller, error) {
	contract, err := bindOracleRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleRegistryCaller{contract: contract}, nil
}

// NewOracleRegistryTransactor creates a new write-only instance of OracleRegistry, bound to a specific deployed contract.
func NewOracleRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleRegistryTransactor, error) {
	contract, err := bindOracleRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleRegistryTransactor{contract: contract}, nil
}

// NewOracleRegistryFilterer creates a new log filterer instance of OracleRegistry, bound to a specific deployed contract.
func NewOracleRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleRegistryFilterer, error) {
	contract, err := bindOracleRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleRegistryFilterer{contract: contract}, nil
}

// bindOracleRegistry binds a generic wrapper to an already deployed contract.
func bindOracleRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleRegistry *OracleRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleRegistry.Contract.OracleRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleRegistry *OracleRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleRegistry.Contract.OracleRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleRegistry *OracleRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleRegistry.Contract.OracleRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleRegistry *OracleRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleRegistry *OracleRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleRegistry *OracleRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetAggregator is a free data retrieval call binding the contract method 0x331b1816.
//
// Solidity: function getAggregator(bytes32 assetId) view returns(address)
func (_OracleRegistry *OracleRegistryCaller) GetAggregator(opts *bind.CallOpts, assetId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _OracleRegistry.contract.Call(opts, &out, "getAggregator", assetId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAggregator is a free data retrieval call binding the contract method 0x331b1816.
//
// Solidity: function getAggregator(bytes32 assetId) view returns(address)
func (_OracleRegistry *OracleRegistrySession) GetAggregator(assetId [32]byte) (common.Address, error) {
	return _OracleRegistry.Contract.GetAggregator(&_OracleRegistry.CallOpts, assetId)
}

// GetAggregator is a free data retrieval call binding the contract method 0x331b1816.
//
// Solidity: function getAggregator(bytes32 assetId) view returns(address)
func (_OracleRegistry *OracleRegistryCallerSession) GetAggregator(assetId [32]byte) (common.Address, error) {
	return _OracleRegistry.Contract.GetAggregator(&_OracleRegistry.CallOpts, assetId)
}

// ListAssets is a free data retrieval call binding the contract method 0xeadeb9a3.
//
// Solidity: function listAssets() view returns(bytes32[])
func (_OracleRegistry *OracleRegistryCaller) ListAssets(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _OracleRegistry.contract.Call(opts, &out, "listAssets")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// ListAssets is a free data retrieval call binding the contract method 0xeadeb9a3.
//
// Solidity: function listAssets() view returns(bytes32[])
func (_OracleRegistry *OracleRegistrySession) ListAssets() ([][32]byte, error) {
	return _OracleRegistry.Contract.ListAssets(&_OracleRegistry.CallOpts)
}

// ListAssets is a free data retrieval call binding the contract method 0xeadeb9a3.
//
// Solidity: function listAssets() view returns(bytes32[])
func (_OracleRegistry *OracleRegistryCallerSession) ListAssets() ([][32]byte, error) {
	return _OracleRegistry.Contract.ListAssets(&_OracleRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OracleRegistry *OracleRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OracleRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OracleRegistry *OracleRegistrySession) Owner() (common.Address, error) {
	return _OracleRegistry.Contract.Owner(&_OracleRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OracleRegistry *OracleRegistryCallerSession) Owner() (common.Address, error) {
	return _OracleRegistry.Contract.Owner(&_OracleRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_OracleRegistry *OracleRegistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OracleRegistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_OracleRegistry *OracleRegistrySession) PendingOwner() (common.Address, error) {
	return _OracleRegistry.Contract.PendingOwner(&_OracleRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_OracleRegistry *OracleRegistryCallerSession) PendingOwner() (common.Address, error) {
	return _OracleRegistry.Contract.PendingOwner(&_OracleRegistry.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OracleRegistry *OracleRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleRegistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OracleRegistry *OracleRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _OracleRegistry.Contract.AcceptOwnership(&_OracleRegistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_OracleRegistry *OracleRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OracleRegistry.Contract.AcceptOwnership(&_OracleRegistry.TransactOpts)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xc2867745.
//
// Solidity: function registerAsset(bytes32 assetId, address aggregator) returns()
func (_OracleRegistry *OracleRegistryTransactor) RegisterAsset(opts *bind.TransactOpts, assetId [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _OracleRegistry.contract.Transact(opts, "registerAsset", assetId, aggregator)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xc2867745.
//
// Solidity: function registerAsset(bytes32 assetId, address aggregator) returns()
func (_OracleRegistry *OracleRegistrySession) RegisterAsset(assetId [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _OracleRegistry.Contract.RegisterAsset(&_OracleRegistry.TransactOpts, assetId, aggregator)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xc2867745.
//
// Solidity: function registerAsset(bytes32 assetId, address aggregator) returns()
func (_OracleRegistry *OracleRegistryTransactorSession) RegisterAsset(assetId [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _OracleRegistry.Contract.RegisterAsset(&_OracleRegistry.TransactOpts, assetId, aggregator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OracleRegistry *OracleRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OracleRegistry *OracleRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _OracleRegistry.Contract.RenounceOwnership(&_OracleRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OracleRegistry *OracleRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OracleRegistry.Contract.RenounceOwnership(&_OracleRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OracleRegistry *OracleRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _OracleRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OracleRegistry *OracleRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OracleRegistry.Contract.TransferOwnership(&_OracleRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OracleRegistry *OracleRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OracleRegistry.Contract.TransferOwnership(&_OracleRegistry.TransactOpts, newOwner)
}

// OracleRegistryAssetRegisteredIterator is returned from FilterAssetRegistered and is used to iterate over the raw logs and unpacked data for AssetRegistered events raised by the OracleRegistry contract.
type OracleRegistryAssetRegisteredIterator struct {
	Event *OracleRegistryAssetRegistered // Event containing the contract specifics and raw log

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
func (it *OracleRegistryAssetRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleRegistryAssetRegistered)
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
		it.Event = new(OracleRegistryAssetRegistered)
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
func (it *OracleRegistryAssetRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleRegistryAssetRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleRegistryAssetRegistered represents a AssetRegistered event raised by the OracleRegistry contract.
type OracleRegistryAssetRegistered struct {
	AssetId    [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAssetRegistered is a free log retrieval operation binding the contract event 0x7b8c7b505365aa1b7f9ce04295e6da7c743d877f121b9debcf6a8a9d1806ce46.
//
// Solidity: event AssetRegistered(bytes32 indexed assetId, address indexed aggregator)
func (_OracleRegistry *OracleRegistryFilterer) FilterAssetRegistered(opts *bind.FilterOpts, assetId [][32]byte, aggregator []common.Address) (*OracleRegistryAssetRegisteredIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _OracleRegistry.contract.FilterLogs(opts, "AssetRegistered", assetIdRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &OracleRegistryAssetRegisteredIterator{contract: _OracleRegistry.contract, event: "AssetRegistered", logs: logs, sub: sub}, nil
}

// WatchAssetRegistered is a free log subscription operation binding the contract event 0x7b8c7b505365aa1b7f9ce04295e6da7c743d877f121b9debcf6a8a9d1806ce46.
//
// Solidity: event AssetRegistered(bytes32 indexed assetId, address indexed aggregator)
func (_OracleRegistry *OracleRegistryFilterer) WatchAssetRegistered(opts *bind.WatchOpts, sink chan<- *OracleRegistryAssetRegistered, assetId [][32]byte, aggregator []common.Address) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _OracleRegistry.contract.WatchLogs(opts, "AssetRegistered", assetIdRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleRegistryAssetRegistered)
				if err := _OracleRegistry.contract.UnpackLog(event, "AssetRegistered", log); err != nil {
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

// ParseAssetRegistered is a log parse operation binding the contract event 0x7b8c7b505365aa1b7f9ce04295e6da7c743d877f121b9debcf6a8a9d1806ce46.
//
// Solidity: event AssetRegistered(bytes32 indexed assetId, address indexed aggregator)
func (_OracleRegistry *OracleRegistryFilterer) ParseAssetRegistered(log types.Log) (*OracleRegistryAssetRegistered, error) {
	event := new(OracleRegistryAssetRegistered)
	if err := _OracleRegistry.contract.UnpackLog(event, "AssetRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleRegistryAssetUpdatedIterator is returned from FilterAssetUpdated and is used to iterate over the raw logs and unpacked data for AssetUpdated events raised by the OracleRegistry contract.
type OracleRegistryAssetUpdatedIterator struct {
	Event *OracleRegistryAssetUpdated // Event containing the contract specifics and raw log

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
func (it *OracleRegistryAssetUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleRegistryAssetUpdated)
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
		it.Event = new(OracleRegistryAssetUpdated)
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
func (it *OracleRegistryAssetUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleRegistryAssetUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleRegistryAssetUpdated represents a AssetUpdated event raised by the OracleRegistry contract.
type OracleRegistryAssetUpdated struct {
	AssetId       [32]byte
	OldAggregator common.Address
	NewAggregator common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAssetUpdated is a free log retrieval operation binding the contract event 0x0ce97bfad6771d8c3b0fc0c6a15cd488df2d92ff4a99e3354fc4f0d596034305.
//
// Solidity: event AssetUpdated(bytes32 indexed assetId, address indexed oldAggregator, address indexed newAggregator)
func (_OracleRegistry *OracleRegistryFilterer) FilterAssetUpdated(opts *bind.FilterOpts, assetId [][32]byte, oldAggregator []common.Address, newAggregator []common.Address) (*OracleRegistryAssetUpdatedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var oldAggregatorRule []interface{}
	for _, oldAggregatorItem := range oldAggregator {
		oldAggregatorRule = append(oldAggregatorRule, oldAggregatorItem)
	}
	var newAggregatorRule []interface{}
	for _, newAggregatorItem := range newAggregator {
		newAggregatorRule = append(newAggregatorRule, newAggregatorItem)
	}

	logs, sub, err := _OracleRegistry.contract.FilterLogs(opts, "AssetUpdated", assetIdRule, oldAggregatorRule, newAggregatorRule)
	if err != nil {
		return nil, err
	}
	return &OracleRegistryAssetUpdatedIterator{contract: _OracleRegistry.contract, event: "AssetUpdated", logs: logs, sub: sub}, nil
}

// WatchAssetUpdated is a free log subscription operation binding the contract event 0x0ce97bfad6771d8c3b0fc0c6a15cd488df2d92ff4a99e3354fc4f0d596034305.
//
// Solidity: event AssetUpdated(bytes32 indexed assetId, address indexed oldAggregator, address indexed newAggregator)
func (_OracleRegistry *OracleRegistryFilterer) WatchAssetUpdated(opts *bind.WatchOpts, sink chan<- *OracleRegistryAssetUpdated, assetId [][32]byte, oldAggregator []common.Address, newAggregator []common.Address) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var oldAggregatorRule []interface{}
	for _, oldAggregatorItem := range oldAggregator {
		oldAggregatorRule = append(oldAggregatorRule, oldAggregatorItem)
	}
	var newAggregatorRule []interface{}
	for _, newAggregatorItem := range newAggregator {
		newAggregatorRule = append(newAggregatorRule, newAggregatorItem)
	}

	logs, sub, err := _OracleRegistry.contract.WatchLogs(opts, "AssetUpdated", assetIdRule, oldAggregatorRule, newAggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleRegistryAssetUpdated)
				if err := _OracleRegistry.contract.UnpackLog(event, "AssetUpdated", log); err != nil {
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

// ParseAssetUpdated is a log parse operation binding the contract event 0x0ce97bfad6771d8c3b0fc0c6a15cd488df2d92ff4a99e3354fc4f0d596034305.
//
// Solidity: event AssetUpdated(bytes32 indexed assetId, address indexed oldAggregator, address indexed newAggregator)
func (_OracleRegistry *OracleRegistryFilterer) ParseAssetUpdated(log types.Log) (*OracleRegistryAssetUpdated, error) {
	event := new(OracleRegistryAssetUpdated)
	if err := _OracleRegistry.contract.UnpackLog(event, "AssetUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleRegistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the OracleRegistry contract.
type OracleRegistryOwnershipTransferStartedIterator struct {
	Event *OracleRegistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *OracleRegistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleRegistryOwnershipTransferStarted)
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
		it.Event = new(OracleRegistryOwnershipTransferStarted)
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
func (it *OracleRegistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleRegistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleRegistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the OracleRegistry contract.
type OracleRegistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_OracleRegistry *OracleRegistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleRegistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OracleRegistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OracleRegistryOwnershipTransferStartedIterator{contract: _OracleRegistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_OracleRegistry *OracleRegistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *OracleRegistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OracleRegistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleRegistryOwnershipTransferStarted)
				if err := _OracleRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_OracleRegistry *OracleRegistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*OracleRegistryOwnershipTransferStarted, error) {
	event := new(OracleRegistryOwnershipTransferStarted)
	if err := _OracleRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OracleRegistry contract.
type OracleRegistryOwnershipTransferredIterator struct {
	Event *OracleRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OracleRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleRegistryOwnershipTransferred)
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
		it.Event = new(OracleRegistryOwnershipTransferred)
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
func (it *OracleRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the OracleRegistry contract.
type OracleRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OracleRegistry *OracleRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OracleRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OracleRegistryOwnershipTransferredIterator{contract: _OracleRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OracleRegistry *OracleRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OracleRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OracleRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleRegistryOwnershipTransferred)
				if err := _OracleRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_OracleRegistry *OracleRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*OracleRegistryOwnershipTransferred, error) {
	event := new(OracleRegistryOwnershipTransferred)
	if err := _OracleRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
