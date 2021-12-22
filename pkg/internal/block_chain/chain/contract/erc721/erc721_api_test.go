package erc721

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"testing"
)

func TestGetInputForCreateAsset(t *testing.T) {
	res, err := GetInputForCreateAsset(
		"angle",
		5,
		common.HexToAddress("0x8b2ff3eaa80c998302fa20f5012f68716c5e710a"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hexutil.Encode(res))
	//Output:0x0324143f000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000050000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a0000000000000000000000000000000000000000000000000000000000000005616e676c65000000000000000000000000000000000000000000000000000000
}

func TestUnPackCreateAssetEven(t *testing.T) {
	str := "0x0000000000000000000000009fea72b0c4ec5ce6a0bcce8837cf75f9fc51db14000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000000000c616e676c65206b696c6c65720000000000000000000000000000000000000000"
	data, err := hexutil.Decode(str)
	if err != nil {
		t.Error(err)
		return
	}
	res1, res2, res3, res4, err := UnPackCreateAssetEven(data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res1)
	//Output: 0x9fea72b0c4ec5ce6a0bcce8837cf75f9fc51db14
	t.Log(res2)
	//Output: 0
	t.Log(res3)
	//Output: angle killer
	t.Log(res4)
	//Output: 5
}

func TestGetInputForApproval(t *testing.T) {
	res, err := GetInputForApproval(
		common.HexToAddress("0x8b2ff3eaa80c998302fa20f5012f68716c5e710a"),
		new(big.Int).SetUint64(1))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hexutil.Encode(res))
	//OutPut:0x095ea7b30000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a0000000000000000000000000000000000000000000000000000000000000001
}

func TestGetInputForTransferFrom(t *testing.T) {
	res, err := GetInputForTransferFrom(
		common.HexToAddress("0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5"),
		common.HexToAddress("0x8b2ff3eaa80c998302fa20f5012f68716c5e710a"),
		new(big.Int).SetUint64(1))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hexutil.Encode(res))
	// Output:0x23b872dd000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d50000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a0000000000000000000000000000000000000000000000000000000000000001
}

func TestParseTransferFrom(t *testing.T) {
	input := "0x23b872dd000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d50000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a0000000000000000000000000000000000000000000000000000000000000001"
	addr1, addr2, value, err := ParseTransferFrom(input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr1)
	t.Log(addr2)
	t.Log(value)
}

func TestParseTransfer(t *testing.T) {
	input := "0xa9059cbb0000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a0000000000000000000000000000000000000000000000000000000000000001"
	addr, value, err := ParseTransfer(input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr)
	t.Log(value)
}

func TestParseApproval(t *testing.T) {
	input := "0x095ea7b30000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a0000000000000000000000000000000000000000000000000000000000000001"
	addr, value, err := ParseApproval(input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr)
	t.Log(value)
}

func TestGetInputForTransfer(t *testing.T) {
	res, err := GetInputForTransfer(
		common.HexToAddress("0x8b2ff3eaa80c998302fa20f5012f68716c5e710a"),
		new(big.Int).SetUint64(1))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hexutil.Encode(res))
	//Output: 0xa9059cbb0000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a0000000000000000000000000000000000000000000000000000000000000001
}
