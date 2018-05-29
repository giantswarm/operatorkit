# Control Flow Primitives

When using `operatorkit` you boot controllers which are configured with
informers. Informers are basically clients working against the Kubernetes API.
Runtime objects are observed and passed towards the controller. The controller
is also configured with a list of resource implementations. These are executed
in the order they are configured. There are a lot of cases in which you want to
control the execution and lifecycle of resource implementations depending on the
system an operator is dealing with. Control flow primitives also help then
[keeping reconciliation loops short](keeping_reconciliation_loops_short.md).

## Cancel Resources

#### Default Resources

Default resources implement `EnsureCreated` and `EnsureDeleted`. In order to
stop processing here it is good enough to simply `return` within the code. The
resources being configured after the canceled resource are then executed and the
canceled resource is executed again on the next reconciliation loop once the
resync period kicks in again.

#### CRUD Resources

CRUD resources implement a couple of steps to help with structuring more complex
use cases. These steps act like a framework and help navigating through the
stormy waters of reconciliation and its implications. In order to cancel
resources within one of its steps you can call
[`resourcecanceledcontext.SetCanceled(ctx)`](https://github.com/giantswarm/kvm-operator/blob/de7e109f4a652b785bbcf4214a1c8e028bf0eed4/service/controller/v12/resource/namespace/current.go#L66).
A convention for CRUD resource cancelation is to do it within `GetCurrentState`
when possible.

## Cancel Reconciliation

In order to cancel the whole reconciliation you can simply call
[`reconciliationcanceledcontext.SetCanceled(ctx)`](https://github.com/giantswarm/kvm-operator/blob/de7e109f4a652b785bbcf4214a1c8e028bf0eed4/service/controller/v12/resource/namespace/current.go#L45)
which will then stop executing all configured resources within the current
reconciliation loop. On the next reconciliation loop all resources are executed
again based on how they were configured. Note that cancelling resources on
delete events will cause the finalizer to be removed. This means the delete
event will not be replayed. When this behaviour is not desired check on how to
[repeat delete events](#repeat-delete-events).

## Repeat Delete Events

There are separate docs about [using finalizers](using_finalizers.md) which
describe a lot of the background. Thus we just touch the control flow aspects of
finalizers here briefly. In order to repeat delete events you can cause the
operatorkit controller to keep finalizers by calling
[`finalizerskeptcontext.SetKept(ctx)`](https://github.com/giantswarm/kvm-operator/blob/de7e109f4a652b785bbcf4214a1c8e028bf0eed4/service/controller/v12/resource/namespace/current.go#L67)
without using errors for control flow.
