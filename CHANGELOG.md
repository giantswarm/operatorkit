# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]

## Updated

- Updated Kubernetes dependencies to v1.18.9.

## [2.0.0] - 2020-08-11

## Added

- Add configurable pause annotation support.

## Updated

- Updated backward incompatible Kubernetes dependencies to v1.18.5.

## [1.2.0] - 2020-06-29

## Added

- Support writing Kubernetes error events when returning microerror with kind and description.

## [1.1.0] - 2020-06-22

## Added

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



[Unreleased]: https://github.com/giantswarm/operatorkit/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/giantswarm/operatorkit/compare/v1.2.0...v2.0.0
[1.2.0]: https://github.com/giantswarm/operatorkit/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/operatorkit/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/giantswarm/operatorkit/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/operatorkit/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/operatorkit/compare/v0.2.1...v1.0.0
[0.2.1]: https://github.com/giantswarm/operatorkit/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/operatorkit/compare/v0.1.0...v0.2.0

[0.1.0]: https://github.com/giantswarm/operatorkit/releases/tag/v0.1.0
