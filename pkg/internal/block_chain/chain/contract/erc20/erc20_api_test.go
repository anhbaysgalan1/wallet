package erc20

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
	"testing"
)

func TestMethodInputs(t *testing.T) {
	me, input, err := MethodPackInputs(
		"transfer",
		common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
		new(big.Int).SetUint64(1111))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(len(me))
	t.Log(hexutil.Encode(me))
	t.Log(len(hexutil.Encode(input)))
	t.Log(hexutil.Encode(input))
	t.Log(input)
	// Output: 0x000000000000000000000000c778417e063141139fce010982780140aa0cd5ab0000000000000000000000000000000000000000000000000000000000000457

	t.Log(len(input))
	res, err := MethodUnPackInputs("transfer", input)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range res {
		t.Log(fmt.Sprint(v))
		// Output: 0xc778417E063141139Fce010982780140Aa0cD5Ab  1111
	}
}

func TestMethodPackInputs(t *testing.T) {
	tests := []struct {
		to    string
		value uint64
	}{
		{
			to:    "0xc778417e063141139fce010982780140aa0cd5ab",
			value: 1111,
		},
		{
			to:    "0x9fea72b0c4ec5ce6a0bcce8837cf75f9fc51db14",
			value: 2222,
		},
	}
	for _, v := range tests {
		id, input, err := MethodPackInputs(
			"transfer",
			common.HexToAddress(v.to),
			new(big.Int).SetUint64(v.value),
		)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(len(hexutil.Encode(id)))
		// Output: 0xa9059cbb 0xa9059cbb
		t.Log(len(hexutil.Encode(input)))
		// OutPut: 0x000000000000000000000000c778417e063141139fce010982780140aa0cd5ab0000000000000000000000000000000000000000000000000000000000000457
		// OutPut: 0x0000000000000000000000009fea72b0c4ec5ce6a0bcce8837cf75f9fc51db1400000000000000000000000000000000000000000000000000000000000008ae

		t.Log(len(input))
	}

}

func TestUsdTxParsing(t *testing.T) {
	tests := []struct {
		to       string
		value    uint64
		byteCode string
	}{
		{
			to:       "0xc778417e063141139fce010982780140aa0cd5ab",
			value:    1111,
			byteCode: "000000000000000000000000c778417e063141139fce010982780140aa0cd5ab0000000000000000000000000000000000000000000000000000000000000457",
		},
		{
			to:       "0x9fea72b0c4ec5ce6a0bcce8837cf75f9fc51db14",
			value:    2222,
			byteCode: "0000000000000000000000009fea72b0c4ec5ce6a0bcce8837cf75f9fc51db1400000000000000000000000000000000000000000000000000000000000008ae",
		},
	}

	for _, v := range tests {
		input := []byte(v.byteCode)
		t.Log(len(input))
	}
}

func TestCheckInputToUsdTransfer(t *testing.T) {
	t.Log(strings.ToLower("0xdAC17F958D2ee523a2206206994597C13D831ec7"))
	t.Log(strings.ToLower("0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5"))
}

func TestParseEthUsdTransfer(t *testing.T) {
	addr, value, err := ParseTransfer("0xa9059cbb0000000000000000000000005d8e7c5a551845d380f961b610173ad9f8532bf800000000000000000000000000000000000000000000000000000000121c8c2f")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr, value)
}

func TestGetInputForTransfer(t *testing.T) {
	res, err := GetInputForTransfer(common.HexToAddress("0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5"), new(big.Int).Mul(big.NewInt(1e+18), new(big.Int).SetUint64(1000)))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hexutil.Encode(res))
	// Output: 0xa9059cbb000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d500000000000000000000000000000000000000000000003635c9adc5dea00000
}

func TestGetInputForApproval(t *testing.T) {
	res, err := GetInputForApproval(common.HexToAddress("0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5"), new(big.Int).Mul(big.NewInt(1e+18), new(big.Int).SetUint64(1000)))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(len(res))
	t.Log(hexutil.Encode(res))
	// Output: 0x095ea7b3000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d500000000000000000000000000000000000000000000003635c9adc5dea00000
}

func TestGetInputForTransferFrom(t *testing.T) {
	res, err := GetInputForTransferFrom(
		common.HexToAddress("0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5"),
		common.HexToAddress("0x8b2ff3eaa80c998302fa20f5012f68716c5e710a"),
		new(big.Int).Mul(big.NewInt(1e+18), new(big.Int).SetUint64(1000)))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(len(res))
	t.Log(hexutil.Encode(res))
	// Output:0x23b872dd000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d50000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a00000000000000000000000000000000000000000000003635c9adc5dea00000
}

func TestParseTransfer(t *testing.T) {
	input := "0xa9059cbb000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d500000000000000000000000000000000000000000000003635c9adc5dea00000"
	addr, amount, err := ParseTransfer(input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr)
	t.Log(amount)
}

func TestParseApproval(t *testing.T) {
	input := "0x095ea7b3000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d500000000000000000000000000000000000000000000003635c9adc5dea00000"
	addr, amount, err := ParseApproval(input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr)
	t.Log(amount)
}

func TestParseTransferFrom(t *testing.T) {
	input := "0x23b872dd000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d50000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a00000000000000000000000000000000000000000000003635c9adc5dea00000"
	addr1, addr2, amount, err := ParseTransferFrom(input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(addr1, addr2, amount)
}
