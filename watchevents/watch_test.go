package watchevents

import (
	"testing"

	"k8s.io/client-go/tools/cache"
)

func TestCompatibleListWatchUsesListThenWatch(t *testing.T) {
	listWatch := compatibleListWatch(&cache.ListWatch{})
	capability, ok := listWatch.(interface {
		IsWatchListSemanticsUnSupported() bool
	})
	if !ok || !capability.IsWatchListSemanticsUnSupported() {
		t.Fatal("Kubernetes 1.12 compatibility requires list-then-watch startup")
	}
}
