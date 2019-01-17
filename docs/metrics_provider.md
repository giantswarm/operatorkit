# Metrics provider

The idea about metrics provider is to have a metrics driven approach of
verifying the operators functionality. An operator reconciles a system were
several resources exist. As a means of measurement and safety net it is a good
practise to implement metrics providers that emit metrics about the managed
system and its resources.

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
