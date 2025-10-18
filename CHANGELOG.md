# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v0.2.2

- Add standard Go project badges (Go Reference, CI, Go Report Card, codecov)
- Add Installation, API Documentation, and License sections to README
- Add Ginkgo v2 to tools.go for consistent test framework dependency tracking
- Improve README structure with horizontal rules for better visual separation

## v0.2.1

- Add GoDoc comments for all exported handler and factory functions
- Improve error handling consistency by replacing errors.Wrapf with errors.Wrap
- Add pkg test suite setup with Ginkgo v2 and Gomega

## v0.2.0

- Add build metadata support (Git version, commit hash, and build timestamp)
- Docker container now exposes build information via environment variables
- Build metadata automatically injected during Docker build process

## v0.1.0

- Initial release
