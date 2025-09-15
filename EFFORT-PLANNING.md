# Effort Planning: E2.2.2-A credential-management

## Effort ID
E2.2.2-A

## Purpose
Implement real credential management and remove placeholder authentication code

## Scope
- Add CLI flags --username and --token to push and build commands
- Create pkg/gitea/credentials.go with full credential management
- Environment variable support (GITEA_USERNAME, GITEA_PASSWORD)
- Config file parsing (~/.idpbuilder/config)
- System keyring integration for secure storage
- Update getRegistryUsername/Password to return real values
- Remove credential-related TODOs from client.go
- Add comprehensive credential tests

## Dependencies
- E2.2.1 cli-commands (complete)

## Estimated Size
540 lines (including CLI flag integration)

## Implementation Details

### New Files to Create
- pkg/gitea/credentials.go (180 lines - includes CLI provider)
- pkg/gitea/credentials_test.go (100 lines)
- pkg/gitea/config.go (70 lines)
- pkg/gitea/config_test.go (50 lines)

### Files to Modify
- pkg/gitea/client.go (update credential functions, ~60 lines)
- pkg/gitea/client_test.go (update tests, ~50 lines)
- pkg/cmd/push.go (add --username and --token flags, ~30 lines)
- pkg/cmd/build.go (add --username and --token flags, ~30 lines)

### Code to Remove
- TODO comments at lines 30, 145, 151 in client.go
- Empty string returns from credential functions

## Success Criteria
- CLI flags --username and --token working on push/build commands
- Real credentials retrieved from CLI/environment/config/keyring (in priority order)
- All credential-related TODOs removed
- Tests passing for all credential scenarios
- No hardcoded credentials or placeholders
- Credential priority: CLI > Environment > Config > Keyring