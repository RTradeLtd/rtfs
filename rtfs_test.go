package rtfs_test

import (
	"context"
	"testing"

	"github.com/RTradeLtd/rtfs"
)

// test variables
const (
	testPIN        = "QmNZiPk974vDsPmQii3YbrMKfi12KTSNM7XMiYyiea4VYZ"
	nodeOneAPIAddr = "192.168.1.101:5001"
	nodeTwoAPIAddr = "192.168.2.101:5001"
)

func TestInitialize(t *testing.T) {
	_, err := rtfs.NewManager("", nodeOneAPIAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDHTFindProvs(t *testing.T) {
	im, err := rtfs.NewManager("", nodeOneAPIAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = im.DHTFindProvs("QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv", "10")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildCustomRequest(t *testing.T) {
	im, err := rtfs.NewManager("", nodeOneAPIAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := im.BuildCustomRequest(context.Background(),
		nodeOneAPIAddr, "dht/findprovs", nil, "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("received %+v\n", resp)
}

func TestPin(t *testing.T) {
	im, err := rtfs.NewManager("", nodeOneAPIAddr, nil)
	if err != nil {
		t.Fatal(err)
	}

	// create pin
	if err = im.Pin(testPIN); err != nil {
		t.Error(err)
		return
	}

	// check if pin was created
	exists, err := im.ParseLocalPinsForHash(testPIN)
	if err != nil {
		t.Error(err)
		return
	}
	if !exists {
		t.Error("pin not found")
		return
	}
}

func TestGetObjectFileSizeInBytes(t *testing.T) {
	im, err := rtfs.NewManager("", nodeOneAPIAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = im.GetObjectFileSizeInBytes(testPIN)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestObjectStat(t *testing.T) {
	im, err := rtfs.NewManager("", nodeOneAPIAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = im.ObjectStat(testPIN)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPubSub(t *testing.T) {
	im, err := rtfs.NewManager("", nodeOneAPIAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = im.PublishPubSubMessage(im.PubTopic, "data"); err != nil {
		t.Error(err)
		return
	}
}
