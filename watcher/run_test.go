package watcher

import "testing"

func TestPrepareBinaryName(t *testing.T) {
	res := prepareBinaryName("goldorf")
	if res != "./goldorf" {
		t.Errorf("Expected ./goldorf but got %s for prepareBinaryName", res)
	}
}
