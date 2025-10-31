# Wave N Architecture Plan

**Wave**: [Wave Number/Name]
**Phase**: [Parent Phase]
**Created**: [Date]
**Architect**: [Agent Name]
**Fidelity Level**: **REAL CODE** (concrete designs, working examples)

---

## Adaptation Notes

### Lessons from Previous Waves

**What Worked Well**:
- [Specific pattern/approach that succeeded]
- [Code structure that was maintainable]
- [Testing strategy that caught bugs early]

**What to Improve**:
- [Pattern that caused issues]
- [Code duplication to eliminate]
- [Performance bottleneck to address]

### Design Changes from Phase Architecture

**Refinements**:
- [How actual implementation informed better design]
- [Constraints discovered during previous wave]
- [New requirements that emerged]

**Code Evolution**:
```python
# Previous wave approach:
def old_approach():
    # Less maintainable pattern
    pass

# This wave's improvement:
def new_approach():
    # More maintainable pattern
    pass
```

---

## Concrete Designs

### Core Classes/Functions

**User Authentication Module**:

```python
# File: src/auth/authenticator.py

from typing import Optional, Dict
from dataclasses import dataclass

@dataclass
class AuthResult:
    """Result of an authentication attempt."""
    success: bool
    user_id: Optional[str]
    session_token: Optional[str]
    error_message: Optional[str]

class Authenticator:
    """Handles user authentication with JWT tokens."""

    def __init__(self, secret_key: str, token_expiry_seconds: int = 3600):
        """Initialize authenticator with secret key.

        Args:
            secret_key: Secret key for signing tokens
            token_expiry_seconds: Token validity duration (default 1 hour)
        """
        self.secret_key = secret_key
        self.token_expiry = token_expiry_seconds

    def authenticate(self, username: str, password: str) -> AuthResult:
        """Authenticate user with username and password.

        Args:
            username: User's username
            password: User's password (will be hashed)

        Returns:
            AuthResult with success status and session token if successful
        """
        # Implementation will be added in effort definition
        pass

    def validate_token(self, token: str) -> Optional[Dict]:
        """Validate a session token and extract user data.

        Args:
            token: JWT session token

        Returns:
            User data dictionary if token valid, None otherwise
        """
        # Implementation will be added in effort definition
        pass
```

### Database Models

```python
# File: src/models/user.py

from sqlalchemy import Column, Integer, String, DateTime, Boolean
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime

Base = declarative_base()

class User(Base):
    """User model with authentication fields."""

    __tablename__ = 'users'

    id = Column(Integer, primary_key=True, autoincrement=True)
    username = Column(String(50), unique=True, nullable=False, index=True)
    email = Column(String(100), unique=True, nullable=False, index=True)
    password_hash = Column(String(255), nullable=False)
    is_active = Column(Boolean, default=True, nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)
    last_login = Column(DateTime, nullable=True)

    def __repr__(self):
        return f"<User(id={self.id}, username='{self.username}')>"

    def to_dict(self) -> dict:
        """Convert user to dictionary for API responses."""
        return {
            'id': self.id,
            'username': self.username,
            'email': self.email,
            'is_active': self.is_active,
            'created_at': self.created_at.isoformat(),
            'last_login': self.last_login.isoformat() if self.last_login else None
        }
```

### API Endpoints

```python
# File: src/api/auth_routes.py

from fastapi import APIRouter, HTTPException, Depends, status
from pydantic import BaseModel, EmailStr
from typing import Optional

router = APIRouter(prefix="/auth", tags=["authentication"])

class LoginRequest(BaseModel):
    """Login request schema."""
    username: str
    password: str

class LoginResponse(BaseModel):
    """Login response schema."""
    access_token: str
    token_type: str = "bearer"
    user_id: str

class RegisterRequest(BaseModel):
    """User registration schema."""
    username: str
    email: EmailStr
    password: str

@router.post("/login", response_model=LoginResponse, status_code=status.HTTP_200_OK)
async def login(request: LoginRequest, auth: Authenticator = Depends(get_authenticator)):
    """Authenticate user and return session token.

    Args:
        request: Login credentials
        auth: Authenticator instance (injected)

    Returns:
        LoginResponse with access token

    Raises:
        HTTPException: 401 if authentication fails
    """
    result = auth.authenticate(request.username, request.password)

    if not result.success:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail=result.error_message or "Invalid credentials"
        )

    return LoginResponse(
        access_token=result.session_token,
        user_id=result.user_id
    )

@router.post("/register", status_code=status.HTTP_201_CREATED)
async def register(request: RegisterRequest, db: Session = Depends(get_db)):
    """Register a new user.

    Args:
        request: Registration data
        db: Database session (injected)

    Returns:
        Success message

    Raises:
        HTTPException: 400 if username/email already exists
    """
    # Implementation details in effort definition
    pass
```

---

## Working Usage Examples

### Authentication Flow Example

```python
# Example usage in application code

from src.auth.authenticator import Authenticator, AuthResult
from src.models.user import User

# Initialize authenticator
auth = Authenticator(secret_key="your-secret-key-here", token_expiry_seconds=3600)

# Authenticate user
result: AuthResult = auth.authenticate("john_doe", "secure_password")

if result.success:
    print(f"Login successful! Token: {result.session_token}")
    # Use token for subsequent requests
else:
    print(f"Login failed: {result.error_message}")

# Validate token in protected routes
user_data = auth.validate_token(result.session_token)
if user_data:
    print(f"Valid session for user: {user_data['user_id']}")
else:
    print("Invalid or expired token")
```

### Database Operations Example

```python
# Example database operations

from sqlalchemy.orm import Session
from src.models.user import User
from datetime import datetime

def create_user(db: Session, username: str, email: str, password_hash: str) -> User:
    """Create a new user in the database."""
    user = User(
        username=username,
        email=email,
        password_hash=password_hash,
        created_at=datetime.utcnow()
    )
    db.add(user)
    db.commit()
    db.refresh(user)
    return user

def get_user_by_username(db: Session, username: str) -> Optional[User]:
    """Retrieve user by username."""
    return db.query(User).filter(User.username == username).first()

def update_last_login(db: Session, user_id: int) -> None:
    """Update user's last login timestamp."""
    user = db.query(User).filter(User.id == user_id).first()
    if user:
        user.last_login = datetime.utcnow()
        db.commit()
```

### API Client Example

```python
# Example API client usage

import requests

BASE_URL = "http://localhost:8000"

# Register new user
register_response = requests.post(
    f"{BASE_URL}/auth/register",
    json={
        "username": "john_doe",
        "email": "john@example.com",
        "password": "SecureP@ssw0rd"
    }
)

if register_response.status_code == 201:
    print("Registration successful!")

# Login
login_response = requests.post(
    f"{BASE_URL}/auth/login",
    json={
        "username": "john_doe",
        "password": "SecureP@ssw0rd"
    }
)

if login_response.status_code == 200:
    token_data = login_response.json()
    access_token = token_data["access_token"]

    # Use token for authenticated requests
    headers = {"Authorization": f"Bearer {access_token}"}
    profile_response = requests.get(
        f"{BASE_URL}/users/me",
        headers=headers
    )
    print(f"User profile: {profile_response.json()}")
```

---

## Interface Definitions

### Authentication Service Interface

```python
# File: src/interfaces/auth_service.py

from typing import Protocol, Optional, Dict
from src.auth.authenticator import AuthResult

class IAuthenticationService(Protocol):
    """Protocol defining authentication service interface."""

    def authenticate(self, username: str, password: str) -> AuthResult:
        """Authenticate user credentials."""
        ...

    def validate_token(self, token: str) -> Optional[Dict]:
        """Validate session token."""
        ...

    def generate_token(self, user_id: str, claims: Dict) -> str:
        """Generate JWT token with custom claims."""
        ...

    def revoke_token(self, token: str) -> bool:
        """Revoke a session token."""
        ...
```

### Repository Interface

```python
# File: src/interfaces/user_repository.py

from typing import Protocol, Optional, List
from src.models.user import User

class IUserRepository(Protocol):
    """Protocol defining user data access interface."""

    def create(self, user: User) -> User:
        """Create a new user."""
        ...

    def get_by_id(self, user_id: int) -> Optional[User]:
        """Retrieve user by ID."""
        ...

    def get_by_username(self, username: str) -> Optional[User]:
        """Retrieve user by username."""
        ...

    def get_by_email(self, email: str) -> Optional[User]:
        """Retrieve user by email."""
        ...

    def update(self, user: User) -> User:
        """Update existing user."""
        ...

    def delete(self, user_id: int) -> bool:
        """Delete user by ID."""
        ...

    def list_active_users(self) -> List[User]:
        """List all active users."""
        ...
```

---

## Testing Strategy

### Unit Test Example

```python
# File: tests/unit/test_authenticator.py

import pytest
from src.auth.authenticator import Authenticator, AuthResult

@pytest.fixture
def authenticator():
    return Authenticator(secret_key="test-secret", token_expiry_seconds=3600)

def test_successful_authentication(authenticator, mock_user_db):
    """Test successful user authentication."""
    result = authenticator.authenticate("valid_user", "correct_password")

    assert result.success is True
    assert result.user_id is not None
    assert result.session_token is not None
    assert result.error_message is None

def test_failed_authentication_invalid_password(authenticator, mock_user_db):
    """Test authentication fails with wrong password."""
    result = authenticator.authenticate("valid_user", "wrong_password")

    assert result.success is False
    assert result.user_id is None
    assert result.session_token is None
    assert "Invalid credentials" in result.error_message

def test_token_validation_success(authenticator):
    """Test valid token returns user data."""
    # Create a valid token
    auth_result = authenticator.authenticate("valid_user", "correct_password")
    token = auth_result.session_token

    # Validate it
    user_data = authenticator.validate_token(token)

    assert user_data is not None
    assert user_data['user_id'] == auth_result.user_id
```

### Integration Test Example

```python
# File: tests/integration/test_auth_api.py

import pytest
from fastapi.testclient import TestClient
from src.main import app

@pytest.fixture
def client():
    return TestClient(app)

def test_login_endpoint_success(client, test_user):
    """Test /auth/login endpoint with valid credentials."""
    response = client.post(
        "/auth/login",
        json={"username": "testuser", "password": "testpass"}
    )

    assert response.status_code == 200
    data = response.json()
    assert "access_token" in data
    assert data["token_type"] == "bearer"
    assert "user_id" in data

def test_login_endpoint_invalid_credentials(client):
    """Test /auth/login endpoint with invalid credentials."""
    response = client.post(
        "/auth/login",
        json={"username": "nonexistent", "password": "wrong"}
    )

    assert response.status_code == 401
    assert "Invalid credentials" in response.json()["detail"]
```

---

## Dependencies

### External Libraries

```python
# requirements.txt additions for this wave

fastapi==0.104.1
sqlalchemy==2.0.23
pyjwt==2.8.0
bcrypt==4.1.1
pydantic[email]==2.5.0
python-multipart==0.0.6
```

### Internal Dependencies

- Previous Wave: User model foundation
- Previous Wave: Database connection setup
- Previous Wave: Configuration management

---

## Next Steps (Wave Implementation Planning)

The wave implementation plan will provide:
- **Exact file lists** for each effort
- **Detailed code specifications** for each function
- **R213 metadata** (effort IDs, branch names, dependencies)
- **Line count estimates** per effort
- **Test coverage requirements**

**Note**: This document provides **real, working code examples**. The implementation plan will break this into specific efforts with exact specifications.
