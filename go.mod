module github.com/ci123chain/ci123chain

go 1.16

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20200410125653-0a3087568c00 // indirect
	github.com/armon/go-metrics v0.3.8
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cespare/cp v1.1.1 // indirect
	github.com/ci123chain/wasm-util v1.0.1
	github.com/confio/ics23/go v0.6.6
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/cosmos/iavl v0.16.0
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.9.21
	github.com/go-redis/redis/v8 v8.11.0
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/mattn/go-isatty v0.0.12
	github.com/pkg/errors v0.9.1
	github.com/pretty66/gosdk v1.0.3
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.23.0
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/umbracle/go-web3 v0.0.0-20210428185842-ec1b314b9425
	github.com/wasmerio/wasmer-go v1.0.3
	gitlab.oneitfarm.com/bifrost/sesdk v1.0.28
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.37.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
)

replace (
	github.com/cosmos/iavl => github.com/ci123chain/iavl v0.16.0-ci0
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	//github.com/tendermint/tendermint => ../tendermint-ci
	github.com/tendermint/tendermint => github.com/ci123chain/tendermint v0.32.7-ci44
	github.com/wasmerio/wasmer-go => github.com/ci123chain/wasmer-go v1.0.3-rc2
)
