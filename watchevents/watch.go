package watchevents

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/tools/cache"
)

const informerSyncTimeout = 30 * time.Second

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
