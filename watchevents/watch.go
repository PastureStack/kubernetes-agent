package watchevents

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/tools/cache"
)

const informerSyncTimeout = 30 * time.Second

// Kubernetes 1.12 predates streaming initial lists. New client-go releases
// enable that behaviour by default, while older API servers can leave the
// request open without returning the initial objects. Explicitly retain the
// established list-then-watch sequence for this compatibility agent.
type listThenWatchCompatibility struct{}

func (listThenWatchCompatibility) IsWatchListSemanticsUnSupported() bool {
	return true
}

func compatibleListWatch(listWatch *cache.ListWatch) cache.ListerWatcher {
	return cache.ToListWatcherWithWatchListSemantics(listWatch, listThenWatchCompatibility{})
}

func runAndSyncController(controller cache.Controller, resource string) (chan struct{}, error) {
	stop := make(chan struct{})
	go controller.Run(stop)

	ctx, cancel := context.WithTimeout(context.Background(), informerSyncTimeout)
	defer cancel()
	if !cache.WaitForCacheSync(ctx.Done(), controller.HasSynced) {
		close(stop)
		return nil, fmt.Errorf("synchronize %s watch cache within %s", resource, informerSyncTimeout)
	}

	return stop, nil
}

func closeWatch(stop *chan struct{}) {
	if *stop == nil {
		return
	}
	close(*stop)
	*stop = nil
}
