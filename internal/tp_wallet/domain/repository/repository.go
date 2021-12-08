package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	blockDto "tp_wallet/internal/block_chain/dto"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/internal/wallet/repository/cache"
	"tp_wallet/internal/wallet/repository/db"
	"tp_wallet/pkg/database/redis"
	"tp_wallet/pkg/redisCache/common"
)

type Repository interface {
	AccountGetAndCreate(ctx context.Context, req *dto.UidReq) (*dto.AccountGetResp, error)
	AddressGetByUid(ctx context.Context, uid uint64) (map[string]struct{}, error)
	BalanceGetByUid(ctx context.Context, uid uint64) (string, error)
	UidGetByAddress(ctx context.Context, addr string) (uint64, error)
	AccountRegister(ctx context.Context, uid uint64, address string) error

	BillH2OGetForAsync(ctx context.Context) ([]*entity.Bill, error)
	BillNftGetForAsync(ctx context.Context) ([]*entity.Bill, error)
	BillDealWith(ctx context.Context, bill *entity.Bill) (err error)
	BillDealWithPending(ctx context.Context, bill *entity.Bill) error
	BillDealWithSuccess(ctx context.Context, bill *entity.Bill) error
	BillDealWithH2OCash(ctx context.Context, bill *entity.Bill) error
	BillDealWithH2OCharge(ctx context.Context, bill *entity.Bill) error
	BillDealWithFailed(ctx context.Context, bill *entity.Bill) error

	TransferH2OForOffline(ctx context.Context, req *dto.TransferForOfflineReq) (*dto.Empty, error)
	TransferH2OCash(ctx context.Context, req *entity.Bill) (*dto.Empty, error)
	TransferNftCash(ctx context.Context, req *entity.Bill) (*dto.Empty, error)
	BalanceSet(ctx context.Context, uid uint64, amount string, isAdd bool) (string, string, error)

	NftOwnerCreate(ctx context.Context, nft *entity.NftOwner) error
	NftOwnerSetByToken(ctx context.Context, nft *entity.NftOwner) error
	NftOwnerGetByAddr(ctx context.Context, address string, page *dto.Page) ([]*entity.NftOwner, error)
	NftOwnerGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftOwner, error)

	NftContractCreate(ctx context.Context, nft *entity.NftContract) error
	NftContractSetByToken(ctx context.Context, nft *entity.NftContract) error
	NftContractGetByAddr(ctx context.Context, address string, page *dto.Page) ([]*entity.NftContract, error)
	NftContractGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftContract, error)

	NonceGetByAddr(ctx context.Context, addr string) (uint64, error)
	NonceIncr(ctx context.Context, addr string, incr int64) error
	UnlockNonceAddr(ctx context.Context, addr string)
	GetAndLockAddr(ctx context.Context, addr string) (uint64, error)
	GetH2OSysExpendAddr(ctx context.Context) (string, uint64, error)
	GetNftSysExpendAddr(ctx context.Context) (string, uint64, error)
	GetTpSysExpendAddr(ctx context.Context, tp string) (string, uint64, error)
}

type RepositoryStruct struct {
	mongo         *mongo.Client
	Cache         cache.WalletCache
	Db            db.WalletDb
	Lock          common.RedisLock
	BlockChainSrv blockDto.BlockChainSrv
}

func NewWalletRepository(r *redis.Client, d *mongo.Client, bSrv blockDto.BlockChainSrv) Repository {
	repo := &RepositoryStruct{}
	repo.Db = db.NewWalletDb(d)
	repo.Cache = cache.NewCache(r)
	repo.Lock = common.NewRedisLock(r)
	repo.BlockChainSrv = bSrv
	repo.mongo = d
	return repo
}
