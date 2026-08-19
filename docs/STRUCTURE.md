# Project Structure

This document describes the purpose of each folder in the repository. It complements the [README](../README.md), which covers installation and usage.

## Overview

This repository implements **SECA Conformance**, a Go-based CLI tool (built as a `go test` binary driven by [Cobra](https://github.com/spf13/cobra)) that validates Cloud Service Provider (CSP) implementations against the [SECA API specification](https://spec.secapi.cloud). Tests are organized into suites per API domain (Compute, Network, Storage, Authorization, Workspace, Region, Usage) and can run either against a real provider or against a local [WireMock](https://wiremock.org/) mock server.

## Top-level layout

```
.
├── cmd/                  CLI entry point and generated report artifacts
├── config/               Editor/IDE configuration
├── docs/                 Project documentation (this folder)
├── internal/             Private application logic (test framework, suites, mocking)
├── pkg/                  Reusable library code (builders, generators, wrappers)
├── secatest/             Empty/legacy placeholder package
├── tools/                Pinned developer tooling (isolated Go module)
├── wiremock/             Local WireMock mock-server setup and static stub mappings
├── dist/                 Compiled binary output (generated, git-ignored)
├── reports/              Allure test result output (generated, git-ignored)
├── Makefile              Build, lint, mock, run and report automation
├── go.mod / go.sum       Main Go module definition
└── README.md             Installation, configuration and usage guide
```

## `cmd/`

The CLI entry point, implemented as a Go test binary (`package main` with `TestMain`), wired up with Cobra.

- **`cmd/conformance`** — Defines the CLI subcommands:
  - `run` — executes the conformance suites (`go test`) and produces an Allure report.
  - `list` — lists the available test scenarios.
  - `report` — serves the Allure report in a browser (`allure serve`).
  - `summary` — prints/writes a JSON or text summary of Allure results.

  Per-domain `*_test.go` files (`compute_test.go`, `network_test.go`, etc.) wire up and conditionally run each suite (e.g. `TestComputeV1Suites` runs Provider Lifecycle, Provider Queries, Instance Constraints, and Instance Error suites).

- **`cmd/reports/results`** — Sample/generated Allure result JSON files from a prior test run; example output rather than source code.

## `internal/`

Private application logic that implements the test framework itself.

- **`internal/conformance/config`** — Global runtime configuration. `parameters.go` defines `ParametersHolder` (provider URLs, client auth/tenant/region, scenario filters, mock settings, retry settings) populated from CLI flags; `clients.go` builds SDK API clients from those parameters.

- **`internal/conformance/params`** — Domain-specific parameter/config structs used to configure individual test suites/scenarios.

- **`internal/conformance/steps`** — Reusable, Gherkin-style test step implementations per resource type (`compute_v1.go`, `network_v1.go`, `storage_v1.go`, `authorization_v1.go`, `workspace_v1.go`, `region_v1.go`), plus generic CRUD/watch step helpers (create/update, get, list, delete, watch, action), assertions, and API call wrappers.

- **`internal/conformance/suites`** — Base suite framework: shared assertion helpers and suite construction logic (e.g. `CreateRegionalTestSuite`), built on the [`ozontech/allure-go`](https://github.com/ozontech/allure-go) suite framework. Subfolders hold the actual conformance test suites per SECA API domain:
  - `authorization/`, `compute/`, `network/`, `region/`, `storage/`, `usage/`, `workspace/` — e.g. `compute/provider_lifecycle_v1.go`, `compute/instance_error_v1.go`.

- **`internal/constants`** — Shared constants: suite/scenario names (used by the `list` command), HTTP condition/operation constants, and general test constants.

- **`internal/mock`** — Core WireMock client integration wrapping [`wiremock/go-wiremock`](https://github.com/wiremock/go-wiremock), including mock parameters and shared constants.
  - **`internal/mock/scenarios`** — Defines a `Scenario` type that configures/finishes/resets WireMock stub scenarios per test. Subfolders (`authorization/`, `clients/`, `compute/`, `network/`, `region/`, `storage/`, `usage/`, `workspace/`) contain scenario definitions per API domain that script mock server behavior for each test case.
  - **`internal/mock/stubs`** — Builds and registers individual WireMock stub rules (request matcher + response) per resource type, via a shared `Configurator`.

- **`internal/report`** — Allure result parsing/aggregation. Reads raw Allure result files, builds a summary (totals + per-scenario results), and renders it as human-readable text; used by the `summary` and `report` CLI commands.

## `pkg/`

Public/reusable library code that wraps the SECA `go-sdk`, kept separate from `internal` since it is intended to be reusable outside this repo.

- **`pkg/builders`** — Generic, generics-based fluent builders for constructing SECA resource payloads per domain (compute, network, storage, authorization, workspace, region, metadata) plus a validator.

- **`pkg/generators`** — Test-data generators: CIDR/subnet/IP generation, resource name generation, reference/URL builders, and sample resource generation.

- **`pkg/wrappers`** — Typed wrapper interfaces (`ResourceWrapper`/`GlobalResourceWrapper` generics) around SDK resource types per domain, simplifying access to metadata/spec/status in tests.

## `secatest/`

Currently empty (including `secatest/generators`). Likely a placeholder/legacy package name from before code was consolidated under `pkg/generators`.

## `tools/`

An isolated Go module (own `go.mod`/`go.sum`) used only to pin developer-tool dependencies (e.g. `golangci-lint`, `gofumpt`), invoked by the Makefile via `-modfile=./tools/go.mod` so these tools don't pollute the main module's dependency graph.

## `wiremock/`

Local mock server setup used to run conformance tests without a real CSP backend.

- **`wiremock/docker-compose.yml`** — Runs the `wiremock/wiremock` Docker image on port 8080, mounting `config/mappings`.
- **`wiremock/config/mappings/{compute,network,storage}`** — Static WireMock stub-mapping JSON files (request matcher + templated JSON response, e.g. `get-sku.json`, `list-skus.json`) used as baseline/default mocks for these API domains. Dynamic, per-scenario stubs are added at runtime via `internal/mock/stubs`.

## `config/`

Editor/IDE configuration.

- **`config/vscode`** — `launch.json` holds VS Code debug launch configurations (e.g. "Run All Suites") that invoke the `cmd/conformance` test binary with typical CLI flags against a local WireMock instance.

## `docs/`

Project documentation, including this file and supporting images (e.g. `report-viewer.png`, referenced from the README).

## Generated / git-ignored directories

These are build or test-run artifacts, not source, and are excluded from version control via `.gitignore`:

- **`dist/`** — Compiled `secatest` binary produced by `make` (`go test -c -o dist/secatest ./cmd/conformance`).
- **`reports/`** — Allure test result output from local runs.
- **`allure-report_*`** — Generated, timestamped Allure HTML report bundles.
- **`.vscode/`** — Local, per-developer VS Code settings (distinct from the versioned `config/vscode`).

## Root files

- **`Makefile`** — Build/dev automation: compiling the binary, managing the WireMock Docker Compose stack (`mock-run`/`mock-start`/`mock-stop`), running the suite (`run`/`test`), producing Allure summaries/reports (`summary`/`report`), formatting/linting (`fmt`/`lint`/`dupl`), and cleanup (`clean`).
- **`go.mod` / `go.sum`** — Main Go module (`github.com/eu-sovereign-cloud/conformance`), depending on `go-sdk` (the SECA API client), `cobra` (CLI), `ozontech/allure-go` (test reporting/suite framework), `wiremock/go-wiremock` (mock client), and `go-playground/validator`.
- **`.golangci.yml` / `.golangci-dupl.yml`** — Linting configuration, including duplicate-code detection.
