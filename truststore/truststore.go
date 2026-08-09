package truststore

import (
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	defaultPlatformCA = "/var/lib/pasturestack/etc/ssl/ca.crt"
	legacyPlatformCA  = "/var/lib/rancher/etc/ssl/ca.crt"
	maxCAFileSize     = 1024 * 1024
)

// ConfigurePlatformCA selects a mounted control-platform CA without changing
// the container's system trust store. Go also loads the system certificate
// directory, so public roots remain available when SSL_CERT_FILE is set.
func ConfigurePlatformCA() (bool, error) {
	if os.Getenv("SSL_CERT_FILE") != "" {
		return false, nil
	}

	preferred, err := approvedPlatformCAPath(os.Getenv("PLATFORM_CA_ROOT"))
	if err != nil {
		return false, err
	}

	candidates := []string{preferred}
	if preferred != legacyPlatformCA {
		candidates = append(candidates, legacyPlatformCA)
	}
	for _, candidate := range candidates {
		configured, err := configureCandidate(candidate)
		if err != nil {
			return false, err
		}
		if configured {
			return true, nil
		}
	}

	return false, nil
}

func approvedPlatformCAPath(configured string) (string, error) {
	switch configured {
	case "", defaultPlatformCA:
		return defaultPlatformCA, nil
	case legacyPlatformCA:
		return legacyPlatformCA, nil
	default:
		return "", fmt.Errorf("PLATFORM_CA_ROOT must be an approved control-platform CA path")
	}
}

func configureCandidate(path string) (bool, error) {
	directory, name := filepath.Split(path)
	root, err := os.OpenRoot(filepath.Clean(directory))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("open control-platform CA directory: %w", err)
	}
	defer root.Close()

	linkInfo, err := root.Lstat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect control-platform CA entry: %w", err)
	}
	if linkInfo.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("control-platform CA must not be a symbolic link")
	}

	file, err := root.Open(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("open control-platform CA: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("inspect control-platform CA: %w", err)
	}
	if !info.Mode().IsRegular() {
		return false, fmt.Errorf("control-platform CA is not a regular file")
	}
	if info.Size() <= 0 || info.Size() > maxCAFileSize {
		return false, fmt.Errorf("control-platform CA size is outside the accepted range")
	}

	contents, err := io.ReadAll(io.LimitReader(file, maxCAFileSize+1))
	if err != nil {
		return false, fmt.Errorf("read control-platform CA: %w", err)
	}
	if len(contents) == 0 || len(contents) > maxCAFileSize {
		return false, fmt.Errorf("control-platform CA size is outside the accepted range")
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(contents) {
		return false, fmt.Errorf("control-platform CA does not contain a valid PEM certificate")
	}
	if err := os.Setenv("SSL_CERT_FILE", path); err != nil {
		return false, fmt.Errorf("configure control-platform CA: %w", err)
	}
	return true, nil
}
