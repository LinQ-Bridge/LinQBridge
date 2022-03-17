# LinQ

## About LinQ
LinQ runs on Mac OS X, Linux, and Windows. Windows and Mac OS X should be considered experimental - it works fine if you're an app developer but isn't recommended for running nodes.

## Nodes
* Full nodes of supported chains
  Current supported Chain: `Ethereum`, `Klaytn`,`Binance Smart Chain`,`PlatON`.

## Building LinQ
LinQ has one main codebase, written in Go.
#### Prerequisites
* `Go`  
  Requires `Go` version >= 1.14
##### Building LinQ
Clone, compile, and build LinQ:
```shell
make linq
```
Make sure LinQ complied successfully:
```shell
./linq help
```
```
   _     _       _____ 
  | |   (_)     |  _  |
  | |    _ _ __ | | | |
  | |   | | '_ \| | | |
  | |___| | | | \ \/' /
  \_____/_|_| |_|\_/\_\
NAME:
   LinQ - MultiChain-NFT-Bridge Service

USAGE:
   linq [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   tool, t  LinQ Tool
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config <path>  Server config file <path> (default: "./conf/config_devnet.json")
   --nodekey value  P2P node key file
   --help, -h       show help
   --version, -v    print the version

COPYRIGHT:
   Copyright 2022 The LinQ Authors

```
Check available configuration options use:
```shell
./linq tool init --help
```
```
   _     _       _____ 
  | |   (_)     |  _  |
  | |    _ _ __ | | | |
  | |   | | '_ \| | | |
  | |___| | | | \ \/' /
  \_____/_|_| |_|\_/\_\
NAME:
   linq tool init - Server will init or update db when true

USAGE:
   linq tool init [command options] [arguments...]

OPTIONS:
   --config <path>   Server config file <path> (default: "./conf/config_devnet.json")
   --genesis <path>  Server genesis.json file <path> (default: "./conf/genesis.json")
```
## Initialize

#### Generate nodeKey for nodes:
```shell
./linq tool nodekey
```
```
LinQ NodeKey successfully generated
The new file is generated in the current directory
Pubkey Address: <new Pubkey Address>
enode: <new enode info>
```
***Note: nodeKey file must be kept properly.***

#### Generate genesis.json file:
If you want to join the cluster, you need to exchange `enode` and `Pubkey Address` with the cluster nodes to generate a new `Node Pubkey Address set`.  
Convert the `Node Pubkey Address set` to concatenate strings with comma(for example:"Address1,Address2,Address3"), and run this to Generate `genesis.json` file:
```shell
./linq tool genesis -pklist <Node Pubkey Address set String>
```
```
LinQ Genesis.json successfully generated
```
***Note: The `genesis.json` information is not suggested being modified***

#### Generate config.json file:
```shell
./linq tool config
```
```
LinQ config.json successfully generated
```
Open `config.json` and edit it according to the instructions below:
``` json
{
  "RunMode": "testnet", // When RunMode = "testnet", the program runs in the test state
  "DBConfig": { // Mysql database config info
    "Debug": false, // Used to switch whether the database log is output
    "URL": "", // Mysql database url
    "Scheme": "", // Mysql database scheme name
    "User": "", // Mysql database user name
    "Password": "" // Mysql database password
  },
  "LinQConfig": { // LinQ running configuration information
    "DefaultBootNodes": [], // Node enode URL information
    "Addr": "0.0.0.0", // Set the IP address of TCP, Default is "0.0.0.0"
    "Port": 30303 // Set the PORT of TCP, Default is 30303
  },
  "Chains": [ // Used to set listening blockchain node information
    {
      "ChainName": "Klaytn", // The chain name
      "ChainID": 1001,  // The chainId
      "ListenSlot": 5, // Monitoring interval(s)
      "BatchSize": 5, // Maximum number of monitoring blocks per time
      "defer": 1,  // Maximum height delay with chain
      "Nodes": [  //Chain rpc or ws url
        {
          "Url": "" //Multiple url can be set
        }
      ],
      // Contract address
      "CCMContract": "", 
      "NFTProxyContract": "",
      "NFTWrapperContract": "",
      "NFTQueryContract": ""
    }
  ]
}
```
#### Initialize LinQ
After completing the `config.json` modification, initialize LinQ with the following command.
```shell
./linq tool init --config <config.json path> --genesis <genesis.json path>
```
## Running LinQ

Run LinQ with:
```shell
./linq --config <config.json path> --nodekey <nodekey path>
```
```
  Starting LinQ Server
  ...
```
