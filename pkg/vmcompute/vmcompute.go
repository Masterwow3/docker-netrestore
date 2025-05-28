//go:build windows
// +build windows

package vmcompute

//https://github.com/microsoft/SDN/blob/master/Kubernetes/windows/hns.psm1

import (
	"encoding/json"
	"fmt"
	"syscall"
	"unsafe"
)

func utf16PtrFromString(s string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return ptr
}

// Policy represents a network policy
type Policy struct {
	Type string `json:"Type"`
	VLAN int    `json:"VLAN"`
}

// OutputEntry represents each entry in the Output array
type OutputEntry struct {
	Policies []Policy `json:"Policies"`
	Name     string   `json:"Name"`
}

// TopLevel represents the top-level JSON structure
type TopLevel struct {
	Success bool          `json:"Success"`
	Output  []OutputEntry `json:"Output"`
}

func findNextFreeVLAN(policies []Policy) (int, error) {
	vlanSet := make(map[int]bool)

	for _, policy := range policies {
		if policy.Type == "VLAN" {
			vlanSet[policy.VLAN] = true
		}
	}

	for i := 1; i <= 4094; i++ {
		if !vlanSet[i] {
			return i, nil
		}
	}
	return -1, fmt.Errorf("no free VLAN ID available")
}

func getAllHNSNetworks() ([]OutputEntry, error) {
	dll := syscall.NewLazyDLL("vmcompute.dll")
	proc := dll.NewProc("HNSCall")

	method := utf16PtrFromString("GET")
	path := utf16PtrFromString("/networks")
	request := utf16PtrFromString("")

	var response *uint16
	r1, _, err := proc.Call(
		uintptr(unsafe.Pointer(method)),
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(request)),
		uintptr(unsafe.Pointer(&response)),
	)
	if r1 != 0 {
		return make([]OutputEntry, 0), err
	}

	if response == nil {
		return make([]OutputEntry, 0), nil
	}

	out := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(response))[:])

	var data TopLevel
	err = json.Unmarshal([]byte(out), &data)
	if err != nil {
		panic(err)
	}

	if !data.Success {
		return make([]OutputEntry, 0), fmt.Errorf("data success is false")
	}

	return data.Output, nil
}

func GetNextHNSNetworkVlanId() (int, error) {
	networks, err := getAllHNSNetworks()
	if err != nil {
		return -1, err
	}
	var allPolicies []Policy
	for _, entry := range networks {
		allPolicies = append(allPolicies, entry.Policies...)
	}

	return findNextFreeVLAN(allPolicies)
}

func HNSNetworkExists(names []string) (bool, error) {
	networks, err := getAllHNSNetworks()
	if err != nil {
		return true, err
	}

	for _, name := range names {
		for _, entry := range networks {
			if entry.Name == name {
				return true, nil
			}
		}
	}

	return false, nil
}

func AddHNSNetwork(request string) error {
	dll := syscall.NewLazyDLL("vmcompute.dll")
	proc := dll.NewProc("HNSCall")

	method := utf16PtrFromString("POST")
	path := utf16PtrFromString("/networks")
	requestPtr := utf16PtrFromString(request)

	var response *uint16
	r1, _, err := proc.Call(
		uintptr(unsafe.Pointer(method)),
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(requestPtr)),
		uintptr(unsafe.Pointer(&response)),
	)
	if r1 != 0 {
		return err
	}

	if response != nil {
		out := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(response))[:])
		fmt.Println("Response:", out)
	}
	return nil
}
