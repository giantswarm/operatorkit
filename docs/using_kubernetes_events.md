# Using Kubernetes Events

Kubernetes events are objects that show you what is happening inside a cluster,
such as what decisions were made by the controllers, in case an error occurs, or
other messages. Most of the core components already create events through the
API Server.

By default, events are retained for one hour but can be changed globally via a
flag in the kube-apiserver `--event-ttl`. Events are visible by either using
`describe` on Kubernetes objects or `get events` to see what events have been
occured.



### How can events be implemented by an operator?

Operatorkit with version
[v1.2.0](https://github.com/giantswarm/operatorkit/releases/tag/v1.2.0) or
higher supports Kubernetes events.

All it takes is to specify a [microerror with a `Kind` and a `Desc`
(description)](https://pkg.go.dev/github.com/giantswarm/microerror?tab=doc#Error).
Both fields have to be set, otherwise the event recorder in operatorkit won't
create an event.

`Kind` is set as the event reason which is a short, machine understandable
string that gives the reason for the transition into the object's current
status.

`Desc` is set as the event message, a human-readable description of the status
of this operation.

When a microerror occurs, the event recorder emits a **event with type warning**
on given Kubernetes objects for each reconciliation loop.

When similar events are popping up they are counted and won't appear in a new
line in events, e.g.:

```yaml
Events:
  Type     Reason      Age              From           Message
  ----     ------      ----             ----           -------
  Warning  EventError  9s (x5 over 9s)  test-operator  Error of an event
```
