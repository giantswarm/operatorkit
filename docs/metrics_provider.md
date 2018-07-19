# Metrics provider

The idea about metrics provider is to have a metrics driven approach of
verifying the operators functionality. An operator reconciles a system were
several resources exist. As a means of measurement and safety net it is a good
practise to implement metrics providers that emit metrics about the managed
system and its resources.

In [`aws-operator` we implement a Prometheus collector](https://github.com/giantswarm/aws-operator/blob/845afd245b0ace5cc6b37a1bd4f5da6c7e1d12d6/service/collector/collector.go)
to emit metrics about e.g. VPCs. These metrics are used for verifying each guest
cluster has a VPC assigned and no VPC is orphaned without having any guest
cluster assigned.

In [`operatorkit` we implement a Prometheus collector](https://github.com/giantswarm/operatorkit/blob/929bed01204f9a210f4589b4f3282e7f8028cce0/informer/collector.go)
to emit metrics about the creation and deletion timestamps of watched runtime
objects. These metrics are used for various purposes within our monitoring and
alerting system.
