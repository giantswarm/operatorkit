[![CircleCI](https://circleci.com/gh/giantswarm/operatorkit.svg?&style=shield&circle-token=5f7e69042df6538d1e9c7ef0dd1387ca4d7a0d55)](https://circleci.com/gh/giantswarm/operatorkit)

# operatorkit

The operatorkit is a library for creating [Kubernetes operators][operators]. It
emerged as extracted common functionality of number of the operators we
developed in Giant Swarm. The goal of the library is to provide common
structure of the operator projects and encapsulate best practices learned
during running operators in production.

## Features

- Reducing boilerplate code.
- CRD/TPR primitives allowing to reliably create, watch custom resources, and
  decode custom objects.
- Independent packages. It is possible to use only a part of the library. 
- Possibility to change behaviour that often is specific to an organization
  like logging and error handling.
- Minimal set of dependencies.

## Current Scope

Project is split into independent packages providing complementary
functionalities allowing to create production grade Kubernetes operators.

- client - provides unified way of creating Kubernetes clients required by
  other packages.
- crd/tpr - provide CRD/TPR primitives, allowing to reliably create custom
  resources, wait for their initialization, and generate endpoints URLs. 
- informer - provides well defined watching functionality for virtually any
  Kubernetes resource. It also provides a custom objects decoding functionality
  reducing error prone boilerplate.
- framework - provides a framework aiming to help writing reliable, robust
  reconciliation loops.

## Future Scope

- Custom object migration. With operators running in production handling many
  custom objects in multiple installations it is essential to be able to change
  the custom resource definition. This is something we are struggling with at
  the moment. We want to provide some generic way to deal with that problem
  nicely.
- Ensuring processing of custom object deletion. Freeing resources managed by
  the operator after custom resource deletion is crucial part of the operator.
  We think this is important that the operator processes custom resource
  deletion even when it was not running during the deletion. Even if the
  probability of such case is small we do not want risk orphaned resources. We
  want to use [finalizers][finalizers] for that. At the moment we are blocked
  with this bug https://github.com/kubernetes/kubernetes/issues/50528.
- Framework is still under heavy lifting. As we create new operators, and run
  current ones in different environments we discover new problems, learn how to
  deal with them and try to move that knowledge to the framework.

## Projects using operatorkit

- https://github.com/giantswarm/aws-operator
- https://github.com/giantswarm/azure-operator
- https://github.com/giantswarm/cert-operator
- https://github.com/giantswarm/flannel-operator
- https://github.com/giantswarm/ingress-operator
- https://github.com/giantswarm/kvm-operator
- more to come

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/cert-operator/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the contribution workflow as well as reporting bugs.

## License

operatorkit is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.

[finalizers]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#finalizers
[operators]: https://coreos.com/operators
