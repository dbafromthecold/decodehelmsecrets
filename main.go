package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type secretData struct {
	Data map[string]string `json:"data"`
}

type releasePayload struct {
	Chart struct {
		Templates []struct {
			Name string `json:"name"`
			Data string `json:"data"`
		} `json:"templates"`
	} `json:"chart"`
}

func usage() {
	fmt.Println("DecodeHelmSecrets - decode Helm secret release templates from a Kubernetes secret")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  DecodeHelmSecrets -secret <secret-name> [-namespace <namespace>]")
	fmt.Println("  DecodeHelmSecrets <secret-name> [-namespace <namespace>]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -secret string     Kubernetes secret name")
	fmt.Println("  -namespace string  Kubernetes namespace")
}

func main() {
	var secretName string
	var namespace string

	flag.StringVar(&secretName, "secret", "", "Kubernetes secret name")
	flag.StringVar(&namespace, "namespace", "", "Kubernetes namespace")
	flag.Usage = usage
	flag.Parse()

	if secretName == "" && flag.NArg() > 0 {
		secretName = flag.Arg(0)
	}

	if secretName == "" {
		usage()
		os.Exit(2)
	}

	rawSecret, err := fetchSecret(secretName, namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	decoded, err := decodeRelease(rawSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding release payload: %v\n", err)
		os.Exit(1)
	}

	if err := printTemplates(decoded); err != nil {
		fmt.Fprintf(os.Stderr, "error printing templates: %v\n", err)
		os.Exit(1)
	}
}

func fetchSecret(name, namespace string) ([]byte, error) {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return nil, fmt.Errorf("kubectl is not installed or not in PATH: %w", err)
	}

	args := []string{"get", "secret", name, "-o", "json"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("kubectl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("kubectl failed: %w: %s", err, bytes.TrimSpace(output))
	}

	return output, nil
}

func decodeRelease(raw []byte) ([]byte, error) {
	var secret secretData
	if err := json.Unmarshal(raw, &secret); err != nil {
		return nil, fmt.Errorf("failed to parse secret JSON: %w", err)
	}

	releaseBase64, ok := secret.Data["release"]
	if !ok {
		return nil, fmt.Errorf("secret data does not contain release field")
	}

	firstDecode, err := base64.StdEncoding.DecodeString(releaseBase64)
	if err != nil {
		return nil, fmt.Errorf("first base64 decode failed: %w", err)
	}

	secondDecode, err := base64.StdEncoding.DecodeString(string(firstDecode))
	if err != nil {
		return nil, fmt.Errorf("second base64 decode failed: %w", err)
	}

	gz, err := gzip.NewReader(bytes.NewReader(secondDecode))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()

	decoded, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress payload: %w", err)
	}

	return decoded, nil
}

func printTemplates(data []byte) error {
	var payload releasePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("failed to parse decoded JSON: %w", err)
	}

	if len(payload.Chart.Templates) == 0 {
		return fmt.Errorf("decoded release payload does not contain any chart templates")
	}

	for _, template := range payload.Chart.Templates {
		if template.Name == "" {
			return fmt.Errorf("found chart template with empty name")
		}

		fmt.Println(template.Name)
		decodedData, err := base64.StdEncoding.DecodeString(template.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to decode template data for %s: %v\n", template.Name, err)
			fmt.Println(template.Data)
		} else {
			fmt.Println(string(decodedData))
		}
		fmt.Println()
	}

	return nil
}
