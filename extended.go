package rtfs

import "context"

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
