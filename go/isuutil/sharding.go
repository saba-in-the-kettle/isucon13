package isuutil

import (
	"encoding/binary"
	"fmt"
	"github.com/cespare/xxhash"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Shards[T Int] struct {
	shards    []*shard
	weightSum int64
}

type shard struct {
	weight Weight
	db     *sqlx.DB
}

type ShardConfig struct {
	// DisplayName を使うことで人間に分かりやすいようにシャードに名前をつけることができます。
	DisplayName string
	MySQLConfig *mysql.Config
	Weight      Weight
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Weight は余剰方式でのシャードの重みを表します。
// 重みが大きいほど、そのシャードに対するクエリが多くなります。
// 全てに均一に分散させたい場合はすべてのシャードの重みを1にしてください。
// 割当しなくないシャードがある場合は、そのシャードの重みを0にしてください。
type Weight int64

// NewShards は baseDSN から各サーバーに接続するためのDBクライアントを作成します。
// ただし、ホスト名のみシャードのものが使われます。
func NewShards[T Int](configs []*ShardConfig) (*Shards[T], error) {
	shards := []*shard{}
	var weightSum int64
	for _, config := range configs {
		db, err := NewIsuconDB(config.MySQLConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to DB: %w", err)
		}
		shards = append(shards, &shard{weight: config.Weight, db: db})
		weightSum += int64(config.Weight)
	}

	return &Shards[T]{
		shards:    shards,
		weightSum: weightSum,
	}, nil
}

// GetShard は shardKey に応じたDBクライアントを返します。
// シャードの選択は余剰方式によって行われます。
func (s *Shards[T]) GetShard(shardKey T) *sqlx.DB {
	key := int64(shardKey)   // 任意の数字
	mod := key % s.weightSum // [0, weightSum) の間に収まるようにする

	var boundary int64
	for _, s := range s.shards {
		boundary += int64(s.weight)
		if mod < boundary {
			return s.db
		}
	}

	panic("unreachable")
}

// GetShardFromIndex は何らかの事情でシャードキーを指定せずにDBにアクセスしたい場合に、シャードのインデックスからシャードを選択します。
// 例えば、すべてのシャードに対してクエリを発行したい場合などに使えます。
func (s *Shards[T]) GetShardFromIndex(index int) *sqlx.DB {
	return s.shards[index].db
}

// StringToShardKey は文字列をシャードキーに変換します。
// 具体的には、文字列の各文字のUnicodeコードポイントを足し合わせたものを返します。
func StringToShardKey(s string) int64 {
	x := xxhash.New()
	_, _ = x.Write([]byte(s))
	return int64(x.Sum64())
}

// IntToShardKey は整数を一様に分散されたシャードキーに変換します。
// 使いたいシャードーキーに偏りがある場合に有効です。
func IntToShardKey[T Int](i T) uint64 {
	x := xxhash.New()
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(i))
	_, _ = x.Write(b[:])
	return x.Sum64()
}
