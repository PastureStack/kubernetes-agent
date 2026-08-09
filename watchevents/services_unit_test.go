package watchevents

import (
	"reflect"
	"testing"
)

func TestServiceSelectorIsStableAndDoesNotMutateSource(t *testing.T) {
	source := map[string]string{"tier": "api", "app": "demo"}
	wantSource := map[string]string{"tier": "api", "app": "demo"}
	want := "app=demo,io.kubernetes.pod.namespace=production,tier=api"

	for i := 0; i < 20; i++ {
		if got := serviceSelector(source, "production"); got != want {
			t.Fatalf("selector = %q, want %q", got, want)
		}
	}
	if !reflect.DeepEqual(source, wantSource) {
		t.Fatalf("source selector was mutated: %#v", source)
	}
}

func TestServiceSelectorPreservesNilSelector(t *testing.T) {
	if got := serviceSelector(nil, "default"); got != "" {
		t.Fatalf("nil selector = %q, want empty", got)
	}
}
