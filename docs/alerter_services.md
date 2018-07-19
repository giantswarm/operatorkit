# Alerter Services

The idea about alerter services is to have a metrics driven approach of
verifying the operators functionality. An operator might reconcile a system in
which several resources should exist or maybe even not. As a means of
measurement and safety net it is a good practise to implement alerter services
that emit metrics about the managed system and its resources.

In [`aws-operator` we implement a Prometheus collector](https://github.com/giantswarm/aws-operator/blob/845afd245b0ace5cc6b37a1bd4f5da6c7e1d12d6/service/collector/collector.go)
to emit metrics about e.g. VPCs. These metrics are used for verifying each guest
cluster has a VPC assigned and no VPC is orphaned without having any guest
cluster assigned.

In [`operatorkit` we implement a Prometheus collector](https://github.com/giantswarm/operatorkit/blob/929bed01204f9a210f4589b4f3282e7f8028cce0/informer/collector.go)
to emit metrics about the creation and deletion timestamps of watched runtime
objects. These metrics are used for various purposes within our monitoring and
alerting system.
