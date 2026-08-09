package eventhandlers

import (
	util "github.com/PastureStack/kubernetes-agent/platformevents/util"
	revents "github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v3"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	if err := util.CreateAndPublishReply(event, cli); err != nil {
		return err
	}
	return nil
}
