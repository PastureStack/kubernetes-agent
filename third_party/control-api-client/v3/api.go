package client

const (
	publishType              = "publish"
	externalServiceEventType = "externalServiceEvent"
)

// RancherClient is the compatibility surface used by Kubernetes Agent. The
// historical name is retained because it is part of the existing Go API; only
// the two control-plane write endpoints used by the agent are exposed.
type RancherClient struct {
	RancherBaseClient

	Publish              PublishOperations
	ExternalServiceEvent ExternalServiceEventOperations
}

type Publish struct {
	Resource `yaml:"-"`

	Data                 map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`
	Name                 string                 `json:"name,omitempty" yaml:"name,omitempty"`
	PreviousId           string                 `json:"previousId,omitempty" yaml:"previous_id,omitempty"`
	PreviousIds          []string               `json:"previousIds,omitempty" yaml:"previous_ids,omitempty"`
	ResourceId           string                 `json:"resourceId,omitempty" yaml:"resource_id,omitempty"`
	ResourceType         string                 `json:"resourceType,omitempty" yaml:"resource_type,omitempty"`
	Time                 int64                  `json:"time,omitempty" yaml:"time,omitempty"`
	Transitioning        string                 `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string                 `json:"transitioningMessage,omitempty" yaml:"transitioning_message,omitempty"`
}

type ExternalServiceEvent struct {
	Resource `yaml:"-"`

	Environment interface{} `json:"environment,omitempty" yaml:"environment,omitempty"`
	EventType   string      `json:"eventType,omitempty" yaml:"event_type,omitempty"`
	ExternalId  string      `json:"externalId,omitempty" yaml:"external_id,omitempty"`
	Kind        string      `json:"kind,omitempty" yaml:"kind,omitempty"`
	Service     interface{} `json:"service,omitempty" yaml:"service,omitempty"`
}

// Service contains the exact compatibility fields emitted by Kubernetes
// Agent. Additional control-plane fields are intentionally not modelled.
type Service struct {
	Resource `yaml:"-"`

	ExternalId string `json:"externalId,omitempty" yaml:"external_id,omitempty"`
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Selector   string `json:"selector,omitempty" yaml:"selector,omitempty"`
	Uuid       string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Vip        string `json:"vip,omitempty" yaml:"vip,omitempty"`
}

// Stack contains the exact compatibility fields emitted for namespace removal.
type Stack struct {
	Resource `yaml:"-"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

type PublishOperations interface {
	Create(*Publish) (*Publish, error)
}

type ExternalServiceEventOperations interface {
	Create(*ExternalServiceEvent) (*ExternalServiceEvent, error)
}

type PublishClient struct {
	rancherClient *RancherClient
}

func (c *PublishClient) Create(value *Publish) (*Publish, error) {
	response := &Publish{}
	err := c.rancherClient.doCreate(publishType, value, response)
	return response, err
}

type ExternalServiceEventClient struct {
	rancherClient *RancherClient
}

func (c *ExternalServiceEventClient) Create(value *ExternalServiceEvent) (*ExternalServiceEvent, error) {
	response := &ExternalServiceEvent{}
	err := c.rancherClient.doCreate(externalServiceEventType, value, response)
	return response, err
}

func constructClient(baseClient *RancherBaseClientImpl) *RancherClient {
	compatibilityClient := &RancherClient{RancherBaseClient: baseClient}
	compatibilityClient.Publish = &PublishClient{rancherClient: compatibilityClient}
	compatibilityClient.ExternalServiceEvent = &ExternalServiceEventClient{rancherClient: compatibilityClient}
	return compatibilityClient
}

func NewRancherClient(opts *ClientOpts) (*RancherClient, error) {
	baseClient := &RancherBaseClientImpl{Types: map[string]Schema{}}
	compatibilityClient := constructClient(baseClient)
	if err := setupRancherBaseClient(baseClient, opts); err != nil {
		return nil, err
	}
	return compatibilityClient, nil
}
