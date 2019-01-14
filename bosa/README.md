# bosa
Airdrop tools for BOS Mainnet.

## Usage

Choose a binary from the `bin` directory of your favorite platform and rename it to `bosa`.

Then you can type `bosa` to see how it works.

```
bosa is a command-line Swiss Army knife for BOS - by BOSCore.

Usage:
  bosa [command]

Available Commands:
  create      Create accounts for BOS Mainnet.
  help        Help about any command
  updateauth  Update auth for EOS Mainnet msig accounts on BOS Mainnet.
  version     Show the program version

Flags:
  -h, --help            help for bosa
  -v, --verbose         Display verbose output (also see 'output.log')
  -w, --write-actions   Write actions to actions.json (default true)

Use "bosa [command] --help" for more information about a command.
```

### Create Accounts

Create accounts for the snapshot is pretty simple:

```
bosa create
```

### UpdateAuth

Because we do a mapping of EOS Mainnet accounts on BOS, we should update their permission after create their accounts using a predefined public key.

```
bosa updateauth
```

### Config

The example `config.yaml` is a sample config file of `bosa`:

```
mainnet: false
testnet_truncate_snapshot: 5025
http_endpoints:
  - http://localhost:8888
snapshot:
  all: ./data/accounts_info_bos_snapshot.airdrop.csv
  normal: ./data/accounts_info_bos_snapshot.airdrop.normal.csv
  msig: ./data/accounts_info_bos_snapshot.airdrop.msig.csv
  msig_json: ./data/accounts_info_bos_snapshot.airdrop.msig.json
creator:
  name: bos.airdrop
  pubkey: EOS7dFpX3B5xCyudno88n4JYN63vnNs9MrFrAg6dMEbqTK5m6w9Ex
  prikey: 5JtrRrtLwDKA2U8kF5KZWaN9Bjrk7fAuY3n4NmiqyjQ5Nd7dCpv
```

If you wanna use another config file, for example `config-mainnet.yaml`, just type:

```
bosa create ./config-mainnet.yaml
bosa updateauth ./config-mainnet.yaml
```

The table below show their meanings:

| Key                       | Meanings                                                     |
| ------------------------- | ------------------------------------------------------------ |
| mainnet                   | If true, will create all the accounts in the snapshot.       |
| testnet_truncate_snapshot | if mainnet is false, will only create `testnet_truncate_snapshot` numbers of accounts in the snapshot. |
| http_endpoints            | A list of endpoints for `bosa` to send transactions.         |
| snapshot.all              | All the snapshots of BOS Mainnet, including `msig` accounts. `msig` accounts are using a predefined key of `creator.pubkey` in the csv, because we will change it's permissions later using `bosa updateauth` |
| snapshot.normal           | Snapshot of BOS Mainnet accounts of non-msig accounts.       |
| snapshot.msig             | Snapshot of BOS Mainnet accounts of msig accounts using a predefined public key. |
| snapshot.msig_json        | Snapshot of BOS Mainnet accounts of msig accounts.           |
| creator.name              | Creator's account name.                                      |
| creator.pubkey            | Creator's public key.                                        |
| creator.prikey            | Creator's private key.                                       |


