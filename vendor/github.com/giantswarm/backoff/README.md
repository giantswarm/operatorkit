[![CircleCI](https://circleci.com/gh/giantswarm/backoff.svg?&style=shield&circle-token=776bf6423e66027e72228034b905a74cdbe871dc)](https://circleci.com/gh/giantswarm/backoff)

# backoff

Backoff is a library abstracting retry functionality using
https://godoc.org/github.com/cenkalti/backoff. We use this library to have a
unified interface and the opportunity to wrap and extend backoff implementations
as we need them.
