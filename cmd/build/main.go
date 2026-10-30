package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var lambdaTargets = []string{
	"presign",
	"intake",
	"analysis",
	"caregiver_list",
	"caregiver_status",
	"caregiver_invite",
	"senior_pair",
	"digest",
	"telegram_hook",
}

func main() {
	cleanFlag := flag.Bool("clean", false, "clean build artifacts")
	frontendFlag := flag.Bool("frontend", false, "build frontend bundle")
	flag.Parse()

	if *cleanFlag {
		if err := cleanArtifacts(); err != nil {
			fmt.Fprintf(os.Stderr, "clean failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("cleaned build artifacts")
		return
	}

	if err := os.MkdirAll("bin", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "create bin dir: %v\n", err)
		os.Exit(1)
	}

	for _, target := range lambdaTargets {
		srcDir := filepath.Join("cmd", target)
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			continue
		}

		if err := buildAndZipLambda(target, srcDir); err != nil {
			fmt.Fprintf(os.Stderr, "failed to build %s: %v\n", target, err)
			os.Exit(1)
		}
		fmt.Printf("packaged bin/%s.zip\n", target)
	}

	if *frontendFlag {
		if err := buildFrontend(); err != nil {
			fmt.Fprintf(os.Stderr, "frontend build failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("frontend built successfully")
	}
}

func cleanArtifacts() error {
	_ = os.RemoveAll("bin")
	_ = os.RemoveAll(filepath.Join("frontend", "build"))
	_ = os.RemoveAll(filepath.Join("frontend", ".next"))
	return nil
}

func buildAndZipLambda(target, srcDir string) error {
	tempDir, err := os.MkdirTemp("", "sentry-build-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	bootstrapPath := filepath.Join(tempDir, "bootstrap")

	cmd := exec.Command("go", "build",
		"-trimpath",
		"-tags", "lambda.norpc",
		"-ldflags", "-s -w",
		"-o", bootstrapPath,
		"./"+srcDir,
	)

	cmd.Env = append(os.Environ(),
		"GOOS=linux",
		"GOARCH=arm64",
		"CGO_ENABLED=0",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}

	zipPath := filepath.Join("bin", target+".zip")
	if err := zipBootstrap(bootstrapPath, zipPath); err != nil {
		return fmt.Errorf("zip binary: %w", err)
	}

	return nil
}

func zipBootstrap(binaryPath, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("create zip: %w", err)
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	info, err := os.Stat(binaryPath)
	if err != nil {
		return fmt.Errorf("stat bootstrap: %w", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("file info header: %w", err)
	}

	header.Name = "bootstrap"
	header.Method = zip.Deflate
	// executable permissions for linux lambda runtime
	header.SetMode(0755)

	writer, err := w.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create header: %w", err)
	}

	f, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("open bootstrap: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(writer, f); err != nil {
		return fmt.Errorf("write to zip: %w", err)
	}

	return nil
}

func buildFrontend() error {
	var cmd *exec.Cmd
	if _, err := exec.LookPath("npm.cmd"); err == nil {
		cmd = exec.Command("npm.cmd", "run", "build")
	} else {
		cmd = exec.Command("npm", "run", "build")
	}
	cmd.Dir = "frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
