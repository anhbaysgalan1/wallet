package orm

type MongoBlockHeight struct {
	Height  uint64 `bson:"height"json:"height"`
	NetWork string `bson:"net_work"json:"net_work"`
}
