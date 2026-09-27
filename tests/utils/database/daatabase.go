package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func CreateDBConnection() (*sql.DB, error) {
	//godotenv.Load()は実行時のカレントディレクトリの.envを探すため
	//repositores/.envをリポジトリから呼んだときに探す為リポジトリからみたパスを指定してあげる
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, err
	}

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASS")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		password,
		host,
		port,
		dbName,
	)

	//main.goではエラーが起きたらコードを終了させているけど
	//こっちではエラーを返させる
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
