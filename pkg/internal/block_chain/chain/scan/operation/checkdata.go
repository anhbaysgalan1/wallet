package operation

import (
	"errors"
	"go.uber.org/zap"
	"strings"
	"tp_wallet/internal/block_chain/chain/contract/erc1155"
	"tp_wallet/internal/block_chain/chain/contract/erc20"
	"tp_wallet/internal/block_chain/chain/contract/erc721"
	"tp_wallet/internal/block_chain/chain/scan/common"
	"tp_wallet/pkg/log"
)

func CheckRacingBoatInput(input string) (string, string, string, uint, error) {
	//判断交易类型
	if len(input) < 11 {
		return "", "", "", 0, errors.New("input length is too low")
	}

	switch strings.ToLower(input[:10]) {
	case common.TransferPackInput:
		addr, nftID, err := erc721.ParseTransfer(input)
		if err != nil {
			log.GetLogger().Error("erc721.ParseTransfer", zap.Error(err))
			return "", "", "", 0, err
		}
		return addr, "", nftID, common.RowingNftTransferCode, nil
	case common.ApprovalPackInput:
		addr, nftID, err := erc721.ParseApproval(input)
		if err != nil {
			log.GetLogger().Error("erc721.ParseApproval", zap.Error(err))
			return "", "", "", 0, err
		}
		return addr, "", nftID, common.RowingNftApprovalCode, nil
	case common.TransferPackFromInput:
		addr1, addr2, nftID, err := erc721.ParseTransferFrom(input)
		if err != nil {
			log.GetLogger().Error("erc721.ParseTransferFrom", zap.Error(err))
			return "", "", "", 0, err
		}
		return addr1, addr2, nftID, common.RowingNftTransferFromCode, nil
	case common.CreateAsset721PackInput:
		return "", "", "", common.RowingNftCreateCode, nil
	default:
		log.GetLogger().Error("unknown rowing boat data:", zap.Error(errors.New(input)))
		return "", "", "", common.UnknownCode, errors.New("unknown rowing boat data")
	}
}

// CheckMaterialInput 返回：from 、to、nftID、amount、Code、error
func CheckMaterialInput(input string) (string, string, string, string, uint, error) {
	//判断交易类型
	if len(input) < 11 {
		return "", "", "", "", 0, errors.New("input length is too low")
	}

	switch strings.ToLower(input[:10]) {
	case common.TransferMaterialInput:
		from, to, nftID, amount, err := erc1155.ParseTransfer(input)
		if err != nil {
			log.GetLogger().Error("erc1155.ParseTransfer", zap.Error(err))
			return "", "", "", "", 0, err
		}
		return from, to, nftID, amount, common.MaterialNftTransferCode, nil
	case common.ExpandMaterialInput:
		nftID, amount, err := erc1155.ParseExpand(input)
		if err != nil {
			log.GetLogger().Error(" erc1155.ParseExpand", zap.Error(err))
			return "", "", "", "", 0, err
		}
		return "", "", nftID, amount, common.MaterialNftExpandCode, nil
	case common.CreateMaterialInput:
		amount, announcer, name, err := erc1155.ParseCreateAsset(input)
		if err != nil {
			log.GetLogger().Error("erc1155.ParseCreateAsset", zap.Error(err))
			return "", "", "", "", 0, err
		}
		return "", announcer, amount, name, common.MaterialNftCreateCode, nil
	default:
		log.GetLogger().Error("unknown rowing boat data:", zap.Error(errors.New(input)))
		return "", "", "", "", common.UnknownCode, errors.New("unknown rowing boat data")
	}
}

func CheckRacerInput(input string) (string, string, string, uint, error) {
	//判断交易类型
	if len(input) < 11 {
		return "", "", "", 0, errors.New("input length is too low")
	}

	switch strings.ToLower(input[:10]) {
	case common.TransferPackInput:
		addr, nftID, err := erc721.ParseTransfer(input)
		if err != nil {
			log.GetLogger().Error("erc721.ParseTransfer", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", nftID, common.RacerNftTransferCode, nil
	case common.ApprovalPackInput:
		addr, nftID, err := erc721.ParseApproval(input)
		if err != nil {
			log.GetLogger().Error("erc721.ParseApproval", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", nftID, common.RacerNftApprovalCode, nil
	case common.TransferPackFromInput:
		addr1, addr2, nftID, err := erc721.ParseTransferFrom(input)
		if err != nil {
			log.GetLogger().Error("erc721.ParseTransferFrom", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr1, addr2, nftID, common.RacerNftTransferFromCode, nil
	case common.CreateAsset721PackInput:
		return "", "", "", common.RacerNftCreateCode, nil
	default:
		log.GetLogger().Error("unknown rowing boat data", zap.Any("hash", input))
		return "", "", "", common.UnknownCode, errors.New("unknown rowing boat data")
	}
}

func CheckH2OInput(input string) (string, string, string, uint, error) {
	//判断交易类型
	if len(input) < 11 {
		return "", "", "", 0, errors.New("input length is too low")
	}

	switch strings.ToLower(input[:10]) {
	case common.TransferPackInput:
		addr, amount, err := erc20.ParseTransfer(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseTransfer", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", amount, common.H2OTransferCode, nil
	case common.ApprovalPackInput:
		addr, amount, err := erc20.ParseApproval(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseApproval", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", amount, common.H2OApprovalCode, nil
	case common.TransferPackFromInput:
		addr1, addr2, amount, err := erc20.ParseTransferFrom(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseTransferFrom", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr1, addr2, amount, common.H2OTransferFromCode, nil
	default:
		log.GetLogger().Error("unknown H2O data", zap.Any("hash", input))
		return "", "", "", common.UnknownCode, errors.New("unknown H2O data")
	}
}

func CheckFFInput(input string) (string, string, string, uint, error) {
	//判断交易类型
	if len(input) < 11 {
		return "", "", "", 0, errors.New("input length is too low")
	}

	switch strings.ToLower(input[:10]) {
	case common.TransferPackInput:
		addr, amount, err := erc20.ParseTransfer(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseTransfer", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", amount, common.FFTransferCode, nil
	case common.ApprovalPackInput:
		addr, amount, err := erc20.ParseApproval(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseApproval", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", amount, common.FFApprovalCode, nil
	case common.TransferPackFromInput:
		addr1, addr2, amount, err := erc20.ParseTransferFrom(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseTransferFrom", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr1, addr2, amount, common.FFTransferFromCode, nil
	default:
		log.GetLogger().Error("unknown FF data", zap.Any("hash", input))
		return "", "", "", common.UnknownCode, errors.New("unknown FF data")
	}
}

func CheckF1Input(input string) (string, string, string, uint, error) {
	//判断交易类型
	if len(input) < 11 {
		return "", "", "", 0, errors.New("input length is too low")
	}

	switch strings.ToLower(input[:10]) {
	case common.TransferPackInput:
		addr, amount, err := erc20.ParseTransfer(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseTransfer", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", amount, common.F1TransferCode, nil
	case common.ApprovalPackInput:
		addr, amount, err := erc20.ParseApproval(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseApproval", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr, "", amount, common.F1ApprovalCode, nil
	case common.TransferPackFromInput:
		addr1, addr2, amount, err := erc20.ParseTransferFrom(input)
		if err != nil {
			log.GetLogger().Error("erc20.ParseTransferFrom", zap.Error(err), zap.Any("hash", input))
			return "", "", "", 0, err
		}
		return addr1, addr2, amount, common.F1TransferFromCode, nil
	default:
		log.GetLogger().Error("unknown F1 data", zap.Any("hash", input))
		return "", "", "", common.UnknownCode, errors.New("unknown F1 data")
	}
}
