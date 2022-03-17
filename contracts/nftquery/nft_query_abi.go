// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nftquery

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

// PolyNFTQueryABI is the input ABI used to generate the binding from.
const PolyNFTQueryABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_limit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getAndCheckTokenUrl\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"ignore\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"getFilterTokensByIndex\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"getOwnerTokensByIndex\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"args\",\"type\":\"bytes\"}],\"name\":\"getTokensByIds\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_limit\",\"type\":\"uint256\"}],\"name\":\"setQueryLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// PolyNFTQueryBin is the compiled bytecode used for deploying new contracts.
var PolyNFTQueryBin = "0x60806040523480156200001157600080fd5b506040516200222b3803806200222b833981810160405260408110156200003757600080fd5b8101908080519060200190929190805190602001909291905050506000620000646200019960201b60201c565b9050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3506000811162000179576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260068152602001807f216c6567616c000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b806001819055506200019182620001a160201b60201c565b5050620003e3565b600033905090565b620001b16200023860201b60201c565b62000224576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b62000235816200029e60201b60201c565b50565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16620002826200019960201b60201c565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141562000326576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526026815260200180620022056026913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b611e1280620003f36000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80638f32d59b116100665780638f32d59b146101f4578063bbc8e44614610214578063c5cc029a1461030e578063c792f10f14610408578063f2fde38b1461056557610093565b806312d459ab14610098578063574ded3b14610188578063715018a6146101b65780638da5cb5b146101c0575b600080fd5b610104600480360360608110156100ae57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506105a9565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561014c578082015181840152602081019050610131565b50505050905090810190601f1680156101795780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6101b46004803603602081101561019e57600080fd5b8101908080359060200190929190505050610828565b005b6101be6108ac565b005b6101c86109e4565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101fc610a0d565b60405180821515815260200191505060405180910390f35b61028a6004803603608081101561022a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190505050610a6b565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102d25780820151818401526020810190506102b7565b50505050905090810190601f1680156102ff5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6103846004803603608081101561032457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190505050610daf565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156103cc5780820151818401526020810190506103b1565b50505050905090810190601f1680156103f95780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6104e16004803603604081101561041e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561045b57600080fd5b82018360208201111561046d57600080fd5b8035906020019184600183028401116401000000008311171561048f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050611194565b60405180831515815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561052957808201518184015260208101905061050e565b50505050905090810190601f1680156105565780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6105a76004803603602081101561057b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611376565b005b600060608060405180602001604052806000815250905060008673ffffffffffffffffffffffffffffffffffffffff16636352211e866040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561061357600080fd5b505afa158015610627573d6000803e3d6000fd5b505050506040513d602081101561063d57600080fd5b810190808051906020019092919050505090508073ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff161415806106b75750600073ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff16145b156106ca57600082935093505050610820565b8673ffffffffffffffffffffffffffffffffffffffff1663c87b56dd866040518263ffffffff1660e01b81526004018082815260200191505060006040518083038186803b15801561071b57600080fd5b505afa15801561072f573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250602081101561075957600080fd5b810190808051604051939291908464010000000082111561077957600080fd5b8382019150602082018581111561078f57600080fd5b82518660018202830111640100000000821117156107ac57600080fd5b8083526020830192505050908051906020019080838360005b838110156107e05780820151818401526020810190506107c5565b50505050905090810190601f16801561080d5780820380516001836020036101000a031916815260200191505b5060405250505091506001829350935050505b935093915050565b610830610a0d565b6108a2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b8060018190555050565b6108b4610a0d565b610926576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff1660008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a360008060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16610a4f6113fc565b73ffffffffffffffffffffffffffffffffffffffff1614905090565b60006060806000841480610a80575060015484115b15610a92576000819250925050610da6565b60008773ffffffffffffffffffffffffffffffffffffffff166370a08231886040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015610afb57600080fd5b505afa158015610b0f573d6000803e3d6000fd5b505050506040513d6020811015610b2557600080fd5b810190808051906020019092919050505090506000811480610b475750808610155b15610b5a57600082935093505050610da6565b6000610b67878784611404565b9050600089905060008a905060008990505b838111610d985760008273ffffffffffffffffffffffffffffffffffffffff16632f745c598d846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060206040518083038186803b158015610bf157600080fd5b505afa158015610c05573d6000803e3d6000fd5b505050506040513d6020811015610c1b57600080fd5b8101908080519060200190929190505050905060608473ffffffffffffffffffffffffffffffffffffffff1663c87b56dd836040518263ffffffff1660e01b81526004018082815260200191505060006040518083038186803b158015610c8157600080fd5b505afa158015610c95573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506020811015610cbf57600080fd5b8101908080516040519392919084640100000000821115610cdf57600080fd5b83820191506020820185811115610cf557600080fd5b8251866001820283011164010000000082111715610d1257600080fd5b8083526020830192505050908051906020019080838360005b83811015610d46578082015181840152602081019050610d2b565b50505050905090810190601f168015610d735780820380516001836020036101000a031916815260200191505b506040525050509050610d87888383611428565b975050508080600101915050610b79565b506001859650965050505050505b94509492505050565b60006060806000841480610dc4575060015484115b15610dd657600081925092505061118b565b60008790506000889050600089905060008273ffffffffffffffffffffffffffffffffffffffff166318160ddd6040518163ffffffff1660e01b815260040160206040518083038186803b158015610e2d57600080fd5b505afa158015610e41573d6000803e3d6000fd5b505050506040513d6020811015610e5757600080fd5b810190808051906020019092919050505090506000811480610e795750808910155b15610e8f5760008596509650505050505061118b565b6000610e9c8a8a84611404565b90505b808a11158015610eae57508181105b1561117d5760008473ffffffffffffffffffffffffffffffffffffffff16634f6ccce78c6040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610f0657600080fd5b505afa158015610f1a573d6000803e3d6000fd5b505050506040513d6020811015610f3057600080fd5b8101908080519060200190929190505050905060018b019a5060008473ffffffffffffffffffffffffffffffffffffffff16636352211e836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610f9c57600080fd5b505afa158015610fb0573d6000803e3d6000fd5b505050506040513d6020811015610fc657600080fd5b810190808051906020019092919050505090508c73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561101a576001830192505050610e9f565b60608773ffffffffffffffffffffffffffffffffffffffff1663c87b56dd846040518263ffffffff1660e01b81526004018082815260200191505060006040518083038186803b15801561106d57600080fd5b505afa158015611081573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525060208110156110ab57600080fd5b81019080805160405193929190846401000000008211156110cb57600080fd5b838201915060208201858111156110e157600080fd5b82518660018202830111640100000000821117156110fe57600080fd5b8083526020830192505050908051906020019080838360005b83811015611132578082015181840152602081019050611117565b50505050905090810190601f16801561115f5780820380516001836020036101000a031916815260200191505b506040525050509050611173898483611428565b9850505050610e9f565b600186975097505050505050505b94509492505050565b60006060600080600060606111a98785611557565b809550819350505060008214806111c1575060015482115b156111d657600081955095505050505061136f565b600088905060005b83811015611361576111f08987611557565b809750819650505060608273ffffffffffffffffffffffffffffffffffffffff1663c87b56dd876040518263ffffffff1660e01b81526004018082815260200191505060006040518083038186803b15801561124b57600080fd5b505afa15801561125f573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250602081101561128957600080fd5b81019080805160405193929190846401000000008211156112a957600080fd5b838201915060208201858111156112bf57600080fd5b82518660018202830111640100000000821117156112dc57600080fd5b8083526020830192505050908051906020019080838360005b838110156113105780820151818401526020810190506112f5565b50505050905090810190601f16801561133d5780820380516001836020036101000a031916815260200191505b506040525050509050611351848783611428565b93505080806001019150506111de565b506001829650965050505050505b9250929050565b61137e610a0d565b6113f0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657281525060200191505060405180910390fd5b6113f9816116b0565b50565b600033905090565b600080600184860103905082811061141d576001830390505b809150509392505050565b606083611434846117f3565b61143d846118d2565b6040516020018084805190602001908083835b602083106114735780518252602082019150602081019050602083039250611450565b6001836020036101000a03801982511681845116808217855250505050505090500183805190602001908083835b602083106114c457805182526020820191506020810190506020830392506114a1565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b6020831061151557805182526020820191506020810190506020830392506114f2565b6001836020036101000a038019825116818451168082178552505050505050905001935050505060405160208183030381529060405293508390509392505050565b6000808351602084011115801561157057506020830183105b6115c5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526023815260200180611dba6023913960400191505060405180910390fd5b600060405160206000600182038760208a0101515b838310156115fa5780821a838601536001830192506001820391506115da565b5050508082016040528151925050507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81111561169f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f56616c75652065786365656473207468652072616e676500000000000000000081525060200191505060405180910390fd5b806020850192509250509250929050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415611736576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526026815260200180611d946026913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82111561188b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f56616c756520657863656564732075696e743235362072616e6765000000000081525060200191505060405180910390fd5b6060604051905060208082526000601f5b828210156118bf5785811a8260208601015360018201915060018103905061189c565b5050604082016040525080915050919050565b60606000825190506118e3816119a8565b836040516020018083805190602001908083835b6020831061191a57805182526020820191506020810190506020830392506118f7565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b6020831061196b5780518252602082019150602081019050602083039250611948565b6001836020036101000a03801982511681845116808217855250505050505090500192505050604051602081830303815290604052915050919050565b606060fd8267ffffffffffffffff1610156119cd576119c682611c81565b9050611c7c565b61ffff8267ffffffffffffffff1611611ab9576119ed60fd60f81b611ca6565b6119f683611cbb565b6040516020018083805190602001908083835b60208310611a2c5780518252602082019150602081019050602083039250611a09565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b60208310611a7d5780518252602082019150602081019050602083039250611a5a565b6001836020036101000a038019825116818451168082178552505050505050905001925050506040516020818303038152906040529050611c7c565b63ffffffff8267ffffffffffffffff1611611ba757611adb60fe60f81b611ca6565b611ae483611d03565b6040516020018083805190602001908083835b60208310611b1a5780518252602082019150602081019050602083039250611af7565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b60208310611b6b5780518252602082019150602081019050602083039250611b48565b6001836020036101000a038019825116818451168082178552505050505050905001925050506040516020818303038152906040529050611c7c565b611bb460ff60f81b611ca6565b611bbd83611d4b565b6040516020018083805190602001908083835b60208310611bf35780518252602082019150602081019050602083039250611bd0565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b60208310611c445780518252602082019150602081019050602083039250611c21565b6001836020036101000a0380198251168184511680821785525050505050509050019250505060405160208183030381529060405290505b919050565b6060806040519050600181528260f81b60208201526021810160405280915050919050565b6060611cb48260f81c611c81565b9050919050565b606080604051905060028082526000601f5b82821015611cf05785811a82602086010153600182019150600181039050611ccd565b5050602282016040525080915050919050565b606080604051905060048082526000601f5b82821015611d385785811a82602086010153600182019150600181039050611d15565b5050602482016040525080915050919050565b606080604051905060088082526000601f5b82821015611d805785811a82602086010153600182019150600181039050611d5d565b505060288201604052508091505091905056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f20616464726573734e65787455696e743235362c206f66667365742065786365656473206d6178696d756da264697066735822122021d4197742b7a66998a148f77fffeb7bfd0184972d84e2c6956dccc7b037da3e64736f6c634300060c00334f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373"

// DeployPolyNFTQuery deploys a new Ethereum contract, binding an instance of PolyNFTQuery to it.
func DeployPolyNFTQuery(auth *bind.TransactOpts, backend bind.ContractBackend, _owner common.Address, _limit *big.Int) (common.Address, *types.Transaction, *PolyNFTQuery, error) {
	parsed, err := abi.JSON(strings.NewReader(PolyNFTQueryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PolyNFTQueryBin), backend, _owner, _limit)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PolyNFTQuery{PolyNFTQueryCaller: PolyNFTQueryCaller{contract: contract}, PolyNFTQueryTransactor: PolyNFTQueryTransactor{contract: contract}, PolyNFTQueryFilterer: PolyNFTQueryFilterer{contract: contract}}, nil
}

// PolyNFTQuery is an auto generated Go binding around an Ethereum contract.
type PolyNFTQuery struct {
	PolyNFTQueryCaller     // Read-only binding to the contract
	PolyNFTQueryTransactor // Write-only binding to the contract
	PolyNFTQueryFilterer   // Log filterer for contract events
}

// PolyNFTQueryCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolyNFTQueryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyNFTQueryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolyNFTQueryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyNFTQueryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolyNFTQueryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyNFTQuerySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolyNFTQuerySession struct {
	Contract     *PolyNFTQuery     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PolyNFTQueryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolyNFTQueryCallerSession struct {
	Contract *PolyNFTQueryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PolyNFTQueryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolyNFTQueryTransactorSession struct {
	Contract     *PolyNFTQueryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PolyNFTQueryRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolyNFTQueryRaw struct {
	Contract *PolyNFTQuery // Generic contract binding to access the raw methods on
}

// PolyNFTQueryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolyNFTQueryCallerRaw struct {
	Contract *PolyNFTQueryCaller // Generic read-only contract binding to access the raw methods on
}

// PolyNFTQueryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolyNFTQueryTransactorRaw struct {
	Contract *PolyNFTQueryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolyNFTQuery creates a new instance of PolyNFTQuery, bound to a specific deployed contract.
func NewPolyNFTQuery(address common.Address, backend bind.ContractBackend) (*PolyNFTQuery, error) {
	contract, err := bindPolyNFTQuery(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PolyNFTQuery{PolyNFTQueryCaller: PolyNFTQueryCaller{contract: contract}, PolyNFTQueryTransactor: PolyNFTQueryTransactor{contract: contract}, PolyNFTQueryFilterer: PolyNFTQueryFilterer{contract: contract}}, nil
}

// NewPolyNFTQueryCaller creates a new read-only instance of PolyNFTQuery, bound to a specific deployed contract.
func NewPolyNFTQueryCaller(address common.Address, caller bind.ContractCaller) (*PolyNFTQueryCaller, error) {
	contract, err := bindPolyNFTQuery(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolyNFTQueryCaller{contract: contract}, nil
}

// NewPolyNFTQueryTransactor creates a new write-only instance of PolyNFTQuery, bound to a specific deployed contract.
func NewPolyNFTQueryTransactor(address common.Address, transactor bind.ContractTransactor) (*PolyNFTQueryTransactor, error) {
	contract, err := bindPolyNFTQuery(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolyNFTQueryTransactor{contract: contract}, nil
}

// NewPolyNFTQueryFilterer creates a new log filterer instance of PolyNFTQuery, bound to a specific deployed contract.
func NewPolyNFTQueryFilterer(address common.Address, filterer bind.ContractFilterer) (*PolyNFTQueryFilterer, error) {
	contract, err := bindPolyNFTQuery(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolyNFTQueryFilterer{contract: contract}, nil
}

// bindPolyNFTQuery binds a generic wrapper to an already deployed contract.
func bindPolyNFTQuery(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PolyNFTQueryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolyNFTQuery *PolyNFTQueryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolyNFTQuery.Contract.PolyNFTQueryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolyNFTQuery *PolyNFTQueryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.PolyNFTQueryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolyNFTQuery *PolyNFTQueryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.PolyNFTQueryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolyNFTQuery *PolyNFTQueryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolyNFTQuery.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolyNFTQuery *PolyNFTQueryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolyNFTQuery *PolyNFTQueryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.contract.Transact(opts, method, params...)
}

// GetAndCheckTokenUrl is a free data retrieval call binding the contract method 0x12d459ab.
//
// Solidity: function getAndCheckTokenUrl(address asset, address user, uint256 tokenId) view returns(bool, string)
func (_PolyNFTQuery *PolyNFTQueryCaller) GetAndCheckTokenUrl(opts *bind.CallOpts, asset common.Address, user common.Address, tokenId *big.Int) (bool, string, error) {
	var out []interface{}
	err := _PolyNFTQuery.contract.Call(opts, &out, "getAndCheckTokenUrl", asset, user, tokenId)

	if err != nil {
		return *new(bool), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// GetAndCheckTokenUrl is a free data retrieval call binding the contract method 0x12d459ab.
//
// Solidity: function getAndCheckTokenUrl(address asset, address user, uint256 tokenId) view returns(bool, string)
func (_PolyNFTQuery *PolyNFTQuerySession) GetAndCheckTokenUrl(asset common.Address, user common.Address, tokenId *big.Int) (bool, string, error) {
	return _PolyNFTQuery.Contract.GetAndCheckTokenUrl(&_PolyNFTQuery.CallOpts, asset, user, tokenId)
}

// GetAndCheckTokenUrl is a free data retrieval call binding the contract method 0x12d459ab.
//
// Solidity: function getAndCheckTokenUrl(address asset, address user, uint256 tokenId) view returns(bool, string)
func (_PolyNFTQuery *PolyNFTQueryCallerSession) GetAndCheckTokenUrl(asset common.Address, user common.Address, tokenId *big.Int) (bool, string, error) {
	return _PolyNFTQuery.Contract.GetAndCheckTokenUrl(&_PolyNFTQuery.CallOpts, asset, user, tokenId)
}

// GetFilterTokensByIndex is a free data retrieval call binding the contract method 0xc5cc029a.
//
// Solidity: function getFilterTokensByIndex(address asset, address ignore, uint256 start, uint256 length) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQueryCaller) GetFilterTokensByIndex(opts *bind.CallOpts, asset common.Address, ignore common.Address, start *big.Int, length *big.Int) (bool, []byte, error) {
	var out []interface{}
	err := _PolyNFTQuery.contract.Call(opts, &out, "getFilterTokensByIndex", asset, ignore, start, length)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

// GetFilterTokensByIndex is a free data retrieval call binding the contract method 0xc5cc029a.
//
// Solidity: function getFilterTokensByIndex(address asset, address ignore, uint256 start, uint256 length) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQuerySession) GetFilterTokensByIndex(asset common.Address, ignore common.Address, start *big.Int, length *big.Int) (bool, []byte, error) {
	return _PolyNFTQuery.Contract.GetFilterTokensByIndex(&_PolyNFTQuery.CallOpts, asset, ignore, start, length)
}

// GetFilterTokensByIndex is a free data retrieval call binding the contract method 0xc5cc029a.
//
// Solidity: function getFilterTokensByIndex(address asset, address ignore, uint256 start, uint256 length) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQueryCallerSession) GetFilterTokensByIndex(asset common.Address, ignore common.Address, start *big.Int, length *big.Int) (bool, []byte, error) {
	return _PolyNFTQuery.Contract.GetFilterTokensByIndex(&_PolyNFTQuery.CallOpts, asset, ignore, start, length)
}

// GetOwnerTokensByIndex is a free data retrieval call binding the contract method 0xbbc8e446.
//
// Solidity: function getOwnerTokensByIndex(address asset, address owner, uint256 start, uint256 length) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQueryCaller) GetOwnerTokensByIndex(opts *bind.CallOpts, asset common.Address, owner common.Address, start *big.Int, length *big.Int) (bool, []byte, error) {
	var out []interface{}
	err := _PolyNFTQuery.contract.Call(opts, &out, "getOwnerTokensByIndex", asset, owner, start, length)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

// GetOwnerTokensByIndex is a free data retrieval call binding the contract method 0xbbc8e446.
//
// Solidity: function getOwnerTokensByIndex(address asset, address owner, uint256 start, uint256 length) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQuerySession) GetOwnerTokensByIndex(asset common.Address, owner common.Address, start *big.Int, length *big.Int) (bool, []byte, error) {
	return _PolyNFTQuery.Contract.GetOwnerTokensByIndex(&_PolyNFTQuery.CallOpts, asset, owner, start, length)
}

// GetOwnerTokensByIndex is a free data retrieval call binding the contract method 0xbbc8e446.
//
// Solidity: function getOwnerTokensByIndex(address asset, address owner, uint256 start, uint256 length) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQueryCallerSession) GetOwnerTokensByIndex(asset common.Address, owner common.Address, start *big.Int, length *big.Int) (bool, []byte, error) {
	return _PolyNFTQuery.Contract.GetOwnerTokensByIndex(&_PolyNFTQuery.CallOpts, asset, owner, start, length)
}

// GetTokensByIds is a free data retrieval call binding the contract method 0xc792f10f.
//
// Solidity: function getTokensByIds(address asset, bytes args) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQueryCaller) GetTokensByIds(opts *bind.CallOpts, asset common.Address, args []byte) (bool, []byte, error) {
	var out []interface{}
	err := _PolyNFTQuery.contract.Call(opts, &out, "getTokensByIds", asset, args)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

// GetTokensByIds is a free data retrieval call binding the contract method 0xc792f10f.
//
// Solidity: function getTokensByIds(address asset, bytes args) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQuerySession) GetTokensByIds(asset common.Address, args []byte) (bool, []byte, error) {
	return _PolyNFTQuery.Contract.GetTokensByIds(&_PolyNFTQuery.CallOpts, asset, args)
}

// GetTokensByIds is a free data retrieval call binding the contract method 0xc792f10f.
//
// Solidity: function getTokensByIds(address asset, bytes args) view returns(bool, bytes)
func (_PolyNFTQuery *PolyNFTQueryCallerSession) GetTokensByIds(asset common.Address, args []byte) (bool, []byte, error) {
	return _PolyNFTQuery.Contract.GetTokensByIds(&_PolyNFTQuery.CallOpts, asset, args)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PolyNFTQuery *PolyNFTQueryCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PolyNFTQuery.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PolyNFTQuery *PolyNFTQuerySession) IsOwner() (bool, error) {
	return _PolyNFTQuery.Contract.IsOwner(&_PolyNFTQuery.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_PolyNFTQuery *PolyNFTQueryCallerSession) IsOwner() (bool, error) {
	return _PolyNFTQuery.Contract.IsOwner(&_PolyNFTQuery.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PolyNFTQuery *PolyNFTQueryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolyNFTQuery.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PolyNFTQuery *PolyNFTQuerySession) Owner() (common.Address, error) {
	return _PolyNFTQuery.Contract.Owner(&_PolyNFTQuery.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PolyNFTQuery *PolyNFTQueryCallerSession) Owner() (common.Address, error) {
	return _PolyNFTQuery.Contract.Owner(&_PolyNFTQuery.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PolyNFTQuery *PolyNFTQueryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyNFTQuery.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PolyNFTQuery *PolyNFTQuerySession) RenounceOwnership() (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.RenounceOwnership(&_PolyNFTQuery.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PolyNFTQuery *PolyNFTQueryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.RenounceOwnership(&_PolyNFTQuery.TransactOpts)
}

// SetQueryLimit is a paid mutator transaction binding the contract method 0x574ded3b.
//
// Solidity: function setQueryLimit(uint256 _limit) returns()
func (_PolyNFTQuery *PolyNFTQueryTransactor) SetQueryLimit(opts *bind.TransactOpts, _limit *big.Int) (*types.Transaction, error) {
	return _PolyNFTQuery.contract.Transact(opts, "setQueryLimit", _limit)
}

// SetQueryLimit is a paid mutator transaction binding the contract method 0x574ded3b.
//
// Solidity: function setQueryLimit(uint256 _limit) returns()
func (_PolyNFTQuery *PolyNFTQuerySession) SetQueryLimit(_limit *big.Int) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.SetQueryLimit(&_PolyNFTQuery.TransactOpts, _limit)
}

// SetQueryLimit is a paid mutator transaction binding the contract method 0x574ded3b.
//
// Solidity: function setQueryLimit(uint256 _limit) returns()
func (_PolyNFTQuery *PolyNFTQueryTransactorSession) SetQueryLimit(_limit *big.Int) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.SetQueryLimit(&_PolyNFTQuery.TransactOpts, _limit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PolyNFTQuery *PolyNFTQueryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PolyNFTQuery.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PolyNFTQuery *PolyNFTQuerySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.TransferOwnership(&_PolyNFTQuery.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PolyNFTQuery *PolyNFTQueryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PolyNFTQuery.Contract.TransferOwnership(&_PolyNFTQuery.TransactOpts, newOwner)
}

// PolyNFTQueryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PolyNFTQuery contract.
type PolyNFTQueryOwnershipTransferredIterator struct {
	Event *PolyNFTQueryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PolyNFTQueryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyNFTQueryOwnershipTransferred)
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
		it.Event = new(PolyNFTQueryOwnershipTransferred)
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
func (it *PolyNFTQueryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyNFTQueryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyNFTQueryOwnershipTransferred represents a OwnershipTransferred event raised by the PolyNFTQuery contract.
type PolyNFTQueryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PolyNFTQuery *PolyNFTQueryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PolyNFTQueryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PolyNFTQuery.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PolyNFTQueryOwnershipTransferredIterator{contract: _PolyNFTQuery.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PolyNFTQuery *PolyNFTQueryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PolyNFTQueryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PolyNFTQuery.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyNFTQueryOwnershipTransferred)
				if err := _PolyNFTQuery.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_PolyNFTQuery *PolyNFTQueryFilterer) ParseOwnershipTransferred(log types.Log) (*PolyNFTQueryOwnershipTransferred, error) {
	event := new(PolyNFTQueryOwnershipTransferred)
	if err := _PolyNFTQuery.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
