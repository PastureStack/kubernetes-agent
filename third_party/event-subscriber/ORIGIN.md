# Origin and local compatibility patch

This minimal module contains only the packages used by Kubernetes Agent. The
source comes from `rancher/event-subscriber` commit
`dddf44b13e15ae38f8bb70a6de0a095a45da91ad` under Apache License 2.0.

The local compatibility patch changes the historical mixed-case Logrus import
to the canonical `github.com/sirupsen/logrus` module path. No event names,
payloads, locking behaviour, or retry behaviour are changed.
