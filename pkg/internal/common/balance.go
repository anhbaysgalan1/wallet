package common

import (
	"strconv"
	"time"
)

const (
	walletSysBalance      = "tp_wallet:wallet_balance_"
	lockAccountBalance    = "lock:tp_wallet:account_balance_"
	LockAccountBalanceTtl = time.Second * 15
)

func KeyLockAccountBalance(uid uint64) string {
	return lockAccountBalance + strconv.FormatUint(uid, 10)
}

func KeyWalletSysBalance(uid uint64) string {
	return walletSysBalance + strconv.FormatUint(uid, 10)
}
