package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompatibilityClientWritesOnlySupportedResources(t *testing.T) {
	requests := make(chan string, 2)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v3":
			response.Header().Set("X-API-Schemas", server.URL+"/v3/schemas")
			response.WriteHeader(http.StatusOK)
		case "/v3/schemas":
			_, _ = response.Write([]byte(`{
				"data":[
					{"id":"publish","collectionMethods":["POST"],"links":{"collection":"` + server.URL + `/v3/publish"}},
					{"id":"externalServiceEvent","collectionMethods":["POST"],"links":{"collection":"` + server.URL + `/v3/externalServiceEvents"}}
				]
			}`))
		case "/v3/publish", "/v3/externalServiceEvents":
			if request.Method != http.MethodPost {
				http.Error(response, "method must be POST", http.StatusMethodNotAllowed)
				return
			}
			var payload map[string]interface{}
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				http.Error(response, "invalid JSON", http.StatusBadRequest)
				return
			}
			requests <- request.URL.Path
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write([]byte(`{}`))
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	compatibilityClient, err := NewRancherClient(&ClientOpts{Url: server.URL})
	if err != nil {
		t.Fatalf("NewRancherClient: %v", err)
	}
	if _, err := compatibilityClient.Publish.Create(&Publish{Name: "reply"}); err != nil {
		t.Fatalf("create publish: %v", err)
	}
	if _, err := compatibilityClient.ExternalServiceEvent.Create(&ExternalServiceEvent{
		EventType:   "service.create",
		Environment: &Stack{Kind: "environment"},
		Service:     Service{Kind: "kubernetesService", Name: "api"},
	}); err != nil {
		t.Fatalf("create external service event: %v", err)
	}

	for _, expected := range []string{"/v3/publish", "/v3/externalServiceEvents"} {
		if actual := <-requests; actual != expected {
			t.Fatalf("request path = %s, want %s", actual, expected)
		}
	}
}

func TestExternalServiceEventWireFields(t *testing.T) {
	payload, err := json.Marshal(ExternalServiceEvent{
		ExternalId:  "kubernetes://namespace-id",
		EventType:   "stack.remove",
		Environment: &Stack{Kind: "environment"},
		Service: Service{
			ExternalId: "service-id",
			Kind:       "kubernetesService",
			Name:       "api",
			Selector:   "app=api",
			Uuid:       "legacy-id",
			Vip:        "10.43.0.5",
		},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, field := range []string{"externalId", "eventType", "environment", "service"} {
		if _, ok := decoded[field]; !ok {
			t.Fatalf("wire field %q is absent: %s", field, payload)
		}
	}
}
