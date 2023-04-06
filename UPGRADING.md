# Upgrade Guide

Go Sync follows semantic versioning where possible, so minor and revision versions should be upgradable with no breaking
changes; however, major version upgrades may need specific steps.

# v1.0.0

### Breaking changes

* Go Sync has been moved from `github.com/ovotech/go-sync` to `github.com/ovotech/go-sync/packages/gosync`
* Adapters have been moved from `github.com/ovotech/go-sync/adapters/*` to `github.com/ovotech/go-sync/packages/*`
* `Init` function signature has been modified to accept any number of optsFns.

### Notable changes

* `New` functions have been deprecated in favour of `Init` functions.
