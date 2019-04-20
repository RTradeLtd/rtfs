package krab_test

import (
	"os"
	"testing"

	ci "github.com/libp2p/go-libp2p-crypto"

	"github.com/RTradeLtd/rtfs/v2/krab"
)

const (
	passphrase  = "password123"
	testKeyName = "suchkeymuchencryptverywow"
	dsPath      = "./testds"
)

func TestKrab(t *testing.T) {
	defer func() {
		if err := os.RemoveAll(dsPath); err != nil {
			t.Fatal(err)
		}
	}()
	km, err := krab.NewKrab(krab.Opts{Passphrase: passphrase, DSPath: dsPath, ReadOnly: false})
	if err != nil {
		t.Fatal(err)
	}
	defer km.Close()
	// this should fail
	if has, err := km.Has(testKeyName); err == nil {
		t.Fatal("key was found when it shouldn't have been")
	} else if has {
		t.Fatal("unexpected error occured")
	}
	// create a key
	pk, _, err := ci.GenerateKeyPair(ci.Ed25519, 256)
	if err != nil {
		t.Fatal(err)
	}
	// store the key
	if err := km.Put(testKeyName, pk); err != nil {
		t.Fatal(err)
	}
	// check if key exists, this should pass
	if has, err := km.Has(testKeyName); err != nil {
		t.Fatal(err)
	} else if !has {
		t.Fatal("key not present when it should be")
	}
	// get the key
	if pk, err := km.Get(testKeyName); err != nil {
		t.Fatal(err)
	} else if pk == nil {
		t.Fatal("empty private key, unexpected error occured")
	}
	// test key list
	if list, err := km.List(); err != nil {
		t.Fatal(err)
	} else if len(list) > 1 {
		t.Fatal("only 1 key should be present")
	} else if list[0] != testKeyName {
		t.Fatal("bad key name recovered")
	}
	// delete the key
	if err := km.Delete(testKeyName); err != nil {
		t.Fatal(err)
	}
	// verify key was deleted
	if has, err := km.Has(testKeyName); err == nil {
		t.Fatal("key was found when it shouldn't have been")
	} else if has {
		t.Fatal("unexpected error occured")
	}
}
