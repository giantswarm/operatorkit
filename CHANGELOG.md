# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]



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



[Unreleased]: https://github.com/giantswarm/operatorkit/compare/v1.0.0...HEAD

[1.0.0]: https://github.com/giantswarm/operatorkit/compare/v0.2.1...1.0.0
[0.2.1]: https://github.com/giantswarm/operatorkit/compare/v0.2.0...0.2.1
[0.2.0]: https://github.com/giantswarm/operatorkit/compare/v0.1.0...0.2.0

[0.1.0]: https://github.com/giantswarm/operatorkit/releases/tag/v0.1.0
