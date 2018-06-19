# Troubleshooting Tips

When developing operators, there might occasionally be moments when Kubernetes
APIs or apiserver don't return error message with full context or when the error
message can be a bit misleading. Purpose of this document is to list some common
pitfalls that might be hard to debug in some cases and where returned error
message is not perfectly clear.

## RBAC rules

When `informer.Watch()` is called, it might fail to unknown error:
`unknown (get clusternetworkconfigs.core.giantswarm.io)`. This error is not very
descriptive and OperatorKit cannot really improve here since the original error
is generic `StatusError`.

Root cause here is though missing RBAC rule: Operator must have watch permission
to used CRD in order to be able to operate on it.

## CRD registration

When creating new CRD for an operator and having everything working, but getting
an error from `streamwatcher.go` that states `unable to decode an event from the
watch stream` when an object for corresponding CR is created, one possible
mistake in this case is missing type registration when using generated CRD
clientset. When type registration is missing, `client-go` decoder cannot
recognize event and therefore fails with message like following:

`ERROR: logging before flag.Parse: E0619 13:28:42.006292       1 streamwatcher.go:109] Unable to decode an event from the watch stream: unable to decode watch event: no kind "ClusterNetworkConfig" is registered for version "core.giantswarm.io/v1alpha1"`

