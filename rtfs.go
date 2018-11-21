package rtfs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	ipfsapi "github.com/RTradeLtd/go-ipfs-api"
)

// IpfsManager is our helper wrapper for IPFS
type IpfsManager struct {
	shell       *ipfsapi.Shell
	keystore    *KeystoreManager
	nodeAPIAddr string
}

// NewManager is used to initialize our Ipfs manager struct
func NewManager(ipfsURL string, keystore *KeystoreManager) (*IpfsManager, error) {
	// set up shell
	sh := newShell(ipfsURL)
	sh.SetTimeout(time.Minute * 5)
	if _, err := sh.ID(); err != nil {
		return nil, fmt.Errorf("failed to connect to ipfs node at '%s': %s", ipfsURL, err.Error())
	}

	// instantiate manager
	return &IpfsManager{
		shell:       sh,
		nodeAPIAddr: ipfsURL,
		keystore:    keystore,
	}, nil
}

// SetTimeout is used to set a timeout for our api client
func (im *IpfsManager) SetTimeout(time time.Duration) {
	im.shell.SetTimeout(time)
}

// PublishToIPNSDetails is used for fine grained control over IPNS record publishing
func (im *IpfsManager) PublishToIPNSDetails(contentHash, keyName string, lifetime, ttl time.Duration, resolve bool) (*ipfsapi.PublishResponse, error) {
	if im.keystore == nil {
		return nil, errors.New("attempting to create ipns entry with dynamic keys keystore is not enabled/generated yet")
	}

	if keyPresent, err := im.keystore.CheckIfKeyExists(keyName); err != nil {
		return nil, err
	} else if !keyPresent {
		return nil, errors.New("attempting to sign with non existent key")
	}

	return im.shell.PublishWithDetails(contentHash, keyName, lifetime, ttl, resolve)
}

// Pin is a wrapper method to pin a hash to the local node,
// but also alert the rest of the local nodes to pin
// after which the pin will be sent to the cluster
func (im *IpfsManager) Pin(hash string) error {
	if err := im.shell.Pin(hash); err != nil {
		return fmt.Errorf("failed to pin %s: %s", hash, err.Error())
	}
	return nil
}

// Add is a wrapper used to add a file to IPFS
// currently until https://github.com/ipfs/go-ipfs/issues/5376 it is added with no pin
// thus a manual pin must be triggered afterwards
func (im *IpfsManager) Add(r io.Reader) (string, error) {
	return im.shell.AddNoPin(r)
}

// Cat is used to get cat an ipfs object
func (im *IpfsManager) Cat(cid string, out interface{}) ([]byte, error) {
	r, err := im.shell.Cat(cid)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

// DagGet is used to get an ipld object
func (im *IpfsManager) DagGet(cid string, out interface{}) error {
	return im.shell.DagGet(cid, out)
}

// DagPut is used to store data as an ipld object
func (im *IpfsManager) DagPut(data interface{}, encoding, kind string) (string, error) {
	return im.shell.DagPut(data, encoding, kind)
}

// GetObjectFileSizeInBytes is used to retrieve the cumulative byte size of an object
func (im *IpfsManager) GetObjectFileSizeInBytes(key string) (int, error) {
	stat, err := im.shell.ObjectStat(key)
	if err != nil {
		return 0, err
	}
	return stat.CumulativeSize, nil
}

// ObjectStat is used to retrieve the stats about an object
func (im *IpfsManager) ObjectStat(key string) (*ipfsapi.ObjectStats, error) {
	return im.shell.ObjectStat(key)
}

// ParseLocalPinsForHash checks whether or not a pin is present
func (im *IpfsManager) ParseLocalPinsForHash(hash string) (bool, error) {
	pins, err := im.shell.Pins()
	if err != nil {
		return false, err
	}
	if info := pins[hash]; info.Type != "" {
		return true, nil
	}
	return false, nil
}

// BuildCustomRequest is used to build a custom request
func (im *IpfsManager) BuildCustomRequest(ctx context.Context, url, commad string, opts map[string]string, args ...string) (*ipfsapi.Response, error) {
	req := ipfsapi.NewRequest(ctx, url, commad, args...)
	if len(opts) > 0 {
		currentOpts := req.Opts
		for k, v := range opts {
			currentOpts[k] = v
		}
		req.Opts = currentOpts
	}
	hc := &http.Client{Timeout: time.Minute * 1}
	resp, err := req.Send(hc)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DHTFindProvs is used to find providers of a given CID
// Currently bugged and wil only fetch 1 provider
func (im *IpfsManager) DHTFindProvs(cid, numProviders string) error {
	var (
		opts = map[string]string{
			"num-providers": numProviders,
		}
		cmd = "dht/findprovs"
		out = DHTFindProvsResponse{}
	)
	resp, err := im.BuildCustomRequest(context.Background(),
		im.nodeAPIAddr, cmd, opts, cid)
	if err != nil {
		return err
	}
	return resp.Decode(&out)
}
