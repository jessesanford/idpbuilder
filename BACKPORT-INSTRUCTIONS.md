# Backport Instructions for SW Engineer
## Assignment: cli-commands

## Your Working Directory
`/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave2/cli-commands`

## Context
The integration testing identified critical fixes for the push command that need to be backported to this effort.

## Fixes to Apply

### File: `pkg/cmd/push/push.go`

1. **Trust Store Initialization** (Lines 74-79):
   ```diff
   - var trustStore certs.TrustStoreManager
   - // Setup certificate trust (unless --insecure)
   - if !insecure {
   -     trustStore, err = certs.NewTrustStoreManager("")
   + // Always create trust store (required by NewGiteaClient)
   + trustStore, err := certs.NewTrustStoreManager("")
   + if err != nil {
   +     return fmt.Errorf("failed to create trust store: %w", err)
   + }
   + if !insecure {
   ```
   **Issue Fixed**: Trust store was nil when --insecure flag was used, causing crash

2. **Tarball Path Handling** (Lines 48-50):
   ```diff
   - // Validate image name
   - if !strings.Contains(image, ":") {
   + // Don't add :latest tag if this is a tarball path
   + if !strings.Contains(image, ".tar") && !strings.Contains(image, ":") {
       image = image + ":latest"
   ```
   **Issue Fixed**: Code was adding ":latest" to tarball filenames

3. **Registry Reference Construction** (Lines 137-146):
   ```diff
   + // Parse registry URL to get host
   + registryHost := strings.TrimPrefix(registryURL, "https://")
   + registryHost = strings.TrimPrefix(registryHost, "http://")
   + 
   + // Construct the full image reference for pushing
   + imageRef := fmt.Sprintf("%s/%s/hello-world:v1", registryHost, strings.ToLower(username))
   + 
   - progress.UpdateMessage("Pushing to registry")
   + progress.UpdateMessage("Pushing to " + imageRef)
   ```
   **Issue Fixed**: Push was trying to use Docker Hub instead of Gitea registry

4. **Use Correct Reference for Push** (Lines 156-161):
   ```diff
   - if err := client.Push(context.Background(), img, image, pushOpts); err != nil {
   + if err := client.Push(context.Background(), img, imageRef, pushOpts); err != nil {
   
   - fmt.Printf("Successfully pushed %s\n", image)
   + fmt.Printf("Successfully pushed %s to %s\n", tarballPath, imageRef)
   ```
   **Issue Fixed**: Using tarball path as image reference instead of registry reference

### File: `pkg/registry/gitea_client_test.go` (if exists in this effort)
Apply the same test assertion fixes as listed for gitea-registry-client effort.

## Process
1. Check if this directory has a git repository
2. If no git repo exists, check if the code needs these fixes
3. Apply ALL the fixes to the push.go file
4. Apply test fixes if the test file exists
5. Build the code to verify it compiles
6. Run tests if available
7. Document your work in sw-engineer-state.yaml
8. Report completion status

## Success Criteria
- All push command fixes applied correctly
- Code compiles without errors
- Tests pass (if available)
- Push command functionality fixed
- Work documented properly

## Notes
- These are CRITICAL fixes - the push command doesn't work without them
- The original effort branch may not exist remotely
- Focus on ensuring the local code has the correct fixes
- If the fixes are already applied, document that fact