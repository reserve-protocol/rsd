// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ReserveDollarEternalStorageABI is the input ABI used to generate the binding from.
const ReserveDollarEternalStorageABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"escapeHatch\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"addBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setAllowed\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newEscapeHatch\",\"type\":\"address\"}],\"name\":\"transferEscapeHatch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"frozenTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"},{\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"setFrozenTime\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"subBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"escapeHatchAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"oldEscapeHatch\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newEscapeHatch\",\"type\":\"address\"}],\"name\":\"EscapeHatchTransferred\",\"type\":\"event\"}]"

// ReserveDollarEternalStorageBin is the compiled bytecode used for deploying new contracts.
const ReserveDollarEternalStorageBin = `608060405234801561001057600080fd5b50604051602080610ac08339810180604052602081101561003057600080fd5b505160008054600160a060020a0319908116331790915560018054600160a060020a0390931692909116919091179055610a518061006f6000396000f3fe608060405234801561001057600080fd5b50600436106100f1576000357c010000000000000000000000000000000000000000000000000000000090048063b06230741161009e578063e30443bc11610078578063e30443bc146102d2578063e3d670d71461030b578063f2fde38b1461033e576100f1565b8063b06230741461022d578063b65dc41314610260578063cf8eeb7e14610299576100f1565b80635c658165116100cf5780635c658165146101a55780638babf203146101f25780638da5cb5b14610225576100f1565b80631554611f146100f657806321e5383a1461012757806333dd1b8a14610162575b600080fd5b6100fe610371565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101606004803603604081101561013d57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020013561038d565b005b6101606004803603606081101561017857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060400135610476565b6101e0600480360360408110156101bb57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516610535565b60408051918252519081900360200190f35b6101606004803603602081101561020857600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610552565b6100fe610666565b6101e06004803603602081101561024357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610682565b6101606004803603604081101561027657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610694565b610160600480360360408110156102af57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610743565b610160600480360360408110156102e857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356107ff565b6101e06004803603602081101561032157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166108ae565b6101606004803603602081101561035457600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166108c0565b60015473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff16331461041357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6f6e6c794f776e65720000000000000000000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8216600090815260026020526040902054610449908263ffffffff6109f716565b73ffffffffffffffffffffffffffffffffffffffff90921660009081526002602052604090209190915550565b60005473ffffffffffffffffffffffffffffffffffffffff1633146104fc57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6f6e6c794f776e65720000000000000000000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff92831660009081526003602090815260408083209490951682529290925291902055565b600360209081526000928352604080842090915290825290205481565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105d857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f6e6f7420617574686f72697a6564000000000000000000000000000000000000604482015290519081900360640190fd5b60015460405173ffffffffffffffffffffffffffffffffffffffff8084169216907f089af7288b55770a7c1dfd40b9d9e464c64031c45326c0916854814b6c16da2890600090a3600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60046020526000908152604090205481565b60005473ffffffffffffffffffffffffffffffffffffffff16331461071a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6f6e6c794f776e65720000000000000000000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff909116600090815260046020526040902055565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107c957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6f6e6c794f776e65720000000000000000000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff8216600090815260026020526040902054610449908263ffffffff610a1016565b60005473ffffffffffffffffffffffffffffffffffffffff16331461088557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6f6e6c794f776e65720000000000000000000000000000000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff909116600090815260026020526040902055565b60026020526000908152604090205481565b60005473ffffffffffffffffffffffffffffffffffffffff163314806108fd575060015473ffffffffffffffffffffffffffffffffffffffff1633145b151561096a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f6e6f7420617574686f72697a6564000000000000000000000000000000000000604482015290519081900360640190fd5b6000805460405173ffffffffffffffffffffffffffffffffffffffff808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600082820183811015610a0957600080fd5b9392505050565b600082821115610a1f57600080fd5b5090039056fea165627a7a723058202efbf6c1454b2837c001af8059908ddf7e3f6ca7f20a109d3d8f8d5c6c5ec1c60029`

// DeployReserveDollarEternalStorage deploys a new Ethereum contract, binding an instance of ReserveDollarEternalStorage to it.
func DeployReserveDollarEternalStorage(auth *bind.TransactOpts, backend bind.ContractBackend, escapeHatchAddress common.Address) (common.Address, *types.Transaction, *ReserveDollarEternalStorage, error) {
	parsed, err := abi.JSON(strings.NewReader(ReserveDollarEternalStorageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ReserveDollarEternalStorageBin), backend, escapeHatchAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReserveDollarEternalStorage{ReserveDollarEternalStorageCaller: ReserveDollarEternalStorageCaller{contract: contract}, ReserveDollarEternalStorageTransactor: ReserveDollarEternalStorageTransactor{contract: contract}, ReserveDollarEternalStorageFilterer: ReserveDollarEternalStorageFilterer{contract: contract}}, nil
}

// ReserveDollarEternalStorage is an auto generated Go binding around an Ethereum contract.
type ReserveDollarEternalStorage struct {
	ReserveDollarEternalStorageCaller     // Read-only binding to the contract
	ReserveDollarEternalStorageTransactor // Write-only binding to the contract
	ReserveDollarEternalStorageFilterer   // Log filterer for contract events
}

// ReserveDollarEternalStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReserveDollarEternalStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReserveDollarEternalStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReserveDollarEternalStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReserveDollarEternalStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReserveDollarEternalStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReserveDollarEternalStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReserveDollarEternalStorageSession struct {
	Contract     *ReserveDollarEternalStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ReserveDollarEternalStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReserveDollarEternalStorageCallerSession struct {
	Contract *ReserveDollarEternalStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// ReserveDollarEternalStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReserveDollarEternalStorageTransactorSession struct {
	Contract     *ReserveDollarEternalStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// ReserveDollarEternalStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReserveDollarEternalStorageRaw struct {
	Contract *ReserveDollarEternalStorage // Generic contract binding to access the raw methods on
}

// ReserveDollarEternalStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReserveDollarEternalStorageCallerRaw struct {
	Contract *ReserveDollarEternalStorageCaller // Generic read-only contract binding to access the raw methods on
}

// ReserveDollarEternalStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReserveDollarEternalStorageTransactorRaw struct {
	Contract *ReserveDollarEternalStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReserveDollarEternalStorage creates a new instance of ReserveDollarEternalStorage, bound to a specific deployed contract.
func NewReserveDollarEternalStorage(address common.Address, backend bind.ContractBackend) (*ReserveDollarEternalStorage, error) {
	contract, err := bindReserveDollarEternalStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReserveDollarEternalStorage{ReserveDollarEternalStorageCaller: ReserveDollarEternalStorageCaller{contract: contract}, ReserveDollarEternalStorageTransactor: ReserveDollarEternalStorageTransactor{contract: contract}, ReserveDollarEternalStorageFilterer: ReserveDollarEternalStorageFilterer{contract: contract}}, nil
}

// NewReserveDollarEternalStorageCaller creates a new read-only instance of ReserveDollarEternalStorage, bound to a specific deployed contract.
func NewReserveDollarEternalStorageCaller(address common.Address, caller bind.ContractCaller) (*ReserveDollarEternalStorageCaller, error) {
	contract, err := bindReserveDollarEternalStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReserveDollarEternalStorageCaller{contract: contract}, nil
}

// NewReserveDollarEternalStorageTransactor creates a new write-only instance of ReserveDollarEternalStorage, bound to a specific deployed contract.
func NewReserveDollarEternalStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*ReserveDollarEternalStorageTransactor, error) {
	contract, err := bindReserveDollarEternalStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReserveDollarEternalStorageTransactor{contract: contract}, nil
}

// NewReserveDollarEternalStorageFilterer creates a new log filterer instance of ReserveDollarEternalStorage, bound to a specific deployed contract.
func NewReserveDollarEternalStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*ReserveDollarEternalStorageFilterer, error) {
	contract, err := bindReserveDollarEternalStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReserveDollarEternalStorageFilterer{contract: contract}, nil
}

// bindReserveDollarEternalStorage binds a generic wrapper to an already deployed contract.
func bindReserveDollarEternalStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReserveDollarEternalStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ReserveDollarEternalStorage.Contract.ReserveDollarEternalStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.ReserveDollarEternalStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.ReserveDollarEternalStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ReserveDollarEternalStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.contract.Transact(opts, method, params...)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCaller) Allowed(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ReserveDollarEternalStorage.contract.Call(opts, out, "allowed", arg0, arg1)
	return *ret0, err
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ReserveDollarEternalStorage.Contract.Allowed(&_ReserveDollarEternalStorage.CallOpts, arg0, arg1)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCallerSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ReserveDollarEternalStorage.Contract.Allowed(&_ReserveDollarEternalStorage.CallOpts, arg0, arg1)
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance(address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCaller) Balance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ReserveDollarEternalStorage.contract.Call(opts, out, "balance", arg0)
	return *ret0, err
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance(address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) Balance(arg0 common.Address) (*big.Int, error) {
	return _ReserveDollarEternalStorage.Contract.Balance(&_ReserveDollarEternalStorage.CallOpts, arg0)
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance(address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCallerSession) Balance(arg0 common.Address) (*big.Int, error) {
	return _ReserveDollarEternalStorage.Contract.Balance(&_ReserveDollarEternalStorage.CallOpts, arg0)
}

// EscapeHatch is a free data retrieval call binding the contract method 0x1554611f.
//
// Solidity: function escapeHatch() constant returns(address)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCaller) EscapeHatch(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ReserveDollarEternalStorage.contract.Call(opts, out, "escapeHatch")
	return *ret0, err
}

// EscapeHatch is a free data retrieval call binding the contract method 0x1554611f.
//
// Solidity: function escapeHatch() constant returns(address)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) EscapeHatch() (common.Address, error) {
	return _ReserveDollarEternalStorage.Contract.EscapeHatch(&_ReserveDollarEternalStorage.CallOpts)
}

// EscapeHatch is a free data retrieval call binding the contract method 0x1554611f.
//
// Solidity: function escapeHatch() constant returns(address)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCallerSession) EscapeHatch() (common.Address, error) {
	return _ReserveDollarEternalStorage.Contract.EscapeHatch(&_ReserveDollarEternalStorage.CallOpts)
}

// FrozenTime is a free data retrieval call binding the contract method 0xb0623074.
//
// Solidity: function frozenTime(address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCaller) FrozenTime(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ReserveDollarEternalStorage.contract.Call(opts, out, "frozenTime", arg0)
	return *ret0, err
}

// FrozenTime is a free data retrieval call binding the contract method 0xb0623074.
//
// Solidity: function frozenTime(address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) FrozenTime(arg0 common.Address) (*big.Int, error) {
	return _ReserveDollarEternalStorage.Contract.FrozenTime(&_ReserveDollarEternalStorage.CallOpts, arg0)
}

// FrozenTime is a free data retrieval call binding the contract method 0xb0623074.
//
// Solidity: function frozenTime(address ) constant returns(uint256)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCallerSession) FrozenTime(arg0 common.Address) (*big.Int, error) {
	return _ReserveDollarEternalStorage.Contract.FrozenTime(&_ReserveDollarEternalStorage.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ReserveDollarEternalStorage.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) Owner() (common.Address, error) {
	return _ReserveDollarEternalStorage.Contract.Owner(&_ReserveDollarEternalStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageCallerSession) Owner() (common.Address, error) {
	return _ReserveDollarEternalStorage.Contract.Owner(&_ReserveDollarEternalStorage.CallOpts)
}

// AddBalance is a paid mutator transaction binding the contract method 0x21e5383a.
//
// Solidity: function addBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactor) AddBalance(opts *bind.TransactOpts, key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.contract.Transact(opts, "addBalance", key, value)
}

// AddBalance is a paid mutator transaction binding the contract method 0x21e5383a.
//
// Solidity: function addBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) AddBalance(key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.AddBalance(&_ReserveDollarEternalStorage.TransactOpts, key, value)
}

// AddBalance is a paid mutator transaction binding the contract method 0x21e5383a.
//
// Solidity: function addBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorSession) AddBalance(key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.AddBalance(&_ReserveDollarEternalStorage.TransactOpts, key, value)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x33dd1b8a.
//
// Solidity: function setAllowed(address from, address to, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactor) SetAllowed(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.contract.Transact(opts, "setAllowed", from, to, value)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x33dd1b8a.
//
// Solidity: function setAllowed(address from, address to, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) SetAllowed(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SetAllowed(&_ReserveDollarEternalStorage.TransactOpts, from, to, value)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x33dd1b8a.
//
// Solidity: function setAllowed(address from, address to, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorSession) SetAllowed(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SetAllowed(&_ReserveDollarEternalStorage.TransactOpts, from, to, value)
}

// SetBalance is a paid mutator transaction binding the contract method 0xe30443bc.
//
// Solidity: function setBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactor) SetBalance(opts *bind.TransactOpts, key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.contract.Transact(opts, "setBalance", key, value)
}

// SetBalance is a paid mutator transaction binding the contract method 0xe30443bc.
//
// Solidity: function setBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) SetBalance(key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SetBalance(&_ReserveDollarEternalStorage.TransactOpts, key, value)
}

// SetBalance is a paid mutator transaction binding the contract method 0xe30443bc.
//
// Solidity: function setBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorSession) SetBalance(key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SetBalance(&_ReserveDollarEternalStorage.TransactOpts, key, value)
}

// SetFrozenTime is a paid mutator transaction binding the contract method 0xb65dc413.
//
// Solidity: function setFrozenTime(address who, uint256 time) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactor) SetFrozenTime(opts *bind.TransactOpts, who common.Address, time *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.contract.Transact(opts, "setFrozenTime", who, time)
}

// SetFrozenTime is a paid mutator transaction binding the contract method 0xb65dc413.
//
// Solidity: function setFrozenTime(address who, uint256 time) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) SetFrozenTime(who common.Address, time *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SetFrozenTime(&_ReserveDollarEternalStorage.TransactOpts, who, time)
}

// SetFrozenTime is a paid mutator transaction binding the contract method 0xb65dc413.
//
// Solidity: function setFrozenTime(address who, uint256 time) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorSession) SetFrozenTime(who common.Address, time *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SetFrozenTime(&_ReserveDollarEternalStorage.TransactOpts, who, time)
}

// SubBalance is a paid mutator transaction binding the contract method 0xcf8eeb7e.
//
// Solidity: function subBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactor) SubBalance(opts *bind.TransactOpts, key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.contract.Transact(opts, "subBalance", key, value)
}

// SubBalance is a paid mutator transaction binding the contract method 0xcf8eeb7e.
//
// Solidity: function subBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) SubBalance(key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SubBalance(&_ReserveDollarEternalStorage.TransactOpts, key, value)
}

// SubBalance is a paid mutator transaction binding the contract method 0xcf8eeb7e.
//
// Solidity: function subBalance(address key, uint256 value) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorSession) SubBalance(key common.Address, value *big.Int) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.SubBalance(&_ReserveDollarEternalStorage.TransactOpts, key, value)
}

// TransferEscapeHatch is a paid mutator transaction binding the contract method 0x8babf203.
//
// Solidity: function transferEscapeHatch(address newEscapeHatch) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactor) TransferEscapeHatch(opts *bind.TransactOpts, newEscapeHatch common.Address) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.contract.Transact(opts, "transferEscapeHatch", newEscapeHatch)
}

// TransferEscapeHatch is a paid mutator transaction binding the contract method 0x8babf203.
//
// Solidity: function transferEscapeHatch(address newEscapeHatch) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) TransferEscapeHatch(newEscapeHatch common.Address) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.TransferEscapeHatch(&_ReserveDollarEternalStorage.TransactOpts, newEscapeHatch)
}

// TransferEscapeHatch is a paid mutator transaction binding the contract method 0x8babf203.
//
// Solidity: function transferEscapeHatch(address newEscapeHatch) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorSession) TransferEscapeHatch(newEscapeHatch common.Address) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.TransferEscapeHatch(&_ReserveDollarEternalStorage.TransactOpts, newEscapeHatch)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.TransferOwnership(&_ReserveDollarEternalStorage.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReserveDollarEternalStorage.Contract.TransferOwnership(&_ReserveDollarEternalStorage.TransactOpts, newOwner)
}

// ReserveDollarEternalStorageEscapeHatchTransferredIterator is returned from FilterEscapeHatchTransferred and is used to iterate over the raw logs and unpacked data for EscapeHatchTransferred events raised by the ReserveDollarEternalStorage contract.
type ReserveDollarEternalStorageEscapeHatchTransferredIterator struct {
	Event *ReserveDollarEternalStorageEscapeHatchTransferred // Event containing the contract specifics and raw log

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
func (it *ReserveDollarEternalStorageEscapeHatchTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveDollarEternalStorageEscapeHatchTransferred)
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
		it.Event = new(ReserveDollarEternalStorageEscapeHatchTransferred)
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
func (it *ReserveDollarEternalStorageEscapeHatchTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveDollarEternalStorageEscapeHatchTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveDollarEternalStorageEscapeHatchTransferred represents a EscapeHatchTransferred event raised by the ReserveDollarEternalStorage contract.
type ReserveDollarEternalStorageEscapeHatchTransferred struct {
	OldEscapeHatch common.Address
	NewEscapeHatch common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterEscapeHatchTransferred is a free log retrieval operation binding the contract event 0x089af7288b55770a7c1dfd40b9d9e464c64031c45326c0916854814b6c16da28.
//
// Solidity: event EscapeHatchTransferred(address indexed oldEscapeHatch, address indexed newEscapeHatch)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageFilterer) FilterEscapeHatchTransferred(opts *bind.FilterOpts, oldEscapeHatch []common.Address, newEscapeHatch []common.Address) (*ReserveDollarEternalStorageEscapeHatchTransferredIterator, error) {

	var oldEscapeHatchRule []interface{}
	for _, oldEscapeHatchItem := range oldEscapeHatch {
		oldEscapeHatchRule = append(oldEscapeHatchRule, oldEscapeHatchItem)
	}
	var newEscapeHatchRule []interface{}
	for _, newEscapeHatchItem := range newEscapeHatch {
		newEscapeHatchRule = append(newEscapeHatchRule, newEscapeHatchItem)
	}

	logs, sub, err := _ReserveDollarEternalStorage.contract.FilterLogs(opts, "EscapeHatchTransferred", oldEscapeHatchRule, newEscapeHatchRule)
	if err != nil {
		return nil, err
	}
	return &ReserveDollarEternalStorageEscapeHatchTransferredIterator{contract: _ReserveDollarEternalStorage.contract, event: "EscapeHatchTransferred", logs: logs, sub: sub}, nil
}

// WatchEscapeHatchTransferred is a free log subscription operation binding the contract event 0x089af7288b55770a7c1dfd40b9d9e464c64031c45326c0916854814b6c16da28.
//
// Solidity: event EscapeHatchTransferred(address indexed oldEscapeHatch, address indexed newEscapeHatch)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageFilterer) WatchEscapeHatchTransferred(opts *bind.WatchOpts, sink chan<- *ReserveDollarEternalStorageEscapeHatchTransferred, oldEscapeHatch []common.Address, newEscapeHatch []common.Address) (event.Subscription, error) {

	var oldEscapeHatchRule []interface{}
	for _, oldEscapeHatchItem := range oldEscapeHatch {
		oldEscapeHatchRule = append(oldEscapeHatchRule, oldEscapeHatchItem)
	}
	var newEscapeHatchRule []interface{}
	for _, newEscapeHatchItem := range newEscapeHatch {
		newEscapeHatchRule = append(newEscapeHatchRule, newEscapeHatchItem)
	}

	logs, sub, err := _ReserveDollarEternalStorage.contract.WatchLogs(opts, "EscapeHatchTransferred", oldEscapeHatchRule, newEscapeHatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveDollarEternalStorageEscapeHatchTransferred)
				if err := _ReserveDollarEternalStorage.contract.UnpackLog(event, "EscapeHatchTransferred", log); err != nil {
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

// ReserveDollarEternalStorageOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ReserveDollarEternalStorage contract.
type ReserveDollarEternalStorageOwnershipTransferredIterator struct {
	Event *ReserveDollarEternalStorageOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ReserveDollarEternalStorageOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReserveDollarEternalStorageOwnershipTransferred)
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
		it.Event = new(ReserveDollarEternalStorageOwnershipTransferred)
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
func (it *ReserveDollarEternalStorageOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReserveDollarEternalStorageOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReserveDollarEternalStorageOwnershipTransferred represents a OwnershipTransferred event raised by the ReserveDollarEternalStorage contract.
type ReserveDollarEternalStorageOwnershipTransferred struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*ReserveDollarEternalStorageOwnershipTransferredIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReserveDollarEternalStorage.contract.FilterLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ReserveDollarEternalStorageOwnershipTransferredIterator{contract: _ReserveDollarEternalStorage.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_ReserveDollarEternalStorage *ReserveDollarEternalStorageFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ReserveDollarEternalStorageOwnershipTransferred, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReserveDollarEternalStorage.contract.WatchLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReserveDollarEternalStorageOwnershipTransferred)
				if err := _ReserveDollarEternalStorage.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
