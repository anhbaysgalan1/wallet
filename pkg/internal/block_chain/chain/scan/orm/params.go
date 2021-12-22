package orm

import "go.mongodb.org/mongo-driver/mongo"

var (
	MonCli *mongo.Client
)

const (
	DatabaseNameForChain       = "scan"
	CollectionNameNumberForBsc = "blocknumber"
	CollectionBsc              = "bsc"
	CollectionH2O              = "h2o"
	CollectionFF               = "ff"
	CollectionF1               = "f1"
	CollectionBNB              = "bnb"
	CollectionRacer            = "racer"
	CollectionCreateRacer      = "createracer"
	CollectionRacingBoat       = "rab"
	CollectionCreateRacingBoat = "createrab"
	CollectionCreateMaterial   = "creatematerial"
	CollectionMaterial         = "material"
)
