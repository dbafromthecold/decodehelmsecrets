# DecodeHelmSecrets

DecodeHelmSecrets is a lightweight command-line utility for extracting and decoding Helm release templates from Kubernetes secrets.

The tool retrieves a Helm release secret via `kubectl`, reads the `data.release` payload, performs the required base64 and gzip decoding steps, and outputs each decoded chart template.

## Features

- Fetch Helm release secret data from Kubernetes
- Decode nested base64 payloads
- Decompress gzip-encoded release content
- Print decoded chart template names and contents
- No external JSON or shell parsing dependencies required at runtime

## Installation

Build the tool locally:

```bash
go build -o DecodeHelmSecrets
```

## Usage

Run the tool with a secret name:

```bash
./DecodeHelmSecrets -secret <secret-name>
```

Optionally specify a namespace:

```bash
./DecodeHelmSecrets -secret <secret-name> -namespace <namespace>
```

The secret name can also be provided as the first positional argument:

```bash
./DecodeHelmSecrets <secret-name>
```

## Requirements

- Go installed for building the utility
- `kubectl` installed and configured for the target Kubernetes cluster
- Access to the Helm release secret in the target namespace

## Notes

This utility is intended for inspecting Helm release payloads stored in Kubernetes secrets. It does not modify cluster resources.
