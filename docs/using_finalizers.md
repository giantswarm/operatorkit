# Using Finalizers

[Kubernetes operators][operators] in general try to ensure a certain state. The
create state is already constantly ensured by the create and update events of
the informer. Deletion events are only executed once though. In order to also
replicate a continuous reconciliation for the delete state a Kubernetes concept
called [finalizers][finalizers] is used. They are actually only simple strings
being tracked in the object metadata of any Kubernetes runtime object. Once a
runtime object is deleted the Kubernetes API sees the list of finalizers applied
is not empty and defers deletion until no finalizers are tracked anymore. Users
are responsible to remove finalizers from their reconciled runtime objects.

The good news is `operatorkit` does all this for you already. As soon as an
operator using operatorkit reconciles a runtime object the object's metadata is
updated by adding an operator specific finalizer. This might look like the
following.

```
apiVersion: provider.giantswarm.io/v1alpha1
kind: KVMConfig
metadata:
...
  finalizers:
  - operatorkit.giantswarm.io/kvm-operator
  - operatorkit.giantswarm.io/kvm-operator-drainer
...
```

## Configuration

#### RESTClient

The first important aspect of configuring operatorkit controllers is the
`RESTClient` dependency. [For reconciled objects of custom generated clientsets](https://github.com/giantswarm/kvm-operator/blob/4794d415a21c6e6d0e2ba7fb6f1ef1591101f4b0/service/controller/cluster.go#L337-L346)
the type specific REST client has to be configured.

```
c := controller.Config{
  ...
	RESTClient: config.G8sClient.ProviderV1alpha1().RESTClient(),
  ...
}
```

[For reconciled objects of custom standard clientsets](https://github.com/giantswarm/kvm-operator/blob/4794d415a21c6e6d0e2ba7fb6f1ef1591101f4b0/service/controller/drainer.go#L69-L76)
the normal core REST client has to be configured.

```
c := controller.Config{
  ...
	RESTClient: config.K8sClient.CoreV1().RESTClient(),
  ...
}
```

#### Name

The second important aspect of configuring operatorkit controllers is the `Name`
setting. In case operators boot multiple controllers for different
reconciliation loops and purposes the configured name is used to manage
finalizers.

[The name setting for the `cluster` controller in the kvm-operator looks like this.](https://github.com/giantswarm/kvm-operator/blob/4794d415a21c6e6d0e2ba7fb6f1ef1591101f4b0/service/controller/cluster.go#L337-L346)

```
c := controller.Config{
  ...
	Name: config.ProjectName,
  ...
}
```

[The name setting for the `drainer` controller in the kvm-operator looks like this.](https://github.com/giantswarm/kvm-operator/blob/4794d415a21c6e6d0e2ba7fb6f1ef1591101f4b0/service/controller/drainer.go#L69-L76)

```
c := controller.Config{
  ...
	Name: config.ProjectName + "-drainer",
  ...
}
```

## Control Flow

The default behaviour for delete events to be replayed is to return an error in
the resource implementations. This behaviour makes sense for unexpected cases in
which the operator should just retry the deletion because of a former faulty
behaviour of the system. Since errors are not a good way to manage control flow
[replaying delete events can be requested by an operatorkit primitive](https://github.com/giantswarm/kvm-operator/blob/4794d415a21c6e6d0e2ba7fb6f1ef1591101f4b0/service/controller/v12/resource/endpoint/current.go#L42).

```
finalizerskeptcontext.SetKept(ctx)
```

[finalizers]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#finalizers
[operators]: https://coreos.com/operators
