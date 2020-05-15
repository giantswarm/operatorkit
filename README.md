[![GoDoc](https://godoc.org/github.com/giantswarm/operatorkit?status.svg)](http://godoc.org/github.com/giantswarm/operatorkit) [![CircleCI](https://circleci.com/gh/giantswarm/operatorkit.svg?&style=shield&circle-token=5f7e69042df6538d1e9c7ef0dd1387ca4d7a0d55)](https://circleci.com/gh/giantswarm/operatorkit)

# operatorkit

Package `operatorkit` implements an opinionated framework for developing
[Kubernetes operators]. It emerged as we extracted common functionality from a
number of the operators we developed at [Giant Swarm][giantswarm]. The goal of
this library is to provide a common structure for operator projects and to
encapsulate best practices we learned while running operators in production.



## Features

- CRD primitives to reliably create, watch and delete custom resources, as well
  as any Kubernetes runtime object.
- Managing [finalizers][finalizers] on reconciled objects, making sure the code
  is executed at least once for each create/delete/update event.
- Guarantees to perform only one successful deletion event reconciliation to
  avoid unnecessary, possibly expensive interactions with third party systems.
- Resource wrapping to gain ability of composing resources like middlewares.
- Control flow primitives that allow cancellation and repetition of resource
  implementations.
- Independent packages. It is possible to use only certain parts of the library
  without being bound to all primitives it provides.
- Ability to change behaviour that is often specific to an organization like
  logging and error handling.



## Docs

- [Control Flow Primitives](docs/control_flow_primitives.md)
- [File Structure](docs/file_structure.md)
- [Keeping Reconciliation Loops Short](docs/keeping_reconciliation_loops_short.md)
- [Managing CR Status Sub Resources](docs/managing_cr_status_sub_resources.md)
- [Metrics Provider](docs/metrics_provider.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Using Finalizers](docs/using_finalizers.md)



## Integration Tests

You can simply create a [`kind`](https://github.com/kubernetes-sigs/kind/)
cluster to run the integration tests.

```
kind create cluster
```

The tests need to figure out how to connect to the Kubernetes cluster. Therefore
we need to set an environment variable pointing to your local kube config.

```
export E2E_KUBECONFIG=~/.kube/config
```

Now you can easily run the integration tests.

```
go test -v -tags=k8srequired ./integration/test/<test-name>
```

Once you did your testing you may want to delete your local test cluster again.

```
kind delete cluster
```



## Projects using operatorkit

[Giant Swarm operators] using `operatorkit`.



## Example

For a detailed state of art implementation, please see
[giantswarm/aws-operator](https://github.com/giantswarm/aws-operator).



## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.



## License

`operatorkit` is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.



[finalizers]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#finalizers
[giantswarm]: https://giantswarm.io
[Giant Swarm operators]: https://github.com/search?p=1&q=topic%3Aoperator+org%3Agiantswarm&type=Repositories
[Kubernetes operators]: https://coreos.com/operators
