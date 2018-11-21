package rtfs_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/RTradeLtd/rtfs"
	ci "github.com/libp2p/go-libp2p-crypto"
)

func TestKeystoreManager_noCustomPath(t *testing.T) {
	rtfs.NewKeystoreManager()
}

func TestKeystoreManager(t *testing.T) {
	defer func() {
		if err := os.RemoveAll("temp"); err != nil {
			t.Fatal(err)
		}
	}()

	km, err := rtfs.NewKeystoreManager("temp")
	if err != nil {
		t.Fatal(err)
	}

	keys, err := km.ListKeyIdentifiers()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range keys {
		_, err := km.GetPrivateKeyByName(v)
		if err != nil {
			t.Error(err)
			continue
		}
		present, err := km.CheckIfKeyExists(v)
		if err != nil {
			t.Error(err)
			continue
		}
		if !present {
			t.Error("key not present when it should be")
			continue
		}
	}

	present, err := km.CheckIfKeyExists("thiskeyshouldreallynotexistwithsucharandomname")
	if err != nil {
		t.Fatal(err)
	}
	if present {
		t.Fatal("key found when it should'nt have been")
	}
	// DO NOT USE 1024 IN PRODUCTION, >= 2048
	pk, _, err := ci.GenerateKeyPair(ci.RSA, 1024)
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}

	hexed := hex.EncodeToString(b)
	err = km.SavePrivateKey(hexed, pk)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetKey(t *testing.T) {
	defer func() {
		if err := os.RemoveAll("temp"); err != nil {
			t.Fatal(err)
		}
	}()

	var (
		k1 = "b6ec4a647a7738ef8eea3b21777ecf41630d6d0ac79dc36739d81e927f910a65"
		k2 = "test1"
	)

	km, err := rtfs.NewKeystoreManager("temp")
	if err != nil {
		t.Fatal(err)
	}

	// Create keys
	km.CreateAndSaveKey(k1, 1, 1)
	km.CreateAndSaveKey(k2, 1, 1)

	// Check if keys exist
	present, err := km.CheckIfKeyExists(k1)
	if err != nil {
		t.Fatal(err)
	}

	if !present {
		t.Error("key not present when it should be")
	}

	pk1, err := km.GetPrivateKeyByName(k1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", pk1.GetPublic())

	present, err = km.CheckIfKeyExists(k2)
	if err != nil {
		t.Fatal(err)
	}

	if !present {
		t.Fatal("key not present when it should be")
	}

	pk2, err := km.GetPrivateKeyByName(k2)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", pk2.GetPublic())
}

func Blah(t *testing.T) {
	defer func() {
		if err := os.RemoveAll("temp"); err != nil {
			t.Fatal(err)
		}
	}()
	km, err := rtfs.NewKeystoreManager("temp")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		name    string
		keyType int
		size    int
	}
	tests := []struct {
		name string
		args args
	}{
		{"EDKey-Success", args{"edkey1", ci.Ed25519, 256}},
		{"RSAKey-Success", args{"rsakey1", ci.RSA, 2048}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk1, err := km.CreateAndSaveKey(
				tt.args.name,
				tt.args.keyType,
				tt.args.size,
			)
			if err != nil {
				t.Fatal(err)
			}
			mnemonic, err := km.ExportKeyAsMnemonic(tt.args.name)
			if err != nil {
				t.Fatal(err)
			}
			pk2, err := rtfs.MnemonicToKey(mnemonic)
			if err != nil {
				t.Fatal(err)
			}
			if valid := pk1.Equals(pk2); !valid {
				t.Fatal("failed to properly recover key")
			}
		})
	}
}
