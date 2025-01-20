module go.sia.tech/walletd

go 1.23.1

toolchain go1.23.2

require (
	github.com/mattn/go-sqlite3 v1.14.24
	go.sia.tech/core v0.9.0
	go.sia.tech/coreutils v0.10.0
	go.sia.tech/jape v0.12.1
	go.sia.tech/web/walletd v0.27.0
	go.uber.org/zap v1.27.0
	golang.org/x/term v0.28.0
	gopkg.in/yaml.v3 v3.0.1
	lukechampine.com/flagg v1.1.1
	lukechampine.com/frand v1.5.1
	lukechampine.com/upnp v0.3.0
)

require (
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	go.etcd.io/bbolt v1.3.11 // indirect
	go.sia.tech/mux v1.3.0 // indirect
	go.sia.tech/web v0.0.0-20240610131903-5611d44a533e // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
)

// replace go.sia.tech/core => ../sia-core
// replace go.sia.tech/coreutils => ../sia-coreutils
replace go.sia.tech/core => github.com/komodoplatform/sia-core v0.0.0-20241122201700-fbeb493c0be1

replace go.sia.tech/coreutils => github.com/komodoplatform/sia-coreutils v0.0.0-20241122153804-5873e645f8b8
