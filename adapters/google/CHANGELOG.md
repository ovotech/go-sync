# Changelog

All notable changes to this adapter will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v1.0.0

### Added

- `WithLogger` and `WithAdminService` ConfigFns.

### Changed

 - Unless `WithAdminService` is passed to the `InitFn`, the adapter will now use default GCP credentials.

### Removed

 - `New` functions have been removed. Use `InitFn` to instantiate a new adapter.
 - `GoogleAuthenticationMechanism` has been removed.
