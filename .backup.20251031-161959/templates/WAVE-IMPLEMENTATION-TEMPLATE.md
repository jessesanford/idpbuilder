# Wave N Implementation Plan

**Wave**: [Wave Number/Name]
**Phase**: [Parent Phase]
**Created**: [Date]
**Planner**: [Code Reviewer Agent Name]
**Fidelity Level**: **EXACT SPECIFICATIONS** (detailed efforts with R213 metadata)

---

## Wave Overview

**Goal**: [1-2 sentence summary from wave architecture]

**Architecture Reference**: See `WAVE-N-ARCHITECTURE.md` for design details

**Total Efforts**: [Number]

---

## Effort Definitions

### Effort 1: [Effort Name]

#### R213 Metadata

```json
{
  "effort_id": "EFFORT_[PHASE]_[WAVE]_001",
  "effort_name": "[Descriptive Name]",
  "branch_name": "effort/[phase]-[wave]-[descriptive-name]",
  "parent_wave": "[WAVE_N]",
  "parent_phase": "[PHASE_N]",
  "depends_on": [],
  "estimated_lines": 250,
  "complexity": "medium",
  "can_parallelize": true
}
```

#### Scope

**Purpose**: [What this effort accomplishes]

**Boundaries**: [What is explicitly OUT of scope]

#### Files to Create/Modify

**New Files**:
- `src/auth/authenticator.py` (150 lines)
- `src/auth/password_hasher.py` (50 lines)
- `tests/unit/test_authenticator.py` (100 lines)

**Modified Files**:
- `src/config.py` (add auth configuration, +20 lines)
- `requirements.txt` (add PyJWT, bcrypt, +2 lines)

**Total Estimated Lines**: 322 lines

#### Exact Code Specifications

**File: src/auth/authenticator.py**

```python
"""User authentication with JWT tokens.

This module provides the Authenticator class for handling user login,
session token generation, and token validation.
"""

from typing import Optional, Dict
from dataclasses import dataclass
import jwt
import bcrypt
from datetime import datetime, timedelta

@dataclass
class AuthResult:
    """Result of an authentication attempt.

    Attributes:
        success: True if authentication succeeded
        user_id: User identifier if successful, None otherwise
        session_token: JWT token if successful, None otherwise
        error_message: Error description if failed, None otherwise
    """
    success: bool
    user_id: Optional[str]
    session_token: Optional[str]
    error_message: Optional[str]


class Authenticator:
    """Handles user authentication with JWT tokens.

    This class provides methods to:
    - Authenticate users with username/password
    - Generate JWT session tokens
    - Validate session tokens
    - Revoke tokens (via blacklist)
    """

    def __init__(self, secret_key: str, token_expiry_seconds: int = 3600):
        """Initialize authenticator.

        Args:
            secret_key: Secret key for signing JWT tokens
            token_expiry_seconds: Token validity duration (default 1 hour)

        Raises:
            ValueError: If secret_key is empty
        """
        if not secret_key:
            raise ValueError("Secret key cannot be empty")

        self.secret_key = secret_key
        self.token_expiry = token_expiry_seconds
        self.algorithm = "HS256"

    def authenticate(self, username: str, password: str, user_repo) -> AuthResult:
        """Authenticate user with credentials.

        Algorithm:
        1. Retrieve user from repository by username
        2. If user not found, return failure (don't reveal why)
        3. Verify password hash using bcrypt
        4. If password invalid, return failure
        5. Generate JWT token with user_id and expiration
        6. Update user's last_login timestamp
        7. Return success with token

        Args:
            username: User's username
            password: User's password (plaintext)
            user_repo: User repository instance

        Returns:
            AuthResult with success status and token if successful
        """
        # Step 1: Retrieve user
        user = user_repo.get_by_username(username)
        if not user:
            return AuthResult(
                success=False,
                user_id=None,
                session_token=None,
                error_message="Invalid credentials"
            )

        # Step 2: Verify password
        if not self._verify_password(password, user.password_hash):
            return AuthResult(
                success=False,
                user_id=None,
                session_token=None,
                error_message="Invalid credentials"
            )

        # Step 3: Generate token
        token = self._generate_token(str(user.id), {"username": user.username})

        # Step 4: Update last login
        user_repo.update_last_login(user.id)

        return AuthResult(
            success=True,
            user_id=str(user.id),
            session_token=token,
            error_message=None
        )

    def validate_token(self, token: str) -> Optional[Dict]:
        """Validate JWT token and extract claims.

        Algorithm:
        1. Decode JWT using secret key
        2. Check if token is expired
        3. Extract user_id from claims
        4. Return claims if valid

        Args:
            token: JWT session token

        Returns:
            Dictionary with user_id and custom claims if valid, None otherwise
        """
        try:
            payload = jwt.decode(token, self.secret_key, algorithms=[self.algorithm])
            return payload
        except jwt.ExpiredSignatureError:
            return None
        except jwt.InvalidTokenError:
            return None

    def _generate_token(self, user_id: str, custom_claims: Dict) -> str:
        """Generate JWT token with expiration.

        Args:
            user_id: User identifier
            custom_claims: Additional claims to include

        Returns:
            JWT token string
        """
        expiration = datetime.utcnow() + timedelta(seconds=self.token_expiry)
        payload = {
            "user_id": user_id,
            "exp": expiration,
            **custom_claims
        }
        return jwt.encode(payload, self.secret_key, algorithm=self.algorithm)

    def _verify_password(self, plaintext: str, hashed: str) -> bool:
        """Verify password against hash.

        Args:
            plaintext: Password to verify
            hashed: Bcrypt password hash

        Returns:
            True if password matches hash
        """
        return bcrypt.checkpw(plaintext.encode('utf-8'), hashed.encode('utf-8'))
```

**Implementation Requirements**:
- Use PyJWT 2.8.0+ for token generation/validation
- Use bcrypt for password hashing (not implemented here, in PasswordHasher)
- Token expiry must be configurable
- Error messages must NOT reveal whether username or password was wrong
- All exceptions during validation must return None (fail-safe)

#### Tests Required

**File: tests/unit/test_authenticator.py**

```python
"""Unit tests for Authenticator class."""

import pytest
from unittest.mock import Mock, MagicMock
from datetime import datetime, timedelta
import jwt
from src.auth.authenticator import Authenticator, AuthResult

@pytest.fixture
def mock_user_repo():
    """Mock user repository."""
    repo = Mock()
    repo.get_by_username = Mock()
    repo.update_last_login = Mock()
    return repo

@pytest.fixture
def mock_user():
    """Mock user object."""
    user = Mock()
    user.id = 123
    user.username = "testuser"
    user.password_hash = "$2b$12$HASH_HERE"  # Bcrypt hash
    return user

@pytest.fixture
def authenticator():
    """Authenticator instance for testing."""
    return Authenticator(secret_key="test-secret-key", token_expiry_seconds=3600)

def test_successful_authentication(authenticator, mock_user_repo, mock_user):
    """Test successful authentication returns token."""
    mock_user_repo.get_by_username.return_value = mock_user
    authenticator._verify_password = Mock(return_value=True)

    result = authenticator.authenticate("testuser", "password123", mock_user_repo)

    assert result.success is True
    assert result.user_id == "123"
    assert result.session_token is not None
    assert result.error_message is None
    mock_user_repo.update_last_login.assert_called_once_with(123)

def test_authentication_fails_user_not_found(authenticator, mock_user_repo):
    """Test authentication fails when user doesn't exist."""
    mock_user_repo.get_by_username.return_value = None

    result = authenticator.authenticate("nonexistent", "password123", mock_user_repo)

    assert result.success is False
    assert result.user_id is None
    assert result.session_token is None
    assert result.error_message == "Invalid credentials"

def test_authentication_fails_wrong_password(authenticator, mock_user_repo, mock_user):
    """Test authentication fails with incorrect password."""
    mock_user_repo.get_by_username.return_value = mock_user
    authenticator._verify_password = Mock(return_value=False)

    result = authenticator.authenticate("testuser", "wrong_password", mock_user_repo)

    assert result.success is False
    assert result.user_id is None
    assert result.session_token is None
    assert result.error_message == "Invalid credentials"

def test_validate_token_success(authenticator):
    """Test token validation with valid token."""
    token = authenticator._generate_token("123", {"username": "testuser"})

    claims = authenticator.validate_token(token)

    assert claims is not None
    assert claims["user_id"] == "123"
    assert claims["username"] == "testuser"

def test_validate_token_expired(authenticator):
    """Test token validation fails for expired token."""
    authenticator.token_expiry = -1  # Token already expired
    token = authenticator._generate_token("123", {"username": "testuser"})

    claims = authenticator.validate_token(token)

    assert claims is None

def test_validate_token_invalid(authenticator):
    """Test token validation fails for malformed token."""
    claims = authenticator.validate_token("invalid.token.here")

    assert claims is None
```

**Test Coverage Requirements**:
- Minimum 90% code coverage
- All success paths tested
- All failure paths tested
- Edge cases (expired tokens, malformed input) tested

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- None (this is a foundational effort)

**Downstream Dependencies** (efforts that depend on this):
- Effort 2: API endpoints (needs Authenticator class)
- Effort 3: Protected routes (needs token validation)

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All tests passing (100% pass rate)
- [ ] Code coverage ≥90%
- [ ] No linting errors (flake8, mypy)
- [ ] Documentation complete (all public methods have docstrings)
- [ ] Line count within estimate (±15%)

---

### Effort 2: [Effort Name]

#### R213 Metadata

```json
{
  "effort_id": "EFFORT_[PHASE]_[WAVE]_002",
  "effort_name": "[Descriptive Name]",
  "branch_name": "effort/[phase]-[wave]-[descriptive-name]",
  "parent_wave": "[WAVE_N]",
  "parent_phase": "[PHASE_N]",
  "depends_on": ["EFFORT_[PHASE]_[WAVE]_001"],
  "estimated_lines": 350,
  "complexity": "medium",
  "can_parallelize": false
}
```

#### Scope

[Same structure as Effort 1]

#### Files to Create/Modify

[Same structure as Effort 1]

#### Exact Code Specifications

[Same structure as Effort 1]

#### Tests Required

[Same structure as Effort 1]

#### Dependencies

[Same structure as Effort 1]

#### Acceptance Criteria

[Same structure as Effort 1]

---

### Effort 3: [Effort Name]

[Same structure as Effort 1 and 2]

---

## Parallelization Strategy

### Parallel Group 1
- Effort 1 (no dependencies)
- Effort 4 (no dependencies)

### Sequential Group 2
- Effort 2 (depends on Effort 1)
- Effort 3 (depends on Effort 2)

### Parallel Group 3
- Effort 5 (depends on Effort 2)
- Effort 6 (depends on Effort 2)

**Rationale**: Efforts with no dependencies can execute in parallel. Dependent efforts must run sequentially.

---

## Wave Size Compliance

**Total Wave Lines**: [Sum of all effort lines]

**Size Limit**: 3500 lines (soft), 4000 lines (hard)

**Status**:
- [ ] Within soft limit (3500 lines)
- [ ] Within hard limit (4000 lines)
- [ ] Requires split plan (>4000 lines)

---

## Integration Strategy

1. **Effort 1** completes → Review → Merge to wave integration branch
2. **Effort 2** starts (depends on Effort 1) → Review → Merge
3. **Effort 3** starts (depends on Effort 2) → Review → Merge
4. **Effort 4-6** complete in parallel → Review → Merge
5. **Wave integration tests** run on wave branch
6. **Wave review** by Architect
7. **Merge to phase branch**

---

## Testing Strategy

### Unit Tests
- Each effort includes unit tests (specified above)
- Target: 90% code coverage per effort

### Integration Tests
- Run after all efforts merged to wave branch
- Test cross-effort interactions
- Verify end-to-end functionality

### Wave-Level Tests
```python
# tests/integration/test_wave_N_integration.py

def test_authentication_flow_end_to_end():
    """Test complete authentication flow across efforts."""
    # Test user registration (Effort 2)
    # Test user login (Effort 1)
    # Test protected endpoint access (Effort 3)
    # Verify token validation (Effort 1)
    pass

def test_error_handling_across_efforts():
    """Test error scenarios span multiple efforts."""
    # Test invalid credentials propagation
    # Test expired token handling
    # Test database connection failures
    pass
```

---

## Risk Mitigation

**High-Risk Efforts**:
- **Effort 2**: Database schema changes (risk: data loss)
  - Mitigation: Create backup before migration, test on staging first

**External Dependencies**:
- PyJWT library: Ensure version compatibility
- Bcrypt library: Performance may vary by platform

**Complexity Hotspots**:
- Token validation logic: High test coverage required
- Password hashing: Security-critical, audit carefully

---

## Next Steps

1. **Orchestrator** creates effort infrastructure (branches, worktrees)
2. **SW Engineers** spawn for parallel efforts
3. **Efforts** execute with continuous size monitoring
4. **Code Reviewers** review after each effort completion
5. **Wave integration** after all efforts merged
6. **Architect review** before phase integration

**Note**: This document provides **exact specifications** for implementation. SW Engineers should follow these precisely while adapting to discoveries during coding.
