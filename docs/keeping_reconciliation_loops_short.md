# Keeping Reconciliation Loops Short

Resource implementations should not block more than necessary and execute as
fast as possible. This is to not block reconciliation of other runtime objects,
because their reconciliation is synchronous. For instance the reconciliation of
a certain cloud provider resource taking half an hour prevents managing other
resources during that time. So keeping reconciliation loops short makes the
operator more efficient in managing multiple resources and the nature of
operators by design allow to simply catch up with any state on the next resync
period.

As an example, resource implementations aim to ensure a certain state of the
systems they manage. Very common states in that regard are `created` and
`deleted`. In order to reach these states it might take time for e.g. the cloud
provider the resource implementation relies on to finish certain state
transitions. In our example the resource implementation then does two things.

- Ask the cloud provider to create a certain resource.
- Wait for the created resource to actually be created.

The first step is often synchronous and relatively fast. The second step can
take a very long time though. Here the implementation should rather check for
the current state of the managed system instead of blocking and waiting forever.
See also real examples in the wild.

- [aws-operator creating AWS CloudFormation stacks](https://github.com/giantswarm/aws-operator/blob/ee3ece0107442b60eb2755f669bf2945a5ab05f5/service/controller/v11/resource/cloudformation/create.go#L31-L57)
