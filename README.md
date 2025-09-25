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


| Parameter                     | Variable                    | Description                                                 |
|-------------------------------|-----------------------------|-------------------------------------------------------------|
| `--provider.region.v1`        | `PROVIDER_REGION_V1`        | URL of a Region V1 provider API implementation              |
| `--provider.authorization.v1` | `PROVIDER_AUTHORIZATION_V1` | URL of a Authorization V1 provider API implementation       |
| `--client.authtoken`          | `CLIENT_AUTH_TOKEN`         | Valid JWT token to access the CSP API's                     |
| `--client.tenant`             | `CLIENT_TENANT`             | Name of the Tenant used in the tests                        |
| `--client.region`             | `CLIENT_REGION`             | Name of the Region used in the regional tests               |
| `--scenario.users`            | `SCENARIO_USERS`            | Comma-separated list of valid CSP users                     | 
| `--scenario.cidr`             | `SCENARIO_CIDR`             | CIDR range availabed in the CSP to create network resources |
| `--scenario.publicips`        | `SCENARIO_PUBLIC_IPS`       | Public IPs range, in CIDR format, to create CSP public IP's |

## Running

To run the conformance tests, set the [configuration](#configuration) variables and run the following command:
```bash
secatest run \
  --provider.region.v1=#REGION_API \
  --provider.authorization.v1=#AUTHORIZATION_API \
  --client.authtoken=#TOKEN \
  --client.region=#REGION \
  --client.tenant=#TENANT \
  --scenario.users=#USERS \
  --scenario.cidr=#CIDR \
  --scenario.publicips=#PUBLIC_IPS
```

## Viewing Result

To see the the result report run the following command:
```bash
secatest report
```

Your default browser will be opened, with the Allure Report viewer:

![Viewer](docs/report-viewer.png)