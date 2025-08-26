package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// Push uploads an image to the registry with automatic signing
func (rc *registryClient) Push(ctx context.Context, image string, auth api.AuthConfig) error {
	ref, err := parseImageReference(image)
	if err != nil {
		return fmt.Errorf("invalid image reference: %w", err)
	}

	authHandler := newAuthHandler(auth)
	authHandler.setHTTPClient(rc.httpClient)

	// Get image data from local store (simplified for scope)
	manifest, layers, err := rc.getLocalImageData(ref)
	if err != nil {
		return fmt.Errorf("failed to get local image: %w", err)
	}

	// Upload layers
	for _, layer := range layers {
		if err := rc.uploadLayer(ctx, ref, layer, authHandler); err != nil {
			return fmt.Errorf("failed to upload layer %s: %w", layer.Digest, err)
		}
	}

	// Upload manifest
	if err := rc.uploadManifest(ctx, ref, manifest, authHandler); err != nil {
		return fmt.Errorf("failed to upload manifest: %w", err)
	}

	// Sign image if security manager available
	if rc.securityMgr != nil {
		if _, err := rc.securityMgr.SignImage(ctx, image, nil); err != nil {
			return fmt.Errorf("failed to sign image: %w", err)
		}
	}

	return nil
}

// Pull downloads an image from the registry with verification
func (rc *registryClient) Pull(ctx context.Context, image string, auth api.AuthConfig) (*api.Image, error) {
	ref, err := parseImageReference(image)
	if err != nil {
		return nil, fmt.Errorf("invalid image reference: %w", err)
	}

	authHandler := newAuthHandler(auth)
	authHandler.setHTTPClient(rc.httpClient)

	// Download manifest
	manifest, err := rc.downloadManifest(ctx, ref, authHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to download manifest: %w", err)
	}

	// Verify signature if security manager available
	if rc.securityMgr != nil {
		if err := rc.securityMgr.VerifySignature(ctx, image, nil); err != nil {
			return nil, fmt.Errorf("signature verification failed: %w", err)
		}
	}

	// Download layers
	layers := make([]*api.LayerInfo, 0, len(manifest.Layers))
	for _, descriptor := range manifest.Layers {
		layer, err := rc.downloadLayer(ctx, ref, descriptor, authHandler)
		if err != nil {
			return nil, fmt.Errorf("failed to download layer %s: %w", descriptor.Digest, err)
		}
		layers = append(layers, layer)
	}

	return &api.Image{
		Name:      ref.Name,
		Tag:       ref.Tag,
		Digest:    manifest.Config.Digest,
		MediaType: manifest.MediaType,
		Layers:    layers,
	}, nil
}

func (rc *registryClient) uploadLayer(ctx context.Context, ref *imageReference, layer *api.LayerInfo, auth *authHandler) error {
	baseURL := rc.getRegistryURL(ref)
	
	headURL := fmt.Sprintf("%s/v2/%s/blobs/%s", baseURL, ref.Name, layer.Digest)
	if headReq, _ := http.NewRequestWithContext(ctx, "HEAD", headURL, nil); headReq != nil {
		auth.authenticate(headReq, nil)
		if resp, err := rc.httpClient.Do(headReq); err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
	}

	uploadURL := fmt.Sprintf("%s/v2/%s/blobs/uploads/", baseURL, ref.Name)
	initReq, _ := http.NewRequestWithContext(ctx, "POST", uploadURL, nil)
	auth.authenticate(initReq, nil)
	
	resp, err := rc.httpClient.Do(initReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("upload init failed: %s", resp.Status)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return fmt.Errorf("no upload location provided")
	}

	putURL := fmt.Sprintf("%s&digest=%s", location, layer.Digest)
	putReq, _ := http.NewRequestWithContext(ctx, "PUT", putURL, bytes.NewReader(layer.Data))
	putReq.Header.Set("Content-Type", "application/octet-stream")
	putReq.Header.Set("Content-Length", strconv.Itoa(len(layer.Data)))
	auth.authenticate(putReq, nil)

	putResp, err := rc.httpClient.Do(putReq)
	if err != nil {
		return err
	}
	defer putResp.Body.Close()

	if putResp.StatusCode != http.StatusCreated {
		return fmt.Errorf("layer upload failed: %s", putResp.Status)
	}
	return nil
}

func (rc *registryClient) downloadLayer(ctx context.Context, ref *imageReference, descriptor *api.Descriptor, auth *authHandler) (*api.LayerInfo, error) {
	blobURL := fmt.Sprintf("%s/v2/%s/blobs/%s", rc.getRegistryURL(ref), ref.Name, descriptor.Digest)
	req, err := http.NewRequestWithContext(ctx, "GET", blobURL, nil)
	if err != nil {
		return nil, err
	}
	auth.authenticate(req, nil)
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("layer download failed: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &api.LayerInfo{
		Digest:    descriptor.Digest,
		Size:      descriptor.Size,
		MediaType: descriptor.MediaType,
		Data:      data,
	}, nil
}

func (rc *registryClient) getLocalImageData(ref *imageReference) (*api.Manifest, []*api.LayerInfo, error) {
	return &api.Manifest{}, []*api.LayerInfo{}, nil
}

func (rc *registryClient) uploadManifest(ctx context.Context, ref *imageReference, manifest *api.Manifest, auth *authHandler) error {
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/%s", rc.getRegistryURL(ref), ref.Name, ref.Tag)
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", manifestURL, bytes.NewReader(manifestBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", manifest.MediaType)
	auth.authenticate(req, nil)
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("manifest upload failed: %s", resp.Status)
	}
	return nil
}

func (rc *registryClient) downloadManifest(ctx context.Context, ref *imageReference, auth *authHandler) (*api.Manifest, error) {
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/%s", rc.getRegistryURL(ref), ref.Name, ref.Tag)
	req, err := http.NewRequestWithContext(ctx, "GET", manifestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	auth.authenticate(req, nil)
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest download failed: %s", resp.Status)
	}
	var manifest api.Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}