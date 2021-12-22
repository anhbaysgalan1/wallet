package operation

import (
	"strings"
	"testing"
)

func TestCheckH2OInput(t *testing.T) {
	tests := []struct {
		input      string
		method     string
		wantAddr1  string
		wantAddr2  string
		wantAmount string
		wantCode   uint
	}{
		{
			input:      "0x23b872dd000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d50000000000000000000000008b2ff3eaa80c998302fa20f5012f68716c5e710a00000000000000000000000000000000000000000000003635c9adc5dea00000",
			method:     "transferFrom",
			wantAddr1:  "0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5",
			wantAddr2:  "0x8b2ff3eaa80c998302fa20f5012f68716c5e710a",
			wantAmount: "1000000000000000000000",
			wantCode:   3,
		},
		{
			input:      "0xa9059cbb000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d500000000000000000000000000000000000000000000003635c9adc5dea00000",
			method:     "transfer",
			wantAddr1:  "0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5",
			wantAddr2:  "",
			wantAmount: "1000000000000000000000",
			wantCode:   1,
		},
		{
			input:      "0x095ea7b3000000000000000000000000d75596573b4e691e2ee7cb3b5618b8ab8618c7d500000000000000000000000000000000000000000000003635c9adc5dea00000",
			method:     "approve",
			wantAddr1:  "0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5",
			wantAddr2:  "",
			wantAmount: "1000000000000000000000",
			wantCode:   2,
		},
	}

	for _, v := range tests {
		addr1, addr2, amount, code, err := CheckH2OInput(v.input)

		if err != nil {
			t.Error(err)
			return
		}

		if code != v.wantCode {
			t.Error(v)
			return
		}

		if strings.ToLower(addr1) != v.wantAddr1 {
			t.Error(v)
			return
		}

		if strings.ToLower(addr2) != v.wantAddr2 {
			t.Error(v)
			return
		}

		if amount != v.wantAmount {
			t.Error(v)
			return
		}
	}
}
