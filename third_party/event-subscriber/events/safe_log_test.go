package events

import (
	"bytes"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSafeLogValueKeepsUntrustedInputOnOneLine(t *testing.T) {
	got := safeLogValue("valid\r\nforged=true\t\x00end")
	if strings.ContainsAny(got, "\r\n\t\x00") {
		t.Fatalf("safeLogValue retained a control character: %q", got)
	}
	if got != `valid\r\nforged=true\t�end` {
		t.Fatalf("safeLogValue = %q", got)
	}
}

func TestSafeLogValuePreservesUnicodeAndBoundsOutput(t *testing.T) {
	input := "臺灣節點-" + strings.Repeat("a", maxLogValueRunes+20)
	got := safeLogValue(input)
	if !strings.HasPrefix(got, "臺灣節點-") || !strings.HasSuffix(got, "…") {
		t.Fatalf("safeLogValue did not preserve legitimate Unicode or truncate safely: %q", got)
	}
}

func TestSafeEventFieldsExcludePayloadAndReplyDestination(t *testing.T) {
	event := &Event{
		Name:         "host.update\nforged=true",
		ID:           "event-1",
		ReplyTo:      "secret.reply.destination",
		ResourceID:   "host-1",
		ResourceType: "host",
		Data:         map[string]interface{}{"token": "must-not-be-logged"},
	}
	fields := safeEventFields(event)
	if len(fields) != 4 {
		t.Fatalf("safeEventFields returned %d fields, want 4", len(fields))
	}
	serialized := ""
	for key, value := range fields {
		serialized += key + "=" + value.(string) + " "
	}
	if strings.Contains(serialized, "must-not-be-logged") || strings.Contains(serialized, event.ReplyTo) {
		t.Fatalf("safe event fields exposed payload or reply destination: %q", serialized)
	}
	if strings.ContainsAny(serialized, "\r\n") {
		t.Fatalf("safe event fields retained a line break: %q", serialized)
	}
}

func TestDoWorkSanitizesIdentifiersAndDoesNotLogPayload(t *testing.T) {
	logger := log.StandardLogger()
	originalOutput := logger.Out
	originalFormatter := logger.Formatter
	originalLevel := logger.Level
	t.Cleanup(func() {
		logger.SetOutput(originalOutput)
		logger.SetFormatter(originalFormatter)
		logger.SetLevel(originalLevel)
	})

	var output bytes.Buffer
	logger.SetOutput(&output)
	logger.SetFormatter(&log.TextFormatter{DisableColors: true, DisableTimestamp: true})
	logger.SetLevel(log.DebugLevel)

	event := &Event{
		Name:       "unhandled\nlevel=error",
		ID:         "event-1",
		ReplyTo:    "must-not-be-logged.reply",
		ResourceID: "host-1",
		Data:       map[string]interface{}{"token": "must-not-be-logged-token"},
	}
	doWork(event, map[string]EventHandler{}, nil, nopLocker(event))

	got := strings.TrimSuffix(output.String(), "\n")
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("untrusted input created additional log records: %q", got)
	}
	for _, line := range lines {
		if strings.ContainsRune(line, '\r') {
			t.Fatalf("operational log retained a carriage return: %q", got)
		}
	}
	if !strings.Contains(got, "unhandled") || !strings.Contains(got, "level=error") {
		t.Fatalf("sanitized identifier was not retained for diagnosis: %q", got)
	}
	if strings.Contains(got, "must-not-be-logged") {
		t.Fatalf("event payload or reply destination reached the log: %q", got)
	}
}
