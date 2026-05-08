package verify

import (
	"fmt"
	"os"

	"github.com/sigstore/sigstore-go/pkg/bundle"
	"github.com/sigstore/sigstore-go/pkg/root"
	"github.com/sigstore/sigstore-go/pkg/tuf"
	"github.com/sigstore/sigstore-go/pkg/verify"
)

// IsSigstoreSigned checks if a file has a valid Sigstore signature.
// It looks for a .sigstore.json bundle file adjacent to the artifact.
func IsSigstoreSigned(path string) error {
	bundlePath := path + ".sigstore.json"
	if _, err := os.Stat(bundlePath); os.IsNotExist(err) {
		return fmt.Errorf("sigstore bundle not found: %s", bundlePath)
	}

	// 1. Load trusted root from Sigstore public good infrastructure via TUF
	tufClient, err := tuf.New(tuf.DefaultOptions())
	if err != nil {
		return fmt.Errorf("failed to initialize TUF client: %w", err)
	}

	trustedRootJSON, err := tufClient.GetTarget("trusted_root.json")
	if err != nil {
		return fmt.Errorf("failed to fetch trusted root: %w", err)
	}

	trustedRoot, err := root.NewTrustedRootFromJSON(trustedRootJSON)
	if err != nil {
		return fmt.Errorf("failed to parse trusted root: %w", err)
	}

	// 2. Configure the verifier with transparency log and timestamp requirements
	verifierConfig := []verify.VerifierOption{
		verify.WithSignedCertificateTimestamps(1),
		verify.WithObserverTimestamps(1),
		verify.WithTransparencyLog(1),
	}

	trustedMaterial := root.TrustedMaterialCollection{trustedRoot}

	sev, err := verify.NewVerifier(trustedMaterial, verifierConfig...)
	if err != nil {
		return fmt.Errorf("failed to create Sigstore verifier: %w", err)
	}

	// 3. Load the bundle
	b, err := bundle.LoadJSONFromPath(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to load Sigstore bundle: %w", err)
	}

	// 4. Open the artifact for verification
	artifactFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open artifact for verification: %w", err)
	}
	defer artifactFile.Close()

	// 5. Build identity policy
	// Accept any valid Sigstore identity from public good instance.
	// TODO: In production, restrict to your specific issuer and SAN identity.
	certID, err := verify.NewShortCertificateIdentity("", ".*", "", ".*")
	if err != nil {
		return fmt.Errorf("failed to create certificate identity policy: %w", err)
	}

	policy := verify.NewPolicy(
		verify.WithArtifact(artifactFile),
		verify.WithCertificateIdentity(certID),
	)

	// 6. Verify the bundle against the artifact and policy
	_, err = sev.Verify(b, policy)
	if err != nil {
		return fmt.Errorf("sigstore verification failed: %w", err)
	}

	return nil
}
