package hostwatch

import (
	"github.com/PastureStack/kubernetes-agent/kubernetesclient"
	"github.com/rancher/go-rancher-metadata/metadata"
	log "github.com/sirupsen/logrus"
)

var (
	maxRetryCount int = 3
)

func statusSync(kClient *kubernetesclient.Client, metadataClient metadata.Client) {
	hosts, err := metadataClient.GetHosts()
	if err != nil {
		log.Errorf("Error reading host list from metadata service: [%v]", err)
		return
	}
	for _, host := range hosts {
		mNodeStatus := host.State
		switch mNodeStatus {
		case ActivatingState:
			cordonUncordon(host, kClient, metadataClient, false)
		case DeactivatingState:
			cordonUncordon(host, kClient, metadataClient, true)
		}
	}
}

func cordonUncordon(host metadata.Host, kClient *kubernetesclient.Client, metadataClient metadata.Client, unschedulable bool) {
	changed := false
	for retryCount := 0; retryCount <= maxRetryCount; retryCount++ {
		node, err := kClient.Node.ByName(host.Hostname)
		if err != nil {
			log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
			continue
		}
		if node.Spec.Unschedulable == unschedulable {
			changed = true
			break
		}
		node.Spec.Unschedulable = unschedulable
		_, err = kClient.Node.ReplaceNode(node)
		if err != nil {
			log.Errorf("Error updating node [%s] with new schedulable state, err :[%v]", host.Hostname, err)
			continue
		}
		changed = true
		break
	}
	if !changed {
		log.Errorf("Failed to cordon/uncordon node: [%s]", host.Hostname)
	}
}
