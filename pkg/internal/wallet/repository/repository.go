package repository

import (
	"context"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.mongodb.org/mongo-driver/mongo"
	"tp_wallet/internal/wallet/adapter/block_chain"
	"tp_wallet/internal/wallet/adapter/props"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/internal/wallet/repository/cache"
	"tp_wallet/internal/wallet/repository/db"
	"tp_wallet/pkg/database/redis"
	"tp_wallet/pkg/redisCache/common"
)

type Repository interface {
	AccountGetAndCreate(ctx context.Context, req *walletPb.UidReq) (*walletPb.AccountGetResp, error)
	AddressGetByUid(ctx context.Context, uid uint64) (map[string]struct{}, error)
	BalanceGetByUid(ctx context.Context, uid uint64) (map[string]entity.Amount, error)
	UidGetByAddress(ctx context.Context, addr string) (uint64, error)
	AccountRegister(ctx context.Context, uid uint64, address, currency string) error

	BillCurrencyGetForAsync(ctx context.Context) ([]*entity.Bill, error)
	BillNftGetForAsync(ctx context.Context) ([]*entity.Bill, error)
	BillDealWith(ctx context.Context, bill *entity.Bill) (err error)
	BillDealWithPending(ctx context.Context, bill *entity.Bill) error
	BillDealWithSuccess(ctx context.Context, bill *entity.Bill) error
	BillDealWithFailed(ctx context.Context, bill *entity.Bill) error

	TransferCurrencyForOffline(ctx context.Context, req *walletPb.TransferForOfflineReq) (*walletPb.Empty, error)
	TransferCurrencyCash(ctx context.Context, req *entity.Bill) (*walletPb.Empty, error)
	TransferNftCash(ctx context.Context, req *entity.Bill) (*walletPb.Empty, error)
	BalanceSet(ctx context.Context, uid uint64, amount, currency string, isAdd bool) (entity.Amount, entity.Amount, error)

	NftOwnerCreate(ctx context.Context, nft *entity.NftOwner) error
	NftOwnerSetByToken(ctx context.Context, nft *entity.NftOwner) error
	NftOwnerGetByAddr(ctx context.Context, address string, page *walletPb.Page) ([]*entity.NftOwner, error)
	NftOwnerGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftOwner, error)

	NftContractCreate(ctx context.Context, nft *entity.NftContract) error
	NftContractSetByToken(ctx context.Context, nft *entity.NftContract) error
	NftContractGetByAddr(ctx context.Context, address string, page *walletPb.Page) ([]*entity.NftContract, error)
	NftContractGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftContract, error)

	NonceGetByAddr(ctx context.Context, addr string) (uint64, error)
	NonceIncr(ctx context.Context, addr string, incr int64) error
	UnlockNonceAddr(ctx context.Context, addr string)
	GetAndLockAddr(ctx context.Context, addr string) (uint64, error)
	GetCurrencySysExpendAddr(ctx context.Context, currency string) (string, uint64, error)
	GetNftSysExpendAddr(ctx context.Context, contractType string) (string, uint64, error)
	GetTpSysExpendAddr(ctx context.Context, tp string) (string, uint64, error)

	NftCreate(ctx context.Context, bill *entity.Bill, nftData entity.NftData) error
}

type RepositoryStruct struct {
	Mongo         *mongo.Client
	Cache         cache.WalletCache
	Db            db.WalletDb
	Lock          common.RedisLock
	BlockChainSrv block_chain.BlockChainSrv
	PropsSrv      props.PropsSrv
}

func NewWalletRepository(r *redis.Client, d *mongo.Client, bSrv block_chain.BlockChainSrv, PropsSrv props.PropsSrv) Repository {
	repo := &RepositoryStruct{}
	repo.Db = db.NewWalletDb(d)
	repo.Cache = cache.NewCache(r)
	repo.Lock = common.NewRedisLock(r)
	repo.BlockChainSrv = bSrv
	repo.Mongo = d
	repo.PropsSrv = PropsSrv
	return repo
}
