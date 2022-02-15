package service

import (
	"fmt"
	"testing"
)

func TestManager_WatchKey(t *testing.T) {
	type args struct {
		key      string
		callback func(key string, value []byte)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test",
			args{
				key: "CoinSummer/Sodium/test",
				callback: func(key string, value []byte) {
					fmt.Println(string(value))
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewManager(&Config{
				Addr:       "127.0.0.1:8500",
				Token:      "",
				Datacenter: "",
				WatchKey:   "",
			})
			if err := s.WatchKey(tt.args.key, tt.args.callback); (err != nil) != tt.wantErr {
				t.Errorf("WatchKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
