package mysql

import (
	"testing"
)

func TestNew(t *testing.T) {
	config := &Config{
		User: "root",
		Pass: "123456",
		Host: "127.0.0.1:3306",
		Db: "test",
	}

	t.Run("testsql", func(t *testing.T) {
		sql, err := New(config)
		t.Logf("err:%s", err)
		t.Log("config: ", sql.config)
		count := -1
		err = sql.Client.Table("test").Count(&count).Error
		t.Logf("count: %d", count)
		t.Logf("count err:%s", err)
	})
}
