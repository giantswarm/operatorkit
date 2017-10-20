[![CircleCI](https://circleci.com/gh/giantswarm/operatorkit.svg?&style=shield&circle-token=5f7e69042df6538d1e9c7ef0dd1387ca4d7a0d55)](https://circleci.com/gh/giantswarm/operatorkit)

# operatorkit

operatorkit package is a library for creating [Kubernetes
operators][operators]. It emerged as we extracted common functionality from
number of the operators we developed at Giant Swarm. The goal of the library is
to provide a common structure of operator projects and to encapsulate best
practices learned while running operators in production.

## Features

- CRD/TPR primitives to reliably create, watch and delete custom resources.
- Provides custom object decoder reducing boilerplate code to the minimum.
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
- crd/tpr - provides CRD/TPR primitives, allowing to reliably create custom
  resources, wait for their initialization, and generate endpoint URLs. 
- informer - provides well defined watching functionality for virtually any
  Kubernetes resource. The informer is deterministic, meaning it does not
  dispatch events twice after the resync period, which saves some cycles. It
  also features rate limiting of the event dispatching. It also provides
  functionality for decoding custom objects, reducing error prone boilerplate
  code.
- framework - provides a framework aiming to help writing reliable, robust
  reconciliation loops. The heart of the framework is a Resouce interface
  guiding the user step by step towards a proper reconciliation that drives the
  current state towards the desired state. The Resource is designed to be
  composable. The framework provides useful Resource wrappers, making it easy
  to retry on failures, emit metrics and log consistent diagnostic messages.

## Future Scope

- Custom object migration. With operators running in production handling many
  custom objects in multiple installations it is essential to be able to change
  the custom resource definition. This is something we are struggling with at
  the moment. We want to provide some generic way to deal with that problem
  nicely.
- Ensuring processing of custom object deletion. Freeing resources managed by
  the operator after custom resource deletion is a crucial part of the
  operator. We think it is important that the operator processes custom
  resource deletion even when it was not running during the deletion. E.g.
  custom object deleted during operator redepolyment. Even if the probability
  of such a case is small we do not want to risk orphaned resources. We plan to
  use [finalizers][finalizers] for that.
- Framework is still under heavy development. As we create new operators, and
  run them in different environments we discover new problems, learn how to
  deal with them and try to move that knowledge to the framework.

## Projects using operatorkit

- https://github.com/giantswarm/aws-operator (not using the framwork yet, but
  we work hard on that)
- https://github.com/giantswarm/azure-operator
- https://github.com/giantswarm/cert-operator
- https://github.com/giantswarm/draughtsman-operator (currently not used in
  production because of goroutine leaks in Kubernetes apiservers)
- https://github.com/giantswarm/endpoint-operator (WIP)
- https://github.com/giantswarm/flannel-operator
- https://github.com/giantswarm/ingress-operator
- https://github.com/giantswarm/kvm-operator
- https://github.com/giantswarm/prometheus-config-controller (WIP)
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
