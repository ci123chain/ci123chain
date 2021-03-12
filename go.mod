module github.com/ci123chain/ci123chain

go 1.13

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20200410125653-0a3087568c00 // indirect
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cespare/cp v1.1.1
	github.com/ci123chain/wasm-util v1.0.1
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/ethereum/go-ethereum v1.9.21
	github.com/go-redis/redis/v8 v8.6.0
	github.com/gogo/protobuf v1.3.0
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/mattn/go-isatty v0.0.10
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/pretty66/gosdk v1.0.3
	github.com/rjeczalik/notify v0.9.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.7.0
	github.com/stumble/gorocksdb v0.0.3 // indirect
	github.com/tanhuiya/fabric-crypto v0.0.0-20191114090500-ee2b23759e39
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/iavl v0.12.4
	github.com/tendermint/tendermint v0.32.3
	github.com/tendermint/tm-db v0.2.0
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/wasmerio/go-ext-wasm v0.3.1
	gitlab.oneitfarm.com/bifrost/cilog v0.1.5
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools v2.2.0+incompatible
)

replace github.com/tendermint/tendermint => github.com/ci123chain/tendermint v0.32.7-rc14
