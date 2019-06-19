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

// MintAndBurnAdminABI is the input ABI used to generate the binding from.
const MintAndBurnAdminABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"proposals\",\"outputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"isMint\",\"type\":\"bool\"},{\"name\":\"time\",\"type\":\"uint256\"},{\"name\":\"closed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"cancelAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"},{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"isMint\",\"type\":\"bool\"}],\"name\":\"cancel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"isMint\",\"type\":\"bool\"}],\"name\":\"propose\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delay\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"reserve\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"},{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"isMint\",\"type\":\"bool\"}],\"name\":\"confirm\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"reserveDollar\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"isMint\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"delayUntil\",\"type\":\"uint256\"}],\"name\":\"ProposalCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"isMint\",\"type\":\"bool\"}],\"name\":\"ProposalConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"isMint\",\"type\":\"bool\"}],\"name\":\"ProposalCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"AllProposalsCancelled\",\"type\":\"event\"}]"

// MintAndBurnAdminBin is the compiled bytecode used for deploying new contracts.
const MintAndBurnAdminBin = `608060405234801561001057600080fd5b50604051602080610e3c8339810180604052602081101561003057600080fd5b505160008054600160a060020a03909216600160a060020a03199283161790556001805490911633179055610dd28061006a6000396000f3fe608060405234801561001057600080fd5b50600436106100a5576000357c0100000000000000000000000000000000000000000000000000000000900480636a42b8f8116100785780636a42b8f81461019d578063cd3293de146101b7578063dba82a45146101e8578063f851a4401461022f576100a5565b8063013cf08b146100aa57806318cb2b181461010b57806334a0f49a146101155780635e77e64e1461015c575b600080fd5b6100c7600480360360208110156100c057600080fd5b5035610237565b6040805173ffffffffffffffffffffffffffffffffffffffff9096168652602086019490945291151584840152606084015215156080830152519081900360a00190f35b610113610294565b005b6101136004803603608081101561012b57600080fd5b5080359073ffffffffffffffffffffffffffffffffffffffff60208201351690604081013590606001351515610353565b6101136004803603606081101561017257600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602081013590604001351515610498565b6101a5610715565b60408051918252519081900360200190f35b6101bf61071b565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610113600480360360808110156101fe57600080fd5b5080359073ffffffffffffffffffffffffffffffffffffffff60208201351690604081013590606001351515610737565b6101bf610a5c565b600280548290811061024557fe5b60009182526020909120600590910201805460018201546002830154600384015460049094015473ffffffffffffffffffffffffffffffffffffffff9093169450909260ff9182169290911685565b60015473ffffffffffffffffffffffffffffffffffffffff16331461031a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6d7573742062652061646d696e00000000000000000000000000000000000000604482015290519081900360640190fd5b6000610327600282610ce9565b506040517f3732302b0efc3e1e883bb80d83c641031dc1e32223cb406c3e4d5de68208c4e990600090a1565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103d957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6d7573742062652061646d696e00000000000000000000000000000000000000604482015290519081900360640190fd5b6103e584848484610a78565b60016002858154811015156103f657fe5b600091825260209182902060046005909202010180549215157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009093169290921790915560408051868152918201849052821515828201525173ffffffffffffffffffffffffffffffffffffffff8516917fc1ea9ad7fe3cfb48a741fc229353411aabb3b135d9446697bf6db7c197a9ac0f919081900360600190a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461051e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6d7573742062652061646d696e00000000000000000000000000000000000000604482015290519081900360640190fd5b6040805160a08101825273ffffffffffffffffffffffffffffffffffffffff858116808352602080840187815286151585870181815261a8c04201606080890182815260006080808c0182815260028054600181018255938190529c517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace600590940293840180547fffffffffffffffffffffffff00000000000000000000000000000000000000001691909d1617909b5596517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf82015593517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad0850180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090811692151592909217905590517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad185015597517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad2909301805490981692151592909217909655955487517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff919091018152928301899052828701529381018390529351919390927fd1d2eb762bbbecbc03b8a9dd22368874018771d0c93d855cd08c5a8fa6086b9692918290030190a250505050565b61a8c081565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60015473ffffffffffffffffffffffffffffffffffffffff1633146107bd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f6d7573742062652061646d696e00000000000000000000000000000000000000604482015290519081900360640190fd5b6107c984848484610a78565b426002858154811015156107d957fe5b90600052602060002090600502016003015410151561085957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f746f6f206561726c790000000000000000000000000000000000000000000000604482015290519081900360640190fd5b600160028581548110151561086a57fe5b600091825260209182902060046005909202010180549215157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff009093169290921790915560408051868152918201849052821515828201525173ffffffffffffffffffffffffffffffffffffffff8516917fc398e86b1dfd2596a48f97df67ac573ef31251ea5b65d73e4096be478df18f57919081900360600190a2600280548590811061091457fe5b600091825260209091206002600590920201015460ff16156109c55760008054604080517f40c10f1900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015260248201879052915191909216926340c10f19926044808201939182900301818387803b1580156109a857600080fd5b505af11580156109bc573d6000803e3d6000fd5b50505050610a56565b60008054604080517f79cc679000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015260248201879052915191909216926379cc6790926044808201939182900301818387803b158015610a3d57600080fd5b505af1158015610a51573d6000803e3d6000fd5b505050505b50505050565b60015473ffffffffffffffffffffffffffffffffffffffff1681565b6002805485908110610a8657fe5b600091825260209091206004600590920201015460ff1615610b0957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f70726f706f73616c20616c726561647920636c6f736564000000000000000000604482015290519081900360640190fd5b8273ffffffffffffffffffffffffffffffffffffffff16600285815481101515610b2f57fe5b600091825260209091206005909102015473ffffffffffffffffffffffffffffffffffffffff1614610bc257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f61646472206d69736d6174636865640000000000000000000000000000000000604482015290519081900360640190fd5b81600285815481101515610bd257fe5b906000526020600020906005020160010154141515610c5257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f76616c7565206d69736d61746368656400000000000000000000000000000000604482015290519081900360640190fd5b801515600285815481101515610c6457fe5b600091825260209091206002600590920201015460ff16151514610a5657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f69734d696e74206d69736d617463686564000000000000000000000000000000604482015290519081900360640190fd5b815481835581811115610d1557600502816005028360005260206000209182019101610d159190610d1a565b505050565b610da391905b80821115610d9f5780547fffffffffffffffffffffffff00000000000000000000000000000000000000001681556000600182018190556002820180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0090811690915560038301919091556004820180549091169055600501610d20565b5090565b9056fea165627a7a72305820828a7b96423c59dd6b7aef88912972d5d833e3ac6ef35c67efa559f6fd0d57db0029`

// DeployMintAndBurnAdmin deploys a new Ethereum contract, binding an instance of MintAndBurnAdmin to it.
func DeployMintAndBurnAdmin(auth *bind.TransactOpts, backend bind.ContractBackend, reserveDollar common.Address) (common.Address, *types.Transaction, *MintAndBurnAdmin, error) {
	parsed, err := abi.JSON(strings.NewReader(MintAndBurnAdminABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MintAndBurnAdminBin), backend, reserveDollar)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MintAndBurnAdmin{MintAndBurnAdminCaller: MintAndBurnAdminCaller{contract: contract}, MintAndBurnAdminTransactor: MintAndBurnAdminTransactor{contract: contract}, MintAndBurnAdminFilterer: MintAndBurnAdminFilterer{contract: contract}}, nil
}

// MintAndBurnAdmin is an auto generated Go binding around an Ethereum contract.
type MintAndBurnAdmin struct {
	MintAndBurnAdminCaller     // Read-only binding to the contract
	MintAndBurnAdminTransactor // Write-only binding to the contract
	MintAndBurnAdminFilterer   // Log filterer for contract events
}

// MintAndBurnAdminCaller is an auto generated read-only Go binding around an Ethereum contract.
type MintAndBurnAdminCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintAndBurnAdminTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MintAndBurnAdminTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintAndBurnAdminFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MintAndBurnAdminFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintAndBurnAdminSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MintAndBurnAdminSession struct {
	Contract     *MintAndBurnAdmin // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MintAndBurnAdminCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MintAndBurnAdminCallerSession struct {
	Contract *MintAndBurnAdminCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// MintAndBurnAdminTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MintAndBurnAdminTransactorSession struct {
	Contract     *MintAndBurnAdminTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// MintAndBurnAdminRaw is an auto generated low-level Go binding around an Ethereum contract.
type MintAndBurnAdminRaw struct {
	Contract *MintAndBurnAdmin // Generic contract binding to access the raw methods on
}

// MintAndBurnAdminCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MintAndBurnAdminCallerRaw struct {
	Contract *MintAndBurnAdminCaller // Generic read-only contract binding to access the raw methods on
}

// MintAndBurnAdminTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MintAndBurnAdminTransactorRaw struct {
	Contract *MintAndBurnAdminTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMintAndBurnAdmin creates a new instance of MintAndBurnAdmin, bound to a specific deployed contract.
func NewMintAndBurnAdmin(address common.Address, backend bind.ContractBackend) (*MintAndBurnAdmin, error) {
	contract, err := bindMintAndBurnAdmin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdmin{MintAndBurnAdminCaller: MintAndBurnAdminCaller{contract: contract}, MintAndBurnAdminTransactor: MintAndBurnAdminTransactor{contract: contract}, MintAndBurnAdminFilterer: MintAndBurnAdminFilterer{contract: contract}}, nil
}

// NewMintAndBurnAdminCaller creates a new read-only instance of MintAndBurnAdmin, bound to a specific deployed contract.
func NewMintAndBurnAdminCaller(address common.Address, caller bind.ContractCaller) (*MintAndBurnAdminCaller, error) {
	contract, err := bindMintAndBurnAdmin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdminCaller{contract: contract}, nil
}

// NewMintAndBurnAdminTransactor creates a new write-only instance of MintAndBurnAdmin, bound to a specific deployed contract.
func NewMintAndBurnAdminTransactor(address common.Address, transactor bind.ContractTransactor) (*MintAndBurnAdminTransactor, error) {
	contract, err := bindMintAndBurnAdmin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdminTransactor{contract: contract}, nil
}

// NewMintAndBurnAdminFilterer creates a new log filterer instance of MintAndBurnAdmin, bound to a specific deployed contract.
func NewMintAndBurnAdminFilterer(address common.Address, filterer bind.ContractFilterer) (*MintAndBurnAdminFilterer, error) {
	contract, err := bindMintAndBurnAdmin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdminFilterer{contract: contract}, nil
}

// bindMintAndBurnAdmin binds a generic wrapper to an already deployed contract.
func bindMintAndBurnAdmin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MintAndBurnAdminABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintAndBurnAdmin *MintAndBurnAdminRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MintAndBurnAdmin.Contract.MintAndBurnAdminCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintAndBurnAdmin *MintAndBurnAdminRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.MintAndBurnAdminTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintAndBurnAdmin *MintAndBurnAdminRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.MintAndBurnAdminTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintAndBurnAdmin *MintAndBurnAdminCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MintAndBurnAdmin.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintAndBurnAdmin *MintAndBurnAdminTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintAndBurnAdmin *MintAndBurnAdminTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_MintAndBurnAdmin *MintAndBurnAdminCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MintAndBurnAdmin.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_MintAndBurnAdmin *MintAndBurnAdminSession) Admin() (common.Address, error) {
	return _MintAndBurnAdmin.Contract.Admin(&_MintAndBurnAdmin.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_MintAndBurnAdmin *MintAndBurnAdminCallerSession) Admin() (common.Address, error) {
	return _MintAndBurnAdmin.Contract.Admin(&_MintAndBurnAdmin.CallOpts)
}

// Delay is a free data retrieval call binding the contract method 0x6a42b8f8.
//
// Solidity: function delay() constant returns(uint256)
func (_MintAndBurnAdmin *MintAndBurnAdminCaller) Delay(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MintAndBurnAdmin.contract.Call(opts, out, "delay")
	return *ret0, err
}

// Delay is a free data retrieval call binding the contract method 0x6a42b8f8.
//
// Solidity: function delay() constant returns(uint256)
func (_MintAndBurnAdmin *MintAndBurnAdminSession) Delay() (*big.Int, error) {
	return _MintAndBurnAdmin.Contract.Delay(&_MintAndBurnAdmin.CallOpts)
}

// Delay is a free data retrieval call binding the contract method 0x6a42b8f8.
//
// Solidity: function delay() constant returns(uint256)
func (_MintAndBurnAdmin *MintAndBurnAdminCallerSession) Delay() (*big.Int, error) {
	return _MintAndBurnAdmin.Contract.Delay(&_MintAndBurnAdmin.CallOpts)
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) constant returns(address addr, uint256 value, bool isMint, uint256 time, bool closed)
func (_MintAndBurnAdmin *MintAndBurnAdminCaller) Proposals(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Addr   common.Address
	Value  *big.Int
	IsMint bool
	Time   *big.Int
	Closed bool
}, error) {
	ret := new(struct {
		Addr   common.Address
		Value  *big.Int
		IsMint bool
		Time   *big.Int
		Closed bool
	})
	out := ret
	err := _MintAndBurnAdmin.contract.Call(opts, out, "proposals", arg0)
	return *ret, err
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) constant returns(address addr, uint256 value, bool isMint, uint256 time, bool closed)
func (_MintAndBurnAdmin *MintAndBurnAdminSession) Proposals(arg0 *big.Int) (struct {
	Addr   common.Address
	Value  *big.Int
	IsMint bool
	Time   *big.Int
	Closed bool
}, error) {
	return _MintAndBurnAdmin.Contract.Proposals(&_MintAndBurnAdmin.CallOpts, arg0)
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 ) constant returns(address addr, uint256 value, bool isMint, uint256 time, bool closed)
func (_MintAndBurnAdmin *MintAndBurnAdminCallerSession) Proposals(arg0 *big.Int) (struct {
	Addr   common.Address
	Value  *big.Int
	IsMint bool
	Time   *big.Int
	Closed bool
}, error) {
	return _MintAndBurnAdmin.Contract.Proposals(&_MintAndBurnAdmin.CallOpts, arg0)
}

// Reserve is a free data retrieval call binding the contract method 0xcd3293de.
//
// Solidity: function reserve() constant returns(address)
func (_MintAndBurnAdmin *MintAndBurnAdminCaller) Reserve(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MintAndBurnAdmin.contract.Call(opts, out, "reserve")
	return *ret0, err
}

// Reserve is a free data retrieval call binding the contract method 0xcd3293de.
//
// Solidity: function reserve() constant returns(address)
func (_MintAndBurnAdmin *MintAndBurnAdminSession) Reserve() (common.Address, error) {
	return _MintAndBurnAdmin.Contract.Reserve(&_MintAndBurnAdmin.CallOpts)
}

// Reserve is a free data retrieval call binding the contract method 0xcd3293de.
//
// Solidity: function reserve() constant returns(address)
func (_MintAndBurnAdmin *MintAndBurnAdminCallerSession) Reserve() (common.Address, error) {
	return _MintAndBurnAdmin.Contract.Reserve(&_MintAndBurnAdmin.CallOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0x34a0f49a.
//
// Solidity: function cancel(uint256 index, address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactor) Cancel(opts *bind.TransactOpts, index *big.Int, addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.contract.Transact(opts, "cancel", index, addr, value, isMint)
}

// Cancel is a paid mutator transaction binding the contract method 0x34a0f49a.
//
// Solidity: function cancel(uint256 index, address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminSession) Cancel(index *big.Int, addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.Cancel(&_MintAndBurnAdmin.TransactOpts, index, addr, value, isMint)
}

// Cancel is a paid mutator transaction binding the contract method 0x34a0f49a.
//
// Solidity: function cancel(uint256 index, address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactorSession) Cancel(index *big.Int, addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.Cancel(&_MintAndBurnAdmin.TransactOpts, index, addr, value, isMint)
}

// CancelAll is a paid mutator transaction binding the contract method 0x18cb2b18.
//
// Solidity: function cancelAll() returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactor) CancelAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintAndBurnAdmin.contract.Transact(opts, "cancelAll")
}

// CancelAll is a paid mutator transaction binding the contract method 0x18cb2b18.
//
// Solidity: function cancelAll() returns()
func (_MintAndBurnAdmin *MintAndBurnAdminSession) CancelAll() (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.CancelAll(&_MintAndBurnAdmin.TransactOpts)
}

// CancelAll is a paid mutator transaction binding the contract method 0x18cb2b18.
//
// Solidity: function cancelAll() returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactorSession) CancelAll() (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.CancelAll(&_MintAndBurnAdmin.TransactOpts)
}

// Confirm is a paid mutator transaction binding the contract method 0xdba82a45.
//
// Solidity: function confirm(uint256 index, address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactor) Confirm(opts *bind.TransactOpts, index *big.Int, addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.contract.Transact(opts, "confirm", index, addr, value, isMint)
}

// Confirm is a paid mutator transaction binding the contract method 0xdba82a45.
//
// Solidity: function confirm(uint256 index, address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminSession) Confirm(index *big.Int, addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.Confirm(&_MintAndBurnAdmin.TransactOpts, index, addr, value, isMint)
}

// Confirm is a paid mutator transaction binding the contract method 0xdba82a45.
//
// Solidity: function confirm(uint256 index, address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactorSession) Confirm(index *big.Int, addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.Confirm(&_MintAndBurnAdmin.TransactOpts, index, addr, value, isMint)
}

// Propose is a paid mutator transaction binding the contract method 0x5e77e64e.
//
// Solidity: function propose(address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactor) Propose(opts *bind.TransactOpts, addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.contract.Transact(opts, "propose", addr, value, isMint)
}

// Propose is a paid mutator transaction binding the contract method 0x5e77e64e.
//
// Solidity: function propose(address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminSession) Propose(addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.Propose(&_MintAndBurnAdmin.TransactOpts, addr, value, isMint)
}

// Propose is a paid mutator transaction binding the contract method 0x5e77e64e.
//
// Solidity: function propose(address addr, uint256 value, bool isMint) returns()
func (_MintAndBurnAdmin *MintAndBurnAdminTransactorSession) Propose(addr common.Address, value *big.Int, isMint bool) (*types.Transaction, error) {
	return _MintAndBurnAdmin.Contract.Propose(&_MintAndBurnAdmin.TransactOpts, addr, value, isMint)
}

// MintAndBurnAdminAllProposalsCancelledIterator is returned from FilterAllProposalsCancelled and is used to iterate over the raw logs and unpacked data for AllProposalsCancelled events raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminAllProposalsCancelledIterator struct {
	Event *MintAndBurnAdminAllProposalsCancelled // Event containing the contract specifics and raw log

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
func (it *MintAndBurnAdminAllProposalsCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintAndBurnAdminAllProposalsCancelled)
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
		it.Event = new(MintAndBurnAdminAllProposalsCancelled)
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
func (it *MintAndBurnAdminAllProposalsCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintAndBurnAdminAllProposalsCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintAndBurnAdminAllProposalsCancelled represents a AllProposalsCancelled event raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminAllProposalsCancelled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterAllProposalsCancelled is a free log retrieval operation binding the contract event 0x3732302b0efc3e1e883bb80d83c641031dc1e32223cb406c3e4d5de68208c4e9.
//
// Solidity: event AllProposalsCancelled()
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) FilterAllProposalsCancelled(opts *bind.FilterOpts) (*MintAndBurnAdminAllProposalsCancelledIterator, error) {

	logs, sub, err := _MintAndBurnAdmin.contract.FilterLogs(opts, "AllProposalsCancelled")
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdminAllProposalsCancelledIterator{contract: _MintAndBurnAdmin.contract, event: "AllProposalsCancelled", logs: logs, sub: sub}, nil
}

// WatchAllProposalsCancelled is a free log subscription operation binding the contract event 0x3732302b0efc3e1e883bb80d83c641031dc1e32223cb406c3e4d5de68208c4e9.
//
// Solidity: event AllProposalsCancelled()
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) WatchAllProposalsCancelled(opts *bind.WatchOpts, sink chan<- *MintAndBurnAdminAllProposalsCancelled) (event.Subscription, error) {

	logs, sub, err := _MintAndBurnAdmin.contract.WatchLogs(opts, "AllProposalsCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintAndBurnAdminAllProposalsCancelled)
				if err := _MintAndBurnAdmin.contract.UnpackLog(event, "AllProposalsCancelled", log); err != nil {
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

// MintAndBurnAdminProposalCancelledIterator is returned from FilterProposalCancelled and is used to iterate over the raw logs and unpacked data for ProposalCancelled events raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminProposalCancelledIterator struct {
	Event *MintAndBurnAdminProposalCancelled // Event containing the contract specifics and raw log

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
func (it *MintAndBurnAdminProposalCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintAndBurnAdminProposalCancelled)
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
		it.Event = new(MintAndBurnAdminProposalCancelled)
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
func (it *MintAndBurnAdminProposalCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintAndBurnAdminProposalCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintAndBurnAdminProposalCancelled represents a ProposalCancelled event raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminProposalCancelled struct {
	Index  *big.Int
	Addr   common.Address
	Value  *big.Int
	IsMint bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterProposalCancelled is a free log retrieval operation binding the contract event 0xc1ea9ad7fe3cfb48a741fc229353411aabb3b135d9446697bf6db7c197a9ac0f.
//
// Solidity: event ProposalCancelled(uint256 index, address indexed addr, uint256 value, bool isMint)
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) FilterProposalCancelled(opts *bind.FilterOpts, addr []common.Address) (*MintAndBurnAdminProposalCancelledIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _MintAndBurnAdmin.contract.FilterLogs(opts, "ProposalCancelled", addrRule)
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdminProposalCancelledIterator{contract: _MintAndBurnAdmin.contract, event: "ProposalCancelled", logs: logs, sub: sub}, nil
}

// WatchProposalCancelled is a free log subscription operation binding the contract event 0xc1ea9ad7fe3cfb48a741fc229353411aabb3b135d9446697bf6db7c197a9ac0f.
//
// Solidity: event ProposalCancelled(uint256 index, address indexed addr, uint256 value, bool isMint)
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) WatchProposalCancelled(opts *bind.WatchOpts, sink chan<- *MintAndBurnAdminProposalCancelled, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _MintAndBurnAdmin.contract.WatchLogs(opts, "ProposalCancelled", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintAndBurnAdminProposalCancelled)
				if err := _MintAndBurnAdmin.contract.UnpackLog(event, "ProposalCancelled", log); err != nil {
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

// MintAndBurnAdminProposalConfirmedIterator is returned from FilterProposalConfirmed and is used to iterate over the raw logs and unpacked data for ProposalConfirmed events raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminProposalConfirmedIterator struct {
	Event *MintAndBurnAdminProposalConfirmed // Event containing the contract specifics and raw log

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
func (it *MintAndBurnAdminProposalConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintAndBurnAdminProposalConfirmed)
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
		it.Event = new(MintAndBurnAdminProposalConfirmed)
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
func (it *MintAndBurnAdminProposalConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintAndBurnAdminProposalConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintAndBurnAdminProposalConfirmed represents a ProposalConfirmed event raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminProposalConfirmed struct {
	Index  *big.Int
	Addr   common.Address
	Value  *big.Int
	IsMint bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterProposalConfirmed is a free log retrieval operation binding the contract event 0xc398e86b1dfd2596a48f97df67ac573ef31251ea5b65d73e4096be478df18f57.
//
// Solidity: event ProposalConfirmed(uint256 index, address indexed addr, uint256 value, bool isMint)
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) FilterProposalConfirmed(opts *bind.FilterOpts, addr []common.Address) (*MintAndBurnAdminProposalConfirmedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _MintAndBurnAdmin.contract.FilterLogs(opts, "ProposalConfirmed", addrRule)
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdminProposalConfirmedIterator{contract: _MintAndBurnAdmin.contract, event: "ProposalConfirmed", logs: logs, sub: sub}, nil
}

// WatchProposalConfirmed is a free log subscription operation binding the contract event 0xc398e86b1dfd2596a48f97df67ac573ef31251ea5b65d73e4096be478df18f57.
//
// Solidity: event ProposalConfirmed(uint256 index, address indexed addr, uint256 value, bool isMint)
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) WatchProposalConfirmed(opts *bind.WatchOpts, sink chan<- *MintAndBurnAdminProposalConfirmed, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _MintAndBurnAdmin.contract.WatchLogs(opts, "ProposalConfirmed", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintAndBurnAdminProposalConfirmed)
				if err := _MintAndBurnAdmin.contract.UnpackLog(event, "ProposalConfirmed", log); err != nil {
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

// MintAndBurnAdminProposalCreatedIterator is returned from FilterProposalCreated and is used to iterate over the raw logs and unpacked data for ProposalCreated events raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminProposalCreatedIterator struct {
	Event *MintAndBurnAdminProposalCreated // Event containing the contract specifics and raw log

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
func (it *MintAndBurnAdminProposalCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintAndBurnAdminProposalCreated)
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
		it.Event = new(MintAndBurnAdminProposalCreated)
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
func (it *MintAndBurnAdminProposalCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintAndBurnAdminProposalCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintAndBurnAdminProposalCreated represents a ProposalCreated event raised by the MintAndBurnAdmin contract.
type MintAndBurnAdminProposalCreated struct {
	Index      *big.Int
	Addr       common.Address
	Value      *big.Int
	IsMint     bool
	DelayUntil *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalCreated is a free log retrieval operation binding the contract event 0xd1d2eb762bbbecbc03b8a9dd22368874018771d0c93d855cd08c5a8fa6086b96.
//
// Solidity: event ProposalCreated(uint256 index, address indexed addr, uint256 value, bool isMint, uint256 delayUntil)
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) FilterProposalCreated(opts *bind.FilterOpts, addr []common.Address) (*MintAndBurnAdminProposalCreatedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _MintAndBurnAdmin.contract.FilterLogs(opts, "ProposalCreated", addrRule)
	if err != nil {
		return nil, err
	}
	return &MintAndBurnAdminProposalCreatedIterator{contract: _MintAndBurnAdmin.contract, event: "ProposalCreated", logs: logs, sub: sub}, nil
}

// WatchProposalCreated is a free log subscription operation binding the contract event 0xd1d2eb762bbbecbc03b8a9dd22368874018771d0c93d855cd08c5a8fa6086b96.
//
// Solidity: event ProposalCreated(uint256 index, address indexed addr, uint256 value, bool isMint, uint256 delayUntil)
func (_MintAndBurnAdmin *MintAndBurnAdminFilterer) WatchProposalCreated(opts *bind.WatchOpts, sink chan<- *MintAndBurnAdminProposalCreated, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _MintAndBurnAdmin.contract.WatchLogs(opts, "ProposalCreated", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintAndBurnAdminProposalCreated)
				if err := _MintAndBurnAdmin.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
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
