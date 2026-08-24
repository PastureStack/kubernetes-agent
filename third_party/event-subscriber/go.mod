module github.com/PastureStack/event-subscriber

go 1.27.0

require (
	github.com/PastureStack/control-api-client v0.0.0
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674
	github.com/sirupsen/logrus v1.10.1
)

require (
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/PastureStack/control-api-client => ../control-api-client
