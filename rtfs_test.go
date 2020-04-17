package rtfs_test

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	ipfsapi "github.com/RTradeLtd/go-ipfs-api/v3"
	"github.com/RTradeLtd/rtfs/v2"
)

// test variables
const (
	testPIN                = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
	ipnsPath               = "/ipns/Qmd2GzQc68XXicmUpJZUadjsTcPUsXgP1iP1Hp6CYaY4xU"
	testDefaultReadme      = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
	testRefsHash           = "QmPS6VssQGyBYjGQSK8ordvXaU1yUoaUmTfmrV7daLeRPH"
	nodeOneAPIAddr         = "127.0.0.1:5001"
	remoteNodeOneMultiAddr = "/ip4/172.218.49.115/tcp/4003/ipfs/QmXow5Vu8YXqvabkptQ7HddvNPpbLhXzmmU53yPCM54EQa"
	remoteNodeTwoMultiAddr = "/ip4/172.218.49.115/tcp/4002/ipfs/QmPvnFXWAz1eSghXD6JKpHxaGjbVo4VhBXY2wdBxKPbne5"
)

type args struct {
	addr    string
	token   string
	timeout time.Duration
	direct  bool
}

func TestInitialize(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFilterLogs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := rtfs.NewManager(nodeOneAPIAddr, "", time.Minute*5)
	if err != nil {
		t.Fatal(err)
	}
	logs, err := m.GetLogs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var nilLogger ipfsapi.Logger
	if logs == nilLogger {
		t.Fatal("logs are nil")
	}
}

func TestSwarmConnect(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			if err := im.SwarmConnect(ctx, remoteNodeOneMultiAddr); err != nil {
				if err := im.SwarmConnect(ctx, remoteNodeTwoMultiAddr); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestCustomRequest(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := im.CustomRequest(context.Background(),
		nodeOneAPIAddr, "dht/findprovs", nil, "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("received %+v\n", resp)
}

func TestPin(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			if err = im.Pin(testPIN); err != nil {
				t.Fatal(err)
			}
			if exists, err := im.CheckPin(testPIN); err != nil {
				t.Fatal(err)
			} else if !exists {
				t.Fatal("pin does not exist")
			}
		})
	}
}

func TestStat(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			if stat, err := im.Stat(testPIN); err != nil {
				t.Fatal(err)
			} else if stat == nil {
				t.Fatal("failed to retrieve oject stats")
			}
		})
	}
}

func TestDagGet(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			var out interface{}
			if err := im.DagGet(testPIN, &out); err != nil {
				t.Fatal(err)
			} else if out == nil {
				t.Fatal("failed to get dag")
			}
		})
	}
}

func TestDagPut(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			type testDag struct {
				Foo string `json:"foo"`
				Bar string `json:"bar"`
			}
			a := testDag{"hello", "world"}
			marshaled, err := json.Marshal(&a)
			if err != nil {
				t.Fatal(err)
			}
			if resp, err := im.DagPut(marshaled, "json", "cbor"); err != nil {
				t.Fatal(err)
			} else if resp != "bafyreiaopeffny6qlthkjaoqri4qz5ru544mfpjfo3rvkgv4qq2zfjvgtm" {
				t.Fatal("failed to generate correct dag object")
			}
		})
	}
}

func TestNodeAddress(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			if im.NodeAddress() != nodeOneAPIAddr {
				t.Fatal("bad node address")
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			file, err := os.Open("./hello.txt")
			if err != nil {
				t.Fatal(err)
			}
			if resp, err := im.Add(file); err != nil {
				t.Fatal(err)
			} else if resp != "QmdDHMP6quqdW7n2a5uHkCPoeM1bqg7d4hFkZVyR7vYjCS" {
				t.Fatal("bad hash generated")
			}
		})
	}
}

func TestPubSub_Success(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			if err = im.PubSubPublish("topic", "data"); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPubSub_Failure(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	// test topic failure
	if err = im.PubSubPublish("", "data"); err == nil {
		t.Fatal("failed to validate topic")
	}
	if err = im.PubSubPublish("topic", ""); err == nil {
		t.Fatal("failed to validate data")
	}
}

func TestPatchLink(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			newHash, err := im.PatchLink(testDefaultReadme, "testPatchLink", testPIN, false)
			if err != nil {
				t.Fatal(err)
			}
			if newHash != "QmT2d72bKhhzXSQ5TJ72mPdc3sTQmdrwuPqqLfybL6uUVc" {
				t.Fatal("failed to correctly link objects")
			}
			templateObject, err := im.NewObject("unixfs-dir")
			if err != nil {
				t.Fatal(err)
			}
			if _, err = im.PatchLink(templateObject, "a/b/c", templateObject, false); err == nil {
				t.Fatal("failed to detect error")
			}
			newHash, err = im.PatchLink(templateObject, "a/b/c", templateObject, true)
			if err != nil {
				t.Fatal(err)
			}
			if newHash != "QmQ5D3xbMWFQRC9BKqbvnSnHri31GqvtWG1G6rE8xAZf1J" {
				t.Fatal("failed to correct patch object")
			}
		})
	}
}
func TestAppendData(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			newHash, err := im.AppendData(testPIN, "hello this is some data")
			if err != nil {
				t.Fatal(err)
			}
			if newHash != "QmSaEmskKhcd45e94giip1cmHdBq3oydhVL2sNV7Tf474u" {
				t.Fatal("failed to correctly append data")
			}
		})
	}
}

func TestSetData(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			newHash, err := im.SetData(testPIN, "hello this is some data")
			if err != nil {
				t.Fatal(err)
			}
			if newHash != "QmVyk9mFoZmUgw5zMh6GkF7FQRKGxzTxMZmQRSJaGJq9FK" {
				t.Fatal("failed to correctly set data")
			}
		})
	}
}

func TestNewObject(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			hash, err := im.NewObject("")
			if err != nil {
				t.Fatal(err)
			}
			if hash != "QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n" {
				t.Fatal("failed to generate new object")
			}
			hash, err = im.NewObject("faketemplate")
			if err == nil {
				t.Fatal("failed to recognize invalid template")
			}
			hash, err = im.NewObject("unixfs-dir")
			if err != nil {
				t.Fatal(err)
			}
			if hash != "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn" {
				t.Fatal("failed to generate unixfs-dir template object")
			}
		})
	}
}

func TestIPNS_Publish_And_Resolve(t *testing.T) {
	t.Skip()
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			var (
				resp         *ipfsapi.PublishResponse
				resolvedHash string
			)
			if tt.name == "Direct" {
				resp, err = im.Publish(testDefaultReadme, "self", time.Hour*24, time.Hour*24, true)
				if err != nil {
					t.Fatal(err)
				}
				resolvedHash, err = im.Resolve(resp.Name)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				resp, err = im.Publish(testDefaultReadme, "self", time.Hour*24, time.Hour*24, true)
				if err != nil {
					t.Fatal(err)
				}
				resolvedHash, err = im.Resolve(resp.Name)
				if err != nil {
					t.Fatal(err)
				}
			}
			if strings.Split(resolvedHash, "/")[2] != testDefaultReadme {
				t.Fatal("failed to resolve correct hash")
			}
		})
	}
}

func TestRTFS_Dedups_And_Calculate_Ref_Size(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			size, refs, err := rtfs.DedupAndCalculatePinSize(testPIN, im)
			if err != nil {
				t.Fatal(err)
			}
			if len(refs) != 7 {
				t.Fatal("invalid refs count")
			}
			if size != 6169 {
				t.Fatal("bad size recovered")
			}
			size2, err := im.DeduplicatedSize(testPIN)
			if err != nil {
				t.Fatal(err)
			}
			if int64(size2) != size {
				t.Fatal("bad sizere covered")
			}
		})
	}
}

func TestRTNS_PinUpdate(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			hash, err := im.Add(strings.NewReader("hello"))
			if err != nil {
				t.Fatal(err)
			}
			newHash, err := im.AppendData(testPIN, hash)
			if err != nil {
				t.Fatal(err)
			}
			if err := im.Pin(testPIN); err != nil {
				t.Fatal(err)
			}
			newPin, err := im.PinUpdate(testPIN, newHash)
			if err != nil {
				t.Fatal(err)
			}
			if newPin != "QmXf3Mfh4bkdh2TuYiVN8EcEVXyAhoDpvriCXgQAj6hCzu" {
				t.Fatal("bad pin update")
			}
		})
	}
}

func TestRefs(t *testing.T) {
	expected := []string{
		"QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V",
		"QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		"QmY5heUM5qgRubMDD1og9fhCPA6QdkMp3QCwd4s7gJsyE7",
		"QmejvEPop4D7YUadeGqYWmZxHhLc4JBUCzJJHWMzdcMe2y",
		"QmXgqKTbzdh83pQtKFb19SpMCpDDcKR2ujqk3pKph9aCNF",
		"QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB",
		"QmQ5vhrL7uv6tuoN9KeVBwd4PwfQkXdVVmDLUZuTNxqgvm",
	}
	tests := []struct {
		name string
		args args
	}{
		{"Non-Direct", args{nodeOneAPIAddr, "", time.Minute * 5, false}},
		{"Direct", args{nodeOneAPIAddr, "hello", time.Minute * 5, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.addr, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			var references []string
			if tt.name == "Direct" {
				references, err = im.Refs(testDefaultReadme, true, false)
			} else {
				references, err = im.Refs(testDefaultReadme, true, false)
			}
			if err != nil {
				t.Fatal(err)
			}
			sort.Strings(expected)
			sort.Strings(references)
			if !reflect.DeepEqual(expected, references) {
				t.Fatal("recovered references not equal to expected")
			}
		})
	}
}

func TestAdd_Dir(t *testing.T) {
	type args struct {
		url     string
		token   string
		timeout time.Duration
		dir     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Pass", args{nodeOneAPIAddr, "", time.Minute * 5, "./beam"}, false},
		{"Fail", args{nodeOneAPIAddr, "", time.Minute * 5, "/root/toor"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := rtfs.NewManager(tt.args.url, tt.args.token, tt.args.timeout)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := im.AddDir(tt.args.dir); (err != nil) != tt.wantErr {
				t.Fatalf("AddDir() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
