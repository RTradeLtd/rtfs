package rtfs

// DefaultFSKeystorePath is the default path to an fs keystore
var DefaultFSKeystorePath = "/ipfs/keystore"

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
