package beam

import (
	"bytes"
	"time"

	"github.com/RTradeLtd/rtfs"
)

// Laser is used to transfer content between two different private networks
type Laser struct {
	src *rtfs.IpfsManager
	dst *rtfs.IpfsManager
}

// NewLaser creates a laser client to beam content between different ipfs networks
func NewLaser(srcURL, dstURL string) (*Laser, error) {
	src, err := rtfs.NewManager(srcURL, nil, time.Minute*10)
	if err != nil {
		return nil, err
	}
	dst, err := rtfs.NewManager(dstURL, nil, time.Minute*10)
	if err != nil {
		return nil, err
	}
	return &Laser{
		src: src,
		dst: dst,
	}, nil
}

// BeamFrom is used to transfer content bewween two different network
func (l *Laser) BeamFrom(source bool, contentHash string) error {
	if source {
		data, err := l.src.Cat(contentHash)
		if err != nil {
			return err
		}
		if _, err = l.dst.Add(bytes.NewReader(data)); err != nil {
			return err
		}
		return nil
	} else {
		data, err := l.dst.Cat(contentHash)
		if err != nil {
			return err
		}
		if _, err = l.src.Add(bytes.NewReader(data)); err != nil {
			return err
		}
		return nil
	}
}
