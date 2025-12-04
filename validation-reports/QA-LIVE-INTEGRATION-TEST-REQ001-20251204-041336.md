# QA LIVE INTEGRATION TEST REPORT: REQ-001

**Validation Level**: PROJECT (Live Integration Test)
**Validated By**: qa-agent-20251204-041336
**Validated At**: 2025-12-04T04:15:30Z
**Validation Result**: PASSED

---

## Requirement Under Test

**REQ-001**: "WHEN the user executes `idpbuilder push <image-name> --registry <registry-url>` with valid credentials and a locally available image THEN the image shall be pushed to the registry and be pullable from that registry"

---

## Test Environment

| Component | Value |
|-----------|-------|
| Registry | gitea.cnoe.localtest.me:8443 |
| Registry Type | Gitea OCI Registry |
| Username | giteaAdmin |
| Test Image | giteaadmin/testapp:v1.0 |
| Binary | ./bin/idpbuilder |
| Binary SHA256 | 889200769da7ac8c35c3761c6d2cc0f9ac94601f6df0277d37a443dd17539629 |

---

## R775 Cryptographic Execution Proof

### Evidence Chain

| Artifact | SHA256 Hash |
|----------|-------------|
| Binary (idpbuilder) | 889200769da7ac8c35c3761c6d2cc0f9ac94601f6df0277d37a443dd17539629 |
| Registry Before | 3a6820b392c60d05f57e3add60aebc65dc4aa4c54521aba7868771078a24064f |
| Registry After | 6aeba5907c7147ded75374be28ec0c2fa9e87ee819ef35a6a1aef7a306955880 |
| Execution Log | 11030cdad02d426e073c35ab1330be90333d5c89251f0047961f9766b386d67c |

---

## Test Execution

### Pre-Conditions Verified

- [x] idpbuilder binary exists and is executable
- [x] Test image (testapp:v1.0) exists locally in Docker
- [x] Gitea registry accessible (returns 401 without auth)
- [x] Valid credentials available

### Registry State Before Push

```json
{"repositories":["giteaadmin/gitea/gitea","giteaadmin/kindest/node"]}
```

Note: `giteaadmin/testapp` does NOT exist in registry.

### Push Command Executed

```bash
./bin/idpbuilder push giteaadmin/testapp:v1.0 \
    --registry https://gitea.cnoe.localtest.me:8443 \
    --username giteaAdmin \
    --password [REDACTED] \
    --insecure
```

**Start Time**: 2025-12-04T04:15:16.789700712Z
**End Time**: 2025-12-04T04:15:17.997548296Z
**Exit Code**: 0

**Output**:
```
gitea.cnoe.localtest.me:8443/giteaadmin/testapp:v1.0@sha256:1eb9d8d452eafa7c6a78e88a5de24b9f290d7b7445e61ec3f6c20561753dfb62
```

### Registry State After Push

```json
{"repositories":["giteaadmin/gitea/gitea","giteaadmin/kindest/node","giteaadmin/testapp"]}
```

Note: `giteaadmin/testapp` NOW EXISTS in registry.

### Tags Verification

```json
{"name":"giteaadmin/testapp","tags":["v1.0"]}
```

### Pull Verification

```bash
docker pull gitea.cnoe.localtest.me:8443/giteaadmin/testapp:v1.0
```

**Output**:
```
v1.0: Pulling from giteaadmin/testapp
Digest: sha256:d3627c1d653bb27e10cb3d6c40f8558d593d231f20700ad4d86afd101b794129
Status: Downloaded newer image for gitea.cnoe.localtest.me:8443/giteaadmin/testapp:v1.0
```

**Pull Exit Code**: 0

---

## REQ-001 Acceptance Criteria Validation

| Criterion | Expected | Actual | Result |
|-----------|----------|--------|--------|
| Command executes | Exit code 0 | Exit code 0 | PASS |
| Image pushed to registry | Image in catalog | Image in catalog | PASS |
| Image pullable from registry | Pull succeeds | Pull succeeds | PASS |
| Output shows reference | Reference printed | Reference printed | PASS |

---

## OBSERVATIONS

### Note on Namespace Requirement

The Gitea registry requires a namespace prefix (e.g., `giteaadmin/`). When testing with plain `testapp:v1.0` (no namespace), the push failed with 404 error because Gitea doesn't support root-level repositories.

**First attempt (failed)**:
```bash
./bin/idpbuilder push testapp:v1.0 --registry https://gitea.cnoe.localtest.me:8443 ...
# Error: 404 Not Found
```

**Second attempt (succeeded)**:
```bash
./bin/idpbuilder push giteaadmin/testapp:v1.0 --registry https://gitea.cnoe.localtest.me:8443 ...
# Success
```

This is expected behavior for Gitea's registry which requires namespace-prefixed repository names.

---

## Validation Result

### REQ-001: **PASSED**

The `idpbuilder push` command:
1. Successfully authenticated with the Gitea registry
2. Pushed the local image to the registry (exit code 0)
3. Registry catalog shows the new image
4. Image is pullable from the registry (verified via docker pull)

---

## Evidence Files

All evidence stored in `.software-factory/evidence/`:

1. `binary-hash-20251204-041336.txt` - Binary SHA256 hash
2. `registry-before-20251204-041336.json` - Registry state before push
3. `registry-after-20251204-041336.json` - Registry state after push
4. `live-push-execution-20251204-041336.txt` - Full execution log with R775 proof

---

## QA Agent Approval

**Status**: APPROVED
**Recommendation**: REQ-001 implementation is functional and meets acceptance criteria

**Signed**: qa-agent-20251204-041336
**Date**: 2025-12-04T04:15:30Z
