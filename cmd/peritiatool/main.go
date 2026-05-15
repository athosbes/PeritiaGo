package main

import (
	"fmt"
	"os"

	"github.com/athosbes/PeritiaGo/internal/verify"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: peritiatool <command> <file>\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  sign   <file> [--id-token TOKEN]  Sign a file with Sigstore\n")
		fmt.Fprintf(os.Stderr, "  verify <file>                     Verify a Sigstore signature\n")
		os.Exit(1)
	}

	command := os.Args[1]
	filePath := os.Args[2]

	switch command {
	case "sign":
		idToken := ""
		// Parse --id-token flag
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "--id-token" && i+1 < len(os.Args) {
				idToken = os.Args[i+1]
				break
			}
		}
		// Also check SIGSTORE_ID_TOKEN env var (used in GitHub Actions)
		if idToken == "" {
			idToken = os.Getenv("SIGSTORE_ID_TOKEN")
		}
		if idToken == "" {
			fmt.Fprintf(os.Stderr, "Error: OIDC ID token required. Provide --id-token or set SIGSTORE_ID_TOKEN env var.\n")
			fmt.Fprintf(os.Stderr, "  GitHub Actions: automatically set via OIDC\n")
			fmt.Fprintf(os.Stderr, "  Local: use 'gcloud auth print-identity-token' or similar\n")
			os.Exit(1)
		}

		fmt.Printf("Signing %s with Sigstore...\n", filePath)
		if err := verify.SignBinary(filePath, idToken); err != nil {
			fmt.Fprintf(os.Stderr, "Signing failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Success: Bundle saved to %s.sigstore.json\n", filePath)

	case "verify":
		fmt.Printf("Verifying Sigstore signature for %s...\n", filePath)
		if err := verify.IsSigstoreSigned(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Verification failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Success: Sigstore signature is valid.")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
