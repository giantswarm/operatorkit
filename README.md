[![CircleCI](https://circleci.com/gh/giantswarm/operatorkit.svg?&style=shield&circle-token=5f7e69042df6538d1e9c7ef0dd1387ca4d7a0d55)](https://circleci.com/gh/giantswarm/operatorkit)

# operatorkit

operatorkit package is a library for creating [Kubernetes
operators][operators]. It emerged as we extracted common functionality from
number of the operators we developed at Giant Swarm. The goal of the library is
to provide a common structure of operator projects and to encapsulate best
practices learned while running operators in production.

## Features

- CRD primitives to reliably create, watch and delete custom resources.
- Managing finalizers on reconciled objects, making sure the code is executed
  at least once for each delete event.
- Independent packages. It is possible to use only certain parts of the
  library.
- Possible to change behaviour that often is specific to an organization
  like logging and error handling.
- Minimal set of dependencies.

## Current Scope

The project is split into independent packages providing complementary
functionality, making it easier create production grade Kubernetes operators.

- client - provides a unified way of creating Kubernetes clients required by
  other packages.
- informer - provides well defined watching functionality for virtually any
  Kubernetes resource. The informer is deterministic, meaning it does not
  dispatch events twice after the resync period, which saves some cycles. It
  also features rate limiting of the event dispatching. It also provides
  functionality for decoding custom objects, reducing error prone boilerplate
  code.
- controller - provides a framework aiming to help writing reliable, robust
  controllers performing reconciliation loops. The heart of the controller is
  a Resource interface. The reconciliation primitive allowing split the
  reconciliation into smaller bits. Controller manages [finalizers][finalizers]
  on reconciled objects, making sure all resources are executed at least once
  during the deletion.

## Projects using operatorkit

- https://github.com/giantswarm/aws-operator
- https://github.com/giantswarm/azure-operator
- https://github.com/giantswarm/cert-operator
- https://github.com/giantswarm/cluster-operator
- https://github.com/giantswarm/chart-operator
- https://github.com/giantswarm/flannel-operator
- https://github.com/giantswarm/ingress-operator
- https://github.com/giantswarm/kvm-operator
- https://github.com/giantswarm/node-operator
- https://github.com/giantswarm/prometheus-config-controller
- more to come

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/cert-operator/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.

## License

operatorkit is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.

[finalizers]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#finalizers
[operators]: https://coreos.com/operators
