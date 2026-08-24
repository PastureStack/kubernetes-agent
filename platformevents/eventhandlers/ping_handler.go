package eventhandlers

import (
	"github.com/PastureStack/control-api-client/v3"
	revents "github.com/PastureStack/event-subscriber/events"
	util "github.com/PastureStack/kubernetes-agent/platformevents/util"
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
