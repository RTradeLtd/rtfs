package beam_test

import (
	"os"
	"testing"
	"time"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/rtfs"
	"github.com/RTradeLtd/rtfs/beam"
)

func TestBeam(t *testing.T) {
	cfg, err := config.LoadConfig("../testenv/config.json")
	if err != nil {
		t.Fatal(err)
	}
	laser, err := beam.NewLaser(cfg.IPFS.APIConnection.Host+":"+cfg.IPFS.APIConnection.Port, cfg.IPFS.APIConnection.Host+":"+cfg.IPFS.APIConnection.Port)
	if err != nil {
		t.Fatal(err)
	}
	rtfsManager, err := rtfs.NewManager(cfg.IPFS.APIConnection.Host+":"+cfg.IPFS.APIConnection.Port, nil, time.Minute*5)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("../testenv/config.json")
	if err != nil {
		t.Fatal(err)
	}
	cid, err := rtfsManager.Add(file)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 2; i++ {
		if err = laser.Beam(i, cid); err != nil {
			t.Fatal(err)
		}
	}
}
