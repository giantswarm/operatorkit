# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Dependencies Upgrade
  - giantswarm/architect@5.11.4
  - golang.org/x/net v0.34.0
  - golang.org/x/sync v0.10.0
  - golang.org/x/sys v0.29.0
  - golang.org/x/term v0.28.0
  - golang.org/x/text v0.21.0
  - google.golang.org/protobuf v1.36.3
  - github.com/giantswarm/micrologger v1.1.2
  - github.com/go-logr/logr v1.4.2
  - github.com/giantswarm/exporterkit v1.2.0
  - github.com/prometheus/client_golang v1.19.0
  - github.com/prometheus/common v0.48.0
  - golang.org/x/oauth2 v0.16.0
  - github.com/giantswarm/to v0.4.2

## [7.2.0] - 2023-11-09

### Changed

- Upgrade go to 1.21
- Upgrade k8s dependencies to 1.28.x

## [7.1.0] - 2022-07-18

### Changed

- Change Reconcile errors total to include controller name. 

## [7.0.1] - 2022-02-07

### Changed

- Export `GetFinalizerName` function.

### Added

- Controller boot log line.

## [7.0.0] - 2021-12-20

### Changed

- Upgrade github.com/giantswarm/backoff v0.2.0 to v1.0.0
- Upgrade github.com/giantswarm/exporterkit v0.2.1 to v1.0.0
- Upgrade github.com/giantswarm/microerror v0.3.0 to v0.4.0
- Upgrade github.com/giantswarm/micrologger v0.5.0 to v0.6.0
- Upgrade github.com/giantswarm/k8sclient v6.0.0 to v7.0.0

## [6.1.0] - 2021-12-17

### Fixed

- Update `k8sclient` to v6.1.0 with CRDClient that was removed in v6.0.0.

## [6.0.0] - 2021-11-12

### Added

- Add new Kubernetes API, `examples.testing.giantswarm.io`, for integration tests without importing `apiextensions`.

### Changed

- Update `k8sclient` to v6.0.0, `controller-runtime` to v0.8.3, and Kubernetes dependencies to v0.20.12.
- Adjust signature of `NewRuntimeObjectFunc` to return `client.Object` instead of `runtime.Object`.

## [5.0.0] - 2021-05-25

### Fixed

- Reduced memory usage of the timestamp collector using server-side filtering for watched resources.

### Changed

- Replaced `github.com/giantswarm/operatorkit/v4/pkg/controller/internal/selector.Selector` with
  `k8s.io/apimachinery/pkg/labels.Selector` in `controller.Config` to streamline the usage of server-side filtering.

## [4.3.1] - 2021-04-06

### Fixed

- Remove usage of self link for Kubernetes 1.20 support.

## [4.3.0] - 2021-03-16

### Added

- Add `Controller.Stop` method to stop controller reconciliation and metrics collection.

### Fixed

- Re-expose `controller.NewSelector()`.
- Only close manager channel once.
- Add `AllowedLabels` to configmap resource to prevent unnecessary updates.

## [4.2.0] - 2021-01-07

### Added

- Add `operatorkit_controller_last_reconciled` metrics.

### Fixed

- Add object context to pause annotation related logs.

## [4.1.0] - 2020-12-18

### Added

- Add `namespace` into controller setting.
- Add `SentryTags` Config field to allow setting custom tags to be sent alongside errors to `sentry.io`.

### Fixed

- Propagate label selectors to timestamp collector

## [4.0.0] - 2020-10-27

### Updated

- Update apiextensions to v3 and replace CAPI with Giant Swarm fork.
- Prepare module v4.

## [3.0.0] - 2020-10-23

### Removed

- Drop `controller.ProcessDelete` and `controller.ProcessUpdate`.

## [2.0.2] - 2020-10-15

### Fixed

- Fix pause logic being triggered by empty values on non-target annotation keys.

## [2.0.1] - 2020-09-24

### Updated

- Updated Kubernetes dependencies to v1.18.9.

## [2.0.0] - 2020-08-11

### Added

- Add configurable pause annotation support.

### Updated

- Updated backward incompatible Kubernetes dependencies to v1.18.5.

## [1.2.0] - 2020-06-29

### Added

- Support writing Kubernetes error events when returning microerror with kind and description.

## [1.1.0] - 2020-06-22

### Added

- Optional support for sentry error collector.

## [1.0.2] - 2020-06-18

### Fixed

- Use local `Selector` introduced in 1.0.1 type in `controller` package.



## [1.0.1] 2020-06-09

### Added

- Add local less rigid `Selector` interface type in `controller.Controller`.
  The new `Selector` interface is backward compatible with previously used
  `apiextensions` implementation. #407

### Fixed

- Fix the issue where `operatorkit_controller_creation_timestamp` and
  `operatorkit_controller_deletion_timestamp` metrics were not emitted for all
  the controllers.



## [1.0.0] 2020-05-18

### Added

- Add `handler.Interface` (not used yet).

### Changed

- Remove resource set concept.
- Remove CRD management. Due to versioning issues throughout the lifecycle of
  operators CRDs must be managed in a different way.
- Use v3 `k8sclient`.



## [0.2.1] 2020-05-06

### Added

- Add `cachekeycontext` package.



## [0.2.0] 2020-03-24

### Changed

- Switch from dep to Go modules.
- Use architect orb.



## [0.1.0] 2020-03-19

### Added

- First release.



[Unreleased]: https://github.com/giantswarm/operatorkit/compare/v7.2.0...HEAD
[7.2.0]: https://github.com/giantswarm/operatorkit/compare/v7.1.0...v7.2.0
[7.1.0]: https://github.com/giantswarm/operatorkit/compare/v7.0.1...v7.1.0
[7.0.1]: https://github.com/giantswarm/operatorkit/compare/v7.0.0...v7.0.1
[7.0.0]: https://github.com/giantswarm/operatorkit/compare/v6.1.0...v7.0.0
[6.1.0]: https://github.com/giantswarm/operatorkit/compare/v6.0.0...v6.1.0
[6.0.0]: https://github.com/giantswarm/operatorkit/compare/v5.0.0...v6.0.0
[5.0.0]: https://github.com/giantswarm/operatorkit/compare/v4.3.1...v5.0.0
[4.3.1]: https://github.com/giantswarm/operatorkit/compare/v4.3.0...v4.3.1
[4.3.0]: https://github.com/giantswarm/operatorkit/compare/v4.2.0...v4.3.0
[4.2.0]: https://github.com/giantswarm/operatorkit/compare/v4.1.0...v4.2.0
[4.1.0]: https://github.com/giantswarm/operatorkit/compare/v4.0.0...v4.1.0
[4.0.0]: https://github.com/giantswarm/operatorkit/compare/v3.0.0...v4.0.0
[3.0.0]: https://github.com/giantswarm/operatorkit/compare/v2.0.2...v3.0.0
[2.0.2]: https://github.com/giantswarm/operatorkit/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/giantswarm/operatorkit/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/giantswarm/operatorkit/compare/v1.2.0...v2.0.0
[1.2.0]: https://github.com/giantswarm/operatorkit/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/operatorkit/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/giantswarm/operatorkit/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/operatorkit/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/operatorkit/compare/v0.2.1...v1.0.0
[0.2.1]: https://github.com/giantswarm/operatorkit/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/operatorkit/compare/v0.1.0...v0.2.0

[0.1.0]: https://github.com/giantswarm/operatorkit/releases/tag/v0.1.0
