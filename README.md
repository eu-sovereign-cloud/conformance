# SECA Conformance Tests

A comprehensive conformance testing tool to validate a Cloud Service Provider (CSP) of the [SECA API specification](https://spec.secapi.cloud).

## Overview

SECA Conformance ensure that CSP implementations comply with the standardized API specification for sovereign cloud services. This tool validates API endpoints, resource lifecycle management, and compliance with SECA standards across multiple cloud providers.

## Requirements

- POSIX compatible environment;
- Make Build Tool;
- Git SCM;
- [Go](https://go.dev/doc/install) 1.24 or higher;
- [Allure Report](https://allurereport.org/docs/install/).

### Installation

```bash
git clone https://github.com/eu-sovereign-cloud/conformance
cd conformance
make
make install
```

## Configuration

The following configurations are required to run the tool. These configurations can be set as command line parameters or environment variables:


| Parameter                     | Variable                    | Description                                                                                                               | Required |
|-------------------------------|-----------------------------|---------------------------------------------------------------------------------------------------------------------------|----------|
| `--provider.region.v1`        | `PROVIDER_REGION_V1`        | URL of a Region V1 provider API implementation                                                                            | True     |
| `--provider.authorization.v1` | `PROVIDER_AUTHORIZATION_V1` | URL of a Authorization V1 provider API implementation                                                                     | True     |
| `--client.auth.token`         | `CLIENT_AUTH_TOKEN`         | Valid JWT token to access the CSP API's                                                                                   | True     |
| `--client.tenant`             | `CLIENT_TENANT`             | Name of the Tenant used in the tests                                                                                      | True     |
| `--client.region`             | `CLIENT_REGION`             | Name of the Region used in the regional tests                                                                             | True     |
| `--scenarios.filter`          | `SCENARIOS_FILTER`          | Regular expression to filter scenarios to run. To know the available scenarios run the [list](#listing-scenarios) command | False    |
| `--scenarios.users`           | `SCENARIOS_USERS`           | Comma-separated list of valid CSP users                                                                                   | True     |
| `--scenarios.cidr`            | `SCENARIOS_CIDR`            | CIDR range available in the CSP to create network resources                                                               | True     |
| `--scenarios.public.ips`      | `SCENARIOS_PUBLIC_IPS`      | Public IPs range, in CIDR format, to create CSP public IP's                                                               | True     |
| `--report.results.path`       | `REPORT_RESULTS_PATH`       | Path to store the tests result reports                                                                                    | False    |

## Running

To run the conformance tests, set the [configuration](#configuration) variables and run the following command:
```bash
secatest run \
  --provider.region.v1=#REGION_API \
  --provider.authorization.v1=#AUTHORIZATION_API \
  --client.auth.token=#TOKEN \
  --client.region=#REGION \
  --client.tenant=#TENANT \
  --scenarios.users=#USERS \
  --scenarios.cidr=#CIDR \
  --scenarios.public.ips=#PUBLIC_IPS
```

## Viewing Result

To see the the result report run the following command:
```bash
secatest report
```

Your default browser will be opened, with the Allure Report viewer:

![Viewer](docs/report-viewer.png)

## Listing Scenarios

To see the list of available test scenarios run the following command:
```bash
secatest list
```
