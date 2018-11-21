package rtfs

import "context"

// DHTFindProvsResponse is a response from the findprovs command
type DHTFindProvsResponse struct {
	ID        string `json:"id,omitempty"`
	Type      int    `json:"type,omitempty"`
	Responses [][]struct {
		ID    string   `json:"id,omitempty"`
		Addrs []string `json:"addrs,omitempty"`
	} `json:"responses,omitempty"`
	Extra string `json:"extra,omitempty"`
}

// DHTFindProvs is used to find providers of a given CID
// Currently bugged and wil only fetch 1 provider
func DHTFindProvs(im Manager, cid, numProviders string) error {
	var (
		opts = map[string]string{
			"num-providers": numProviders,
		}
		cmd = "dht/findprovs"
		out = DHTFindProvsResponse{}
	)
	resp, err := im.CustomRequest(context.Background(),
		im.NodeAddress(), cmd, opts, cid)
	if err != nil {
		return err
	}
	return resp.Decode(&out)
}
