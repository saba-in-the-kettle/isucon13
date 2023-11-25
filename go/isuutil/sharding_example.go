package isuutil

import (
	"fmt"
	"log"
	"os"
)

func exampleSharding() {
	baseDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s&multiStatements=%t",
		os.Getenv("ISUCON_DB_USER"),
		os.Getenv("ISUCON_DB_PASSWORD"),
		os.Getenv("ISUCON_DB_HOST"),
		os.Getenv("ISUCON_DB_PORT"),
		os.Getenv("ISUCON_DB_NAME"),
		"Asia%2FTokyo",
		false,
	)

	s2Config, err := OverrideAddr(baseDSN, "192.168.0.12")
	if err != nil {
		log.Fatalf("failed to override addr: %v", err)
	}
	s3Config, err := OverrideAddr(baseDSN, "192.168.0.13")
	if err != nil {
		log.Fatalf("failed to override addr: %v", err)
	}
	userShardConfigs := []*ShardConfig{
		{DisplayName: "s2", MySQLConfig: s2Config, Weight: 1},
		{DisplayName: "s3", MySQLConfig: s3Config, Weight: 4},
	}
	userShards, err := NewShards[int64](userShardConfigs)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// -------------------ここまでがmain関数で実行しておくやつ---------------------

	// こんな感じでユーザに対応するシャードを取得できる
	var userID int64 = 123456
	userShard := userShards.GetShard(userID)
	fmt.Println(userShard)
}
