package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// GetManifest retrieves and parses an image manifest
func (rc *registryClient) GetManifest(ctx context.Context, image string) (*api.Manifest, error) {
	ref, err := parseImageReference(image)
	if err != nil {
		return nil, fmt.Errorf("invalid image reference: %w", err)
	}

	authHandler := newAuthHandler(rc.defaultAuth)
	authHandler.setHTTPClient(rc.httpClient)

	return rc.downloadManifest(ctx, ref, authHandler)
}

// ListTags returns all available tags for a repository
func (rc *registryClient) ListTags(ctx context.Context, repository string) ([]string, error) {
	baseURL := rc.getRegistryBaseURL()
	url := fmt.Sprintf("%s/v2/%s/tags/list", baseURL, repository)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	authHandler := newAuthHandler(rc.defaultAuth)
	authHandler.setHTTPClient(rc.httpClient)
	authHandler.authenticate(req, nil)

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tags list failed: %s", resp.Status)
	}

	var result struct {
		Tags []string `json:"tags"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode tags: %w", err)
	}

	return result.Tags, nil
}

func (rc *registryClient) CheckManifestExists(ctx context.Context, image string) (bool, error) {
	ref, err := parseImageReference(image)
	if err != nil {
		return false, fmt.Errorf("invalid image reference: %w", err)
	}
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/%s", rc.getRegistryURL(ref), ref.Name, ref.Tag)
	req, err := http.NewRequestWithContext(ctx, "HEAD", manifestURL, nil)
	if err != nil {
		return false, err
	}
	authHandler := newAuthHandler(rc.defaultAuth)
	authHandler.setHTTPClient(rc.httpClient)
	authHandler.authenticate(req, nil)
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}