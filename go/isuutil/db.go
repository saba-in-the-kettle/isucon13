package isuutil

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	driverName = "mysql"
)

// NewIsuconDB はISUCON用にカスタマイズされたsqlxのDBクライアントを返します。
// 再起動試験対策済み。
func NewIsuconDB(config *mysql.Config) (*sqlx.DB, error) {

	return newIsuconDB(config)
}

// NewIsuconDBFromDSN はISUCON用にカスタマイズされたsqlxのDBクライアントを返します。
// 再起動試験対策済み。
func NewIsuconDBFromDSN(dsn string) (*sqlx.DB, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	return newIsuconDB(config)
}

func newIsuconDB(config *mysql.Config) (*sqlx.DB, error) {
	// ISUCONにおける必須の設定項目たち
	config.ParseTime = true
	config.InterpolateParams = true

	// OpenTelemetryのSQLインストルメンテーションを有効にすることで、Jaegerから発行されているSQLを見れるようにしている
	stdDb, err := sql.Open(driverName, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open to DB: %w", err)
	}

	dbx := sqlx.NewDb(stdDb, driverName)

	// コネクション数はデフォルトでは無制限になっている。
	// 数十から数百くらいで要調整。
	dbx.SetMaxOpenConns(20)
	dbx.SetMaxIdleConns(20)
	dbx.SetConnMaxLifetime(5 * time.Minute)

	// 再起動試験対策
	// Pingして接続が確立するまで待つ
	for {
		if err := dbx.Ping(); err == nil {
			break
		} else {
			fmt.Println(err)
			time.Sleep(time.Second * 2)
		}
	}
	fmt.Println("ISUCON DB ready")

	return dbx, nil
}

// CreateIndexIfNotExists はMySQLのインデックスが存在しない場合に、インデックスを作成します。
// 既に存在する場合はエラーを無視します。
// ISUCONではinitializeのタイミングで、DROP TABLEするのではなくTRUNCATEする場合があります。
// その場合はインデックスは消されず残ってしまうので、Duplicateエラーが発生します。
func CreateIndexIfNotExists(db *sqlx.DB, query string) error {
	_, err := db.Exec(query)

	// 既に存在する場合はエラーになるが、それ以外のエラーはそのまま返す
	var mysqlErr *mysql.MySQLError
	if err != nil {
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1061 || mysqlErr.Number == 1060 {
				fmt.Println("detected already existing index, but it's ok")
				return nil
			}
		}
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

// OverrideAddr はDSNのアドレスを上書きします。
// addr は 127.0.0.1:3306 のような形式で指定してください。
func OverrideAddr(basDSN string, addr string) (*mysql.Config, error) {
	mysqlCfg, err := mysql.ParseDSN(basDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	mysqlCfg.Addr = addr
	return mysqlCfg, nil
}
