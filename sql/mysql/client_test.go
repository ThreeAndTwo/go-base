package mysql

import (
	"testing"
)

func TestNew(t *testing.T) {
	config := &Config{
		user: "root",
		pass: "123456",
		host: "127.0.0.1:3306",
		db: "test",
	}

	t.Run("testsql", func(t *testing.T) {
		sql, err := New(config)
		t.Logf("err:%s", err)
		t.Log("config: ", sql.config)
		count := -1
		err = sql.Client.Table("defi_sync").Count(&count).Error
		t.Logf("count: %d", count)
		t.Logf("count err:%s", err)
	})
}
