# Metrics provider

The idea about metrics provider is to have a metrics driven approach of
verifying the operators functionality. An operator reconciles a system where
several resources exist. As a means of measurement and as a safety net, it is a good
practise to implement metrics providers that emit metrics about the managed
system and its resources.

## Examples

In [`aws-operator` we implement a Prometheus collector](https://github.com/giantswarm/aws-operator/tree/master/service/collector)
to emit metrics about e.g. VPCs. These metrics are used for verifying each
tenant cluster has a VPC assigned and no VPC is orphaned without having any
tenant cluster assigned.

In [`azure-operator` we implement a Prometheus collector](https://github.com/giantswarm/azure-operator/tree/master/service/collector)
to emit metrics about e.g. ARM Deployments. These metrics are used for alerting
on deployments being in failed state or stuck in upgrading.

In [`operatorkit` we implement a Prometheus collector](https://github.com/giantswarm/operatorkit/tree/master/informer/collector)
to emit metrics about the creation and deletion timestamps of watched runtime
objects. These metrics are used for various purposes within our monitoring and
alerting system.


## Prometheus collectors
We make use of [Prometheus collector interface](https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector) 
(which is part of the prometheus client library) to collect metrics from our 
operators.
This interface allows us to implement new collectors by implementing the 
`Collect()` method, which is used to send prometheus metrics to a shared channel.
All the collectors [need to be registered](https://godoc.org/github.com/prometheus/client_golang/prometheus#Registerer) 
so that they are included in metrics collection.

Each time that the metrics endpoint is scraped by Prometheus, all the 
registered collectors will have their `Collect()` method called, and metrics 
will be exposed on the response. We use the [exporterkit](https://github.com/giantswarm/exporterkit) 
library to make exposing metrics easier. 

[Here you can find a commented example](https://godoc.org/github.com/prometheus/client_golang/prometheus#ex-Collector) on the prometheus client docs.
