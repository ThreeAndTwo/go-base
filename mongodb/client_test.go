package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type fields struct {
	config *Config
}

func TestMongo(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "normal connect",
			fields: fields{
				config: &Config{
					URI:      "mongodb://localhost:27017",
					Username: "admin",
					Password: "123456",
				},
			},
			want: true,
		},
		{
			name: "URI is null",
			fields: fields{
				config: &Config{
					URI:      "",
					Username: "admin",
					Password: "123456",
				},
			},
			want: false,
		},
		{
			name: "AUTH is null",
			fields: fields{
				config: &Config{
					URI:      "mongodb://localhost:27017",
					Username: "admin",
					Password: "",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.fields.config)
			if err != nil {
				t.Errorf("new mongo error %s", err)
				return
			}

			collection := client.Database("testing").Collection("numbers")
			res, err := collection.InsertOne(context.TODO(), bson.D{{"name", "pi"}, {"value", 3.14159}})
			if (err == nil) == tt.want {
				t.Logf("res: %s", res)
			}
		})
	}
}
