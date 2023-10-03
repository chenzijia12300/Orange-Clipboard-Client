package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"orangeadd.com/clipboard-client/common/resource"
)

type ClipboardModel struct {
	ID         int
	Msg        string
	MsgType    int
	CreateTime int64
}

var DB *sql.DB

const CreateSQL = `
			CREATE TABLE IF NOT EXISTS clipboard (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				msg TEXT,
				msg_type INTEGER,
				create_time INTEGER
        )`

const (
	MsgTextType  = 1
	MsgImageType = 2
)

func InitDB() error {
	db, err := sql.Open("sqlite3", "clipboard.sqlite")
	if err != nil {
		resource.Logger.Error("初始化sqlite失败", zap.Error(err))
		return err
	}
	DB = db
	_, err = DB.Exec(CreateSQL)
	if err != nil {
		resource.Logger.Error("初始化clipboard失败", zap.Error(err))
		return err
	}
	return nil
}

func Query(limit, offset int) []ClipboardModel {
	clipboardModels := make([]ClipboardModel, 0)
	rows, err := DB.Query("SELECT * FROM clipboard LIMIT ?,?", limit, offset)
	if err != nil {
		resource.Logger.Error("查询历史剪贴数据失败", zap.Error(err))
		return clipboardModels
	}
	for rows.Next() {
		var (
			id         int
			msg        string
			msgType    int
			createTime int64
		)
		err := rows.Scan(&id, &msg, &msgType, &createTime)
		if err != nil {
			resource.Logger.Error("row scan error", zap.Error(err))
			return clipboardModels
		}
		clipboardModels = append(clipboardModels, ClipboardModel{
			ID:         id,
			Msg:        msg,
			MsgType:    msgType,
			CreateTime: createTime,
		})
	}
	return clipboardModels
}

func Insert(data ClipboardModel) int64 {
	result, err := DB.Exec("INSERT INTO clipboard(msg,msg_type,create_time) VALUES (?,?,?)", data.Msg, data.MsgType, data.CreateTime)
	if err != nil {
		resource.Logger.Error("insert error", zap.Error(err))
		return -1
	}
	id, _ := result.LastInsertId()
	return id
}
