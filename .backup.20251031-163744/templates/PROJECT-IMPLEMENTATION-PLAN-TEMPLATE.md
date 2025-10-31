# PROJECT IMPLEMENTATION PLAN: [Your Project Name]

**Created**: [Date - YYYY-MM-DD]
**Project Type**: [Web Application / API Service / CLI Tool / Library / Mobile App]
**Target Platforms**: [Linux / macOS / Windows / Cloud / etc.]
**Repository**: [URL if exists, or "New Repository"]

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Project Overview](#project-overview)
3. [Technology Stack](#technology-stack)
4. [Architecture Overview](#architecture-overview)
5. [Success Criteria](#success-criteria)
6. [Implementation Phases](#implementation-phases)
7. [Testing Strategy](#testing-strategy)
8. [Deployment Strategy](#deployment-strategy)
9. [Risk Management](#risk-management)
10. [Appendices](#appendices)

---

## Executive Summary

### Project Vision
[2-3 sentences describing the project's purpose and value proposition]

**Example**:
> TaskFlow is a modern task management system designed for remote teams. It provides real-time collaboration, intuitive task organization, and seamless integration with existing workflows. The system aims to reduce task management overhead by 50% while improving team visibility.

### Key Metrics
- **Total Phases**: [Number]
- **Total Waves**: [Number]
- **Estimated Efforts**: [Number]
- **Effort Size Limit**: 800 lines (per R220/R221)
- **Estimated Duration**: [X weeks]
- **Target Test Coverage**: 80%+

### Quick Start Guide
This implementation plan is designed for Software Factory 3.0. To begin:
1. Ensure SF 3.0 is initialized: `bash utilities/init-software-factory-v3.sh --project-name [name]`
2. Place this file in project root as `PROJECT-IMPLEMENTATION-PLAN.md`
3. Run: `/init-software-factory` (for interactive setup)
4. Run: `/continue-software-factory` (to start orchestration)

---

## Project Overview

### Problem Statement
[Describe the problem this project solves. Be specific about pain points.]

**Example**:
> Remote teams struggle with task visibility and coordination across time zones. Current solutions are either too complex (enterprise tools) or too simple (spreadsheets), lacking real-time updates and intuitive organization. This creates communication overhead and missed deadlines.

### Solution Approach
[Describe how this project addresses the problem]

**Example**:
> TaskFlow provides a lightweight, real-time task management system that:
> - Uses WebSocket connections for instant updates across all clients
> - Organizes tasks in flexible hierarchies (projects → tasks → subtasks)
> - Integrates with common tools via webhooks and APIs
> - Requires minimal training due to intuitive UI design

### Target Users
[Define who will use this system and their requirements]

**Example**:
> - **Primary**: Remote development teams (5-20 members)
> - **Secondary**: Project managers tracking multiple projects
> - **Requirements**: Cross-platform access, mobile support, API integration capabilities

### Key Features
1. **[Feature 1]**: [Brief description and value]
2. **[Feature 2]**: [Brief description and value]
3. **[Feature 3]**: [Brief description and value]
4. **[Feature 4]**: [Brief description and value]

**Example**:
1. **Real-time Task Updates**: WebSocket-based instant synchronization across all clients
2. **Flexible Task Hierarchy**: Support for projects, tasks, subtasks, and tags
3. **User Authentication**: JWT-based secure authentication with role-based access control
4. **API Integration**: RESTful API and webhooks for third-party integrations

### Non-Goals (Out of Scope)
- [What this project will NOT do]
- [Features explicitly excluded]

**Example**:
- No time tracking functionality (focus on task management only)
- No built-in chat/messaging (integrate with existing tools instead)
- No mobile native apps (web-first, mobile-responsive)

---

## Technology Stack

### Backend
- **Language**: [e.g., Python 3.11+]
- **Framework**: [e.g., FastAPI 0.104+]
- **Database**: [e.g., PostgreSQL 15+]
- **ORM**: [e.g., SQLAlchemy 2.0]
- **Real-time**: [e.g., WebSockets via FastAPI]
- **Testing**: [e.g., Pytest, pytest-asyncio]

### Frontend
- **Language**: [e.g., TypeScript 5+]
- **Framework**: [e.g., Next.js 14 (App Router)]
- **UI Library**: [e.g., React 18]
- **Styling**: [e.g., TailwindCSS 3]
- **State Management**: [e.g., React Context + useReducer]
- **Testing**: [e.g., Jest, React Testing Library, Playwright]

### Infrastructure
- **Build System**: [e.g., npm/pnpm for frontend, Poetry for backend]
- **Containerization**: [e.g., Docker, Docker Compose]
- **CI/CD**: [e.g., GitHub Actions]
- **Deployment**: [e.g., Docker containers on AWS/GCP/Azure]
- **Monitoring**: [e.g., Prometheus, Grafana]

### Development Tools
- **Version Control**: Git
- **Code Quality**: [e.g., ESLint, Prettier, Black, mypy]
- **Pre-commit Hooks**: [e.g., husky, pre-commit]
- **Documentation**: [e.g., OpenAPI/Swagger, Storybook]

---

## Architecture Overview

### System Architecture Pattern
[Describe the chosen architecture pattern]

**Example**:
> Three-tier architecture with:
> - **Presentation Layer**: Next.js frontend (server-rendered React)
> - **Application Layer**: FastAPI backend (RESTful API + WebSocket)
> - **Data Layer**: PostgreSQL database with SQLAlchemy ORM

### Component Diagram
```
[User Browser]
      ↓
[Next.js Frontend (Port 3000)]
      ↓ (HTTP/WS)
[FastAPI Backend (Port 8000)]
      ↓
[PostgreSQL Database (Port 5432)]
```

### Key Components

#### 1. Backend API Service
**Responsibilities**:
- RESTful API endpoints (authentication, task CRUD)
- WebSocket connection management
- Business logic validation
- Database operations via ORM

**Key Modules**:
- `api/routes/` - API endpoints
- `api/websocket/` - WebSocket handlers
- `models/` - SQLAlchemy models
- `services/` - Business logic layer
- `auth/` - Authentication & authorization

#### 2. Frontend Web Application
**Responsibilities**:
- User interface rendering
- State management (local + server sync)
- WebSocket client connection
- API consumption

**Key Modules**:
- `app/` - Next.js app router pages
- `components/` - Reusable UI components
- `hooks/` - Custom React hooks
- `services/` - API client functions
- `contexts/` - Global state management

#### 3. Database Schema
**Core Entities**:
- Users (authentication, profiles)
- Projects (top-level organization)
- Tasks (individual work items)
- Tags (flexible categorization)
- Activity Log (audit trail)

### Data Flow

#### Task Creation Flow
```
User Action (Frontend)
  → API Request (POST /api/tasks)
  → Backend Validation
  → Database Insert
  → WebSocket Broadcast
  → All Connected Clients Update
```

### External Dependencies
- **Authentication**: [e.g., JWT tokens (self-managed)]
- **Email**: [e.g., SendGrid for notifications]
- **Storage**: [e.g., S3 for file attachments if needed]
- **Monitoring**: [e.g., Sentry for error tracking]

---

## Success Criteria

### Functional Requirements
- [ ] Users can register and authenticate securely
- [ ] Users can create, read, update, delete tasks
- [ ] Real-time updates work across all connected clients
- [ ] Task hierarchy (projects → tasks → subtasks) functions correctly
- [ ] API endpoints respond within 200ms (p95)
- [ ] WebSocket connections remain stable under load

### Quality Requirements
- [ ] All efforts ≤800 lines (measured by `tools/line-counter.sh`)
- [ ] Test coverage ≥80% (backend and frontend)
- [ ] All code passes code review
- [ ] Zero critical security vulnerabilities
- [ ] Documentation complete (API docs, setup guide, architecture)

### Performance Requirements
- [ ] API response time <200ms (p95)
- [ ] WebSocket message latency <100ms
- [ ] Support 100 concurrent WebSocket connections
- [ ] Database queries optimized (indexed, <50ms)

### Deployment Requirements
- [ ] Docker containers build successfully
- [ ] Docker Compose orchestration works locally
- [ ] Production deployment automated
- [ ] Health checks and monitoring operational
- [ ] Rollback procedure documented and tested

---

## Implementation Phases

### Phase 1: Foundation

**Goal**: Establish backend API, database, and authentication

**Duration Estimate**: 2 weeks

#### Wave 1.1: Database & Authentication

**Purpose**: Set up database models and user authentication system

**Efforts**:

##### Effort 1.1.1 - Database Models & Migrations
- **Description**: Create SQLAlchemy models for User, Project, Task, Tag
- **Estimated Lines**: ~350 lines
- **Files**:
  - `backend/models/user.py`
  - `backend/models/project.py`
  - `backend/models/task.py`
  - `backend/models/base.py`
  - `backend/alembic/versions/001_initial.py`
- **Definition of Done**:
  - [ ] All models defined with relationships
  - [ ] Alembic migrations created and tested
  - [ ] Model validation rules implemented
  - [ ] Database schema creates successfully

##### Effort 1.1.2 - User Authentication
- **Description**: Implement JWT-based authentication (register, login, logout)
- **Estimated Lines**: ~400 lines
- **Files**:
  - `backend/auth/jwt.py`
  - `backend/auth/password.py`
  - `backend/api/routes/auth.py`
  - `backend/middleware/auth.py`
- **Definition of Done**:
  - [ ] User registration endpoint working
  - [ ] Login returns valid JWT token
  - [ ] Password hashing with bcrypt
  - [ ] Protected endpoints verify JWT
  - [ ] Token refresh mechanism working

##### Effort 1.1.3 - Core API Endpoints (Tasks)
- **Description**: CRUD endpoints for tasks
- **Estimated Lines**: ~350 lines
- **Files**:
  - `backend/api/routes/tasks.py`
  - `backend/services/task_service.py`
  - `backend/schemas/task.py`
- **Definition of Done**:
  - [ ] GET /api/tasks (list with filters)
  - [ ] POST /api/tasks (create)
  - [ ] PUT /api/tasks/{id} (update)
  - [ ] DELETE /api/tasks/{id} (delete)
  - [ ] Input validation with Pydantic
  - [ ] User-specific task filtering

#### Wave 1.2: Frontend Foundation

**Purpose**: Set up Next.js application with authentication UI

**Efforts**:

##### Effort 1.2.1 - Next.js App Setup
- **Description**: Initialize Next.js with TailwindCSS and app router
- **Estimated Lines**: ~250 lines
- **Files**:
  - `frontend/app/layout.tsx`
  - `frontend/app/page.tsx`
  - `frontend/tailwind.config.js`
  - `frontend/next.config.js`
  - `frontend/tsconfig.json`
- **Definition of Done**:
  - [ ] Next.js 14 app router configured
  - [ ] TailwindCSS integrated and working
  - [ ] Environment variables configured
  - [ ] Basic layout and homepage render
  - [ ] TypeScript strict mode enabled

##### Effort 1.2.2 - Authentication UI
- **Description**: Login, register pages and auth context
- **Estimated Lines**: ~400 lines
- **Files**:
  - `frontend/app/login/page.tsx`
  - `frontend/app/register/page.tsx`
  - `frontend/contexts/AuthContext.tsx`
  - `frontend/components/ProtectedRoute.tsx`
- **Definition of Done**:
  - [ ] Login form with validation
  - [ ] Register form with validation
  - [ ] Auth context manages JWT token
  - [ ] Protected route wrapper component
  - [ ] Auto-redirect on auth state change

##### Effort 1.2.3 - Task List UI
- **Description**: Task list display and basic CRUD UI
- **Estimated Lines**: ~450 lines
- **Files**:
  - `frontend/app/tasks/page.tsx`
  - `frontend/components/TaskList.tsx`
  - `frontend/components/TaskForm.tsx`
  - `frontend/services/taskApi.ts`
- **Definition of Done**:
  - [ ] Task list fetches from API
  - [ ] Create task form working
  - [ ] Edit task inline
  - [ ] Delete with confirmation modal
  - [ ] Loading and error states

#### Wave 1.3: Real-time Features

**Purpose**: Add WebSocket support for real-time updates

**Efforts**:

##### Effort 1.3.1 - WebSocket Backend
- **Description**: WebSocket connection handling and broadcasting
- **Estimated Lines**: ~300 lines
- **Files**:
  - `backend/websocket/manager.py`
  - `backend/websocket/handlers.py`
  - `backend/api/routes/ws.py`
- **Definition of Done**:
  - [ ] WebSocket endpoint accepts connections
  - [ ] Connection manager tracks active clients
  - [ ] Broadcast task updates to all clients
  - [ ] Handle client disconnections gracefully

##### Effort 1.3.2 - WebSocket Frontend
- **Description**: WebSocket client integration in frontend
- **Estimated Lines**: ~350 lines
- **Files**:
  - `frontend/hooks/useWebSocket.ts`
  - `frontend/contexts/WebSocketContext.tsx`
  - `frontend/services/websocketClient.ts`
- **Definition of Done**:
  - [ ] WebSocket connection established
  - [ ] Automatic reconnection on disconnect
  - [ ] Task list updates in real-time
  - [ ] Connection status indicator in UI

---

### Phase 2: Advanced Features

**Goal**: Add project organization, tags, and search

**Duration Estimate**: 2 weeks

#### Wave 2.1: Projects & Organization

**Purpose**: Multi-project support and task hierarchy

**Efforts**:

##### Effort 2.1.1 - Project API
- **Description**: CRUD endpoints for projects
- **Estimated Lines**: ~300 lines
- **Files**:
  - `backend/api/routes/projects.py`
  - `backend/services/project_service.py`
  - `backend/schemas/project.py`
- **Definition of Done**:
  - [ ] Project CRUD endpoints working
  - [ ] Tasks associated with projects
  - [ ] Project-level permissions
  - [ ] Project statistics (task counts, etc.)

##### Effort 2.1.2 - Project UI
- **Description**: Project management UI
- **Estimated Lines**: ~400 lines
- **Files**:
  - `frontend/app/projects/page.tsx`
  - `frontend/components/ProjectList.tsx`
  - `frontend/components/ProjectCard.tsx`
- **Definition of Done**:
  - [ ] Project list view
  - [ ] Project detail page
  - [ ] Create/edit project forms
  - [ ] Task filter by project

##### Effort 2.1.3 - Task Hierarchy
- **Description**: Subtask support
- **Estimated Lines**: ~350 lines
- **Files**:
  - `backend/models/subtask.py`
  - `backend/api/routes/subtasks.py`
  - `frontend/components/SubtaskList.tsx`
- **Definition of Done**:
  - [ ] Subtask model and API
  - [ ] Parent-child relationships
  - [ ] Nested task display UI
  - [ ] Drag-and-drop reordering

#### Wave 2.2: Search & Filtering

**Purpose**: Advanced search and filtering capabilities

**Efforts**:

##### Effort 2.2.1 - Backend Search
- **Description**: Full-text search and filtering
- **Estimated Lines**: ~300 lines
- **Files**:
  - `backend/services/search_service.py`
  - `backend/api/routes/search.py`
- **Definition of Done**:
  - [ ] Full-text search on tasks
  - [ ] Filter by status, priority, assignee
  - [ ] Date range filtering
  - [ ] Sort options (date, priority, etc.)

##### Effort 2.2.2 - Search UI
- **Description**: Search interface and filters
- **Estimated Lines**: ~350 lines
- **Files**:
  - `frontend/components/SearchBar.tsx`
  - `frontend/components/FilterPanel.tsx`
- **Definition of Done**:
  - [ ] Search input with autocomplete
  - [ ] Filter dropdowns for all criteria
  - [ ] Active filter badges
  - [ ] Clear filters button

##### Effort 2.2.3 - Tags System
- **Description**: Tagging system for tasks
- **Estimated Lines**: ~300 lines
- **Files**:
  - `backend/models/tag.py`
  - `backend/api/routes/tags.py`
  - `frontend/components/TagSelector.tsx`
- **Definition of Done**:
  - [ ] Tag CRUD operations
  - [ ] Many-to-many task-tag relationship
  - [ ] Tag autocomplete in UI
  - [ ] Filter tasks by tag

---

### Phase 3: Production Ready

**Goal**: Testing, deployment, monitoring, and documentation

**Duration Estimate**: 1 week

#### Wave 3.1: Testing & Quality

**Purpose**: Comprehensive test coverage and quality assurance

**Efforts**:

##### Effort 3.1.1 - Backend Tests
- **Description**: Comprehensive backend test suite
- **Estimated Lines**: ~500 lines
- **Files**:
  - `backend/tests/test_auth.py`
  - `backend/tests/test_tasks.py`
  - `backend/tests/test_projects.py`
  - `backend/tests/test_websocket.py`
- **Definition of Done**:
  - [ ] Unit tests for all services
  - [ ] API endpoint integration tests
  - [ ] WebSocket connection tests
  - [ ] Coverage ≥80%

##### Effort 3.1.2 - Frontend Tests
- **Description**: Frontend component and integration tests
- **Estimated Lines**: ~600 lines
- **Files**:
  - `frontend/__tests__/components/TaskList.test.tsx`
  - `frontend/__tests__/pages/tasks.test.tsx`
  - `frontend/__tests__/e2e/auth.test.ts`
- **Definition of Done**:
  - [ ] Component unit tests (Jest + RTL)
  - [ ] Integration tests for pages
  - [ ] E2E tests with Playwright
  - [ ] Coverage ≥80%

##### Effort 3.1.3 - Performance Testing
- **Description**: Load testing and performance validation
- **Estimated Lines**: ~200 lines
- **Files**:
  - `tests/load/locustfile.py`
  - `tests/performance/benchmark.py`
- **Definition of Done**:
  - [ ] Load test with 100 concurrent users
  - [ ] API response time validation
  - [ ] WebSocket connection stability test
  - [ ] Database query performance checks

#### Wave 3.2: Deployment & Operations

**Purpose**: Deployment automation and operational monitoring

**Efforts**:

##### Effort 3.2.1 - Docker & Deployment
- **Description**: Containerization and deployment configs
- **Estimated Lines**: ~300 lines
- **Files**:
  - `backend/Dockerfile`
  - `frontend/Dockerfile`
  - `docker-compose.yml`
  - `docker-compose.prod.yml`
  - `.github/workflows/deploy.yml`
- **Definition of Done**:
  - [ ] Backend Docker image builds
  - [ ] Frontend Docker image builds
  - [ ] Docker Compose for local dev
  - [ ] Production deployment config
  - [ ] CI/CD pipeline for deployment

##### Effort 3.2.2 - Monitoring & Logging
- **Description**: Application monitoring and logging setup
- **Estimated Lines**: ~250 lines
- **Files**:
  - `backend/middleware/logging.py`
  - `backend/monitoring/metrics.py`
  - `docker-compose.monitoring.yml`
- **Definition of Done**:
  - [ ] Structured logging configured
  - [ ] Prometheus metrics exposed
  - [ ] Health check endpoints
  - [ ] Error tracking (Sentry or similar)

##### Effort 3.2.3 - Documentation
- **Description**: Complete project documentation
- **Estimated Lines**: ~300 lines (markdown)
- **Files**:
  - `README.md`
  - `docs/SETUP.md`
  - `docs/API.md`
  - `docs/ARCHITECTURE.md`
  - `docs/DEPLOYMENT.md`
- **Definition of Done**:
  - [ ] README with overview and quick start
  - [ ] Setup guide for local development
  - [ ] API documentation (OpenAPI/Swagger)
  - [ ] Architecture diagrams and explanations
  - [ ] Deployment runbook

---

## Testing Strategy

### Test Pyramid

```
           E2E Tests (10%)
         /              \
    Integration Tests (30%)
   /                      \
  Unit Tests (60%)
```

### Backend Testing

**Unit Tests** (Pytest):
- Services layer logic
- Model validation
- Utility functions
- Coverage target: 80%+

**Integration Tests** (Pytest + TestClient):
- API endpoint functionality
- Database operations
- Authentication flow
- WebSocket connections

**Load Tests** (Locust):
- 100 concurrent users
- API response times
- WebSocket stability

### Frontend Testing

**Component Tests** (Jest + React Testing Library):
- Individual component rendering
- User interactions
- State management
- Coverage target: 80%+

**Integration Tests** (Jest):
- Page-level functionality
- API integration
- Context providers

**E2E Tests** (Playwright):
- Critical user journeys
- Authentication flows
- Task CRUD operations
- Real-time updates

### Test Automation

**Pre-commit Hooks**:
- Linting (ESLint, Black)
- Type checking (TypeScript, mypy)
- Unit tests (fast subset)

**CI Pipeline** (GitHub Actions):
- All unit tests
- Integration tests
- Build verification
- Test coverage reports

---

## Deployment Strategy

### Local Development

```bash
# Start all services
docker-compose up

# Backend: http://localhost:8000
# Frontend: http://localhost:3000
# Database: localhost:5432
```

### Staging Environment

- Auto-deploy from `develop` branch
- Full integration testing
- Performance validation
- Security scanning

### Production Environment

**Infrastructure**:
- Cloud platform: [AWS/GCP/Azure]
- Container orchestration: [Docker/Kubernetes]
- Database: Managed PostgreSQL
- CDN: [CloudFlare/CloudFront]

**Deployment Process**:
1. Merge to `main` triggers CI/CD
2. Run full test suite
3. Build Docker images
4. Push to container registry
5. Deploy to production
6. Run smoke tests
7. Monitor metrics

**Rollback Procedure**:
- Keep previous 5 images
- One-command rollback script
- Database migration rollback plan

---

## Risk Management

### Technical Risks

#### Risk 1: WebSocket Connection Stability
- **Probability**: Medium
- **Impact**: High (core feature failure)
- **Mitigation**:
  - Implement automatic reconnection with exponential backoff
  - Add WebSocket connection health checks
  - Fallback to polling if WebSocket fails
  - Comprehensive connection error handling

#### Risk 2: Database Performance Under Load
- **Probability**: Medium
- **Impact**: Medium (slow response times)
- **Mitigation**:
  - Database indexing on commonly queried fields
  - Query optimization and monitoring
  - Connection pooling (SQLAlchemy)
  - Caching layer (Redis) if needed

#### Risk 3: Effort Size Overflow
- **Probability**: Medium
- **Impact**: High (violates R220/R221)
- **Mitigation**:
  - Use `tools/line-counter.sh` frequently during development
  - Plan for aggressive splitting in complex features
  - Code Reviewer monitors size proactively
  - Split plan ready before 600 lines reached

### Process Risks

#### Risk 4: Integration Conflicts
- **Probability**: Low
- **Impact**: Medium (delays integration)
- **Mitigation**:
  - Feature branches kept small and focused
  - Regular merges to wave branches
  - Architect reviews wave-level design
  - Clear interface contracts between efforts

#### Risk 5: Test Coverage Gaps
- **Probability**: Medium
- **Impact**: Medium (quality issues)
- **Mitigation**:
  - TDD approach (tests before code)
  - Coverage reports in CI
  - Code Reviewer validates test quality
  - Critical paths have E2E tests

### External Dependency Risks

#### Risk 6: Third-party Service Outages
- **Probability**: Low
- **Impact**: Low (nice-to-have features)
- **Mitigation**:
  - Core features work without external services
  - Graceful degradation for non-critical features
  - Retry logic with exponential backoff
  - Status page monitoring

---

## Appendices

### A. Glossary

- **Effort**: A unit of work ≤800 lines, completed by one SW Engineer
- **Wave**: A group of related efforts that integrate together
- **Phase**: A major project milestone containing multiple waves
- **Integration Container**: A wave or phase integration point
- **Convergence**: The iterate-fix-reintegrate cycle until all tests pass
- **R220/R221**: Rules defining 800 line effort size limit

### B. Sizing Estimates

**Total Estimated Lines**: ~6,000 lines
**Total Efforts**: 21 efforts
**Average Effort Size**: ~285 lines
**Largest Effort**: 600 lines (Frontend Tests)
**All efforts within limit**: ✅ (target: ≤800 lines)

### C. Development Timeline

**Phase 1**: 2 weeks (3 waves, 9 efforts)
**Phase 2**: 2 weeks (2 waves, 6 efforts)
**Phase 3**: 1 week (2 waves, 6 efforts)
**Total Duration**: ~5 weeks

**Parallelization Opportunities**:
- Wave 1.1: Efforts can run in parallel (independent backend endpoints)
- Wave 1.2: Efforts can run in parallel (independent UI components)
- Wave 2.1 & 2.2: Some efforts can overlap

### D. References

- **Architecture Docs**: [To be created by Architect agent]
- **API Specification**: [Generated from OpenAPI schema]
- **State Machine**: `state-machines/software-factory-3.0-state-machine.json`
- **Rules Library**: `rule-library/` (particularly R220, R221, R287, R373, R420)
- **Line Counter**: `tools/line-counter.sh` (for size measurement)

### E. Assumptions

1. **Environment**: Developers have Docker, Node.js 18+, Python 3.11+ installed
2. **Resources**: Sufficient API credits for full development cycle (~$200-300)
3. **Access**: Write access to target repository
4. **Skills**: Team familiar with FastAPI, Next.js, or willing to learn
5. **Timeline**: No hard deadlines, quality over speed
6. **Scope**: Features listed are complete scope (no scope creep)

### F. Configuration Files

**Backend**:
- `backend/pyproject.toml` - Python dependencies (Poetry)
- `backend/alembic.ini` - Database migration config
- `backend/.env.example` - Environment variables template

**Frontend**:
- `frontend/package.json` - Node dependencies
- `frontend/.env.local.example` - Environment variables template
- `frontend/tsconfig.json` - TypeScript configuration

**Infrastructure**:
- `docker-compose.yml` - Local development
- `docker-compose.prod.yml` - Production deployment
- `.github/workflows/` - CI/CD pipelines

---

## How to Use This Plan

### For Orchestrator Agent
1. Read this plan completely
2. Validate structure (phases → waves → efforts)
3. Spawn Architect for master architecture creation
4. Execute phases sequentially, waves sequentially, efforts in parallel where possible
5. Track progress in `orchestrator-state-v3.json`

### For Architect Agent
1. Use this plan to create `PROJECT-ARCHITECTURE.md`
2. Validate technology choices
3. Review wave-level architecture
4. Ensure consistency across phases

### For Code Reviewer Agent
1. Create effort plans based on wave descriptions
2. Execute R420 protocol before each effort plan
3. Validate effort size estimates (≤800 lines)
4. Review integrated code for quality and compliance

### For Software Engineer Agents
1. Read assigned effort plan
2. Follow implementation guidelines
3. Write tests first (TDD)
4. Measure size with `tools/line-counter.sh`
5. Request review when complete

---

**This implementation plan is ready for Software Factory 3.0 orchestration. Begin by running `/init-software-factory` followed by `/continue-software-factory`.**
