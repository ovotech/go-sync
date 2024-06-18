# Changelog

All notable changes to Go Sync will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v1.0.0

### Removed

 - `New` functions are no longer recommended to instantiate adapters, and have been replaced with `InitFn`.

### Changed

 - `InitFn` signature has been updated to allow adapters to take an arbitrary number of `ConfigFn`.
