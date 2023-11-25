package isuutil

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/sony/sonyflake"
	"math/rand"
)

var sf *sonyflake.Sonyflake

func init() {
	var err error
	sf, err = sonyflake.New(sonyflake.Settings{})
	if err != nil {
		panic(err)
	}
}

// GenerateStringID はIDとして使えるソート可能な文字列を生成します。
// MySQLのインデックスは辞書順のほうがパフォーマンスが高いため、MySQLのテーブルのIDとして使う場合はこちらを使うことを推奨します。
func GenerateStringID() string {
	return xid.New().String()
}

// GenerateUUID はIDとして使えるUUID文字列を生成します。
// 一様に分布されたIDを生成したいときに使ってください。
// なお、MySQLのテーブルのIDとして使いたい場合は、辞書順に並んでいる GenerateStringID を使った方がパフォーマンスが高いです。
func GenerateUUID() string {
	return uuid.NewString()
}

// GenerateIntID はIDとして使える64bit整数を生成します。
// このIDは単調増加なので、AUTO INCREMENTの代わりに使うことができます。
// 大体 `488152078426835572` くらいの値が生成されます。int32には収まりきらないので注意してください。
func GenerateIntID() int64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	// Snoflakeの最終bitはMachine IDなので固定になってしまう。そこで、ランダムな値を足すことで分散させる。
	// (modをとってシャーディングをしたいときに有効)
	return int64(id) + rand.Int63n(10)
}
