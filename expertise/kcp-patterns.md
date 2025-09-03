# KCP Patterns and Best Practices Expertise Module

## ┌─ KCP Core Architecture Patterns [R001-R010] ─┐
### │ Multi-Tenant Workspace Isolation             │
└─────────────────────────────────────────────────┘

**Rule R001**: All KCP resources MUST implement workspace-aware controllers
**Rule R002**: Logical clusters MUST be properly scoped and isolated
**Rule R003**: Cross-workspace resource access MUST use proper authorization

### Workspace-Aware Controller Pattern

```go
// Correct KCP controller pattern with workspace awareness
type WorkspaceAwareController struct {
    client dynamic.Interface
    informer cache.SharedIndexInformer
    workspaceClient kcpkubernetesclient.Interface
}

func (c *WorkspaceAwareController) processWorkQueueItem() bool {
    obj, shutdown := c.workQueue.Get()
    if shutdown {
        return false
    }
    defer c.workQueue.Done(obj)

    req := obj.(reconcile.Request)
    
    // CRITICAL: Extract workspace context from request
    workspace, _, err := cache.SplitMetaNamespaceKey(req.NamespacedName.String())
    if err != nil {
        utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", req.NamespacedName.String()))
        return true
    }
    
    // Use workspace-scoped client
    workspaceConfig := rest.CopyConfig(c.config)
    workspaceConfig.Host = workspaceConfig.Host + "/clusters/" + workspace
    workspaceClient, err := dynamic.NewForConfig(workspaceConfig)
    if err != nil {
        return false
    }
    
    return c.reconcileInWorkspace(req, workspaceClient, workspace)
}
```

### Anti-Pattern: Global Resource Access
```go
// WRONG - Direct global resource access without workspace context
func (c *Controller) getResource(name string) error {
    // This ignores workspace boundaries
    return c.globalClient.Get(context.TODO(), name, &resource)
}

// CORRECT - Workspace-aware resource access
func (c *Controller) getResourceInWorkspace(workspace, name string) error {
    workspaceClient := c.getWorkspaceClient(workspace)
    return workspaceClient.Get(context.TODO(), name, &resource)
}
```

## ┌─ Logical Cluster Management [R011-R020] ─┐
### │ Cluster Lifecycle and State Management   │
└─────────────────────────────────────────────┘

**Rule R011**: Logical clusters MUST have proper lifecycle management
**Rule R012**: Cluster state transitions MUST be atomic and consistent
**Rule R013**: Cross-cluster communication MUST use KCP's native mechanisms

### Logical Cluster Lifecycle Pattern

```go
// LogicalClusterReconciler manages cluster lifecycle
type LogicalClusterReconciler struct {
    client.Client
    Scheme *runtime.Scheme
    WorkspaceClient kcpkubernetesclient.Interface
}

func (r *LogicalClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    var cluster corev1alpha1.LogicalCluster
    if err := r.Get(ctx, req.NamespacedName, &cluster); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Implement state machine for cluster lifecycle
    switch cluster.Status.Phase {
    case corev1alpha1.LogicalClusterPhaseScheduling:
        return r.reconcileScheduling(ctx, &cluster)
    case corev1alpha1.LogicalClusterPhaseInitializing:
        return r.reconcileInitializing(ctx, &cluster)
    case corev1alpha1.LogicalClusterPhaseReady:
        return r.reconcileReady(ctx, &cluster)
    case corev1alpha1.LogicalClusterPhaseDeleting:
        return r.reconcileDeleting(ctx, &cluster)
    default:
        return r.reconcileUnknown(ctx, &cluster)
    }
}

func (r *LogicalClusterReconciler) reconcileScheduling(ctx context.Context, cluster *corev1alpha1.LogicalCluster) (ctrl.Result, error) {
    // CRITICAL: Ensure workspace isolation during scheduling
    workspace := cluster.Annotations[corev1alpha1.LogicalClusterWorkspaceAnnotation]
    if workspace == "" {
        return ctrl.Result{}, fmt.Errorf("cluster must have workspace annotation")
    }
    
    // Update status atomically
    cluster.Status.Phase = corev1alpha1.LogicalClusterPhaseInitializing
    cluster.Status.Conditions = []corev1alpha1.Condition{
        {
            Type:   corev1alpha1.LogicalClusterScheduled,
            Status: corev1alpha1.ConditionTrue,
            LastTransitionTime: metav1.Now(),
            Reason: "Scheduled",
            Message: fmt.Sprintf("Cluster scheduled in workspace %s", workspace),
        },
    }
    
    return ctrl.Result{RequeueAfter: time.Second * 30}, r.Status().Update(ctx, cluster)
}
```

## ┌─ API Export and Binding Patterns [R021-R030] ─┐
### │ Cross-Workspace Resource Sharing             │
└─────────────────────────────────────────────────┘

**Rule R021**: API exports MUST define clear resource boundaries
**Rule R022**: API bindings MUST validate workspace permissions
**Rule R023**: Exported APIs MUST implement proper versioning

### API Export Pattern

```go
// APIExport defines resources available for cross-workspace access
apiVersion: apis.kcp.io/v1alpha1
kind: APIExport
metadata:
  name: my-service-api
  namespace: default
spec:
  latestResourceSchemas:
  - name: myservices.example.com
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              endpoint:
                type: string
              port:
                type: integer
                minimum: 1
                maximum: 65535
          status:
            type: object
            properties:
              phase:
                type: string
                enum: ["Pending", "Ready", "Failed"]
  permissionClaims:
  - group: ""
    resource: secrets
    resourceNames: ["my-service-secret"]
    verb: get
```

### API Binding Validation Pattern

```go
type APIBindingController struct {
    client.Client
    WorkspaceClient kcpkubernetesclient.Interface
}

func (c *APIBindingController) ValidateBinding(ctx context.Context, binding *apisv1alpha1.APIBinding) error {
    // R022: Validate workspace permissions
    sourceWorkspace := binding.Spec.Reference.Export.Path
    targetWorkspace := binding.ClusterName
    
    // Check if target workspace has permission to bind to source
    if !c.hasBindPermission(sourceWorkspace, targetWorkspace) {
        return fmt.Errorf("workspace %s does not have permission to bind from %s", 
            targetWorkspace, sourceWorkspace)
    }
    
    // Validate API compatibility
    export, err := c.getAPIExport(ctx, sourceWorkspace, binding.Spec.Reference.Export.Name)
    if err != nil {
        return fmt.Errorf("failed to get API export: %w", err)
    }
    
    return c.validateSchemaCompatibility(export, binding)
}
```

## ┌─ KCP-Specific Testing Patterns [R031-R040] ─┐
### │ Multi-Workspace Test Strategies              │
└─────────────────────────────────────────────────┘

**Rule R031**: Tests MUST create isolated workspace environments
**Rule R032**: Test cleanup MUST properly handle workspace resources
**Rule R033**: Integration tests MUST validate cross-workspace interactions

### Workspace-Isolated Test Pattern

```go
func TestWorkspaceIsolation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Create test server with multiple workspaces
    server := framework.SharedKcpServer(t)
    
    // Create isolated workspaces for test
    workspace1, err := server.NewOrganizationFixture(t, "test-org-1")
    require.NoError(t, err)
    
    workspace2, err := server.NewOrganizationFixture(t, "test-org-2")
    require.NoError(t, err)
    
    // Test resource isolation
    t.Run("resources are isolated between workspaces", func(t *testing.T) {
        // Create resource in workspace1
        resource1 := &corev1.ConfigMap{
            ObjectMeta: metav1.ObjectMeta{
                Name: "test-config",
                Namespace: "default",
            },
            Data: map[string]string{"key": "value1"},
        }
        
        client1, err := workspace1.GetClient()
        require.NoError(t, err)
        
        err = client1.Create(ctx, resource1)
        require.NoError(t, err)
        
        // Verify resource is NOT visible in workspace2
        client2, err := workspace2.GetClient()
        require.NoError(t, err)
        
        var retrieved corev1.ConfigMap
        err = client2.Get(ctx, client.ObjectKey{
            Name: "test-config",
            Namespace: "default",
        }, &retrieved)
        
        require.True(t, errors.IsNotFound(err), "resource should not be found in different workspace")
    })
    
    // Test cleanup
    t.Cleanup(func() {
        workspace1.Cleanup(t)
        workspace2.Cleanup(t)
    })
}
```

## ┌─ Performance Patterns for KCP [R041-R050] ─┐
### │ Scaling and Resource Management              │
└─────────────────────────────────────────────────┘

**Rule R041**: Controllers MUST implement workspace-aware caching
**Rule R042**: Cross-workspace operations MUST be batched efficiently
**Rule R043**: Resource watches MUST be scoped to relevant workspaces

### Workspace-Aware Caching Pattern

```go
type WorkspaceAwareCache struct {
    cache map[string]cache.Store  // workspace -> resources
    mu    sync.RWMutex
    informerFactory dynamicinformer.DynamicSharedInformerFactory
}

func (c *WorkspaceAwareCache) GetResourceInWorkspace(workspace, gvr, namespace, name string) (runtime.Object, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    workspaceCache, exists := c.cache[workspace]
    if !exists {
        return nil, fmt.Errorf("workspace %s not found in cache", workspace)
    }
    
    key := fmt.Sprintf("%s/%s", namespace, name)
    obj, exists, err := workspaceCache.GetByKey(key)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, errors.NewNotFound(schema.GroupResource{Resource: gvr}, name)
    }
    
    return obj.(runtime.Object), nil
}

func (c *WorkspaceAwareCache) StartWorkspaceInformer(workspace string, gvr schema.GroupVersionResource) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Create workspace-scoped informer
    workspaceInformer := c.informerFactory.ForResource(gvr).Informer()
    
    // Configure to watch only this workspace
    workspaceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            c.handleWorkspaceEvent(workspace, "add", obj)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            c.handleWorkspaceEvent(workspace, "update", newObj)
        },
        DeleteFunc: func(obj interface{}) {
            c.handleWorkspaceEvent(workspace, "delete", obj)
        },
    })
    
    go workspaceInformer.Run(c.ctx.Done())
    return nil
}
```

## ┌─ Pattern Detection Queries [R051-R055] ─┐
### │ Automated KCP Pattern Validation        │
└───────────────────────────────────────────────┘

### Detect Missing Workspace Context
```bash
# Query: Find controllers that don't handle workspace context
grep -r "client.Client" --include="*.go" | grep -v "workspace" | head -10
```

### Detect Improper Resource Access
```bash
# Query: Find direct global client usage
grep -r "\.Get\|\.List\|\.Create" --include="*.go" | grep -v "workspace\|cluster" | head -10
```

### Detect Missing Multi-Tenancy
```bash
# Query: Find resources without proper isolation
grep -r "type.*Controller struct" -A 10 --include="*.go" | grep -v "workspace\|tenant"
```

## ┌─ Validation Criteria [R056-R060] ─┐
### │ KCP Implementation Requirements   │
└─────────────────────────────────────────┘

**Validation R056**: Every controller MUST have workspace-aware client initialization
**Validation R057**: All resource operations MUST include workspace context
**Validation R058**: Cross-workspace access MUST use proper API bindings
**Validation R059**: Tests MUST demonstrate workspace isolation
**Validation R060**: Performance MUST scale with workspace count

### KCP Implementation Checklist

- [ ] Controller implements workspace-aware client pattern
- [ ] All resource operations include workspace context
- [ ] Logical cluster lifecycle is properly managed
- [ ] API exports/bindings are correctly implemented
- [ ] Tests validate workspace isolation
- [ ] Caching is workspace-aware for performance
- [ ] Cross-workspace operations use KCP mechanisms
- [ ] Resource cleanup handles workspace boundaries
- [ ] Error handling includes workspace context
- [ ] Monitoring includes workspace-specific metrics

## ┌─ Common KCP Anti-Patterns to Avoid [R061-R065] ─┐
### │ Mistakes That Break Multi-Tenancy              │
└───────────────────────────────────────────────────┘

**Anti-Pattern R061**: Using global clients for workspace-specific operations
**Anti-Pattern R062**: Hard-coding workspace names or paths
**Anti-Pattern R063**: Bypassing KCP's authorization mechanisms
**Anti-Pattern R064**: Creating cross-workspace references without API bindings
**Anti-Pattern R065**: Ignoring workspace deletion in resource cleanup

```go
// WRONG: Global client ignores workspace boundaries
func BadController() {
    client := kubernetes.NewForConfig(config)  // Global client
    pods, _ := client.CoreV1().Pods("").List()  // Lists across ALL workspaces
}

// CORRECT: Workspace-aware client
func GoodController(workspace string) {
    config.Host = config.Host + "/clusters/" + workspace
    client := kubernetes.NewForConfig(config)  // Workspace-scoped client
    pods, _ := client.CoreV1().Pods("default").List()  // Lists only in workspace
}
```

This expertise module provides comprehensive patterns for KCP development, focusing on multi-tenant workspace isolation, proper resource management, and scalable controller patterns essential for KCP-based systems.