[![CircleCI](https://circleci.com/gh/giantswarm/operatorkit.svg?&style=shield&circle-token=5f7e69042df6538d1e9c7ef0dd1387ca4d7a0d55)](https://circleci.com/gh/giantswarm/operatorkit)

# operatorkit

Package `operatorkit` implements an opinionated framework for developing
[Kubernetes operators][operators]. It emerged as we extracted common
functionality from a number of the operators we developed at [Giant
Swarm][giantswarm]. The goal of this library is to provide a common structure
for operator projects and to encapsulate best practices we learned while running
operators in production.

## Features

- CRD primitives to reliably create, watch and delete custom resources, as well
  as any Kubernetes runtime object.
- Managing [finalizers][finalizers] on reconciled objects, making sure the code
  is executed at least once for each create/delete/update event.
- Guarantees to perform only one successful deletion event reconciliation to
  avoid unnecessary, possibly expensive interactions with third party systems.
- A deterministic informer implementation that guarantees the expected behaviour
  of configured resync periods and rate waits.
- Convenient client library helpers for simpler client creation.
- Resource wrapping to gain ability of composing resources like middlewares.
- Control flow primitives that allow cancellation and repetition of resource
  implementations.
- Independent packages. It is possible to use only certain parts of the library
  without being bound to all primitives it provides.
- Ability to change behaviour that is often specific to an organization like
  logging and error handling.
- Minimal set of dependencies.

## Docs

- [Control Flow Primitives](docs/control_flow_primitives.md)
- [File Structure](docs/file_structure.md)
- [Keeping Reconciliation Loops Short](docs/keeping_reconciliation_loops_short.md)
- [Managing CR Status Sub Resources](docs/managing_cr_status_sub_resources.md)
- [Metrics Provider](docs/metrics_provider.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Using Finalizers](docs/using_finalizers.md)

## Current Scope

The project is split into independent packages providing complementary
functionality, making it easier to create production grade [Kubernetes
operators][operators].

- `client`: provides a unified way of creating Kubernetes clients required by
  other packages.
- `informer`: provides well defined watching functionality for virtually any
  Kubernetes resource. The informer is deterministic, meaning it does not
  dispatch events twice after the resync period, which saves some cycles. It
  also features rate limiting of the event dispatching. It also provides
  functionality for decoding custom objects, reducing error prone boilerplate
  code.
- `controller`: provides a framework aiming to help writing reliable, robust
  controllers performing reconciliation loops. The heart of the controller is a
  Resource interface. The reconciliation primitive allows splitting the
  reconciliation into smaller parts. Controller manages [finalizers][finalizers]
  on reconciled objects, making sure all resources are executed at least once
  during the deletion.

## Projects using operatorkit

- https://github.com/giantswarm/app-operator
- https://github.com/giantswarm/aws-operator
- https://github.com/giantswarm/azure-operator
- https://github.com/giantswarm/cert-operator
- https://github.com/giantswarm/cluster-operator
- https://github.com/giantswarm/chart-operator
- https://github.com/giantswarm/flannel-operator
- https://github.com/giantswarm/kvm-operator
- https://github.com/giantswarm/node-operator
- https://github.com/giantswarm/prometheus-config-controller
- https://github.com/giantswarm/release-operator
- more to come

## Example
An example `memcached-operator` is provided in the repository [giantswarm/memchached-operator](https://github.com/giantswarm/memcached-operator). 
This is a simplified example to illustrate certain primitives. For the detail implementation in the state of art, please see [giantswarm/aws-operator](https://github.com/giantswarm/aws-operator).

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/cert-operator/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.

## License

`operatorkit` is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.

[finalizers]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#finalizers
[giantswarm]: https://giantswarm.io
[operators]: https://coreos.com/operators
