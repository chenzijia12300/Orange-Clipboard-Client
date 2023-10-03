package db

import (
	"orangeadd.com/clipboard-client/common/resource"
	"os"
	"testing"
	"time"
)

// 测试函数
func TestMain(m *testing.M) {
	resource.InitLog()
	code := m.Run()
	os.Exit(code)
}

func TestInitDB(t *testing.T) {
	InitDB()
}

func TestInsert(t *testing.T) {
	InitDB()
	id := Insert(ClipboardModel{
		Msg:        "Test",
		MsgType:    1,
		CreateTime: time.Now().Unix(),
	})
	t.Log(id)
}

func TestQuery(t *testing.T) {
	InitDB()
	models := Query(0, 10)
	for _, model := range models {
		t.Logf("model:%+v", model)
	}
}
