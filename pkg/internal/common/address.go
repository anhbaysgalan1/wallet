package common

import "strconv"

const (
	addressUidToAddr = "tp_wallet:uid_to_addr_"
	AddressAddrToUid = "tp_wallet:addr_to_uid"
)

func KeyAddressUidToAddr(uid uint64) string {
	return addressUidToAddr + strconv.FormatUint(uid, 10)
}
