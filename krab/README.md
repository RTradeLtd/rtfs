# krab

`krab` is used to implement a valid ipfs keystore (ie, adheres to the `Keystore` interface) backed by vault. While the existing ipfs keystore solutions work, they are limited to on-disk storage (with no support for encryption) and in-memory storage.