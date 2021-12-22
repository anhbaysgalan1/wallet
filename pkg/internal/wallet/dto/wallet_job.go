package dto

type WalletJob interface {
	Run()
	Close()
	//JobWalletTransferToBlock(ctx context.Context)
	//JobBlockTransferToWallet(ctx context.Context)
}

type BlockTransferAck struct {
	FromAddr string
	ToAddr   string
	Hash     string
	Amount   string
	Gas      string
}

func (b BlockTransferAck) Check() bool {
	if len(b.Hash) == 0 || len(b.FromAddr) == 0 || len(b.ToAddr) == 0 || len(b.Amount) == 0 {
		return false
	}
	return true
}
