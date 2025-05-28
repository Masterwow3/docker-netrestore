package netrestore

import (
	"Masterwow3/docker-netrestore/pkg/vmcompute"
	"Masterwow3/docker-netrestore/pkg/windows_admin"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

type DockerNetworkValue struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	EnableIPv6  bool                   `json:"enableIPv6"`
	NetworkType string                 `json:"networkType"`
	IpamV4Info  string                 `json:"ipamV4Info"`
	Generic     map[string]interface{} `json:"generic"`
}

type Output struct {
	ID       string   `json:"ID"`
	Name     string   `json:"Name"`
	IPv6     bool     `json:"IPv6"`
	Type     string   `json:"Type"`
	Policies []Policy `json:"Policies"`
	Subnets  []Subnet `json:"Subnets"`
}

type MacPool struct {
	EndMacAddress   string `json:"EndMacAddress"`
	StartMacAddress string `json:"StartMacAddress"`
}

type Policy struct {
	Type string `json:"Type"`
	VLAN int    `json:"VLAN"`
}

type Subnet struct {
	GatewayAddress string `json:"GatewayAddress"`
	AddressPrefix  string `json:"AddressPrefix"`
}

type IpamV4InfoEntry struct {
	IPAMData string `json:"IPAMData"`
	PoolID   string `json:"PoolID"`
}

type IPAMData struct {
	AddressSpace string `json:"AddressSpace"`
	Gateway      string `json:"Gateway"`
	Pool         string `json:"Pool"`
}

type DaemonConfig struct {
	DataRoot string `json:"data-root"`
}

func getDockerNetworkDBPath() (string, error) {
	//DEBUG
	//return "local-kv.db", nil

	file, err := os.Open("C:\\ProgramData\\Docker\\config\\daemon.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var config DaemonConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return "", err
	}

	if config.DataRoot == "" {
		config.DataRoot = `C:\ProgramData\docker`
	}
	return path.Join(config.DataRoot, `network\files\local-kv.db`), nil
}

func FixNetwork() error {
	elevated, err := windows_admin.IsElevated()
	if err != nil {
		return err
	}
	if !elevated {
		return fmt.Errorf("run command with admin privileges")
	}

	dockerNetworkDB, err := getDockerNetworkDBPath()
	if err != nil {
		return err
	}

	options := &bbolt.Options{
		ReadOnly: true,
		Timeout:  3 * time.Second,
	}
	db, err := bbolt.Open(dockerNetworkDB, 0600, options)
	if err != nil {
		return err
	}
	defer db.Close()

	prefix := "docker/network/v1.0/network/"

	err = db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			return b.ForEach(func(k, v []byte) error {
				key := string(k)
				if strings.HasPrefix(key, prefix) {
					err := addHNSNetwork(v[8:])
					if err != nil {
						panic(err)
					}
				}
				return nil
			})
		})
	})
	if err != nil {
		return err
	}

	return nil
}

func addHNSNetwork(data []byte) error {
	var raw DockerNetworkValue
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if raw.NetworkType != "nat" {
		return fmt.Errorf("only the network type nat is supported")
	}

	exists, err := vmcompute.HNSNetworkExists([]string{raw.Name, raw.ID})
	if err != nil {
		return err
	}
	fmt.Printf("%s (ID: %s) exists: %v\n", raw.Name, raw.ID, exists)
	if exists {
		return nil // skip
	}

	var ipamEntries []IpamV4InfoEntry
	err = json.Unmarshal([]byte(raw.IpamV4Info), &ipamEntries)
	if err != nil {
		return err
	}

	var ipamData IPAMData
	err = json.Unmarshal([]byte(ipamEntries[0].IPAMData), &ipamData)
	if err != nil {
		return err
	}

	hnsid := ""
	if gen, ok := raw.Generic["com.docker.network.generic"].(map[string]interface{}); ok {
		if val, ok := gen["com.docker.network.windowsshim.hnsid"].(string); ok {
			hnsid = val
		}
	}

	nextVlandId, err := vmcompute.GetNextHNSNetworkVlanId()
	if err != nil {
		return err
	}

	result := Output{
		ID:   hnsid,
		Name: raw.ID,
		IPv6: raw.EnableIPv6,
		Type: raw.NetworkType,
		Policies: []Policy{
			{
				Type: "VLAN",
				VLAN: nextVlandId,
			},
		},
		Subnets: []Subnet{
			{
				GatewayAddress: strings.Split(ipamData.Gateway, "/")[0],
				AddressPrefix:  ipamData.Pool,
			},
		},
	}

	outJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}

	vmcompute.AddHNSNetwork(string(outJSON))

	return nil
}
