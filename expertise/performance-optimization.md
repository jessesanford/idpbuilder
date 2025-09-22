# Performance Optimization and Scaling Patterns Expertise Module

## ┌─ Controller Performance Patterns [R251-R260] ─┐
### │ Efficient Kubernetes Controller Design       │
└─────────────────────────────────────────────────────┘

**Rule R251**: Controllers MUST implement efficient caching strategies
**Rule R288**: Reconciliation loops MUST avoid unnecessary API calls
**Rule R288**: Watch predicates MUST filter irrelevant events
**Rule R254**: Batch operations MUST be used for bulk updates
**Rule R255**: Resource quotas and limits MUST be properly configured

### Optimized Controller Caching Pattern

```go
// R251: Efficient caching with workspace-aware stores
type OptimizedController struct {
    client.Client
    cache           cache.Cache
    workspaceCache  map[string]*WorkspaceCache
    cacheMutex      sync.RWMutex
    recorder        record.EventRecorder
    metrics         *ControllerMetrics
}

type WorkspaceCache struct {
    store       cache.Store
    informer    cache.SharedIndexInformer
    stopCh      chan struct{}
    lastUpdate  time.Time
}

func (r *OptimizedController) SetupWithManager(mgr ctrl.Manager) error {
    // R288: Efficient event filtering with predicates
    return ctrl.NewControllerManagedBy(mgr).
        For(&myapiv1.MyResource{}).
        Owns(&corev1.ConfigMap{}).
        Owns(&appsv1.Deployment{}).
        WithEventFilter(r.buildOptimizedPredicates()).
        WithOptions(controller.Options{
            MaxConcurrentReconciles: r.calculateOptimalConcurrency(),
            RateLimiter: workqueue.NewMaxOfRateLimiter(
                workqueue.NewItemExponentialFailureRateLimiter(1*time.Second, 10*time.Second),
                &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(50), 100)},
            ),
        }).
        Complete(r)
}

func (r *OptimizedController) buildOptimizedPredicates() predicate.Predicate {
    return predicate.Funcs{
        UpdateFunc: func(e event.UpdateEvent) bool {
            oldResource := e.ObjectOld.(*myapiv1.MyResource)
            newResource := e.ObjectNew.(*myapiv1.MyResource)
            
            // R288: Filter out status-only updates and unnecessary reconciliations
            if oldResource.Generation == newResource.Generation {
                // Only reconcile on meaningful status changes
                return !reflect.DeepEqual(oldResource.Status.Phase, newResource.Status.Phase) ||
                       oldResource.DeletionTimestamp != newResource.DeletionTimestamp
            }
            
            // Reconcile on spec changes
            return !reflect.DeepEqual(oldResource.Spec, newResource.Spec)
        },
        CreateFunc: func(e event.CreateEvent) bool {
            return true // Always process new resources
        },
        DeleteFunc: func(e event.DeleteEvent) bool {
            return false // Handled by finalizers, not delete events
        },
        GenericFunc: func(e event.GenericEvent) bool {
            return false // No generic events needed
        },
    }
}

func (r *OptimizedController) calculateOptimalConcurrency() int {
    // Calculate based on system resources and expected load
    cpus := runtime.NumCPU()
    // Conservative approach: 2 * CPUs, with minimum of 1 and maximum of 10
    concurrency := cpus * 2
    if concurrency < 1 {
        concurrency = 1
    } else if concurrency > 10 {
        concurrency = 10
    }
    return concurrency
}

// R288: Avoid unnecessary API calls through caching
func (r *OptimizedController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    startTime := time.Now()
    defer func() {
        r.metrics.ReconcileDuration.WithLabelValues(req.Name).Observe(time.Since(startTime).Seconds())
    }()

    // Use cached client to reduce API server load
    var resource myapiv1.MyResource
    if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
        if errors.IsNotFound(err) {
            return ctrl.Result{}, nil
        }
        r.metrics.ReconcileErrors.WithLabelValues("get_resource").Inc()
        return ctrl.Result{}, err
    }

    // Check if we've processed this generation recently
    if r.isRecentlyProcessed(&resource) {
        return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
    }

    return r.reconcileOptimized(ctx, &resource)
}

func (r *OptimizedController) reconcileOptimized(ctx context.Context, resource *myapiv1.MyResource) (ctrl.Result, error) {
    // R288: Batch related operations to minimize API calls
    desired := r.buildDesiredState(resource)
    current, err := r.getCurrentStateBatch(ctx, resource)
    if err != nil {
        return ctrl.Result{}, err
    }

    changes := r.calculateChanges(current, desired)
    if len(changes) == 0 {
        // No changes needed, update status if necessary
        return r.updateStatusIfNeeded(ctx, resource, current)
    }

    // R254: Apply changes in batches
    if err := r.applyChangesBatch(ctx, resource, changes); err != nil {
        return ctrl.Result{RequeueAfter: time.Second * 30}, err
    }

    return r.updateStatusIfNeeded(ctx, resource, desired)
}
```

### Batch Operations Pattern

```go
// R254: Batch operations for improved performance
func (r *OptimizedController) applyChangesBatch(ctx context.Context, resource *myapiv1.MyResource, changes []Change) error {
    // Group changes by operation type for batching
    creates := make([]client.Object, 0)
    updates := make([]client.Object, 0)
    deletes := make([]client.Object, 0)

    for _, change := range changes {
        switch change.Type {
        case ChangeTypeCreate:
            creates = append(creates, change.Object)
        case ChangeTypeUpdate:
            updates = append(updates, change.Object)
        case ChangeTypeDelete:
            deletes = append(deletes, change.Object)
        }
    }

    // Execute batched operations
    if len(creates) > 0 {
        if err := r.createBatch(ctx, creates); err != nil {
            return fmt.Errorf("batch create failed: %w", err)
        }
    }

    if len(updates) > 0 {
        if err := r.updateBatch(ctx, updates); err != nil {
            return fmt.Errorf("batch update failed: %w", err)
        }
    }

    if len(deletes) > 0 {
        if err := r.deleteBatch(ctx, deletes); err != nil {
            return fmt.Errorf("batch delete failed: %w", err)
        }
    }

    return nil
}

func (r *OptimizedController) createBatch(ctx context.Context, objects []client.Object) error {
    // Use goroutines for concurrent creation with controlled parallelism
    semaphore := make(chan struct{}, 5) // Limit to 5 concurrent operations
    errCh := make(chan error, len(objects))
    
    for _, obj := range objects {
        go func(o client.Object) {
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            if err := r.Create(ctx, o); err != nil {
                if !errors.IsAlreadyExists(err) {
                    errCh <- err
                    return
                }
            }
            errCh <- nil
        }(obj)
    }

    // Collect results
    var errs []error
    for i := 0; i < len(objects); i++ {
        if err := <-errCh; err != nil {
            errs = append(errs, err)
        }
    }

    if len(errs) > 0 {
        return fmt.Errorf("batch create errors: %v", errs)
    }

    return nil
}
```

## ┌─ Memory and Resource Management [R261-R270] ─┐
### │ Efficient Resource Utilization               │
└─────────────────────────────────────────────────────┘

**Rule R261**: Memory usage MUST be monitored and bounded
**Rule R262**: Object pools MUST be used for frequently allocated objects
**Rule R263**: Goroutine leaks MUST be prevented with proper cleanup
**Rule R264**: CPU usage MUST be optimized for high-throughput scenarios
**Rule R265**: Garbage collection impact MUST be minimized

### Memory-Optimized Object Management

```go
// R261, R262: Memory monitoring and object pooling
type ResourceManager struct {
    client          client.Client
    objectPools     *ObjectPools
    memoryMonitor   *MemoryMonitor
    resourceTracker *ResourceTracker
    metrics         *ResourceMetrics
}

type ObjectPools struct {
    requestPool    sync.Pool
    responsePool   sync.Pool
    bufferPool     sync.Pool
    stringBuilder  sync.Pool
}

func NewOptimizedResourceManager(client client.Client) *ResourceManager {
    return &ResourceManager{
        client: client,
        objectPools: &ObjectPools{
            requestPool: sync.Pool{
                New: func() interface{} {
                    return &ReconcileRequest{
                        Changes: make([]Change, 0, 10), // Pre-allocate capacity
                    }
                },
            },
            responsePool: sync.Pool{
                New: func() interface{} {
                    return &ReconcileResponse{
                        Events: make([]Event, 0, 5),
                    }
                },
            },
            bufferPool: sync.Pool{
                New: func() interface{} {
                    return make([]byte, 0, 1024) // 1KB buffer
                },
            },
            stringBuilder: sync.Pool{
                New: func() interface{} {
                    return &strings.Builder{}
                },
            },
        },
        memoryMonitor:   NewMemoryMonitor(),
        resourceTracker: NewResourceTracker(),
    }
}

// R262: Use object pools to reduce allocations
func (rm *ResourceManager) processResource(ctx context.Context, resource *myapiv1.MyResource) error {
    // Get request object from pool
    req := rm.objectPools.requestPool.Get().(*ReconcileRequest)
    defer func() {
        // Reset and return to pool
        req.Reset()
        rm.objectPools.requestPool.Put(req)
    }()

    // Get response object from pool  
    resp := rm.objectPools.responsePool.Get().(*ReconcileResponse)
    defer func() {
        resp.Reset()
        rm.objectPools.responsePool.Put(resp)
    }()

    // Use pooled buffer for temporary data
    buffer := rm.objectPools.bufferPool.Get().([]byte)
    defer func() {
        buffer = buffer[:0] // Reset length but keep capacity
        rm.objectPools.bufferPool.Put(buffer)
    }()

    return rm.reconcileWithPools(ctx, resource, req, resp, buffer)
}

// R261: Memory monitoring and bounded usage
type MemoryMonitor struct {
    maxMemoryUsage uint64
    currentUsage   uint64
    lastGC         time.Time
    gcTrigger      uint64
    mu             sync.RWMutex
}

func (mm *MemoryMonitor) checkMemoryUsage() error {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    mm.mu.Lock()
    mm.currentUsage = m.Alloc
    mm.mu.Unlock()
    
    // R265: Force GC if memory usage is high
    if m.Alloc > mm.gcTrigger && time.Since(mm.lastGC) > time.Minute {
        runtime.GC()
        mm.lastGC = time.Now()
    }
    
    if m.Alloc > mm.maxMemoryUsage {
        return fmt.Errorf("memory usage %d bytes exceeds limit %d bytes", m.Alloc, mm.maxMemoryUsage)
    }
    
    return nil
}

// R263: Prevent goroutine leaks with proper cleanup
type GoroutineManager struct {
    activeGoroutines sync.WaitGroup
    ctx              context.Context
    cancel           context.CancelFunc
    maxGoroutines    int32
    currentCount     int32
}

func (gm *GoroutineManager) executeAsync(fn func(context.Context) error) error {
    // Check if we're at the limit
    if atomic.LoadInt32(&gm.currentCount) >= gm.maxGoroutines {
        return fmt.Errorf("maximum goroutines limit reached: %d", gm.maxGoroutines)
    }
    
    gm.activeGoroutines.Add(1)
    atomic.AddInt32(&gm.currentCount, 1)
    
    go func() {
        defer func() {
            gm.activeGoroutines.Done()
            atomic.AddInt32(&gm.currentCount, -1)
        }()
        
        if err := fn(gm.ctx); err != nil {
            // Log error but don't panic
            log.Error(err, "async operation failed")
        }
    }()
    
    return nil
}

func (gm *GoroutineManager) shutdown(timeout time.Duration) error {
    // Cancel context to signal shutdown
    gm.cancel()
    
    // Wait for all goroutines to finish with timeout
    done := make(chan struct{})
    go func() {
        gm.activeGoroutines.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-time.After(timeout):
        return fmt.Errorf("shutdown timeout: %d goroutines still running", atomic.LoadInt32(&gm.currentCount))
    }
}
```

## ┌─ Caching and Data Access Optimization [R271-R280] ─┐
### │ Efficient Data Retrieval and Storage             │
└─────────────────────────────────────────────────────────┘

**Rule R271**: Frequently accessed data MUST be cached appropriately
**Rule R272**: Cache invalidation MUST be handled correctly
**Rule R273**: Database queries MUST be optimized for performance
**Rule R274**: Pagination MUST be implemented for large result sets
**Rule R275**: Connection pools MUST be configured for optimal throughput

### Multi-Level Caching Strategy

```go
// R271: Multi-level caching for optimal performance
type CacheManager struct {
    l1Cache    *LRUCache          // In-memory fast cache
    l2Cache    *DistributedCache  // Redis/similar for shared caching
    etcdCache  *EtcdCache         // Direct etcd caching for K8s objects
    metrics    *CacheMetrics
}

type CacheKey struct {
    Workspace string
    Kind      string
    Namespace string
    Name      string
    Version   string
}

func (ck CacheKey) String() string {
    return fmt.Sprintf("%s/%s/%s/%s@%s", ck.Workspace, ck.Kind, ck.Namespace, ck.Name, ck.Version)
}

type LRUCache struct {
    cache      *lru.Cache
    maxSize    int
    ttl        time.Duration
    mu         sync.RWMutex
    metrics    *CacheMetrics
    timestamps map[string]time.Time
}

func NewLRUCache(maxSize int, ttl time.Duration) *LRUCache {
    cache, _ := lru.New(maxSize)
    lc := &LRUCache{
        cache:      cache,
        maxSize:    maxSize,
        ttl:        ttl,
        timestamps: make(map[string]time.Time),
    }
    
    // Start cleanup goroutine
    go lc.cleanup()
    return lc
}

func (lc *LRUCache) Get(key string) (interface{}, bool) {
    lc.mu.RLock()
    defer lc.mu.RUnlock()
    
    // Check TTL
    if timestamp, exists := lc.timestamps[key]; exists {
        if time.Since(timestamp) > lc.ttl {
            // Item expired, remove it
            lc.cache.Remove(key)
            delete(lc.timestamps, key)
            lc.metrics.Misses.Inc()
            return nil, false
        }
    }
    
    if value, exists := lc.cache.Get(key); exists {
        lc.metrics.Hits.Inc()
        return value, true
    }
    
    lc.metrics.Misses.Inc()
    return nil, false
}

func (lc *LRUCache) Set(key string, value interface{}) {
    lc.mu.Lock()
    defer lc.mu.Unlock()
    
    lc.cache.Add(key, value)
    lc.timestamps[key] = time.Now()
}

// R272: Proper cache invalidation
func (cm *CacheManager) InvalidateResource(workspace, kind, namespace, name string) {
    pattern := fmt.Sprintf("%s/%s/%s/%s@*", workspace, kind, namespace, name)
    
    // Invalidate from all cache levels
    cm.l1Cache.InvalidatePattern(pattern)
    cm.l2Cache.InvalidatePattern(pattern)
    cm.etcdCache.InvalidatePattern(pattern)
}

func (cm *CacheManager) GetResource(ctx context.Context, key CacheKey, fetcher func() (interface{}, error)) (interface{}, error) {
    keyStr := key.String()
    
    // Try L1 cache first (fastest)
    if value, found := cm.l1Cache.Get(keyStr); found {
        return value, nil
    }
    
    // Try L2 cache (shared)
    if value, found := cm.l2Cache.Get(ctx, keyStr); found {
        // Store in L1 for faster future access
        cm.l1Cache.Set(keyStr, value)
        return value, nil
    }
    
    // Fetch from source
    value, err := fetcher()
    if err != nil {
        return nil, err
    }
    
    // Store in both cache levels
    cm.l1Cache.Set(keyStr, value)
    cm.l2Cache.Set(ctx, keyStr, value, time.Hour) // L2 cache with longer TTL
    
    return value, nil
}
```

### Optimized Database Access Pattern

```go
// R273, R274: Optimized queries with pagination
type ResourceStore struct {
    db          *sql.DB
    queryCache  map[string]*sql.Stmt
    connPool    *ConnectionPool
    mu          sync.RWMutex
}

type QueryOptions struct {
    Limit     int
    Offset    int
    OrderBy   string
    Filters   map[string]interface{}
    Workspace string
}

func (rs *ResourceStore) ListResources(ctx context.Context, opts QueryOptions) (*ResourceList, error) {
    // R274: Implement pagination for large result sets
    if opts.Limit == 0 {
        opts.Limit = 100 // Default page size
    }
    if opts.Limit > 1000 {
        opts.Limit = 1000 // Maximum page size
    }
    
    // Build optimized query with proper indexing
    query := rs.buildOptimizedQuery(opts)
    
    // Use prepared statement from cache
    stmt, err := rs.getOrCreateStatement(query)
    if err != nil {
        return nil, fmt.Errorf("failed to prepare statement: %w", err)
    }
    
    // Execute with timeout
    queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    rows, err := stmt.QueryContext(queryCtx, rs.buildQueryArgs(opts)...)
    if err != nil {
        return nil, fmt.Errorf("query execution failed: %w", err)
    }
    defer rows.Close()
    
    return rs.scanResourceList(rows, opts)
}

func (rs *ResourceStore) buildOptimizedQuery(opts QueryOptions) string {
    var queryBuilder strings.Builder
    
    // Base query with proper indexing
    queryBuilder.WriteString(`
        SELECT r.id, r.name, r.namespace, r.spec, r.status, r.created_at, r.updated_at
        FROM resources r
        WHERE r.workspace = $1
    `)
    
    argCount := 1
    
    // Add filters with parameterized queries to prevent SQL injection
    for field, _ := range opts.Filters {
        argCount++
        queryBuilder.WriteString(fmt.Sprintf(" AND r.%s = $%d", field, argCount))
    }
    
    // Add ordering with index utilization
    if opts.OrderBy != "" {
        queryBuilder.WriteString(" ORDER BY r." + opts.OrderBy)
    } else {
        queryBuilder.WriteString(" ORDER BY r.created_at DESC")
    }
    
    // Add pagination
    queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", opts.Limit, opts.Offset))
    
    return queryBuilder.String()
}

// R275: Optimized connection pool configuration
type ConnectionPool struct {
    db          *sql.DB
    maxOpen     int
    maxIdle     int
    maxLifetime time.Duration
    maxIdleTime time.Duration
}

func NewConnectionPool(dsn string) (*ConnectionPool, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    pool := &ConnectionPool{
        db:          db,
        maxOpen:     25,  // Based on expected load
        maxIdle:     10,  // Keep some connections ready
        maxLifetime: time.Hour,
        maxIdleTime: time.Minute * 15,
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(pool.maxOpen)
    db.SetMaxIdleConns(pool.maxIdle)
    db.SetConnMaxLifetime(pool.maxLifetime)
    db.SetConnMaxIdleTime(pool.maxIdleTime)
    
    return pool, nil
}
```

## ┌─ Concurrency and Parallelism [R281-R290] ─┐
### │ Efficient Concurrent Processing           │
└─────────────────────────────────────────────────┘

**Rule R281**: Concurrent operations MUST be properly synchronized
**Rule R282**: Worker pools MUST be sized appropriately for workload
**Rule R283**: Rate limiting MUST prevent resource exhaustion
**Rule R284**: Circuit breakers MUST handle downstream failures
**Rule R285**: Backpressure MUST be implemented for high-load scenarios

### High-Performance Worker Pool Pattern

```go
// R281, R282: Properly sized and synchronized worker pool
type WorkerPool struct {
    workerCount   int
    jobQueue      chan Job
    workers       []*Worker
    wg            sync.WaitGroup
    ctx           context.Context
    cancel        context.CancelFunc
    rateLimiter   *RateLimiter
    circuitBreaker *CircuitBreaker
    metrics       *WorkerMetrics
}

type Job interface {
    Execute(ctx context.Context) error
    Priority() int
    ID() string
}

type Worker struct {
    id          int
    pool        *WorkerPool
    jobChannel  chan Job
    quit        chan struct{}
    metrics     *WorkerMetrics
}

func NewWorkerPool(workerCount int, bufferSize int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    pool := &WorkerPool{
        workerCount:    workerCount,
        jobQueue:       make(chan Job, bufferSize),
        workers:        make([]*Worker, workerCount),
        ctx:            ctx,
        cancel:         cancel,
        rateLimiter:    NewRateLimiter(100, time.Second), // 100 ops/sec
        circuitBreaker: NewCircuitBreaker("worker-pool", 10, time.Minute),
        metrics:        NewWorkerMetrics(),
    }
    
    // Start workers
    for i := 0; i < workerCount; i++ {
        worker := &Worker{
            id:         i,
            pool:       pool,
            jobChannel: make(chan Job, 1),
            quit:       make(chan struct{}),
            metrics:    pool.metrics,
        }
        pool.workers[i] = worker
        pool.wg.Add(1)
        go worker.start()
    }
    
    // Start job dispatcher
    go pool.dispatch()
    
    return pool
}

func (w *Worker) start() {
    defer w.pool.wg.Done()
    
    for {
        select {
        case job := <-w.jobChannel:
            w.executeJob(job)
        case <-w.quit:
            return
        case <-w.pool.ctx.Done():
            return
        }
    }
}

func (w *Worker) executeJob(job Job) {
    startTime := time.Now()
    w.metrics.ActiveJobs.Inc()
    defer func() {
        w.metrics.ActiveJobs.Dec()
        w.metrics.JobDuration.WithLabelValues(job.ID()).Observe(time.Since(startTime).Seconds())
    }()
    
    // R283: Rate limiting to prevent resource exhaustion
    if err := w.pool.rateLimiter.Wait(w.pool.ctx); err != nil {
        w.metrics.JobsRateLimited.Inc()
        return
    }
    
    // R284: Circuit breaker for downstream failures
    err := w.pool.circuitBreaker.Execute(func() error {
        return job.Execute(w.pool.ctx)
    })
    
    if err != nil {
        w.metrics.JobErrors.WithLabelValues(job.ID()).Inc()
        log.Error(err, "Job execution failed", "worker", w.id, "job", job.ID())
    } else {
        w.metrics.JobsCompleted.WithLabelValues(job.ID()).Inc()
    }
}

// R285: Backpressure implementation
func (wp *WorkerPool) Submit(job Job) error {
    select {
    case wp.jobQueue <- job:
        wp.metrics.JobsQueued.Inc()
        return nil
    case <-time.After(time.Second * 5): // Backpressure timeout
        wp.metrics.JobsRejected.Inc()
        return fmt.Errorf("job queue full, rejecting job %s", job.ID())
    case <-wp.ctx.Done():
        return fmt.Errorf("worker pool shutting down")
    }
}

func (wp *WorkerPool) dispatch() {
    for {
        select {
        case job := <-wp.jobQueue:
            // Find least loaded worker
            worker := wp.findOptimalWorker()
            select {
            case worker.jobChannel <- job:
                // Job dispatched successfully
            case <-time.After(time.Millisecond * 100):
                // Worker overloaded, try another
                select {
                case wp.jobQueue <- job: // Put back in queue
                case <-wp.ctx.Done():
                    return
                }
            }
        case <-wp.ctx.Done():
            return
        }
    }
}
```

### Rate Limiting and Circuit Breaker Pattern

```go
// R283: Advanced rate limiting with token bucket
type RateLimiter struct {
    limiter *rate.Limiter
    burst   int
}

func NewRateLimiter(rps int, per time.Duration) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Every(per/time.Duration(rps)), rps),
        burst:   rps,
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    return rl.limiter.Wait(ctx)
}

// R284: Circuit breaker implementation
type CircuitBreaker struct {
    name           string
    maxFailures    int
    timeout        time.Duration
    state          CircuitState
    failureCount   int
    lastFailureTime time.Time
    mu             sync.RWMutex
    metrics        *CircuitBreakerMetrics
}

type CircuitState int

const (
    StateClosed CircuitState = iota
    StateOpen
    StateHalfOpen
)

func NewCircuitBreaker(name string, maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        name:        name,
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       StateClosed,
        metrics:     NewCircuitBreakerMetrics(name),
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.canExecute() {
        cb.metrics.RequestsRejected.Inc()
        return fmt.Errorf("circuit breaker %s is open", cb.name)
    }
    
    err := fn()
    cb.onResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        // Check if timeout has passed
        if time.Since(cb.lastFailureTime) > cb.timeout {
            return true // Allow one request to test
        }
        return false
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

func (cb *CircuitBreaker) onResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failureCount++
        cb.lastFailureTime = time.Now()
        cb.metrics.Failures.Inc()
        
        if cb.failureCount >= cb.maxFailures {
            cb.state = StateOpen
            cb.metrics.StateChanges.WithLabelValues("open").Inc()
        }
    } else {
        cb.metrics.Successes.Inc()
        
        if cb.state == StateHalfOpen {
            cb.state = StateClosed
            cb.failureCount = 0
            cb.metrics.StateChanges.WithLabelValues("closed").Inc()
        }
    }
}
```

## ┌─ Monitoring and Observability [R291-R300] ─┐
### │ Performance Metrics and Profiling         │
└─────────────────────────────────────────────────┘

**Rule R291**: Performance metrics MUST be exposed for monitoring
**Rule R292**: Resource usage MUST be tracked and alerted on
**Rule R293**: Profiling MUST be available for performance analysis
**Rule R294**: Distributed tracing MUST be implemented for complex operations
**Rule R295**: Log levels MUST be configurable for performance optimization

### Comprehensive Performance Monitoring

```go
// R291: Performance metrics exposure
type PerformanceMonitor struct {
    registry       prometheus.Registerer
    httpDuration   *prometheus.HistogramVec
    grpcDuration   *prometheus.HistogramVec
    cacheHitRate   *prometheus.CounterVec
    memoryUsage    prometheus.Gauge
    cpuUsage       prometheus.Gauge
    goroutineCount prometheus.Gauge
    activeRequests prometheus.Gauge
}

func NewPerformanceMonitor(registry prometheus.Registerer) *PerformanceMonitor {
    pm := &PerformanceMonitor{
        registry: registry,
        httpDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "http_request_duration_seconds",
                Help: "HTTP request duration in seconds",
                Buckets: prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms to 16s
            },
            []string{"method", "endpoint", "status_code"},
        ),
        grpcDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "grpc_request_duration_seconds", 
                Help: "gRPC request duration in seconds",
                Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
            },
            []string{"method", "status_code"},
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "cache_requests_total",
                Help: "Total cache requests",
            },
            []string{"cache_name", "result"}, // result: hit/miss
        ),
        memoryUsage: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "memory_usage_bytes",
                Help: "Current memory usage in bytes",
            },
        ),
        cpuUsage: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "cpu_usage_percent",
                Help: "Current CPU usage percentage",
            },
        ),
        goroutineCount: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "goroutines_total",
                Help: "Current number of goroutines",
            },
        ),
        activeRequests: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "active_requests_total",
                Help: "Current number of active requests",
            },
        ),
    }
    
    // Register metrics
    registry.MustRegister(
        pm.httpDuration,
        pm.grpcDuration,
        pm.cacheHitRate,
        pm.memoryUsage,
        pm.cpuUsage,
        pm.goroutineCount,
        pm.activeRequests,
    )
    
    // Start background monitoring
    go pm.startResourceMonitoring()
    
    return pm
}

// R292: Resource usage tracking
func (pm *PerformanceMonitor) startResourceMonitoring() {
    ticker := time.NewTicker(time.Second * 15)
    defer ticker.Stop()
    
    for range ticker.C {
        // Memory usage
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        pm.memoryUsage.Set(float64(m.Alloc))
        
        // Goroutine count
        pm.goroutineCount.Set(float64(runtime.NumGoroutine()))
        
        // CPU usage (simplified - use proper CPU monitoring in production)
        pm.updateCPUUsage()
    }
}

// R293: HTTP middleware with profiling support
func (pm *PerformanceMonitor) HTTPMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        pm.activeRequests.Inc()
        defer pm.activeRequests.Dec()
        
        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        // R294: Distributed tracing support
        ctx := r.Context()
        span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
        defer span.Finish()
        
        // Add tracing headers
        span.SetTag("http.method", r.Method)
        span.SetTag("http.url", r.URL.String())
        
        next.ServeHTTP(wrapped, r.WithContext(ctx))
        
        // Record metrics
        duration := time.Since(start)
        pm.httpDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
            fmt.Sprintf("%d", wrapped.statusCode),
        ).Observe(duration.Seconds())
        
        // Add timing to trace
        span.SetTag("http.status_code", wrapped.statusCode)
        span.SetTag("http.duration_ms", float64(duration.Nanoseconds())/1e6)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### Application Performance Profiling

```go
// R293: Built-in profiling support
func EnableProfiling(port int) {
    // CPU profiling endpoint
    http.HandleFunc("/debug/pprof/profile", func(w http.ResponseWriter, r *http.Request) {
        duration := 30 * time.Second
        if d := r.URL.Query().Get("seconds"); d != "" {
            if parsed, err := time.ParseDuration(d + "s"); err == nil {
                duration = parsed
            }
        }
        
        pprof.Profile(w, r)
    })
    
    // Memory profiling
    http.HandleFunc("/debug/pprof/heap", pprof.Index)
    http.HandleFunc("/debug/pprof/goroutine", pprof.Index)
    http.HandleFunc("/debug/pprof/block", pprof.Index)
    http.HandleFunc("/debug/pprof/mutex", pprof.Index)
    
    // Custom metrics endpoint
    http.HandleFunc("/metrics", promhttp.Handler())
    
    // Performance analysis endpoint
    http.HandleFunc("/debug/performance", func(w http.ResponseWriter, r *http.Request) {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        stats := map[string]interface{}{
            "goroutines":    runtime.NumGoroutine(),
            "memory_alloc":  m.Alloc,
            "memory_sys":    m.Sys,
            "gc_cycles":     m.NumGC,
            "gc_pause_ns":   m.PauseNs[(m.NumGC+255)%256],
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(stats)
    })
    
    log.Printf("Profiling enabled on port %d", port)
    go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
```

## ┌─ Pattern Detection Queries [R301-R305] ─┐
### │ Performance Anti-Pattern Detection      │
└───────────────────────────────────────────────┘

### Detect Inefficient Loops
```bash
# Query: Find potentially inefficient loops in controllers
grep -r "for.*range" --include="*.go" -A 5 | grep -E "client\.Get\|client\.List\|client\.Create" | head -10
```

### Detect Missing Pagination
```bash
# Query: Find List operations without pagination
grep -r "\.List\(" --include="*.go" | grep -v "client\.Limit\|ListOptions\|Continue" | head -10
```

### Detect Goroutine Leaks
```bash
# Query: Find goroutines without proper cleanup
grep -r "go func" --include="*.go" -A 10 | grep -L "defer\|context\.Done\|cancel" | head -10
```

## ┌─ Validation Criteria [R306-R310] ─┐
### │ Performance Implementation Quality │
└─────────────────────────────────────────┘

**Validation R306**: All controllers MUST implement efficient caching
**Validation R307**: Resource usage MUST stay within defined limits  
**Validation R308**: Concurrent operations MUST be properly managed
**Validation R309**: Performance metrics MUST be exposed and monitored
**Validation R310**: Profiling MUST be available for production debugging

### Performance Implementation Checklist

- [ ] Controllers use efficient predicate filtering
- [ ] Caching is implemented at appropriate levels
- [ ] Object pools are used for frequently allocated objects
- [ ] Database queries are optimized with proper indexing
- [ ] Pagination is implemented for large result sets
- [ ] Rate limiting prevents resource exhaustion
- [ ] Circuit breakers handle downstream failures
- [ ] Worker pools are properly sized
- [ ] Memory usage is monitored and bounded
- [ ] Performance metrics are exposed via Prometheus
- [ ] Distributed tracing is implemented
- [ ] Profiling endpoints are available
- [ ] Goroutine leaks are prevented
- [ ] Resource limits are enforced and monitored

This comprehensive performance optimization module provides patterns and techniques for building high-performance, scalable Kubernetes controllers and applications that can handle production workloads efficiently while maintaining observability and reliability.