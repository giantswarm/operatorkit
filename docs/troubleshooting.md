# Troubleshooting

When developing operators, there might occasionally be moments when Kubernetes
APIs or apiserver don't return error message with full context or when the error
message can be a bit misleading. Purpose of this document is to list some common
pitfalls that might be hard to debug in some cases and where returned error
message is not perfectly clear.

## RBAC rules

When `informer.Watch()` is called, it might fail with an "unknown error". The
original error is
[StatusError](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/errors/errors.go#L39)
with `ErrStatus.Reason`:
[Forbidden](https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/types.go#L569).
As such when logged it's not very descriptive with its `unknown` message, but
there is a test for this error in OperatorKit that should catch and produce an
informative log message referring to possible root cause which is a missing RBAC
authorization rule. The operator must have watch permission for the used CRD in
order to be able to operate on it.

Raw `StatusError` log message:

```
unknown (get clusternetworkconfigs.core.giantswarm.io)
```

An improved OperatorKit log message (when CRD name is
`clusternetworkconfigs.core.giantswarm.io`):

```
controller might be missing RBAC rule for clusternetworkconfigs.core.giantswarm.io CRD
```

## CRD registration

When creating new CRD for an operator and having everything working, but getting
an error from `streamwatcher.go` that states `unable to decode an event from the
watch stream` when an object for corresponding CR is created, one possible
mistake in this case is missing type registration when using generated CRD
clientset. Here is an example for Giant Swarm [apiextensions core type
registration](https://github.com/giantswarm/apiextensions/blob/master/pkg/apis/core/v1alpha1/register.go#L17).
When type registration is missing, `client-go` decoder cannot recognize event
and therefore fails with message like following:

```
ERROR: logging before flag.Parse: E0619 13:28:42.006292       1 streamwatcher.go:109] Unable to decode an event from the watch stream: unable to decode watch event: no kind "ClusterNetworkConfig" is registered for version "core.giantswarm.io/v1alpha1"
```
