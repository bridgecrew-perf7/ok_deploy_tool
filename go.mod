module poly-bridge

go 1.14

require (
	github.com/astaxie/beego v1.12.1
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/goleveldb v1.0.0
	github.com/ethereum/go-ethereum v1.9.25
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/joeqian10/neo-gogogo v0.0.0-20201214075916-44b70d175579
	github.com/okex/exchain v0.18.2
	github.com/okex/exchain-go-sdk v0.18.0
	github.com/ontio/ontology v1.11.1-0.20200812075204-26cf1fa5dd47
	github.com/ontio/ontology-crypto v1.0.9
	github.com/ontio/ontology-go-sdk v1.11.4
	github.com/polynetwork/poly v1.3.1
	github.com/polynetwork/poly-go-sdk v0.0.0-20210114035303-84e1615f4ad4
	github.com/polynetwork/poly-io-test v0.0.0-20200819093740-8cf514b07750 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.33.9
	github.com/urfave/cli v1.22.4
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.2-exchain3
	github.com/ethereum/go-ethereum v1.9.25 => github.com/ethereum/go-ethereum v1.9.15
	github.com/joeqian10/neo-gogogo => github.com/blockchain-develop/neo-gogogo v0.0.0-20210126025041-8d21ec4f0324
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.3-exchain
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-exchain2
)
