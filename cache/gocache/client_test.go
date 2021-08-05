package gocache

import (
	"reflect"
	"testing"
	"time"
)

var client *Client

func init() {
	client = New(1 * time.Minute)
}
func TestClient(t *testing.T) {
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_gocache_client_01",
			args: args{
				key: "test01",
				val: []byte("hello world"),
			},
			wantErr: false,
		},
		{
			name: "test_gocache_client_02",
			args: args{
				key: "test02",
				val: struct {
					name string
					age  int
				}{"xiaoming", 18},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.Set(tt.args.key, tt.args.val)
			got, err := client.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.args.val) {
				t.Errorf("Get() got = %v, want %v", got, tt.args.val)
			}
		})
	}
}
