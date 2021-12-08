package repository

import (
	"context"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/internal/wallet/entity"
)

func (repo RepositoryStruct) NftOwnerCreate(ctx context.Context, nft *entity.NftOwner) error {
	return repo.Db.NftOwnerCreate(ctx, nft)
}
func (repo RepositoryStruct) NftOwnerSetByToken(ctx context.Context, nft *entity.NftOwner) error {
	return repo.Db.NftOwnerSetByToken(ctx, nft)
}
func (repo RepositoryStruct) NftOwnerGetByAddr(ctx context.Context, address string, page *dto.Page) ([]*entity.NftOwner, error) {
	return repo.Db.NftOwnerGetByAddr(ctx, address, page)
}
func (repo RepositoryStruct) NftOwnerGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftOwner, error) {
	return repo.Db.NftOwnerGetByNftToken(ctx, contractToken, token)
}
func (repo RepositoryStruct) NftOwnerGetByGameToken(ctx context.Context, contractToken, token string) (*entity.NftOwner, error) {
	return repo.Db.NftOwnerGetByGameToken(ctx, contractToken, token)
}

func (repo RepositoryStruct) NftContractCreate(ctx context.Context, nft *entity.NftContract) error {
	return repo.Db.NftContractCreate(ctx, nft)
}
func (repo RepositoryStruct) NftContractSetByToken(ctx context.Context, nft *entity.NftContract) error {
	return repo.Db.NftContractSetByToken(ctx, nft)
}
func (repo RepositoryStruct) NftContractGetByAddr(ctx context.Context, address string, page *dto.Page) ([]*entity.NftContract, error) {
	return repo.Db.NftContractGetByAddr(ctx, address, page)
}
func (repo RepositoryStruct) NftContractGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftContract, error) {
	return repo.Db.NftContractGetByNftToken(ctx, contractToken, token)
}
