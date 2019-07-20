module github.com/RTradeLtd/rtfs/v2

go 1.12

require (
	github.com/RTradeLtd/config v2.0.5+incompatible // indirect
	github.com/RTradeLtd/config/v2 v2.1.1
	github.com/RTradeLtd/entropy-mnemonics v0.0.0-20170316012907-7b01a644a636
	github.com/RTradeLtd/go-ipfs-api v0.0.0-20190522213636-8e3700e602fd
	github.com/RTradeLtd/krab v1.0.0
	github.com/ipfs/go-ds-badger v0.0.5 // indirect
	github.com/libp2p/go-libp2p-crypto v0.1.0
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8 // indirect
)

replace github.com/dgraph-io/badger v2.0.0-rc.2+incompatible => github.com/dgraph-io/badger v1.6.0
