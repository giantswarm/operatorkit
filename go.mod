module github.com/giantswarm/operatorkit

go 1.14

require (
	github.com/getsentry/sentry-go v0.7.0
	github.com/giantswarm/apiextensions v0.4.18
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/exporterkit v0.2.0
	github.com/giantswarm/k8sclient/v3 v3.1.2
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/to v0.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.2
)
