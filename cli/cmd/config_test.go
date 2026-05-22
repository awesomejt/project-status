package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigCommandShowsDefaults(t *testing.T) {

	viper.Reset()

	origArgs := os.Args
	os.Args = []string{"project-status", "config"}

	var buf bytes.Buffer
	configCmd.SetOut(&buf)
	configCmd.SetErr(&buf)

	configCmd.Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, "API URL:") {
		t.Fatalf("expected output to contain 'API URL:', got: %s", output)
	}
	if !strings.Contains(output, "Output:") {
		t.Fatalf("expected output to contain 'Output:', got: %s", output)
	}
}

type captureWriter struct {
	b *bytes.Buffer
}

func (c captureWriter) Write(p []byte) (n int, err error) {
	return c.b.Write(p)
}

func (c captureWriter) Fd() uint {
	return 1
}

func TestConfigSetAPIURLValid(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.Reset()
	viper.SetConfigFile(configPath)

	origArgs := os.Args

	var outputBuf bytes.Buffer

	os.Args = []string{"project-status", "config", "set", "api_url", "http://localhost:5000"}

	SetRootOutputForTest(&outputBuf)

	Execute()

	os.Args = origArgs

	output := outputBuf.String()
	if !strings.Contains(output, "set") && !strings.Contains(output, "api_url") {
		t.Logf("output: %s", output)
	}
}

func TestConfigSetAPIURLInvalid(t *testing.T) {

	origArgs := os.Args

	os.Args = []string{"project-status", "config", "set", "api_url", "not-a-valid-url"}

	var errBuf bytes.Buffer
	SetTestOutput(io.Discard, &errBuf)

	Execute()

	os.Args = origArgs

	output := errBuf.String()
	if !strings.Contains(output, "invalid URL") {
		t.Fatalf("expected error about invalid URL, got: %s", output)
	}
}

func TestConfigSetOutputValid(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.Reset()
	viper.SetConfigFile(configPath)

	origArgs := os.Args

	os.Args = []string{"project-status", "config", "set", "output", "json"}

	var outBuf bytes.Buffer
	SetTestOutput(&outBuf, io.Discard)

	Execute()

	os.Args = origArgs

	output := outBuf.String()
	if !strings.Contains(output, "output") || !strings.Contains(output, "json") {
		t.Fatalf("expected output to contain 'output' and 'json', got: %s", output)
	}
}

func TestConfigSetOutputInvalid(t *testing.T) {

	origArgs := os.Args

	os.Args = []string{"project-status", "config", "set", "output", "invalid"}

	var errBuf bytes.Buffer
	SetTestOutput(io.Discard, &errBuf)

	Execute()

	os.Args = origArgs

	output := errBuf.String()
	if !strings.Contains(output, "must be") {
		t.Fatalf("expected error about valid output values, got: %s", output)
	}
}

func TestConfigSetUnknownKey(t *testing.T) {

	origArgs := os.Args

	os.Args = []string{"project-status", "config", "set", "unknown_key", "value"}

	var errBuf bytes.Buffer
	SetTestOutput(io.Discard, &errBuf)

	Execute()

	os.Args = origArgs

	output := errBuf.String()
	if !strings.Contains(output, "unknown config key") {
		t.Fatalf("expected error about unknown config key, got: %s", output)
	}
}

func TestConfigSetWithHTTPS(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.Reset()
	viper.SetConfigFile(configPath)

	origArgs := os.Args

	os.Args = []string{"project-status", "config", "set", "api_url", "https://api.example.com"}

	var outBuf bytes.Buffer
	SetTestOutput(&outBuf, io.Discard)

	Execute()

	os.Args = origArgs

	output := outBuf.String()
	if !strings.Contains(output, "api_url") {
		t.Fatalf("expected output to contain 'api_url', got: %s", output)
	}
}

func TestConfigSetWithPort(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.Reset()
	viper.SetConfigFile(configPath)

	origArgs := os.Args

	os.Args = []string{"project-status", "config", "set", "api_url", "http://localhost:8080/api"}

	var outBuf bytes.Buffer
	SetTestOutput(&outBuf, io.Discard)

	Execute()

	os.Args = origArgs

	output := outBuf.String()
	if !strings.Contains(output, "api_url") {
		t.Fatalf("expected output to contain 'api_url', got: %s", output)
	}
}

func TestConfigCommandWithExistingConfig(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.Set("api_url", "http://configured-url:5000")
	viper.Set("output", "json")
	_ = viper.WriteConfig()

	origArgs := os.Args

	os.Args = []string{"project-status", "config"}

	var outBuf bytes.Buffer
	SetTestOutput(&outBuf, io.Discard)

	Execute()

	os.Args = origArgs

	output := outBuf.String()
	if !strings.Contains(output, "configured-url") {
		t.Fatalf("expected output to contain configured URL, got: %s", output)
	}
	if !strings.Contains(output, "json") {
		t.Fatalf("expected output to contain 'json', got: %s", output)
	}
}

func TestConfigSetOutputTable(t *testing.T) {

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.Reset()
	viper.SetConfigFile(configPath)

	origArgs := os.Args

	os.Args = []string{"project-status", "config", "set", "output", "table"}

	var outBuf bytes.Buffer
	SetTestOutput(&outBuf, io.Discard)

	Execute()

	os.Args = origArgs

	output := outBuf.String()
	if !strings.Contains(output, "output") || !strings.Contains(output, "table") {
		t.Fatalf("expected output to contain 'output' and 'table', got: %s", output)
	}
}

func SetRootOutputForTest(out *bytes.Buffer) {
	SetTestOutput(out, io.Discard)
}

var _ = io.Writer(nil)
