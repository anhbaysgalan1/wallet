package dingding

import "testing"

func TestDingdingSend(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "TestDingdingSend",
			args:    args{"这是一条测试数据 生产环境-卡上余额不足，请及时充值"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DingdingSend(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("DingdingSend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DingdingSend() got = %v, want %v", got, tt.want)
			}
		})
	}
}
