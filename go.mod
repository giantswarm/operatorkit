module github.com/giantswarm/operatorkit/v2

go 1.14

require (
	github.com/getsentry/sentry-go v0.7.0
	github.com/giantswarm/apiextensions/v2 v2.2.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/exporterkit v0.2.0
	github.com/giantswarm/k8sclient/v4 v4.0.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/to v0.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
	sigs.k8s.io/controller-runtime v0.6.1
)
