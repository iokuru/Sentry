package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DeploymentOutputs struct {
	ApiUrl                 string
	UserPoolId             string
	UserPoolClientId       string
	Region                 string
	CloudfrontUrl          string
	FrontendBucket         string
	CloudfrontDistribution string
	UploadsBucket          string
}

func main() {
	syncOnly := flag.Bool("sync-only", false, "only sync frontend and generate config without running terraform apply")
	flag.Parse()

	if !*syncOnly {
		log.Println("Step 1/6: Building arm64 lambda packages and frontend assets...")
		if err := runCmd("go", "run", "./cmd/build", "-frontend"); err != nil {
			log.Fatalf("build failed: %v", err)
		}

		log.Println("Step 2/6: Running terraform plan...")
		if err := runCmd("terraform", "-chdir=infra", "plan", "-out=tfplan"); err != nil {
			log.Fatalf("terraform plan failed: %v", err)
		}

		log.Println("Step 3/6: Applying terraform plan...")
		if err := runCmd("terraform", "-chdir=infra", "apply", "-auto-approve", "tfplan"); err != nil {
			log.Fatalf("terraform apply failed: %v", err)
		}
	}

	log.Println("Step 4/6: Reading terraform outputs...")
	outputs, err := readOutputs()
	if err != nil {
		log.Fatalf("failed to read terraform outputs: %v", err)
	}

	log.Println("Step 5/6: Generating runtime /config.json...")
	configData := map[string]string{
		"api_url":             outputs.ApiUrl,
		"user_pool_id":        outputs.UserPoolId,
		"user_pool_client_id": outputs.UserPoolClientId,
		"region":              outputs.Region,
	}

	configBytes, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		log.Fatalf("marshal config: %v", err)
	}

	configPath := filepath.Join("frontend", "build", "config.json")
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		log.Fatalf("write config.json: %v", err)
	}
	_ = os.WriteFile(filepath.Join("frontend", "public", "config.json"), configBytes, 0644)

	var frontendHost, netlifySite, allowedOrigin string
	if data, err := os.ReadFile("inputs.local.json"); err == nil {
		var in struct {
			FrontendHost    string `json:"frontend_host"`
			NetlifySiteName string `json:"netlify_site_name"`
			AllowedOrigin   string `json:"allowed_origin"`
		}
		if err := json.Unmarshal(data, &in); err == nil {
			frontendHost = strings.ToLower(strings.TrimSpace(in.FrontendHost))
			netlifySite = in.NetlifySiteName
			allowedOrigin = strings.ToLower(strings.TrimSpace(in.AllowedOrigin))
		}
	}

	if frontendHost == "netlify" {
		log.Println("Step 6/6: Deploying frontend bundle to Netlify...")
		uploadsHost := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", outputs.UploadsBucket, outputs.Region)
		cspHeader := fmt.Sprintf("Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https:; connect-src 'self' %s %s; font-src 'self'; worker-src 'self' blob:; frame-ancestors 'none'; object-src 'none'; base-uri 'self';", outputs.ApiUrl, uploadsHost)
		headersContent := fmt.Sprintf("/*\n  X-Content-Type-Options: nosniff\n  X-Frame-Options: DENY\n  Referrer-Policy: strict-origin-when-cross-origin\n  Strict-Transport-Security: max-age=31536000; includeSubDomains; preload\n  %s\n", cspHeader)
		_ = os.WriteFile(filepath.Join("frontend", "build", "_headers"), []byte(headersContent), 0644)

		netPath, err := exec.LookPath("netlify")
		if err != nil {
			log.Fatalf("netlify cli not found on PATH: %v", err)
		}
		deployArgs := []string{"deploy", "--prod", "--dir=" + filepath.Join("frontend", "build")}
		if netlifySite != "" {
			deployArgs = append(deployArgs, "--site="+netlifySite)
		}
		if err := runCmd(netPath, deployArgs...); err != nil {
			log.Fatalf("netlify deploy failed: %v", err)
		}
	} else {
		log.Println("Step 6/6: Uploading frontend to S3 and invalidating CloudFront cache...")
		if outputs.FrontendBucket != "" {
			s3Target := fmt.Sprintf("s3://%s", outputs.FrontendBucket)
			if err := runCmd("aws", "s3", "sync", filepath.Join("frontend", "build"), s3Target, "--delete"); err != nil {
				log.Printf("warning: aws s3 sync: %v", err)
			}
		}

		if outputs.CloudfrontDistribution != "" {
			if err := runCmd("aws", "cloudfront", "create-invalidation", "--distribution-id", outputs.CloudfrontDistribution, "--paths", "/*"); err != nil {
				log.Printf("warning: cloudfront invalidation: %v", err)
			}
		}
	}

	frontendUrl := outputs.CloudfrontUrl
	if frontendHost == "netlify" {
		if allowedOrigin != "" {
			frontendUrl = allowedOrigin
		} else if netlifySite != "" {
			frontendUrl = fmt.Sprintf("https://%s.netlify.app", netlifySite)
		}
	}

	fmt.Println("\n=======================================================")
	fmt.Println("Sentry deployment successfully completed!")
	fmt.Printf("Frontend URL: %s\n", frontendUrl)
	fmt.Printf("API URL:      %s\n", outputs.ApiUrl)
	fmt.Printf("User Pool:    %s\n", outputs.UserPoolId)
	fmt.Printf("App Client:   %s\n", outputs.UserPoolClientId)
	fmt.Printf("AWS Region:   %s\n", outputs.Region)
	fmt.Println("=======================================================")
}

func readOutputs() (*DeploymentOutputs, error) {
	out := &DeploymentOutputs{}
	var err error

	out.ApiUrl, err = getOutput("api_url")
	if err != nil {
		return nil, err
	}
	out.UserPoolId, _ = getOutput("user_pool_id")
	out.UserPoolClientId, _ = getOutput("user_pool_client_id")
	out.Region, _ = getOutput("region")
	out.CloudfrontUrl, _ = getOutput("cloudfront_url")
	out.FrontendBucket, _ = getOutput("frontend_bucket")
	out.CloudfrontDistribution, _ = getOutput("cloudfront_distribution_id")
	out.UploadsBucket, _ = getOutput("uploads_bucket")

	return out, nil
}

func getOutput(name string) (string, error) {
	cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", name)
	res, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("terraform output %s: %w", name, err)
	}
	return strings.TrimSpace(string(res)), nil
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
