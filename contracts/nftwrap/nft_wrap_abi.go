// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nftwrap

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PolyNFTWrapperABI is the input ABI used to generate the binding from.
const PolyNFTWrapperABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"fromAsset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"PolyWrapperLock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"txHash\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"efee\",\"type\":\"uint256\"}],\"name\":\"PolyWrapperSpeedUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"chainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"extractFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromAsset\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collector\",\"type\":\"address\"}],\"name\":\"setFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_lockProxy\",\"type\":\"address\"}],\"name\":\"setLockProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"txHash\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"speedUp\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// PolyNFTWrapperBin is the compiled bytecode used for deploying new contracts.
var PolyNFTWrapperBin = "0x60806040523480156200001157600080fd5b506040516200284338038062002843833981810160405260408110156200003757600080fd5b810190808051906020019092919080519060200190929190505050600062000064620001bb60201b60201c565b9050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35060008060146101000a81548160ff0219169083151502179055506001808190555060008114156200019b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260068152602001807f216c6567616c000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b620001ac82620001c360201b60201c565b80600281905550505062000406565b600033905090565b620001d36200025a60201b60201c565b62000246576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b6200025781620002c060201b60201c565b50565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16620002a4620001bb60201b60201c565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141562000348576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806200281d6026913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b61240780620004166000396000f3fe6080604052600436106100e85760003560e01c80638da5cb5b1161008a578063a42dce8011610059578063a42dce80146103c1578063c415b95c14610412578063d3ed7c7614610469578063f2fde38b1461050c576100e8565b80638da5cb5b146102b95780638f32d59b146103105780639a8a05921461033f5780639d4dc0211461036a576100e8565b80635c975abb116100c65780635c975abb1461020b5780636f2b6ee61461023a578063715018a61461028b5780638456cb59146102a2576100e8565b80630985b87f146100ed5780631745399d146101a35780633f4ba83a146101f4575b600080fd5b6101a1600480360360e081101561010357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803567ffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019092919050505061055d565b005b3480156101af57600080fd5b506101f2600480360360208110156101c657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610820565b005b34801561020057600080fd5b50610209610a6c565b005b34801561021757600080fd5b50610220610af0565b604051808215151515815260200191505060405180910390f35b34801561024657600080fd5b506102896004803603602081101561025d57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b06565b005b34801561029757600080fd5b506102a0610d43565b005b3480156102ae57600080fd5b506102b7610e7c565b005b3480156102c557600080fd5b506102ce610f00565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561031c57600080fd5b50610325610f29565b604051808215151515815260200191505060405180910390f35b34801561034b57600080fd5b50610354610f87565b6040518082815260200191505060405180910390f35b34801561037657600080fd5b5061037f610f8d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103cd57600080fd5b50610410600480360360208110156103e457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610fb3565b005b34801561041e57600080fd5b50610427611114565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61050a6004803603606081101561047f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156104bc57600080fd5b8201836020820111156104ce57600080fd5b803590602001918460018302840111640100000000831117156104f057600080fd5b90919293919293908035906020019092919050505061113a565b005b34801561051857600080fd5b5061055b6004803603602081101561052f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506112e0565b005b6001806000828254019250508190555060006001549050600060149054906101000a900460ff16156105f7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b6002548767ffffffffffffffff161415801561061e575060008767ffffffffffffffff1614155b610690576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600a8152602001807f21746f436861696e49640000000000000000000000000000000000000000000081525060200191505060405180910390fd5b61069a8484611366565b6106a688888888611447565b3373ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff167f3a15d8cf4b167dd8963989f8038f2333a4889f74033bb53bfb767a5cced072e2898989898989604051808767ffffffffffffffff1667ffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828152602001965050505050505060405180910390a36001548114610816576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f5265656e7472616e637947756172643a207265656e7472616e742063616c6c0081525060200191505060405180910390fd5b5050505050505050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146108e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f21666565436f6c6c6563746f720000000000000000000000000000000000000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610964573373ffffffffffffffffffffffffffffffffffffffff166108fc479081150290604051600060405180830381858888f1935050505015801561095e573d6000803e3d6000fd5b50610a69565b610a68600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015610a0757600080fd5b505afa158015610a1b573d6000803e3d6000fd5b505050506040513d6020811015610a3157600080fd5b81019080805190602001909291905050508373ffffffffffffffffffffffffffffffffffffffff166116359092919063ffffffff16565b5b50565b610a74610f29565b610ae6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b610aee6116ed565b565b60008060149054906101000a900460ff16905090565b610b0e610f29565b610b80576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610bba57600080fd5b80600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600073ffffffffffffffffffffffffffffffffffffffff16600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d798f8816040518163ffffffff1660e01b815260040160206040518083038186803b158015610c7b57600080fd5b505afa158015610c8f573d6000803e3d6000fd5b505050506040513d6020811015610ca557600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff161415610d40576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600e8152602001807f6e6f74206c6f636b2070726f787900000000000000000000000000000000000081525060200191505060405180910390fd5b50565b610d4b610f29565b610dbd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a360008060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b610e84610f29565b610ef6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b610efe6117f5565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16610f6b6118ff565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b60025481565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610fbb610f29565b61102d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156110d0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f656d74707920616464726573730000000000000000000000000000000000000081525060200191505060405180910390fd5b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6001806000828254019250508190555060006001549050600060149054906101000a900460ff16156111d4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b6111de8583611366565b3373ffffffffffffffffffffffffffffffffffffffff16848460405180838380828437808301925050509250505060405180910390208673ffffffffffffffffffffffffffffffffffffffff167ff6579aef3e0d086d986c5d6972659f8a0d8602ef7945b054be1b88e088773ef6856040518082815260200191505060405180910390a460015481146112d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f5265656e7472616e637947756172643a207265656e7472616e742063616c6c0081525060200191505060405180910390fd5b5050505050565b6112e8610f29565b61135a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b61136381611907565b50565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561141557803414611410576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f696e73756666696369656e74206574686572000000000000000000000000000081525060200191505060405180910390fd5b611443565b6114423330838573ffffffffffffffffffffffffffffffffffffffff16611a4b909392919063ffffffff16565b5b5050565b61144f61235d565b604051806040016040528084604051602001808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140191505060405160208183030381529060405281526020018567ffffffffffffffff16815250905060606114c982611b38565b90508573ffffffffffffffffffffffffffffffffffffffff1663b88d4fde33600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1686856040518563ffffffff1660e01b8152600401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156115c65780820151818401526020810190506115ab565b50505050905090810190601f1680156115f35780820380516001836020036101000a031916815260200191505b5095505050505050600060405180830381600087803b15801561161557600080fd5b505af1158015611629573d6000803e3d6000fd5b50505050505050505050565b6116e88363a9059cbb60e01b8484604051602401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611c1c565b505050565b600060149054906101000a900460ff1661176f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f5061757361626c653a206e6f742070617573656400000000000000000000000081525060200191505060405180910390fd5b60008060146101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6117b26118ff565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1565b600060149054906101000a900460ff1615611878576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260108152602001807f5061757361626c653a207061757365640000000000000000000000000000000081525060200191505060405180910390fd5b6001600060146101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586118bc6118ff565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561198d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806123826026913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b611b32846323b872dd60e01b858585604051602401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611c1c565b50505050565b606080611b488360000151611e51565b611b558460200151611f27565b6040516020018083805190602001908083835b60208310611b8b5780518252602082019150602081019050602083039250611b68565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b60208310611bdc5780518252602082019150602081019050602083039250611bb9565b6001836020036101000a03801982511681845116808217855250505050505090500192505050604051602081830303815290604052905080915050919050565b611c2582611f6f565b611c97576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f5361666545524332303a2063616c6c20746f206e6f6e2d636f6e74726163740081525060200191505060405180910390fd5b600060608373ffffffffffffffffffffffffffffffffffffffff16836040518082805190602001908083835b60208310611ce65780518252602082019150602081019050602083039250611cc3565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611d48576040519150601f19603f3d011682016040523d82523d6000602084013e611d4d565b606091505b509150915081611dc5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656481525060200191505060405180910390fd5b600081511115611e4b57808060200190516020811015611de457600080fd5b8101908080519060200190929190505050611e4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a8152602001806123a8602a913960400191505060405180910390fd5b5b50505050565b6060600082519050611e6281611fba565b836040516020018083805190602001908083835b60208310611e995780518252602082019150602081019050602083039250611e76565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b60208310611eea5780518252602082019150602081019050602083039250611ec7565b6001836020036101000a03801982511681845116808217855250505050505090500192505050604051602081830303815290604052915050919050565b606080604051905060088082526000601f5b82821015611f5c5785811a82602086010153600182019150600181039050611f39565b5050602882016040525080915050919050565b60008060007fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47060001b9050833f91506000801b8214158015611fb15750808214155b92505050919050565b606060fd8267ffffffffffffffff161015611fdf57611fd882612293565b905061228e565b61ffff8267ffffffffffffffff16116120cb57611fff60fd60f81b6122b8565b612008836122cd565b6040516020018083805190602001908083835b6020831061203e578051825260208201915060208101905060208303925061201b565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b6020831061208f578051825260208201915060208101905060208303925061206c565b6001836020036101000a03801982511681845116808217855250505050505090500192505050604051602081830303815290604052905061228e565b63ffffffff8267ffffffffffffffff16116121b9576120ed60fe60f81b6122b8565b6120f683612315565b6040516020018083805190602001908083835b6020831061212c5780518252602082019150602081019050602083039250612109565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b6020831061217d578051825260208201915060208101905060208303925061215a565b6001836020036101000a03801982511681845116808217855250505050505090500192505050604051602081830303815290604052905061228e565b6121c660ff60f81b6122b8565b6121cf83611f27565b6040516020018083805190602001908083835b6020831061220557805182526020820191506020810190506020830392506121e2565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b602083106122565780518252602082019150602081019050602083039250612233565b6001836020036101000a0380198251168184511680821785525050505050509050019250505060405160208183030381529060405290505b919050565b6060806040519050600181528260f81b60208201526021810160405280915050919050565b60606122c68260f81c612293565b9050919050565b606080604051905060028082526000601f5b828210156123025785811a826020860101536001820191506001810390506122df565b5050602282016040525080915050919050565b606080604051905060048082526000601f5b8282101561234a5785811a82602086010153600182019150600181039050612327565b5050602482016040525080915050919050565b604051806040016040528060608152602001600067ffffffffffffffff168152509056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f20616464726573735361666545524332303a204552433230206f7065726174696f6e20646964206e6f742073756363656564a2646970667358221220d72df2a94bde2917c6e5e9624bfe5a9b74c31f91cc1deeff93a5fff14a7d14c164736f6c634300060200334f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373"

// DeployPolyNFTWrapper deploys a new Ethereum contract, binding an instance of PolyNFTWrapper to it.
func DeployPolyNFTWrapper(auth *bind.TransactOpts, backend bind.ContractBackend, _owner common.Address, _chainId *big.Int) (common.Address, *types.Transaction, *PolyNFTWrapper, error) {
	parsed, err := abi.JSON(strings.NewReader(PolyNFTWrapperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PolyNFTWrapperBin), backend, _owner, _chainId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PolyNFTWrapper{PolyNFTWrapperCaller: PolyNFTWrapperCaller{contract: contract}, PolyNFTWrapperTransactor: PolyNFTWrapperTransactor{contract: contract}, PolyNFTWrapperFilterer: PolyNFTWrapperFilterer{contract: contract}}, nil
}

// PolyNFTWrapper is an auto generated Go binding around an Ethereum contract.
type PolyNFTWrapper struct {
	PolyNFTWrapperCaller     // Read-only binding to the contract
	PolyNFTWrapperTransactor // Write-only binding to the contract
	PolyNFTWrapperFilterer   // Log filterer for contract events
}

// PolyNFTWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolyNFTWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyNFTWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolyNFTWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyNFTWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolyNFTWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyNFTWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolyNFTWrapperSession struct {
	Contract     *PolyNFTWrapper   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PolyNFTWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolyNFTWrapperCallerSession struct {
	Contract *PolyNFTWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// PolyNFTWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolyNFTWrapperTransactorSession struct {
	Contract     *PolyNFTWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// PolyNFTWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolyNFTWrapperRaw struct {
	Contract *PolyNFTWrapper // Generic contract binding to access the raw methods on
}

// PolyNFTWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolyNFTWrapperCallerRaw struct {
	Contract *PolyNFTWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// PolyNFTWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolyNFTWrapperTransactorRaw struct {
	Contract *PolyNFTWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolyNFTWrapper creates a new instance of PolyNFTWrapper, bound to a specific deployed contract.
func NewPolyNFTWrapper(address common.Address, backend bind.ContractBackend) (*PolyNFTWrapper, error) {
	contract, err := bindPolyNFTWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapper{PolyNFTWrapperCaller: PolyNFTWrapperCaller{contract: contract}, PolyNFTWrapperTransactor: PolyNFTWrapperTransactor{contract: contract}, PolyNFTWrapperFilterer: PolyNFTWrapperFilterer{contract: contract}}, nil
}

// NewPolyNFTWrapperCaller creates a new read-only instance of PolyNFTWrapper, bound to a specific deployed contract.
func NewPolyNFTWrapperCaller(address common.Address, caller bind.ContractCaller) (*PolyNFTWrapperCaller, error) {
	contract, err := bindPolyNFTWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperCaller{contract: contract}, nil
}

// NewPolyNFTWrapperTransactor creates a new write-only instance of PolyNFTWrapper, bound to a specific deployed contract.
func NewPolyNFTWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*PolyNFTWrapperTransactor, error) {
	contract, err := bindPolyNFTWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperTransactor{contract: contract}, nil
}

// NewPolyNFTWrapperFilterer creates a new log filterer instance of PolyNFTWrapper, bound to a specific deployed contract.
func NewPolyNFTWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*PolyNFTWrapperFilterer, error) {
	contract, err := bindPolyNFTWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperFilterer{contract: contract}, nil
}

// bindPolyNFTWrapper binds a generic wrapper to an already deployed contract.
func bindPolyNFTWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PolyNFTWrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolyNFTWrapper *PolyNFTWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolyNFTWrapper.Contract.PolyNFTWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolyNFTWrapper *PolyNFTWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.PolyNFTWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolyNFTWrapper *PolyNFTWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.PolyNFTWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolyNFTWrapper *PolyNFTWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolyNFTWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolyNFTWrapper *PolyNFTWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolyNFTWrapper *PolyNFTWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.contract.Transact(opts, method, params...)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (_PolyNFTWrapper *PolyNFTWrapperCaller) ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyNFTWrapper.contract.Call(opts, &out, "chainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (_PolyNFTWrapper *PolyNFTWrapperSession) ChainId() (*big.Int, error) {
	return _PolyNFTWrapper.Contract.ChainId(&_PolyNFTWrapper.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (_PolyNFTWrapper *PolyNFTWrapperCallerSession) ChainId() (*big.Int, error) {
	return _PolyNFTWrapper.Contract.ChainId(&_PolyNFTWrapper.CallOpts)
}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperCaller) FeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolyNFTWrapper.contract.Call(opts, &out, "feeCollector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperSession) FeeCollector() (common.Address, error) {
	return _PolyNFTWrapper.Contract.FeeCollector(&_PolyNFTWrapper.CallOpts)
}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperCallerSession) FeeCollector() (common.Address, error) {
	return _PolyNFTWrapper.Contract.FeeCollector(&_PolyNFTWrapper.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PolyNFTWrapper *PolyNFTWrapperCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PolyNFTWrapper.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PolyNFTWrapper *PolyNFTWrapperSession) IsOwner() (bool, error) {
	return _PolyNFTWrapper.Contract.IsOwner(&_PolyNFTWrapper.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PolyNFTWrapper *PolyNFTWrapperCallerSession) IsOwner() (bool, error) {
	return _PolyNFTWrapper.Contract.IsOwner(&_PolyNFTWrapper.CallOpts)
}

// LockProxy is a free data retrieval call binding the contract method 0x9d4dc021.
//
// Solidity: function lockProxy() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperCaller) LockProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolyNFTWrapper.contract.Call(opts, &out, "lockProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LockProxy is a free data retrieval call binding the contract method 0x9d4dc021.
//
// Solidity: function lockProxy() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperSession) LockProxy() (common.Address, error) {
	return _PolyNFTWrapper.Contract.LockProxy(&_PolyNFTWrapper.CallOpts)
}

// LockProxy is a free data retrieval call binding the contract method 0x9d4dc021.
//
// Solidity: function lockProxy() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperCallerSession) LockProxy() (common.Address, error) {
	return _PolyNFTWrapper.Contract.LockProxy(&_PolyNFTWrapper.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolyNFTWrapper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperSession) Owner() (common.Address, error) {
	return _PolyNFTWrapper.Contract.Owner(&_PolyNFTWrapper.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PolyNFTWrapper *PolyNFTWrapperCallerSession) Owner() (common.Address, error) {
	return _PolyNFTWrapper.Contract.Owner(&_PolyNFTWrapper.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PolyNFTWrapper *PolyNFTWrapperCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PolyNFTWrapper.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PolyNFTWrapper *PolyNFTWrapperSession) Paused() (bool, error) {
	return _PolyNFTWrapper.Contract.Paused(&_PolyNFTWrapper.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PolyNFTWrapper *PolyNFTWrapperCallerSession) Paused() (bool, error) {
	return _PolyNFTWrapper.Contract.Paused(&_PolyNFTWrapper.CallOpts)
}

// ExtractFee is a paid mutator transaction binding the contract method 0x1745399d.
//
// Solidity: function extractFee(address token) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) ExtractFee(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "extractFee", token)
}

// ExtractFee is a paid mutator transaction binding the contract method 0x1745399d.
//
// Solidity: function extractFee(address token) returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) ExtractFee(token common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.ExtractFee(&_PolyNFTWrapper.TransactOpts, token)
}

// ExtractFee is a paid mutator transaction binding the contract method 0x1745399d.
//
// Solidity: function extractFee(address token) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) ExtractFee(token common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.ExtractFee(&_PolyNFTWrapper.TransactOpts, token)
}

// Lock is a paid mutator transaction binding the contract method 0x0985b87f.
//
// Solidity: function lock(address fromAsset, uint64 toChainId, address toAddress, uint256 tokenId, address feeToken, uint256 fee, uint256 id) payable returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) Lock(opts *bind.TransactOpts, fromAsset common.Address, toChainId uint64, toAddress common.Address, tokenId *big.Int, feeToken common.Address, fee *big.Int, id *big.Int) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "lock", fromAsset, toChainId, toAddress, tokenId, feeToken, fee, id)
}

// Lock is a paid mutator transaction binding the contract method 0x0985b87f.
//
// Solidity: function lock(address fromAsset, uint64 toChainId, address toAddress, uint256 tokenId, address feeToken, uint256 fee, uint256 id) payable returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) Lock(fromAsset common.Address, toChainId uint64, toAddress common.Address, tokenId *big.Int, feeToken common.Address, fee *big.Int, id *big.Int) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.Lock(&_PolyNFTWrapper.TransactOpts, fromAsset, toChainId, toAddress, tokenId, feeToken, fee, id)
}

// Lock is a paid mutator transaction binding the contract method 0x0985b87f.
//
// Solidity: function lock(address fromAsset, uint64 toChainId, address toAddress, uint256 tokenId, address feeToken, uint256 fee, uint256 id) payable returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) Lock(fromAsset common.Address, toChainId uint64, toAddress common.Address, tokenId *big.Int, feeToken common.Address, fee *big.Int, id *big.Int) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.Lock(&_PolyNFTWrapper.TransactOpts, fromAsset, toChainId, toAddress, tokenId, feeToken, fee, id)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) Pause() (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.Pause(&_PolyNFTWrapper.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) Pause() (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.Pause(&_PolyNFTWrapper.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) RenounceOwnership() (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.RenounceOwnership(&_PolyNFTWrapper.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.RenounceOwnership(&_PolyNFTWrapper.TransactOpts)
}

// SetFeeCollector is a paid mutator transaction binding the contract method 0xa42dce80.
//
// Solidity: function setFeeCollector(address collector) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) SetFeeCollector(opts *bind.TransactOpts, collector common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "setFeeCollector", collector)
}

// SetFeeCollector is a paid mutator transaction binding the contract method 0xa42dce80.
//
// Solidity: function setFeeCollector(address collector) returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) SetFeeCollector(collector common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.SetFeeCollector(&_PolyNFTWrapper.TransactOpts, collector)
}

// SetFeeCollector is a paid mutator transaction binding the contract method 0xa42dce80.
//
// Solidity: function setFeeCollector(address collector) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) SetFeeCollector(collector common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.SetFeeCollector(&_PolyNFTWrapper.TransactOpts, collector)
}

// SetLockProxy is a paid mutator transaction binding the contract method 0x6f2b6ee6.
//
// Solidity: function setLockProxy(address _lockProxy) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) SetLockProxy(opts *bind.TransactOpts, _lockProxy common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "setLockProxy", _lockProxy)
}

// SetLockProxy is a paid mutator transaction binding the contract method 0x6f2b6ee6.
//
// Solidity: function setLockProxy(address _lockProxy) returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) SetLockProxy(_lockProxy common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.SetLockProxy(&_PolyNFTWrapper.TransactOpts, _lockProxy)
}

// SetLockProxy is a paid mutator transaction binding the contract method 0x6f2b6ee6.
//
// Solidity: function setLockProxy(address _lockProxy) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) SetLockProxy(_lockProxy common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.SetLockProxy(&_PolyNFTWrapper.TransactOpts, _lockProxy)
}

// SpeedUp is a paid mutator transaction binding the contract method 0xd3ed7c76.
//
// Solidity: function speedUp(address feeToken, bytes txHash, uint256 fee) payable returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) SpeedUp(opts *bind.TransactOpts, feeToken common.Address, txHash []byte, fee *big.Int) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "speedUp", feeToken, txHash, fee)
}

// SpeedUp is a paid mutator transaction binding the contract method 0xd3ed7c76.
//
// Solidity: function speedUp(address feeToken, bytes txHash, uint256 fee) payable returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) SpeedUp(feeToken common.Address, txHash []byte, fee *big.Int) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.SpeedUp(&_PolyNFTWrapper.TransactOpts, feeToken, txHash, fee)
}

// SpeedUp is a paid mutator transaction binding the contract method 0xd3ed7c76.
//
// Solidity: function speedUp(address feeToken, bytes txHash, uint256 fee) payable returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) SpeedUp(feeToken common.Address, txHash []byte, fee *big.Int) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.SpeedUp(&_PolyNFTWrapper.TransactOpts, feeToken, txHash, fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.TransferOwnership(&_PolyNFTWrapper.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.TransferOwnership(&_PolyNFTWrapper.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTWrapper.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PolyNFTWrapper *PolyNFTWrapperSession) Unpause() (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.Unpause(&_PolyNFTWrapper.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PolyNFTWrapper *PolyNFTWrapperTransactorSession) Unpause() (*types.Transaction, error) {
	return _PolyNFTWrapper.Contract.Unpause(&_PolyNFTWrapper.TransactOpts)
}

// PolyNFTWrapperOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PolyNFTWrapper contract.
type PolyNFTWrapperOwnershipTransferredIterator struct {
	Event *PolyNFTWrapperOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PolyNFTWrapperOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyNFTWrapperOwnershipTransferred)
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
		it.Event = new(PolyNFTWrapperOwnershipTransferred)
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
func (it *PolyNFTWrapperOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyNFTWrapperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyNFTWrapperOwnershipTransferred represents a OwnershipTransferred event raised by the PolyNFTWrapper contract.
type PolyNFTWrapperOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PolyNFTWrapperOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PolyNFTWrapper.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperOwnershipTransferredIterator{contract: _PolyNFTWrapper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PolyNFTWrapperOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PolyNFTWrapper.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyNFTWrapperOwnershipTransferred)
				if err := _PolyNFTWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) ParseOwnershipTransferred(log types.Log) (*PolyNFTWrapperOwnershipTransferred, error) {
	event := new(PolyNFTWrapperOwnershipTransferred)
	if err := _PolyNFTWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// PolyNFTWrapperPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the PolyNFTWrapper contract.
type PolyNFTWrapperPausedIterator struct {
	Event *PolyNFTWrapperPaused // Event containing the contract specifics and raw log

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
func (it *PolyNFTWrapperPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyNFTWrapperPaused)
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
		it.Event = new(PolyNFTWrapperPaused)
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
func (it *PolyNFTWrapperPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyNFTWrapperPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyNFTWrapperPaused represents a Paused event raised by the PolyNFTWrapper contract.
type PolyNFTWrapperPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) FilterPaused(opts *bind.FilterOpts) (*PolyNFTWrapperPausedIterator, error) {

	logs, sub, err := _PolyNFTWrapper.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperPausedIterator{contract: _PolyNFTWrapper.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *PolyNFTWrapperPaused) (event.Subscription, error) {

	logs, sub, err := _PolyNFTWrapper.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyNFTWrapperPaused)
				if err := _PolyNFTWrapper.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) ParsePaused(log types.Log) (*PolyNFTWrapperPaused, error) {
	event := new(PolyNFTWrapperPaused)
	if err := _PolyNFTWrapper.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	return event, nil
}

// PolyNFTWrapperPolyWrapperLockIterator is returned from FilterPolyWrapperLock and is used to iterate over the raw logs and unpacked data for PolyWrapperLock events raised by the PolyNFTWrapper contract.
type PolyNFTWrapperPolyWrapperLockIterator struct {
	Event *PolyNFTWrapperPolyWrapperLock // Event containing the contract specifics and raw log

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
func (it *PolyNFTWrapperPolyWrapperLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyNFTWrapperPolyWrapperLock)
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
		it.Event = new(PolyNFTWrapperPolyWrapperLock)
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
func (it *PolyNFTWrapperPolyWrapperLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyNFTWrapperPolyWrapperLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyNFTWrapperPolyWrapperLock represents a PolyWrapperLock event raised by the PolyNFTWrapper contract.
type PolyNFTWrapperPolyWrapperLock struct {
	FromAsset common.Address
	Sender    common.Address
	ToChainId uint64
	ToAddress common.Address
	TokenId   *big.Int
	FeeToken  common.Address
	Fee       *big.Int
	Id        *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPolyWrapperLock is a free log retrieval operation binding the contract event 0x3a15d8cf4b167dd8963989f8038f2333a4889f74033bb53bfb767a5cced072e2.
//
// Solidity: event PolyWrapperLock(address indexed fromAsset, address indexed sender, uint64 toChainId, address toAddress, uint256 tokenId, address feeToken, uint256 fee, uint256 id)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) FilterPolyWrapperLock(opts *bind.FilterOpts, fromAsset []common.Address, sender []common.Address) (*PolyNFTWrapperPolyWrapperLockIterator, error) {

	var fromAssetRule []interface{}
	for _, fromAssetItem := range fromAsset {
		fromAssetRule = append(fromAssetRule, fromAssetItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyNFTWrapper.contract.FilterLogs(opts, "PolyWrapperLock", fromAssetRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperPolyWrapperLockIterator{contract: _PolyNFTWrapper.contract, event: "PolyWrapperLock", logs: logs, sub: sub}, nil
}

// WatchPolyWrapperLock is a free log subscription operation binding the contract event 0x3a15d8cf4b167dd8963989f8038f2333a4889f74033bb53bfb767a5cced072e2.
//
// Solidity: event PolyWrapperLock(address indexed fromAsset, address indexed sender, uint64 toChainId, address toAddress, uint256 tokenId, address feeToken, uint256 fee, uint256 id)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) WatchPolyWrapperLock(opts *bind.WatchOpts, sink chan<- *PolyNFTWrapperPolyWrapperLock, fromAsset []common.Address, sender []common.Address) (event.Subscription, error) {

	var fromAssetRule []interface{}
	for _, fromAssetItem := range fromAsset {
		fromAssetRule = append(fromAssetRule, fromAssetItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyNFTWrapper.contract.WatchLogs(opts, "PolyWrapperLock", fromAssetRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyNFTWrapperPolyWrapperLock)
				if err := _PolyNFTWrapper.contract.UnpackLog(event, "PolyWrapperLock", log); err != nil {
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

// ParsePolyWrapperLock is a log parse operation binding the contract event 0x3a15d8cf4b167dd8963989f8038f2333a4889f74033bb53bfb767a5cced072e2.
//
// Solidity: event PolyWrapperLock(address indexed fromAsset, address indexed sender, uint64 toChainId, address toAddress, uint256 tokenId, address feeToken, uint256 fee, uint256 id)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) ParsePolyWrapperLock(log types.Log) (*PolyNFTWrapperPolyWrapperLock, error) {
	event := new(PolyNFTWrapperPolyWrapperLock)
	if err := _PolyNFTWrapper.contract.UnpackLog(event, "PolyWrapperLock", log); err != nil {
		return nil, err
	}
	return event, nil
}

// PolyNFTWrapperPolyWrapperSpeedUpIterator is returned from FilterPolyWrapperSpeedUp and is used to iterate over the raw logs and unpacked data for PolyWrapperSpeedUp events raised by the PolyNFTWrapper contract.
type PolyNFTWrapperPolyWrapperSpeedUpIterator struct {
	Event *PolyNFTWrapperPolyWrapperSpeedUp // Event containing the contract specifics and raw log

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
func (it *PolyNFTWrapperPolyWrapperSpeedUpIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyNFTWrapperPolyWrapperSpeedUp)
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
		it.Event = new(PolyNFTWrapperPolyWrapperSpeedUp)
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
func (it *PolyNFTWrapperPolyWrapperSpeedUpIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyNFTWrapperPolyWrapperSpeedUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyNFTWrapperPolyWrapperSpeedUp represents a PolyWrapperSpeedUp event raised by the PolyNFTWrapper contract.
type PolyNFTWrapperPolyWrapperSpeedUp struct {
	FeeToken common.Address
	TxHash   common.Hash
	Sender   common.Address
	Efee     *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPolyWrapperSpeedUp is a free log retrieval operation binding the contract event 0xf6579aef3e0d086d986c5d6972659f8a0d8602ef7945b054be1b88e088773ef6.
//
// Solidity: event PolyWrapperSpeedUp(address indexed feeToken, bytes indexed txHash, address indexed sender, uint256 efee)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) FilterPolyWrapperSpeedUp(opts *bind.FilterOpts, feeToken []common.Address, txHash [][]byte, sender []common.Address) (*PolyNFTWrapperPolyWrapperSpeedUpIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}
	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyNFTWrapper.contract.FilterLogs(opts, "PolyWrapperSpeedUp", feeTokenRule, txHashRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperPolyWrapperSpeedUpIterator{contract: _PolyNFTWrapper.contract, event: "PolyWrapperSpeedUp", logs: logs, sub: sub}, nil
}

// WatchPolyWrapperSpeedUp is a free log subscription operation binding the contract event 0xf6579aef3e0d086d986c5d6972659f8a0d8602ef7945b054be1b88e088773ef6.
//
// Solidity: event PolyWrapperSpeedUp(address indexed feeToken, bytes indexed txHash, address indexed sender, uint256 efee)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) WatchPolyWrapperSpeedUp(opts *bind.WatchOpts, sink chan<- *PolyNFTWrapperPolyWrapperSpeedUp, feeToken []common.Address, txHash [][]byte, sender []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}
	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyNFTWrapper.contract.WatchLogs(opts, "PolyWrapperSpeedUp", feeTokenRule, txHashRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyNFTWrapperPolyWrapperSpeedUp)
				if err := _PolyNFTWrapper.contract.UnpackLog(event, "PolyWrapperSpeedUp", log); err != nil {
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

// ParsePolyWrapperSpeedUp is a log parse operation binding the contract event 0xf6579aef3e0d086d986c5d6972659f8a0d8602ef7945b054be1b88e088773ef6.
//
// Solidity: event PolyWrapperSpeedUp(address indexed feeToken, bytes indexed txHash, address indexed sender, uint256 efee)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) ParsePolyWrapperSpeedUp(log types.Log) (*PolyNFTWrapperPolyWrapperSpeedUp, error) {
	event := new(PolyNFTWrapperPolyWrapperSpeedUp)
	if err := _PolyNFTWrapper.contract.UnpackLog(event, "PolyWrapperSpeedUp", log); err != nil {
		return nil, err
	}
	return event, nil
}

// PolyNFTWrapperUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the PolyNFTWrapper contract.
type PolyNFTWrapperUnpausedIterator struct {
	Event *PolyNFTWrapperUnpaused // Event containing the contract specifics and raw log

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
func (it *PolyNFTWrapperUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyNFTWrapperUnpaused)
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
		it.Event = new(PolyNFTWrapperUnpaused)
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
func (it *PolyNFTWrapperUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyNFTWrapperUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyNFTWrapperUnpaused represents a Unpaused event raised by the PolyNFTWrapper contract.
type PolyNFTWrapperUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) FilterUnpaused(opts *bind.FilterOpts) (*PolyNFTWrapperUnpausedIterator, error) {

	logs, sub, err := _PolyNFTWrapper.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &PolyNFTWrapperUnpausedIterator{contract: _PolyNFTWrapper.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *PolyNFTWrapperUnpaused) (event.Subscription, error) {

	logs, sub, err := _PolyNFTWrapper.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyNFTWrapperUnpaused)
				if err := _PolyNFTWrapper.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PolyNFTWrapper *PolyNFTWrapperFilterer) ParseUnpaused(log types.Log) (*PolyNFTWrapperUnpaused, error) {
	event := new(PolyNFTWrapperUnpaused)
	if err := _PolyNFTWrapper.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	return event, nil
}
