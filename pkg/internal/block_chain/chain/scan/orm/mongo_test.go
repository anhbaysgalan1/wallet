package orm

import (
	"context"
	"testing"
	"time"
	"tp_wallet/internal/block_chain/chain/scan/common"
)

func TestGetHeightByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mh, err := GetHeightByMongo(CollectionBsc)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mh)
}

func TestSetHeightByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	var mh MongoBlockHeight
	mh.Height = 1000
	mh.NetWork = CollectionBsc
	err = SetHeightByMongo(mh)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUpdateHeightByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	var mh MongoBlockHeight
	mh.Height = 33275170
	mh.NetWork = CollectionBsc
	err = UpdateHeightByMongo(mh)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetTransactionForH20ByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	var mh common.H20TransactionForPush
	mh.To = "1"
	mh.Hash = "2"
	mh.Amount = "1"
	mh.From = "3"
	mh.Nonce = "1"
	mh.Contract = "4"

	err = SetH20ByMongo(mh)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetTransactionForRacerByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	var mh common.Nft721TransactionForPush
	mh.To = "1"
	mh.Hash = "2"
	mh.From = "3"
	mh.Nonce = "1"
	mh.Contract = "4"
	mh.NftToken = "1"

	err = SetRacerNftByMongo(mh)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetTransactionForRacingBoatByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	var mh common.Nft721TransactionForPush
	mh.To = "1"
	mh.Hash = "2"
	mh.From = "3"
	mh.Nonce = "1"
	mh.Contract = "4"
	mh.NftToken = "1"

	err = SetRowingNftByMongo(mh)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetH2OByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mh, err := GetH2OByMongo("2")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mh)
}

func TestGetRacerNftByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mh, err := GetRacerNftByMongo("2")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mh)
}

func TestGetRowingNftByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mh, err := GetRowingNftByMongo("2")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mh)
}

func TestSetRacerCreateNftByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	var mt common.Nft721CreateForPush
	mt.NftToken = "1"
	mt.To = "1"
	mt.Hash = "2"
	mt.From = "3"
	mt.Nonce = "0x1111"
	mt.Contract = "4"
	mt.Status = "success"
	mt.StarRating = "5"
	mt.PropsName = "shan dian"
	mt.BlockNumber = "0x1111"

	err = SetRacerCreateNftByMongo(mt)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetRowingCreateNftByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	var mt common.Nft721CreateForPush
	mt.NftToken = "1"
	mt.To = "1"
	mt.Hash = "2"
	mt.From = "3"
	mt.Nonce = "0x1111"
	mt.Contract = "4"
	mt.Status = "success"
	mt.StarRating = "5"
	mt.PropsName = "shan dian"
	mt.BlockNumber = "0x1111"

	err = SetRowingCreateNftByMongo(mt)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetRacerCreateNftByNftId(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mt, err := GetRacerCreateNftByNftId("1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mt)
}

func TestGetRowingCreateNftByNftId(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mt, err := GetRowingCreateNftByNftId("1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mt)
}

func TestGetRacerCreateNftByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mt, err := GetRacerCreateNftByMongo("2")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mt)
}

func TestGetRowingCreateNftByMongo(t *testing.T) {
	cli, err := NewMongoClient("", "", "mongodb://172.18.6.63:27017")
	if err != nil {
		t.Error(err)
		return
	}
	MonCli = cli
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer cli.Disconnect(ctx)

	mt, err := GetRowingCreateNftByMongo("2")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mt)
}
