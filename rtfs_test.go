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

	"github.com/RTradeLtd/rtfs"
)

// test variables
const (
	testPIN             = "QmNZiPk974vDsPmQii3YbrMKfi12KTSNM7XMiYyiea4VYZ"
	ipnsPath            = "/ipns/Qmd2GzQc68XXicmUpJZUadjsTcPUsXgP1iP1Hp6CYaY4xU"
	testDefaultReadme   = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
	testRefsHash        = "QmPS6VssQGyBYjGQSK8ordvXaU1yUoaUmTfmrV7daLeRPH"
	nodeOneAPIAddr      = "192.168.1.101:5001"
	nodeTwoAPIAddr      = "192.168.2.101:5001"
	remoteNodeMultiAddr = "/ip4/172.218.49.115/tcp/4003/ipfs/Qmct4NniSeuCZ58mSpa7USsJRjCPzL4wTwqmjfa6ANTkMX"
)

func TestInitialize(t *testing.T) {
	_, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSwarmConnect(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	if err = im.SwarmConnect(ctx, remoteNodeMultiAddr); err != nil {
		t.Fatal(err)
	}
}

func TestCustomRequest(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
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
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}

	// create pin
	if err = im.Pin(testPIN); err != nil {
		t.Error(err)
		return
	}

	// check if pin was created
	exists, err := im.CheckPin(testPIN)
	if err != nil {
		t.Error(err)
		return
	}
	if !exists {
		t.Error("pin not found")
		return
	}
}

func TestStat(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	_, err = im.Stat(testPIN)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDagGet(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	var out interface{}
	if err = im.DagGet(testPIN, &out); err != nil {
		t.Fatal(err)
	}
}

func TestDagPut(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
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
	} else if resp == "" {
		t.Fatal("unexpected error occured")
	}
}

func TestNodeAddress(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	if im.NodeAddress() != nodeOneAPIAddr {
		t.Fatal("bad node address retrieved")
	}
}

func TestAdd(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("./Makefile")
	if err != nil {
		t.Fatal(err)
	}
	if resp, err := im.Add(file); err != nil {
		t.Fatal(err)
	} else if resp == "" {
		t.Fatal("unexpected error occured")
	}
}

func TestPubSub_Success(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	if err = im.PubSubPublish("topic", "data"); err != nil {
		t.Fatal(err)
	}
}

func TestPubSub_Failure(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
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
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	newHash, err := im.PatchLink(testDefaultReadme, "testPatchLink", testPIN, false)
	if err != nil {
		t.Fatal(err)
	}
	if newHash != "Qmaga5gbbcihFVvZefTJnKJEfadvgvtPeDnhcbqSHVAnTQ" {
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
}

func TestAppendData(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}

	newHash, err := im.AppendData(testPIN, "hello this is some data")
	if err != nil {
		t.Fatal(err)
	}
	if newHash != "Qmd1SksxuY1aQqcStKv3HTNx9CnTsKhkhu9SqEaR4yrdK6" {
		t.Fatal("failed to correctly append data")
	}
}

func TestSetData(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}

	newHash, err := im.SetData(testPIN, "hello this is some data")
	if err != nil {
		t.Fatal(err)
	}
	if newHash != "QmdfQDSAZXtxvbypJgXXiz3PiC3jwVwujNSbZn5Tkvzq8S" {
		t.Fatal("failed to correctly set data")
	}
}

func TestNewObject(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
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
}

func TestIPNS_Publish_And_Resolve(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := im.Publish(testDefaultReadme, "self", time.Hour*24, time.Hour*24, true)
	if err != nil {
		t.Fatal(err)
	}
	resolvedHash, err := im.Resolve(resp.Name)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Split(resolvedHash, "/")[2] != testDefaultReadme {
		t.Fatal("failed to resolve correct hash")
	}
}

func TestRTFS_Dedups_And_Calculate_Ref_Size(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	size, refs, err := rtfs.DedupAndCalculatePinSize(testRefsHash, im)
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) == 0 {
		t.Fatal("invalid refs count")
	}
	if size != 15729672 {
		t.Fatal("bad size recovered")
	}
}

func TestRTNS_PinUpdate(t *testing.T) {
	var (
		oldPin          = "zb2rheJDzFsa7AsCnSxKimX8eF5wkjriJqeGBamjQF79vr14R"
		newPin          = "QmbB6M914rwm9ZezVd2u8Y2k4g5TRoWWxP3PYKkDipCzpT"
		expectedNewPath = "/ipfs/" + newPin
	)
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	// pin the content first
	if err := im.Pin(oldPin); err != nil {
		t.Fatal(err)
	}
	if err := im.Pin(newPin); err != nil {
		t.Fatal(err)
	}
	newPath, err := im.PinUpdate(oldPin, newPin)
	if err != nil {
		t.Fatal(err)
	}
	if newPath != expectedNewPath {
		t.Fatal("failed to correctly get new path")
	}
}

func TestRefs(t *testing.T) {
	im, err := rtfs.NewManager(nodeOneAPIAddr, "", 5*time.Minute, false)
	if err != nil {
		t.Fatal(err)
	}
	references, err := im.Refs(testDefaultReadme, true, false)
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{
		"QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V",
		"QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		"QmY5heUM5qgRubMDD1og9fhCPA6QdkMp3QCwd4s7gJsyE7",
		"QmejvEPop4D7YUadeGqYWmZxHhLc4JBUCzJJHWMzdcMe2y",
		"QmXgqKTbzdh83pQtKFb19SpMCpDDcKR2ujqk3pKph9aCNF",
		"QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB",
		"QmQ5vhrL7uv6tuoN9KeVBwd4PwfQkXdVVmDLUZuTNxqgvm",
	}
	sort.Strings(expected)
	sort.Strings(references)
	if !reflect.DeepEqual(expected, references) {
		t.Fatal("recovered references not equal to expected")
	}
}
