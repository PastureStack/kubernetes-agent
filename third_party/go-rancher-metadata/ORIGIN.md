# Origin and local compatibility patch

This minimal module contains only the package used by Kubernetes Agent. The
source comes from `rancher/go-rancher-metadata` commit
`d2103caca5873119ff423d29cba09b4d03cd69b8` under Apache License 2.0.

The local compatibility patch changes the historical mixed-case Logrus import
to the canonical `github.com/sirupsen/logrus` module path. The metadata API,
schema, polling behaviour, and callback behaviour are unchanged.
