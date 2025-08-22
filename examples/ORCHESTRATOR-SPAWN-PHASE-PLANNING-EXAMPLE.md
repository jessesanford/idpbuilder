# Example: Orchestrator Spawning Code Reviewer for Phase Planning

## Scenario
The orchestrator needs to create detailed implementation plans for Phase 2 after Phase 1 is complete.

## Orchestrator's Task Command

```markdown
Task @agent-code-reviewer:

You are acting as the SENIOR PROJECT MAINTAINER to create a detailed implementation plan for Phase 2.

## Context
- Project: E-Commerce Platform
- Phase 2: Business Logic & Services
- Phase 1 Status: COMPLETE (APIs and data models implemented)
- Duration Estimate: 5 days
- Base Branch: phase1-integration

## Your Role: Senior Maintainer

You have deep knowledge of this codebase and must make critical technical decisions for Phase 2.

### ⚠️ CRITICAL RULES - YOU MUST FOLLOW THESE

1. **DO NOT WRITE COMPLETE CODE**
   - Only provide 10-30 line snippets for genuinely complex algorithms
   - SW Engineers will implement the rest

2. **ENFORCE REUSE FROM PHASE 1**
   - We already have:
     - pkg/api/types - All API types and contracts
     - pkg/models - Database models
     - pkg/validators - Input validation
     - pkg/errors - Error handling
   - You MUST mandate reuse of these

3. **MAKE TECHNICAL DECISIONS**
   - Choose specific libraries with exact versions
   - Define patterns that extend Phase 1
   - Specify interface contracts

## Required Reading

1. Read PROJECT-IMPLEMENTATION-PLAN.md for overall structure
2. Read PHASE1-SPECIFIC-IMPL-PLAN.md to understand what was built
3. Read phase-plans/PHASEX-SPECIFIC-IMPL-PLAN-TEMPLATE.md for format
4. Read orchestrator-state.yaml for current progress

## Phase 2 Overview

Phase 2 implements business logic and service layer:
- Wave 2.1: Core business services (Order, Inventory, Payment)
- Wave 2.2: Integration services (Email, SMS, Shipping)
- Wave 2.3: Background workers and schedulers

## Specific Requirements

### Libraries to Consider
- Message queue: Choose between RabbitMQ or Redis
- Payment processing: Stripe or PayPal SDK
- Email service: SendGrid or AWS SES
- Background jobs: Consider existing from Phase 1 or new

### Must Reuse from Phase 1
- ALL models from pkg/models
- ALL validators from pkg/validators
- Error handling patterns from pkg/errors
- API contracts from pkg/api/types

### Complex Algorithms Needing Snippets
1. Inventory reservation with distributed locking
2. Payment retry logic with exponential backoff
3. Order state machine transitions
4. Dynamic pricing calculation

## Deliverable

Create PHASE2-SPECIFIC-IMPL-PLAN.md with:

1. **Library Selections**
   ```yaml
   core_libraries:
     - name: "github.com/stripe/stripe-go/v74"
       version: "v74.30.0"
       reason: "Official Stripe SDK, well-maintained"
       usage: "Payment processing in pkg/services/payment"
   ```

2. **Reuse Enforcement**
   ```yaml
   must_reuse_from_phase1:
     - "pkg/models/order.go - Order model"
     - "pkg/models/inventory.go - Inventory model"
     - "pkg/validators/payment.go - Payment validation"
   
   forbidden_duplications:
     - "DO NOT create new Order type"
     - "DO NOT reimplement validation"
   ```

3. **Critical Snippets (Max 30 lines each)**
   ```go
   // COMPLEX: Distributed inventory reservation
   func reserveInventoryWithLock(itemID string, qty int) error {
       // [15-20 lines of critical locking logic]
       // SW Engineer implements the rest
   }
   ```

4. **Wave/Effort Breakdown**
   - For each effort:
     - What to build
     - What to reuse from Phase 1
     - Estimated lines
     - Test requirements

Remember:
- You're the maintainer who knows what worked in Phase 1
- Guide the implementation, don't write it all
- Enforce consistency and reuse
- Make the hard technical choices

Output: PHASE2-SPECIFIC-IMPL-PLAN.md following the template
```

## What the Code Reviewer Should Produce

The Code Reviewer (as Senior Maintainer) will create a detailed plan that:

1. **Makes Concrete Decisions**:
   - "Use Stripe SDK v74.30.0 for payments"
   - "Use Redis for distributed locking, not Redlock"
   - "Extend Phase 1's error patterns, don't create new"

2. **Provides Critical Code Only**:
   ```go
   // Just the complex distributed lock acquisition
   func acquireDistributedLock(key string, ttl time.Duration) (*Lock, error) {
       // 20 lines of tricky Redis Lua script logic
       script := redis.NewScript(`
           -- Complex Lua for atomic lock with owner tracking
           -- [critical logic here]
       `)
       // SW Engineer implements timeout, retry, cleanup
   }
   ```

3. **Enforces Reuse**:
   ```markdown
   E2.1.1 Order Service MUST:
   - Import "pkg/models/order.go" for Order type
   - Use "pkg/api/types/order_api.go" for API contracts
   - Extend "pkg/validators/order.go" validation
   - NOT create new Order struct
   - NOT duplicate validation logic
   ```

4. **Structures Work**:
   ```markdown
   Wave 2.1: Core Services (3 efforts)
   - E2.1.1: Order Service (400 lines)
   - E2.1.2: Inventory Service (350 lines)  
   - E2.1.3: Payment Service (450 lines)
   ```

## Benefits of This Approach

1. **Consistency**: Phase 2 builds on Phase 1 properly
2. **No Duplication**: Reuse is enforced
3. **Technical Leadership**: Hard decisions are made upfront
4. **Focused Implementation**: SW Engineers know exactly what to build
5. **Complex Parts Handled**: Critical algorithms provided
6. **Clear Boundaries**: Maintainer guides, engineers implement

## Common Mistakes to Avoid

### ❌ Maintainer Writing Too Much
```go
// BAD: 200 lines of complete service implementation
type OrderService struct {
    // ... entire implementation ...
}
```

### ✅ Maintainer Providing Key Logic Only
```go
// GOOD: Just the complex state transition
func (o *Order) transitionState(from, to State) error {
    // 15 lines of critical state machine logic
    // with complex validation rules
}
// SW Engineer implements the service using this
```

### ❌ Vague Library Choices
"Use some message queue library"

### ✅ Specific Technical Decisions
"Use github.com/rabbitmq/amqp091-go v1.9.0 - mature, handles reconnection"

## Result

The orchestrator receives a detailed Phase 2 plan that:
- Builds properly on Phase 1
- Makes all technical decisions
- Provides critical complex snippets
- Enforces reuse and prevents duplication
- Guides without over-implementing

The orchestrator can then break this into efforts and task SW Engineers with confidence.