package events

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxLogValueRunes = 512

// safeLogValue keeps untrusted event and response values on one printable log
// line. It is for operational logging only; protocol replies retain their
// original values.
func safeLogValue(value interface{}) string {
	runes := []rune(fmt.Sprint(value))
	var builder strings.Builder
	limit := len(runes)
	if limit > maxLogValueRunes {
		limit = maxLogValueRunes
	}
	for _, r := range runes[:limit] {
		switch {
		case r == '\r':
			builder.WriteString(`\r`)
		case r == '\n':
			builder.WriteString(`\n`)
		case r == '\t':
			builder.WriteString(`\t`)
		case r == utf8.RuneError || unicode.IsControl(r):
			builder.WriteRune('�')
		default:
			builder.WriteRune(r)
		}
	}
	if len(runes) > limit {
		builder.WriteString("…")
	}
	sanitized := strings.ReplaceAll(builder.String(), "\r", `\r`)
	sanitized = strings.ReplaceAll(sanitized, "\n", `\n`)
	return sanitized
}

func safeEventFields(event *Event) map[string]interface{} {
	if event == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"eventName":    safeLogValue(event.Name),
		"eventId":      safeLogValue(event.ID),
		"resourceId":   safeLogValue(event.ResourceID),
		"resourceType": safeLogValue(event.ResourceType),
	}
}
