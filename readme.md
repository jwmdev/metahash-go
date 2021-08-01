# MetaHash API GO library

An unofficial Golang library for [#MetaHash](https://metahash.org ) blockchain.
This repository was folked  from the reposity created by [xboston](https://github.com/jwmdev/metahash-go) and I contributed to the original repository. However, due to major changes I made, I decided to create a new repository to avoid introducing breaking changes to the original repository and affect systems that are already depending on the original repository.


![metahash-go](https://raw.githubusercontent.com/jwmdev/metahash-go/master/media/metahash-go.png)

### Requirements

- GO go 1.12+

### Installation

```bash
go get -u github.com/jwmdev/metahash-go
```

### Information

- [Missing #MetaHash API documentation](https://github.com/xboston/metahash-api)
- [Original repository](https://github.com/xboston/metahash-go)
- [Knowledge base](https://developers.metahash.org)
- [Testpage portal](http://testpage.metahash.org/)

## Refence links:
- [Metahash API links](https://talk.metahash.org/t/metahash-api/36)
- [Metahash Node API](https://github.com/metahashorg/metahash-fullnode-client/wiki/Usage)
- [MetaHash Developers Portal](https://metahash.readme.io/docs/as-a-developer)
- [MetaHash Wiki](https://github.com/metahashorg/MetaHash/wiki)

### Methods

- [x] fetch-balance
- [x] fetch-balances
- [x] fetch-history
- [x] get-tx
- [x] get-last-txs
- [x] get-blocks
- [x] get-block-by-number
- [x] get-block-by-hash
- [x] get-dump-block-by-number
- [x] get-dump-block-by-hash
- [x] get-count-blocks
- [x] get-last-node-stat-result
- [x] get-last-node-stat-trust
- [x] get-all-last-nodes-count
- [x] get-nodes-raiting
- [x] get-address-delegations
- [x] get-forging-sum-all
- [x] get-forging-sum
- [x] get-common-balance
- [ ] status
- [ ] mhc_send
- [ ] getinfo


### Extra Methods
- [ ] [generateKey](https://developers.metahash.org/hc/en-us/articles/360002712193-Getting-started-with-Metahash-network)
- [ ] getNonce
- [x] [metahashSupply](https://github.com/metahashorg/MetaHash/wiki/MetaHash-Supply)

### Usage
You can find usage examples in the [examples](https://github.com/jwmdev/metahash-go/tree/master/examples) folder.

## License

This package is released under the MIT license.
