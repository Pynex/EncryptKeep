// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

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

// KeeperMetaData contains all meta data concerning the Keeper contract.
var KeeperMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"activeIdsForUser\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"changeData\",\"inputs\":[{\"name\":\"_id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_newData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"nextDataId\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeData\",\"inputs\":[{\"name\":\"_id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"storeData\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"storeMetaData\",\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"userData\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userMetaData\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"CannotChangeNonExistentData\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"CannotRemoveNonExistentData\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"CannotStoreExistingData\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]}]",
}

// KeeperABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperMetaData.ABI instead.
var KeeperABI = KeeperMetaData.ABI

// Keeper is an auto generated Go binding around an Ethereum contract.
type Keeper struct {
	KeeperCaller     // Read-only binding to the contract
	KeeperTransactor // Write-only binding to the contract
	KeeperFilterer   // Log filterer for contract events
}

// KeeperCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperSession struct {
	Contract     *Keeper           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperCallerSession struct {
	Contract *KeeperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// KeeperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperTransactorSession struct {
	Contract     *KeeperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRaw struct {
	Contract *Keeper // Generic contract binding to access the raw methods on
}

// KeeperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperCallerRaw struct {
	Contract *KeeperCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperTransactorRaw struct {
	Contract *KeeperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeper creates a new instance of Keeper, bound to a specific deployed contract.
func NewKeeper(address common.Address, backend bind.ContractBackend) (*Keeper, error) {
	contract, err := bindKeeper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Keeper{KeeperCaller: KeeperCaller{contract: contract}, KeeperTransactor: KeeperTransactor{contract: contract}, KeeperFilterer: KeeperFilterer{contract: contract}}, nil
}

// NewKeeperCaller creates a new read-only instance of Keeper, bound to a specific deployed contract.
func NewKeeperCaller(address common.Address, caller bind.ContractCaller) (*KeeperCaller, error) {
	contract, err := bindKeeper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperCaller{contract: contract}, nil
}

// NewKeeperTransactor creates a new write-only instance of Keeper, bound to a specific deployed contract.
func NewKeeperTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperTransactor, error) {
	contract, err := bindKeeper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperTransactor{contract: contract}, nil
}

// NewKeeperFilterer creates a new log filterer instance of Keeper, bound to a specific deployed contract.
func NewKeeperFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperFilterer, error) {
	contract, err := bindKeeper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperFilterer{contract: contract}, nil
}

// bindKeeper binds a generic wrapper to an already deployed contract.
func bindKeeper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Keeper *KeeperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Keeper.Contract.KeeperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Keeper *KeeperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Keeper.Contract.KeeperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Keeper *KeeperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Keeper.Contract.KeeperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Keeper *KeeperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Keeper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Keeper *KeeperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Keeper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Keeper *KeeperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Keeper.Contract.contract.Transact(opts, method, params...)
}

// ActiveIdsForUser is a free data retrieval call binding the contract method 0x0592f5d5.
//
// Solidity: function activeIdsForUser(address , uint256 ) view returns(uint256)
func (_Keeper *KeeperCaller) ActiveIdsForUser(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Keeper.contract.Call(opts, &out, "activeIdsForUser", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveIdsForUser is a free data retrieval call binding the contract method 0x0592f5d5.
//
// Solidity: function activeIdsForUser(address , uint256 ) view returns(uint256)
func (_Keeper *KeeperSession) ActiveIdsForUser(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Keeper.Contract.ActiveIdsForUser(&_Keeper.CallOpts, arg0, arg1)
}

// ActiveIdsForUser is a free data retrieval call binding the contract method 0x0592f5d5.
//
// Solidity: function activeIdsForUser(address , uint256 ) view returns(uint256)
func (_Keeper *KeeperCallerSession) ActiveIdsForUser(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Keeper.Contract.ActiveIdsForUser(&_Keeper.CallOpts, arg0, arg1)
}

// NextDataId is a free data retrieval call binding the contract method 0x63ee461d.
//
// Solidity: function nextDataId(address ) view returns(uint256)
func (_Keeper *KeeperCaller) NextDataId(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Keeper.contract.Call(opts, &out, "nextDataId", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextDataId is a free data retrieval call binding the contract method 0x63ee461d.
//
// Solidity: function nextDataId(address ) view returns(uint256)
func (_Keeper *KeeperSession) NextDataId(arg0 common.Address) (*big.Int, error) {
	return _Keeper.Contract.NextDataId(&_Keeper.CallOpts, arg0)
}

// NextDataId is a free data retrieval call binding the contract method 0x63ee461d.
//
// Solidity: function nextDataId(address ) view returns(uint256)
func (_Keeper *KeeperCallerSession) NextDataId(arg0 common.Address) (*big.Int, error) {
	return _Keeper.Contract.NextDataId(&_Keeper.CallOpts, arg0)
}

// UserData is a free data retrieval call binding the contract method 0x3c05eca1.
//
// Solidity: function userData(address , uint256 ) view returns(bytes)
func (_Keeper *KeeperCaller) UserData(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Keeper.contract.Call(opts, &out, "userData", arg0, arg1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// UserData is a free data retrieval call binding the contract method 0x3c05eca1.
//
// Solidity: function userData(address , uint256 ) view returns(bytes)
func (_Keeper *KeeperSession) UserData(arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	return _Keeper.Contract.UserData(&_Keeper.CallOpts, arg0, arg1)
}

// UserData is a free data retrieval call binding the contract method 0x3c05eca1.
//
// Solidity: function userData(address , uint256 ) view returns(bytes)
func (_Keeper *KeeperCallerSession) UserData(arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	return _Keeper.Contract.UserData(&_Keeper.CallOpts, arg0, arg1)
}

// UserMetaData is a free data retrieval call binding the contract method 0x6192b6b0.
//
// Solidity: function userMetaData(address ) view returns(bytes)
func (_Keeper *KeeperCaller) UserMetaData(opts *bind.CallOpts, arg0 common.Address) ([]byte, error) {
	var out []interface{}
	err := _Keeper.contract.Call(opts, &out, "userMetaData", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// UserMetaData is a free data retrieval call binding the contract method 0x6192b6b0.
//
// Solidity: function userMetaData(address ) view returns(bytes)
func (_Keeper *KeeperSession) UserMetaData(arg0 common.Address) ([]byte, error) {
	return _Keeper.Contract.UserMetaData(&_Keeper.CallOpts, arg0)
}

// UserMetaData is a free data retrieval call binding the contract method 0x6192b6b0.
//
// Solidity: function userMetaData(address ) view returns(bytes)
func (_Keeper *KeeperCallerSession) UserMetaData(arg0 common.Address) ([]byte, error) {
	return _Keeper.Contract.UserMetaData(&_Keeper.CallOpts, arg0)
}

// ChangeData is a paid mutator transaction binding the contract method 0xf2836502.
//
// Solidity: function changeData(uint256 _id, bytes _newData) payable returns()
func (_Keeper *KeeperTransactor) ChangeData(opts *bind.TransactOpts, _id *big.Int, _newData []byte) (*types.Transaction, error) {
	return _Keeper.contract.Transact(opts, "changeData", _id, _newData)
}

// ChangeData is a paid mutator transaction binding the contract method 0xf2836502.
//
// Solidity: function changeData(uint256 _id, bytes _newData) payable returns()
func (_Keeper *KeeperSession) ChangeData(_id *big.Int, _newData []byte) (*types.Transaction, error) {
	return _Keeper.Contract.ChangeData(&_Keeper.TransactOpts, _id, _newData)
}

// ChangeData is a paid mutator transaction binding the contract method 0xf2836502.
//
// Solidity: function changeData(uint256 _id, bytes _newData) payable returns()
func (_Keeper *KeeperTransactorSession) ChangeData(_id *big.Int, _newData []byte) (*types.Transaction, error) {
	return _Keeper.Contract.ChangeData(&_Keeper.TransactOpts, _id, _newData)
}

// RemoveData is a paid mutator transaction binding the contract method 0xa94840bb.
//
// Solidity: function removeData(uint256 _id) payable returns()
func (_Keeper *KeeperTransactor) RemoveData(opts *bind.TransactOpts, _id *big.Int) (*types.Transaction, error) {
	return _Keeper.contract.Transact(opts, "removeData", _id)
}

// RemoveData is a paid mutator transaction binding the contract method 0xa94840bb.
//
// Solidity: function removeData(uint256 _id) payable returns()
func (_Keeper *KeeperSession) RemoveData(_id *big.Int) (*types.Transaction, error) {
	return _Keeper.Contract.RemoveData(&_Keeper.TransactOpts, _id)
}

// RemoveData is a paid mutator transaction binding the contract method 0xa94840bb.
//
// Solidity: function removeData(uint256 _id) payable returns()
func (_Keeper *KeeperTransactorSession) RemoveData(_id *big.Int) (*types.Transaction, error) {
	return _Keeper.Contract.RemoveData(&_Keeper.TransactOpts, _id)
}

// StoreData is a paid mutator transaction binding the contract method 0xac5c8535.
//
// Solidity: function storeData(bytes _data) payable returns()
func (_Keeper *KeeperTransactor) StoreData(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _Keeper.contract.Transact(opts, "storeData", _data)
}

// StoreData is a paid mutator transaction binding the contract method 0xac5c8535.
//
// Solidity: function storeData(bytes _data) payable returns()
func (_Keeper *KeeperSession) StoreData(_data []byte) (*types.Transaction, error) {
	return _Keeper.Contract.StoreData(&_Keeper.TransactOpts, _data)
}

// StoreData is a paid mutator transaction binding the contract method 0xac5c8535.
//
// Solidity: function storeData(bytes _data) payable returns()
func (_Keeper *KeeperTransactorSession) StoreData(_data []byte) (*types.Transaction, error) {
	return _Keeper.Contract.StoreData(&_Keeper.TransactOpts, _data)
}

// StoreMetaData is a paid mutator transaction binding the contract method 0xd33e9b27.
//
// Solidity: function storeMetaData(bytes _data) payable returns()
func (_Keeper *KeeperTransactor) StoreMetaData(opts *bind.TransactOpts, _data []byte) (*types.Transaction, error) {
	return _Keeper.contract.Transact(opts, "storeMetaData", _data)
}

// StoreMetaData is a paid mutator transaction binding the contract method 0xd33e9b27.
//
// Solidity: function storeMetaData(bytes _data) payable returns()
func (_Keeper *KeeperSession) StoreMetaData(_data []byte) (*types.Transaction, error) {
	return _Keeper.Contract.StoreMetaData(&_Keeper.TransactOpts, _data)
}

// StoreMetaData is a paid mutator transaction binding the contract method 0xd33e9b27.
//
// Solidity: function storeMetaData(bytes _data) payable returns()
func (_Keeper *KeeperTransactorSession) StoreMetaData(_data []byte) (*types.Transaction, error) {
	return _Keeper.Contract.StoreMetaData(&_Keeper.TransactOpts, _data)
}
