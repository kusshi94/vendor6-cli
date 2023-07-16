package cmd

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/kusshi94/vendor6-cli/pkg/infra"
)

// ipToVendor returns the vendor name for the given IPv6 address.
func ipToVendor(userInput string, db *infra.OUIDb) string {
	// Parse user input as IPv6 address
	nip, err := netip.ParseAddr(userInput)
	if err != nil {
		return fmt.Sprintf("%s is not Valid IPv6 Address", userInput)
	}

	// Chek if the address is IPv6
	if !nip.Is6() {
		return fmt.Sprintf("%s is not IPv6 Adress", nip.String())
	}

	// Get IID from IPv6 address
	iid := getIID(nip)

	// Check if IID is EUI-64
	if !iid.isEUI64() {
		return fmt.Sprintf("%s is not EUI-64 Address", userInput)
	}

	// Get MAC address from IID
	mac := getMAC(iid)

	// Get OUI information from MAC address
	oui := db.Lookup(mac)
	// Return vendor name
	if oui != nil {
		return oui.Company
	}
	return "OUI Not Found"
}

type IID [8]byte

// getIID returns IID from IPv6 address.
func getIID(nip netip.Addr) IID {
	nip16 := nip.As16()
	return IID(nip16[8:])
}

// isEUI64 returns true if IID is EUI-64.
func (iid IID) isEUI64() bool {
	// Check if IID is EUI-64
	return iid[3] == 0xff && iid[4] == 0xfe
}

// getMAC returns MAC address from EUI-64 IID.
func getMAC(eui64 IID) net.HardwareAddr {
	mac := make(net.HardwareAddr, 6)
	mac[0] = eui64[0] ^ 0x02
	mac[1] = eui64[1]
	mac[2] = eui64[2]
	mac[3] = eui64[5]
	mac[4] = eui64[6]
	mac[5] = eui64[7]
	return mac
}
