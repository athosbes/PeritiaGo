package verify

import (
	"fmt"
	"os"
	"time"

	"github.com/sigstore/sigstore-go/pkg/root"
	"github.com/sigstore/sigstore-go/pkg/sign"
	"github.com/sigstore/sigstore-go/pkg/tuf"
	"google.golang.org/protobuf/encoding/protojson"
)

// SignBinary signs a binary file using Sigstore and produces a .sigstore.json bundle.
// idToken is the OIDC identity token (e.g. from GitHub Actions or gcloud).
func SignBinary(path string, idToken string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read artifact: %w", err)
	}

	content := &sign.PlainData{Data: data}

	keypair, err := sign.NewEphemeralKeypair(nil)
	if err != nil {
		return fmt.Errorf("failed to create ephemeral keypair: %w", err)
	}

	opts := sign.BundleOptions{}

	// Setup TUF client for Sigstore public good instance
	tufClient, err := tuf.New(tuf.DefaultOptions())
	if err != nil {
		return fmt.Errorf("failed to initialize TUF client: %w", err)
	}

	opts.TrustedRoot, err = root.GetTrustedRoot(tufClient)
	if err != nil {
		return fmt.Errorf("failed to get trusted root: %w", err)
	}

	signingConfig, err := root.GetSigningConfig(tufClient)
	if err != nil {
		return fmt.Errorf("failed to get signing config: %w", err)
	}

	// Use Fulcio for certificate issuance with the OIDC token
	if idToken != "" {
		fulcioService, err := root.SelectService(
			signingConfig.FulcioCertificateAuthorityURLs(),
			sign.FulcioAPIVersions,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to select Fulcio service: %w", err)
		}

		opts.CertificateProvider = sign.NewFulcio(&sign.FulcioOptions{
			BaseURL: fulcioService.URL,
			Timeout: 30 * time.Second,
			Retries: 1,
		})
		opts.CertificateProviderOptions = &sign.CertificateProviderOptions{
			IDToken: idToken,
		}
	}

	// Add Rekor transparency log
	rekorServices, err := root.SelectServices(
		signingConfig.RekorLogURLs(),
		signingConfig.RekorLogURLsConfig(),
		sign.RekorAPIVersions,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to select Rekor services: %w", err)
	}
	for _, svc := range rekorServices {
		opts.TransparencyLogs = append(opts.TransparencyLogs, sign.NewRekor(&sign.RekorOptions{
			BaseURL: svc.URL,
			Timeout: 90 * time.Second,
			Retries: 1,
			Version: svc.MajorAPIVersion,
		}))
	}

	// Add timestamp authority
	tsaServices, err := root.SelectServices(
		signingConfig.TimestampAuthorityURLs(),
		signingConfig.TimestampAuthorityURLsConfig(),
		sign.TimestampAuthorityAPIVersions,
		time.Now(),
	)
	if err == nil {
		for _, svc := range tsaServices {
			opts.TimestampAuthorities = append(opts.TimestampAuthorities, sign.NewTimestampAuthority(&sign.TimestampAuthorityOptions{
				URL:     svc.URL,
				Timeout: 30 * time.Second,
				Retries: 1,
			}))
		}
	}

	// Sign
	bndl, err := sign.Bundle(content, keypair, opts)
	if err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}

	bundleJSON, err := protojson.Marshal(bndl)
	if err != nil {
		return fmt.Errorf("failed to serialize bundle: %w", err)
	}

	bundlePath := path + ".sigstore.json"
	if err := os.WriteFile(bundlePath, bundleJSON, 0644); err != nil {
		return fmt.Errorf("failed to write bundle: %w", err)
	}

	return nil
}
