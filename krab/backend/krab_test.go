package backend_test

import (
	"context"
	"testing"

	"github.com/RTradeLtd/config"
	pb "github.com/RTradeLtd/grpc/krab"
	"github.com/RTradeLtd/rtfs/krab/backend"
	ci "github.com/libp2p/go-libp2p-crypto"
)

const (
	testCfgPath = "../../testenv/config.json"
	testKeyName = "testkey"
)

func TestKrab(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	// create our server, and listen for connections
	go func() {
		if err := backend.NewServer(cfg.Krab.URL, "tcp", cfg); err != nil {
			t.Fatal(err)
		}
	}()
	// create our client to connect to the server
	client, err := backend.NewClient(cfg.Endpoints)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	// generate a private key to store
	pk, _, err := ci.GenerateKeyPair(ci.Ed25519, 256)
	if err != nil {
		t.Fatal(err)
	}
	// convert the private key to bytes
	pkBytes, err := pk.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	// create a request to store the private key
	putReq := &pb.KeyPut{
		Name:       testKeyName,
		PrivateKey: pkBytes,
	}
	// store the key
	if resp, err := client.ServiceClient.PutPrivateKey(context.Background(), putReq); err != nil {
		t.Fatal(err)
	} else if resp.Status == "" {
		t.Fatal("failed to properly set status")
	}
	// create a request to get the private key
	getReq := &pb.KeyGet{
		Name: testKeyName,
	}
	if resp, err := client.ServiceClient.GetPrivateKey(context.Background(), getReq); err != nil {
		t.Fatal(err)
	} else if resp.Status == "" {
		t.Fatal("failed to properly set status")
	} else {
		// convert the recovered key to verify we recovered the oorrect one
		pk2, err := ci.UnmarshalPrivateKey(resp.PrivateKey)
		if err != nil {
			t.Fatal(err)
		}
		if ok := pk.Equals(pk2); !ok {
			t.Fatal("failed to recover correct private key")
		}
	}
}
