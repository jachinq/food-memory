package storage

import "testing"

func TestSafeJoinRejectsTraversal(t *testing.T) {
	_, err := SafeJoin("/app/uploads", "/uploads/../../../etc/passwd")
	if err == nil {
		t.Fatal("expected path traversal to fail")
	}
}
