package cmd

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/kusshi94/vendor6-cli/pkg/infra"
)

func ipToVendor(userInput string, db *infra.OUIDb) string {
	// IPv6アドレスをパース
	nip, err := netip.ParseAddr(userInput)
	if err != nil {
		return fmt.Sprintf("%s is not Valid IPv6 Address", userInput)
	}

	// IPv6かどうか判定
	if !nip.Is6() {
		return fmt.Sprintf("%s is not IPv6 Adress", nip.String())
	}

	// IIDを抽出
	iid := getIID(nip)

	// EUI-64判定
	if !iid.isEUI64() {
		return fmt.Sprintf("%s is not EUI-64 Address", userInput)
	}

	// MACアドレスを取得
	mac := getMAC(iid)

	// MACアドレスからOUI情報を取得
	oui := db.Lookup(mac)
	// OUI情報を返す
	if oui != nil {
		return oui.Company
	}
	return "OUI Not Found"
}

type IID [8]byte

// IPアドレスからIIDを取り出す
func getIID(nip netip.Addr) IID {
	nip16 := nip.As16()
	return IID(nip16[8:])
}

func (iid IID) isEUI64() bool {
	// 中央にff:feが来るか？
	return iid[3] == 0xff && iid[4] == 0xfe
}

func getMAC(iid IID) net.HardwareAddr {
	mac := make(net.HardwareAddr, 6)
	mac[0] = iid[0] ^ 0x02
	mac[1] = iid[1]
	mac[2] = iid[2]
	mac[3] = iid[5]
	mac[4] = iid[6]
	mac[5] = iid[7]
	return mac
}
