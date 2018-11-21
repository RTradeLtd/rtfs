package rtfs_test

import (
	"testing"
	"time"

	"github.com/RTradeLtd/rtfs"
)

func TestDHTFindProvs(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, nil, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	err = rtfs.DHTFindProvs(im, "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv", "10")
	if err != nil {
		t.Fatal(err)
	}
}
