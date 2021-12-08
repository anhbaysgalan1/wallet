package vo

type NftData struct {
	GameName      string `json:"game_name,omitempty" bson:"game_name"`
	NftGameToken  string `json:"nft_game_token,omitempty" bson:"nft_game_token"`
	NftBlockToken string `json:"nft_block_token,omitempty" bson:"nft_block_token"`
	Level         uint64 `json:"level,omitempty" bson:"level"`
}

func (nd NftData) IsEmpty() bool {
	return nd == NftData{}
}
