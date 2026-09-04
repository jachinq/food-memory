package service

import (
	"reflect"
	"testing"
)

func TestSplitNames(t *testing.T) {
	got := splitNames("鸡肉，豆腐、辣椒, 豆腐")
	want := []string{"鸡肉", "豆腐", "辣椒"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestRound1(t *testing.T) {
	if round1(4.26) != 4.3 {
		t.Fatalf("round1(4.26)=%v", round1(4.26))
	}
}
