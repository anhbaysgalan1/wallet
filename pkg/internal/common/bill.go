package common

import (
	"strconv"
	"time"
)

const (
	KeyAllHashToBillId = "tp_wallet:all_hash_to_bill_id"
	TransferTime       = time.Second * 30 // 每30s处理一次数据库未上链订单
	TransferQueueLimit = 500              // 每30s处理的最多订单数

	keyQueueBillToPending    = "tp_wallet:queue_bill_to_pending_"
	KeyQueueBillToPendingTtl = time.Minute

	lockCurrencyCash = "cash_currency_"
	LockCurrencyTtl  = time.Second * 3

	lockNftCash = "cash_nft_"
	LockNftTtl  = time.Second * 3
)

func KeyLockCurrencyCash(uid uint64) string {
	return lockCurrencyCash + strconv.FormatUint(uid, 10)
}

func KeyLockNftCash(nftToken string) string {
	return lockNftCash + nftToken
}

func KeyQueueBillToPending(id string) string {
	return keyQueueBillToPending + "id"
}
