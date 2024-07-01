package ts2go_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/armsnyder/ts2go"
)

func TestGenerate(t *testing.T) {
	const testdataDir = "internal/testdata"

	dirs, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, dir := range dirs {
		t.Run(
			dir.Name(),
			func(t *testing.T) {
				const sourceFile = "source.ts"
				const expectedFile = "expected.go"

				source, err := os.ReadFile(filepath.Join(testdataDir, dir.Name(), sourceFile))
				if err != nil {
					t.Fatalf("failed to read source.ts: %v", err)
				}

				expected, err := os.ReadFile(filepath.Join(testdataDir, dir.Name(), expectedFile))
				if err != nil {
					t.Fatalf("failed to read expected.go: %v", err)
				}

				generated := &bytes.Buffer{}

				if err := ts2go.Generate(bytes.NewReader(source), generated); err != nil {
					t.Fatalf("failed to generate code: %v", err)
				}

				if !bytes.Equal(generated.Bytes(), expected) {
					t.Errorf("Generated code does not match expected output")
					logDiff(t, expected, generated.Bytes())
				}
			},
		)
	}
}

func logDiff(t *testing.T, expected, generated []byte) {
	t.Helper()

	path, err := exec.LookPath("diff")
	if err != nil {
		t.Log("Install diff(1) to see a prettier diff")
		t.Logf("\n=== Expected ===\n%s\n\n=== Generated ===\n%s", expected, generated)
		return
	}

	tmpDir := t.TempDir()
	expectedPath := filepath.Join(tmpDir, "expected")
	generatedPath := filepath.Join(tmpDir, "generated")

	if err := os.WriteFile(expectedPath, expected, 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := os.WriteFile(generatedPath, generated, 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cmd := exec.Command(path, "-u", "--label", "expected", "--label", "generated", expectedPath, generatedPath)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	cmd.Stderr = buf

	if err := cmd.Run(); err != nil {
		t.Logf("\n%s", buf.String())
	}
}
