package isuutil

import "testing"

func TestIntToShardKey(t *testing.T) {
	var i int64 = 1
	t.Log(IntToShardKey(i))
}
