# Changelog

All notable changes to `zazu-go` are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1]

Version alignment: the whole SDK family now releases in lockstep with zazu-ruby. No functional changes since [0.1.0].

## [0.1.0]

Initial release.

### Added

- `zazu.Client` built on `net/http` (functional options, context-first API)
- Services: `Accounts`, `Beneficiaries`, `CheckoutSessions`, `Customers`, `Entity`, `Invoices`, `PaymentLinks`, `TransferDrafts`, `WebhookEndpoints`
- Cursor-based `Page` with `Next(ctx)` (max 100 records per page)
- `*zazu.Error` mirroring the shared SDK error taxonomy
- Cassette-replay test harness driven by the Ruby SDK's release tarball
