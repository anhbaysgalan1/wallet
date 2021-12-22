package orm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"tp_wallet/internal/block_chain/chain/scan/common"
)

func NewMongoClient(uname, passWard, url string) (*mongo.Client, error) {
	opt := options.Client().ApplyURI(url)
	if uname != "" || passWard != "" {
		opt.Auth = &options.Credential{
			Username: uname,
			Password: passWard,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetHeightByMongo(network string) (MongoBlockHeight, error) {
	if MonCli == nil {
		return MongoBlockHeight{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"net_work", network}}
	var ma MongoBlockHeight
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionNameNumberForBsc).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return MongoBlockHeight{}, err
	}
	return ma, nil
}

func SetHeightByMongo(mh MongoBlockHeight) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionNameNumberForBsc).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetHeightByMongo:", err)
		return err
	}
	return nil
}

func UpdateHeightByMongo(mh MongoBlockHeight) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"net_work", mh.NetWork}}
	update := bson.D{{"$set",
		bson.D{
			{"height", mh.Height},
		},
	}}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionNameNumberForBsc).UpdateOne(
		context.Background(), filter, update)
	if err != nil {
		fmt.Println("UpdateHeightByMongo:", err)
		return err
	}
	return nil
}

func GetH2OByMongo(hash string) (common.H20TransactionForPush, error) {
	if MonCli == nil {
		return common.H20TransactionForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"hash", hash}}
	var ma common.H20TransactionForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionH2O).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.H20TransactionForPush{}, err
	}
	return ma, nil
}

func SetH20ByMongo(mh common.H20TransactionForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionH2O).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetH20ByMongo:", err)
		return err
	}
	return nil
}

func SetFFByMongo(mh common.FFTransactionForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionFF).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetFFByMongo:", err)
		return err
	}
	return nil
}

func SetF1ByMongo(mh common.F1TransactionForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionF1).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetF1ByMongo:", err)
		return err
	}
	return nil
}

func SetBNBByMongo(mh common.BNBTransactionForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionBNB).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetBNBByMongo:", err)
		return err
	}
	return nil
}

func SetMaterialByMongo(mh common.Nft1155TransactionForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionMaterial).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetMaterialByMongo:", err)
		return err
	}
	return nil
}

func SetMaterialCreateByMongo(mh common.Nft1155CreateForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateMaterial).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetMaterialCreateByMongo:", err)
		return err
	}
	return nil
}

func GetRowingNftByMongo(hash string) (common.Nft721TransactionForPush, error) {
	if MonCli == nil {
		return common.Nft721TransactionForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"hash", hash}}
	var ma common.Nft721TransactionForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionRacingBoat).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.Nft721TransactionForPush{}, err
	}
	return ma, nil
}

func GetRowingCreateNftByMongo(hash string) (common.Nft721CreateForPush, error) {
	if MonCli == nil {
		return common.Nft721CreateForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"hash", hash}}
	var ma common.Nft721CreateForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateRacingBoat).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.Nft721CreateForPush{}, err
	}
	return ma, nil
}

func GetRowingCreateNftByMongoContract(contract, nftToken string) (common.Nft721CreateForPush, error) {
	if MonCli == nil {
		return common.Nft721CreateForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"contract", contract}, {"nft_token", nftToken}}
	var ma common.Nft721CreateForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateRacingBoat).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.Nft721CreateForPush{}, err
	}
	return ma, nil
}

func SetRowingCreateNftByMongo(mh common.Nft721CreateForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateRacingBoat).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetRowingCreateNftByMongo:", err)
		return err
	}
	return nil
}

func GetRacerCreateNftByMongo(hash string) (common.Nft721CreateForPush, error) {
	if MonCli == nil {
		return common.Nft721CreateForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"hash", hash}}
	var ma common.Nft721CreateForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateRacer).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.Nft721CreateForPush{}, err
	}
	return ma, nil
}

// GetRacerCreateNftByNftId 通过nftId 拿到nft创建时候的属性
func GetRacerCreateNftByNftId(nft string) (common.Nft721CreateForPush, error) {
	if MonCli == nil {
		return common.Nft721CreateForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"nft", nft}}
	var ma common.Nft721CreateForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateRacer).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.Nft721CreateForPush{}, err
	}
	return ma, nil
}

// GetNftByIdAndContract 通过nftId 和 contract 拿到nft创建时候的属性
func GetNftByIdAndContract(contract, nft string) (common.Nft721CreateForPush, error) {
	switch contract {
	case common.RacerContractAddress:
		return GetRacerCreateNftByNftId(nft)
	case common.RacingBoatContractAddress:
		return GetRowingCreateNftByNftId(nft)
	default:
		return common.Nft721CreateForPush{}, errors.New("error:contract is wrong!")
	}
}

// GetRowingCreateNftByNftId 通过nftId 拿到nft创建时候的属性
func GetRowingCreateNftByNftId(nft string) (common.Nft721CreateForPush, error) {
	if MonCli == nil {
		return common.Nft721CreateForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"nft", nft}}
	var ma common.Nft721CreateForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateRacingBoat).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.Nft721CreateForPush{}, err
	}
	return ma, nil
}

func SetRacerCreateNftByMongo(mh common.Nft721CreateForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionCreateRacer).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetRacerCreateNftByMongo:", err)
		return err
	}
	return nil
}

func SetRowingNftByMongo(mh common.Nft721TransactionForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionRacingBoat).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetTransactionForRacingBoatByMongo:", err)
		return err
	}
	return nil
}

func GetRacerNftByMongo(hash string) (common.Nft721TransactionForPush, error) {
	if MonCli == nil {
		return common.Nft721TransactionForPush{}, errors.New("error:mongo.Client is nil")
	}

	filter := bson.D{{"hash", hash}}
	var ma common.Nft721TransactionForPush
	err := MonCli.Database(DatabaseNameForChain).Collection(CollectionRacer).FindOne(
		context.Background(), filter).Decode(&ma)
	if err != nil {
		return common.Nft721TransactionForPush{}, err
	}
	return ma, nil
}

func SetRacerNftByMongo(mh common.Nft721TransactionForPush) error {
	if MonCli == nil {
		return errors.New("error:mongo.Client is nil")
	}

	_, err := MonCli.Database(DatabaseNameForChain).Collection(CollectionRacer).InsertOne(
		context.Background(), mh)
	if err != nil {
		fmt.Println("SetRacerNftByMongo:", err)
		return err
	}
	return nil
}
