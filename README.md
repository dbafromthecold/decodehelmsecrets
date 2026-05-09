# DecodeHelmSecrets

A simple Go utility that converts the `decodehelmsecrets.sh` workflow into a native Go app.

This tool fetches a Helm secret from Kubernetes using `kubectl`, extracts the `data.release` payload, decodes it twice from base64, decompresses the gzip payload, and prints each Helm chart template name and rendered data.

## Usage

Build the tool:

```bash
go build -o DecodeHelmSecrets
```

Run it against a secret:

```bash
./DecodeHelmSecrets -secret <secret-name> [-namespace <namespace>]
```

You can also use the secret name as the first positional argument:

```bash
./DecodeHelmSecrets <secret-name>
```

## Requirements

- `kubectl` installed and configured for the target cluster
- Access to the secret in the target namespace

## Notes

This utility is a direct Go rewrite of the original shell script and removes the dependency on `jq`, `gunzip`, and shell string parsing.
