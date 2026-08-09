package truststore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigureCandidate(t *testing.T) {
	t.Setenv("SSL_CERT_FILE", "")
	caPath := filepath.Join(t.TempDir(), "ca.crt")
	writeTestCertificate(t, caPath)

	configured, err := configureCandidate(caPath)
	if err != nil {
		t.Fatalf("configure CA: %v", err)
	}
	if !configured {
		t.Fatal("expected mounted CA to be configured")
	}
	if got := os.Getenv("SSL_CERT_FILE"); got != caPath {
		t.Fatalf("SSL_CERT_FILE = %q, want %q", got, caPath)
	}
}

func TestConfigurePlatformCAHonoursExplicitSSLFile(t *testing.T) {
	t.Setenv("SSL_CERT_FILE", "/already/configured.pem")
	t.Setenv("PLATFORM_CA_ROOT", filepath.Join(t.TempDir(), "missing.pem"))

	configured, err := ConfigurePlatformCA()
	if err != nil {
		t.Fatalf("configure CA: %v", err)
	}
	if configured {
		t.Fatal("explicit SSL_CERT_FILE must not be replaced")
	}
}

func TestConfigurePlatformCARejectsInvalidPEM(t *testing.T) {
	t.Setenv("SSL_CERT_FILE", "")
	caPath := filepath.Join(t.TempDir(), "ca.crt")
	if err := os.WriteFile(caPath, []byte("not a certificate"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := configureCandidate(caPath); err == nil {
		t.Fatal("expected invalid PEM to be rejected")
	}
}

func TestApprovedPlatformCAPath(t *testing.T) {
	for _, test := range []struct {
		name       string
		configured string
		want       string
	}{
		{name: "default when unset", configured: "", want: defaultPlatformCA},
		{name: "preferred mount", configured: defaultPlatformCA, want: defaultPlatformCA},
		{name: "legacy mount", configured: legacyPlatformCA, want: legacyPlatformCA},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := approvedPlatformCAPath(test.configured)
			if err != nil {
				t.Fatalf("approvedPlatformCAPath: %v", err)
			}
			if got != test.want {
				t.Fatalf("approvedPlatformCAPath(%q) = %q, want %q", test.configured, got, test.want)
			}
		})
	}
}

func TestApprovedPlatformCAPathRejectsUntrustedVariants(t *testing.T) {
	for _, configured := range []string{
		"../../tmp/attacker-ca.crt",
		defaultPlatformCA + "/../attacker-ca.crt",
		defaultPlatformCA + "\n/attacker-ca.crt",
		"/tmp/attacker-ca.crt",
	} {
		t.Run(configured, func(t *testing.T) {
			if got, err := approvedPlatformCAPath(configured); err == nil || got != "" {
				t.Fatalf("approvedPlatformCAPath(%q) = %q, %v; want rejection", configured, got, err)
			}
		})
	}
}

func TestConfigurePlatformCARejectsUnapprovedPath(t *testing.T) {
	t.Setenv("SSL_CERT_FILE", "")
	caPath := filepath.Join(t.TempDir(), "ca.crt")
	writeTestCertificate(t, caPath)
	t.Setenv("PLATFORM_CA_ROOT", caPath)

	configured, err := ConfigurePlatformCA()
	if err == nil {
		t.Fatal("expected unapproved CA path to be rejected")
	}
	if configured {
		t.Fatal("unapproved CA path must not be configured")
	}
	if got := os.Getenv("SSL_CERT_FILE"); got != "" {
		t.Fatalf("SSL_CERT_FILE = %q, want empty", got)
	}
}

func TestConfigurePlatformCARejectsEscapingSymlink(t *testing.T) {
	t.Setenv("SSL_CERT_FILE", "")
	root := t.TempDir()
	realPath := filepath.Join(t.TempDir(), "real-ca.crt")
	writeTestCertificate(t, realPath)
	symlinkPath := filepath.Join(root, "ca.crt")
	if err := os.Symlink(realPath, symlinkPath); err != nil {
		t.Skipf("create symlink: %v", err)
	}
	configured, err := configureCandidate(symlinkPath)
	if err == nil {
		t.Fatal("expected escaping symlink to be rejected")
	}
	if configured {
		t.Fatal("escaping symlink must not be configured")
	}
}

func writeTestCertificate(t *testing.T, path string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test root"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	contents := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}
}
