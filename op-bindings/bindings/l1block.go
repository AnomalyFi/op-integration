// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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
)

// L1BlockMetaData contains all meta data concerning the L1Block contract.
var L1BlockMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"DEPOSITOR_ACCOUNT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"basefee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"batcherHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1FeeOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1FeeScalar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"number\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_number\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"_basefee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"_batcherHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_l1FeeOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_l1FeeScalar\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_justification\",\"type\":\"bytes\"}],\"name\":\"setL1BlockValues\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104c1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638381f58a11610081578063b80777ea1161005b578063b80777ea146101a4578063e591b282146101c4578063e81b2c6d1461020457600080fd5b80638381f58a1461017e5780638b239f73146101925780639e8c49661461019b57600080fd5b80635cf24969116100b25780635cf249691461013357806364ca23ef1461013c5780637f122dcf1461016957600080fd5b806309bd5a60146100ce57806354fd4d50146100ea575b600080fd5b6100d760025481565b6040519081526020015b60405180910390f35b6101266040518060400160405280600581526020017f312e312e3000000000000000000000000000000000000000000000000000000081525081565b6040516100e1919061034e565b6100d760015481565b6003546101509067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100e1565b61017c6101773660046103de565b61020d565b005b6000546101509067ffffffffffffffff1681565b6100d760055481565b6100d760065481565b6000546101509068010000000000000000900467ffffffffffffffff1681565b6101df73deaddeaddeaddeaddeaddeaddeaddeaddead000181565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e1565b6100d760045481565b3373deaddeaddeaddeaddeaddeaddeaddeaddead0001146102b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603b60248201527f4c31426c6f636b3a206f6e6c7920746865206465706f7369746f72206163636f60448201527f756e742063616e20736574204c3120626c6f636b2076616c7565730000000000606482015260840160405180910390fd5b50506000805467ffffffffffffffff98891668010000000000000000027fffffffffffffffffffffffffffffffff00000000000000000000000000000000909116998916999099179890981790975560019490945560029290925560038054919094167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009190911617909255600491909155600555600655565b600060208083528351808285015260005b8181101561037b5785810183015185820160400152820161035f565b8181111561038d576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b803567ffffffffffffffff811681146103d957600080fd5b919050565b6000806000806000806000806000806101208b8d0312156103fe57600080fd5b6104078b6103c1565b995061041560208c016103c1565b985060408b0135975060608b0135965061043160808c016103c1565b955060a08b0135945060c08b0135935060e08b013592506101008b013567ffffffffffffffff8082111561046457600080fd5b818d0191508d601f83011261047857600080fd5b81358181111561048757600080fd5b8e602082850101111561049957600080fd5b6020830194508093505050509295989b9194979a509295985056fea164736f6c634300080f000a",
}

// L1BlockABI is the input ABI used to generate the binding from.
// Deprecated: Use L1BlockMetaData.ABI instead.
var L1BlockABI = L1BlockMetaData.ABI

// L1BlockBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use L1BlockMetaData.Bin instead.
var L1BlockBin = L1BlockMetaData.Bin

// DeployL1Block deploys a new Ethereum contract, binding an instance of L1Block to it.
func DeployL1Block(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *L1Block, error) {
	parsed, err := L1BlockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(L1BlockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &L1Block{L1BlockCaller: L1BlockCaller{contract: contract}, L1BlockTransactor: L1BlockTransactor{contract: contract}, L1BlockFilterer: L1BlockFilterer{contract: contract}}, nil
}

// L1Block is an auto generated Go binding around an Ethereum contract.
type L1Block struct {
	L1BlockCaller     // Read-only binding to the contract
	L1BlockTransactor // Write-only binding to the contract
	L1BlockFilterer   // Log filterer for contract events
}

// L1BlockCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1BlockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1BlockTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1BlockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1BlockFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1BlockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1BlockSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1BlockSession struct {
	Contract     *L1Block          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1BlockCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1BlockCallerSession struct {
	Contract *L1BlockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// L1BlockTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1BlockTransactorSession struct {
	Contract     *L1BlockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// L1BlockRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1BlockRaw struct {
	Contract *L1Block // Generic contract binding to access the raw methods on
}

// L1BlockCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1BlockCallerRaw struct {
	Contract *L1BlockCaller // Generic read-only contract binding to access the raw methods on
}

// L1BlockTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1BlockTransactorRaw struct {
	Contract *L1BlockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1Block creates a new instance of L1Block, bound to a specific deployed contract.
func NewL1Block(address common.Address, backend bind.ContractBackend) (*L1Block, error) {
	contract, err := bindL1Block(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1Block{L1BlockCaller: L1BlockCaller{contract: contract}, L1BlockTransactor: L1BlockTransactor{contract: contract}, L1BlockFilterer: L1BlockFilterer{contract: contract}}, nil
}

// NewL1BlockCaller creates a new read-only instance of L1Block, bound to a specific deployed contract.
func NewL1BlockCaller(address common.Address, caller bind.ContractCaller) (*L1BlockCaller, error) {
	contract, err := bindL1Block(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1BlockCaller{contract: contract}, nil
}

// NewL1BlockTransactor creates a new write-only instance of L1Block, bound to a specific deployed contract.
func NewL1BlockTransactor(address common.Address, transactor bind.ContractTransactor) (*L1BlockTransactor, error) {
	contract, err := bindL1Block(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1BlockTransactor{contract: contract}, nil
}

// NewL1BlockFilterer creates a new log filterer instance of L1Block, bound to a specific deployed contract.
func NewL1BlockFilterer(address common.Address, filterer bind.ContractFilterer) (*L1BlockFilterer, error) {
	contract, err := bindL1Block(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1BlockFilterer{contract: contract}, nil
}

// bindL1Block binds a generic wrapper to an already deployed contract.
func bindL1Block(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(L1BlockABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Block *L1BlockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Block.Contract.L1BlockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Block *L1BlockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Block.Contract.L1BlockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Block *L1BlockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Block.Contract.L1BlockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Block *L1BlockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Block.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Block *L1BlockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Block.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Block *L1BlockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Block.Contract.contract.Transact(opts, method, params...)
}

// DEPOSITORACCOUNT is a free data retrieval call binding the contract method 0xe591b282.
//
// Solidity: function DEPOSITOR_ACCOUNT() view returns(address)
func (_L1Block *L1BlockCaller) DEPOSITORACCOUNT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "DEPOSITOR_ACCOUNT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DEPOSITORACCOUNT is a free data retrieval call binding the contract method 0xe591b282.
//
// Solidity: function DEPOSITOR_ACCOUNT() view returns(address)
func (_L1Block *L1BlockSession) DEPOSITORACCOUNT() (common.Address, error) {
	return _L1Block.Contract.DEPOSITORACCOUNT(&_L1Block.CallOpts)
}

// DEPOSITORACCOUNT is a free data retrieval call binding the contract method 0xe591b282.
//
// Solidity: function DEPOSITOR_ACCOUNT() view returns(address)
func (_L1Block *L1BlockCallerSession) DEPOSITORACCOUNT() (common.Address, error) {
	return _L1Block.Contract.DEPOSITORACCOUNT(&_L1Block.CallOpts)
}

// Basefee is a free data retrieval call binding the contract method 0x5cf24969.
//
// Solidity: function basefee() view returns(uint256)
func (_L1Block *L1BlockCaller) Basefee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "basefee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Basefee is a free data retrieval call binding the contract method 0x5cf24969.
//
// Solidity: function basefee() view returns(uint256)
func (_L1Block *L1BlockSession) Basefee() (*big.Int, error) {
	return _L1Block.Contract.Basefee(&_L1Block.CallOpts)
}

// Basefee is a free data retrieval call binding the contract method 0x5cf24969.
//
// Solidity: function basefee() view returns(uint256)
func (_L1Block *L1BlockCallerSession) Basefee() (*big.Int, error) {
	return _L1Block.Contract.Basefee(&_L1Block.CallOpts)
}

// BatcherHash is a free data retrieval call binding the contract method 0xe81b2c6d.
//
// Solidity: function batcherHash() view returns(bytes32)
func (_L1Block *L1BlockCaller) BatcherHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "batcherHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BatcherHash is a free data retrieval call binding the contract method 0xe81b2c6d.
//
// Solidity: function batcherHash() view returns(bytes32)
func (_L1Block *L1BlockSession) BatcherHash() ([32]byte, error) {
	return _L1Block.Contract.BatcherHash(&_L1Block.CallOpts)
}

// BatcherHash is a free data retrieval call binding the contract method 0xe81b2c6d.
//
// Solidity: function batcherHash() view returns(bytes32)
func (_L1Block *L1BlockCallerSession) BatcherHash() ([32]byte, error) {
	return _L1Block.Contract.BatcherHash(&_L1Block.CallOpts)
}

// Hash is a free data retrieval call binding the contract method 0x09bd5a60.
//
// Solidity: function hash() view returns(bytes32)
func (_L1Block *L1BlockCaller) Hash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "hash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Hash is a free data retrieval call binding the contract method 0x09bd5a60.
//
// Solidity: function hash() view returns(bytes32)
func (_L1Block *L1BlockSession) Hash() ([32]byte, error) {
	return _L1Block.Contract.Hash(&_L1Block.CallOpts)
}

// Hash is a free data retrieval call binding the contract method 0x09bd5a60.
//
// Solidity: function hash() view returns(bytes32)
func (_L1Block *L1BlockCallerSession) Hash() ([32]byte, error) {
	return _L1Block.Contract.Hash(&_L1Block.CallOpts)
}

// L1FeeOverhead is a free data retrieval call binding the contract method 0x8b239f73.
//
// Solidity: function l1FeeOverhead() view returns(uint256)
func (_L1Block *L1BlockCaller) L1FeeOverhead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "l1FeeOverhead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L1FeeOverhead is a free data retrieval call binding the contract method 0x8b239f73.
//
// Solidity: function l1FeeOverhead() view returns(uint256)
func (_L1Block *L1BlockSession) L1FeeOverhead() (*big.Int, error) {
	return _L1Block.Contract.L1FeeOverhead(&_L1Block.CallOpts)
}

// L1FeeOverhead is a free data retrieval call binding the contract method 0x8b239f73.
//
// Solidity: function l1FeeOverhead() view returns(uint256)
func (_L1Block *L1BlockCallerSession) L1FeeOverhead() (*big.Int, error) {
	return _L1Block.Contract.L1FeeOverhead(&_L1Block.CallOpts)
}

// L1FeeScalar is a free data retrieval call binding the contract method 0x9e8c4966.
//
// Solidity: function l1FeeScalar() view returns(uint256)
func (_L1Block *L1BlockCaller) L1FeeScalar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "l1FeeScalar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L1FeeScalar is a free data retrieval call binding the contract method 0x9e8c4966.
//
// Solidity: function l1FeeScalar() view returns(uint256)
func (_L1Block *L1BlockSession) L1FeeScalar() (*big.Int, error) {
	return _L1Block.Contract.L1FeeScalar(&_L1Block.CallOpts)
}

// L1FeeScalar is a free data retrieval call binding the contract method 0x9e8c4966.
//
// Solidity: function l1FeeScalar() view returns(uint256)
func (_L1Block *L1BlockCallerSession) L1FeeScalar() (*big.Int, error) {
	return _L1Block.Contract.L1FeeScalar(&_L1Block.CallOpts)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint64)
func (_L1Block *L1BlockCaller) Number(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "number")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint64)
func (_L1Block *L1BlockSession) Number() (uint64, error) {
	return _L1Block.Contract.Number(&_L1Block.CallOpts)
}

// Number is a free data retrieval call binding the contract method 0x8381f58a.
//
// Solidity: function number() view returns(uint64)
func (_L1Block *L1BlockCallerSession) Number() (uint64, error) {
	return _L1Block.Contract.Number(&_L1Block.CallOpts)
}

// SequenceNumber is a free data retrieval call binding the contract method 0x64ca23ef.
//
// Solidity: function sequenceNumber() view returns(uint64)
func (_L1Block *L1BlockCaller) SequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "sequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// SequenceNumber is a free data retrieval call binding the contract method 0x64ca23ef.
//
// Solidity: function sequenceNumber() view returns(uint64)
func (_L1Block *L1BlockSession) SequenceNumber() (uint64, error) {
	return _L1Block.Contract.SequenceNumber(&_L1Block.CallOpts)
}

// SequenceNumber is a free data retrieval call binding the contract method 0x64ca23ef.
//
// Solidity: function sequenceNumber() view returns(uint64)
func (_L1Block *L1BlockCallerSession) SequenceNumber() (uint64, error) {
	return _L1Block.Contract.SequenceNumber(&_L1Block.CallOpts)
}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() view returns(uint64)
func (_L1Block *L1BlockCaller) Timestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "timestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() view returns(uint64)
func (_L1Block *L1BlockSession) Timestamp() (uint64, error) {
	return _L1Block.Contract.Timestamp(&_L1Block.CallOpts)
}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() view returns(uint64)
func (_L1Block *L1BlockCallerSession) Timestamp() (uint64, error) {
	return _L1Block.Contract.Timestamp(&_L1Block.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_L1Block *L1BlockCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L1Block.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_L1Block *L1BlockSession) Version() (string, error) {
	return _L1Block.Contract.Version(&_L1Block.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_L1Block *L1BlockCallerSession) Version() (string, error) {
	return _L1Block.Contract.Version(&_L1Block.CallOpts)
}

// SetL1BlockValues is a paid mutator transaction binding the contract method 0x7f122dcf.
//
// Solidity: function setL1BlockValues(uint64 _number, uint64 _timestamp, uint256 _basefee, bytes32 _hash, uint64 _sequenceNumber, bytes32 _batcherHash, uint256 _l1FeeOverhead, uint256 _l1FeeScalar, bytes _justification) returns()
func (_L1Block *L1BlockTransactor) SetL1BlockValues(opts *bind.TransactOpts, _number uint64, _timestamp uint64, _basefee *big.Int, _hash [32]byte, _sequenceNumber uint64, _batcherHash [32]byte, _l1FeeOverhead *big.Int, _l1FeeScalar *big.Int, _justification []byte) (*types.Transaction, error) {
	return _L1Block.contract.Transact(opts, "setL1BlockValues", _number, _timestamp, _basefee, _hash, _sequenceNumber, _batcherHash, _l1FeeOverhead, _l1FeeScalar, _justification)
}

// SetL1BlockValues is a paid mutator transaction binding the contract method 0x7f122dcf.
//
// Solidity: function setL1BlockValues(uint64 _number, uint64 _timestamp, uint256 _basefee, bytes32 _hash, uint64 _sequenceNumber, bytes32 _batcherHash, uint256 _l1FeeOverhead, uint256 _l1FeeScalar, bytes _justification) returns()
func (_L1Block *L1BlockSession) SetL1BlockValues(_number uint64, _timestamp uint64, _basefee *big.Int, _hash [32]byte, _sequenceNumber uint64, _batcherHash [32]byte, _l1FeeOverhead *big.Int, _l1FeeScalar *big.Int, _justification []byte) (*types.Transaction, error) {
	return _L1Block.Contract.SetL1BlockValues(&_L1Block.TransactOpts, _number, _timestamp, _basefee, _hash, _sequenceNumber, _batcherHash, _l1FeeOverhead, _l1FeeScalar, _justification)
}

// SetL1BlockValues is a paid mutator transaction binding the contract method 0x7f122dcf.
//
// Solidity: function setL1BlockValues(uint64 _number, uint64 _timestamp, uint256 _basefee, bytes32 _hash, uint64 _sequenceNumber, bytes32 _batcherHash, uint256 _l1FeeOverhead, uint256 _l1FeeScalar, bytes _justification) returns()
func (_L1Block *L1BlockTransactorSession) SetL1BlockValues(_number uint64, _timestamp uint64, _basefee *big.Int, _hash [32]byte, _sequenceNumber uint64, _batcherHash [32]byte, _l1FeeOverhead *big.Int, _l1FeeScalar *big.Int, _justification []byte) (*types.Transaction, error) {
	return _L1Block.Contract.SetL1BlockValues(&_L1Block.TransactOpts, _number, _timestamp, _basefee, _hash, _sequenceNumber, _batcherHash, _l1FeeOverhead, _l1FeeScalar, _justification)
}
