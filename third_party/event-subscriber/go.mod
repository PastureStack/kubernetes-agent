module github.com/PastureStack/event-subscriber

go 1.26.0

toolchain go1.27.3

require (
	github.com/PastureStack/control-api-client v0.0.0
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674
	github.com/sirupsen/logrus v1.10.2
)

require (
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
)

replace github.com/PastureStack/control-api-client => ../control-api-client
