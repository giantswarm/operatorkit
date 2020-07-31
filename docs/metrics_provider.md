# Metrics provider

The idea about metrics providers is to have a metrics driven approach of
verifying the operators functionality. An operator constantly reconciles the
current state of the system towards a more desired state eventually. As a means
of measurement and as a safety net, it is a good practise to implement metrics
providers that emit metrics about the managed system and its resources. The
collector implementation should be separate from the operator project itself. As
we learned over time, combining operators and collectors may cause more problems
than it is supposed to solve.



## Examples

In Azure we have the [`azure-operator`
project](https://github.com/giantswarm/azure-operator) and its [separate
associated collector](https://github.com/giantswarm/azure-collector) to emit
metrics about e.g. ARM Deployments. These metrics are used for alerting on on
deployments being in failed state or stuck in upgrading.

In [`operatorkit` we implement a Prometheus
collector](https://github.com/giantswarm/operatorkit/tree/master/informer/collector)
to emit metrics about the creation and deletion timestamps of watched runtime
objects. These metrics are used for various purposes within our monitoring and
alerting system. Note that in libraries we combine business logic and metrics
provider as long as we can ensure high coherence within the resulting data.



## Prometheus collectors

We make use of [Prometheus collector interface](https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector)
(which is part of the prometheus client library) to collect metrics from our
operators. This interface allows us to implement new collectors by implementing
the `Collect()` method, which is used to send prometheus metrics to a shared
channel. All the collectors [need to be
registered](https://godoc.org/github.com/prometheus/client_golang/prometheus#Registerer)
so that they are included in metrics collection.

Each time that the metrics endpoint is scraped by Prometheus, all the
registered collectors will have their `Collect()` method called, and metrics
will be exposed on the response. We use the [exporterkit](https://github.com/giantswarm/exporterkit)
library to make exposing metrics easier.

[Here you can find a commented example](https://godoc.org/github.com/prometheus/client_golang/prometheus#ex-Collector) on the prometheus client docs.
