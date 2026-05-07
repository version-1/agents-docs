package matcher

import (
	"reflect"
	"testing"
)

func TestCacheKeyPatternsNormalizesAndSortsPatterns(t *testing.T) {
	got := CacheKeyPatterns([]string{
		" b/ ",
		"",
		"./a",
		"c/",
		"b",
	})
	want := []string{"a", "b", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalized patterns mismatch: got %v want %v", got, want)
	}
}

func TestCacheKeyPatternsIsOrderIndependent(t *testing.T) {
	a := CacheKeyPatterns([]string{"./a/", "b"})
	b := CacheKeyPatterns([]string{"b", "a"})
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("cache key patterns should match: %v != %v", a, b)
	}
}
