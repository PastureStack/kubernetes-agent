package watchevents

import (
	"testing"

	"k8s.io/client-go/tools/cache"
)

func TestListWatchDoesNotDisableStreamingInitialLists(t *testing.T) {
	listWatch := &cache.ListWatch{}
	capability, ok := any(listWatch).(interface {
		IsWatchListSemanticsUnSupported() bool
	})
	if ok && capability.IsWatchListSemanticsUnSupported() {
		t.Fatal("current Kubernetes support must not force legacy list-then-watch startup")
	}
}
