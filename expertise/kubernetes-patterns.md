# Kubernetes Controller Patterns Expertise Module

## ┌─ Controller Architecture Patterns [R066-R075] ─┐
### │ Standard Kubernetes Controller Design          │
└───────────────────────────────────────────────────┘

**Rule R066**: Controllers MUST follow the controller-runtime patterns
**Rule R067**: Reconcile functions MUST be idempotent and stateless
**Rule R068**: Status updates MUST be separated from spec changes
**Rule R069**: Controllers MUST implement proper error handling with exponential backoff
**Rule R070**: Finalizers MUST be used for proper resource cleanup

### Standard Controller Pattern

```go
// Reconciler implements the core controller pattern
type MyResourceReconciler struct {
    client.Client
    Log    logr.Logger
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mygroup.example.com,resources=myresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mygroup.example.com,resources=myresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mygroup.example.com,resources=myresources/finalizers,verbs=update

func (r *MyResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("myresource", req.NamespacedName)

    // Fetch the resource instance
    var resource myapiv1.MyResource
    if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
        if errors.IsNotFound(err) {
            // Resource was deleted, nothing to do
            log.Info("Resource not found, likely deleted")
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, err
    }

    // Handle deletion
    if !resource.ObjectMeta.DeletionTimestamp.IsZero() {
        return r.reconcileDelete(ctx, &resource)
    }

    // Ensure finalizer is present
    if !controllerutil.ContainsFinalizer(&resource, myapiv1.MyResourceFinalizer) {
        controllerutil.AddFinalizer(&resource, myapiv1.MyResourceFinalizer)
        return ctrl.Result{}, r.Update(ctx, &resource)
    }

    // Main reconciliation logic
    return r.reconcileNormal(ctx, &resource)
}

func (r *MyResourceReconciler) reconcileNormal(ctx context.Context, resource *myapiv1.MyResource) (ctrl.Result, error) {
    log := r.Log.WithValues("myresource", resource.Name, "namespace", resource.Namespace)
    
    // R067: Idempotent reconciliation - check current state
    currentState, err := r.getCurrentState(ctx, resource)
    if err != nil {
        return ctrl.Result{}, err
    }
    
    desiredState := r.getDesiredState(resource)
    
    // Only make changes if current state differs from desired
    if !r.statesEqual(currentState, desiredState) {
        if err := r.applyDesiredState(ctx, resource, desiredState); err != nil {
            // R069: Proper error handling
            r.updateStatusWithError(ctx, resource, err)
            return ctrl.Result{RequeueAfter: time.Minute}, err
        }
    }
    
    // R068: Update status separately from spec
    return r.updateStatus(ctx, resource, currentState)
}
```

### Status Update Pattern

```go
func (r *MyResourceReconciler) updateStatus(ctx context.Context, resource *myapiv1.MyResource, state StateInfo) (ctrl.Result, error) {
    // Create a copy for status update to avoid conflicts
    statusResource := resource.DeepCopy()
    
    // Update status fields
    statusResource.Status.Phase = state.Phase
    statusResource.Status.Conditions = r.buildConditions(state)
    statusResource.Status.ObservedGeneration = resource.Generation
    
    // R068: Status subresource update
    if err := r.Status().Update(ctx, statusResource); err != nil {
        return ctrl.Result{}, err
    }
    
    // Determine requeue strategy based on state
    switch state.Phase {
    case myapiv1.PhaseProgressing:
        return ctrl.Result{RequeueAfter: time.Second * 30}, nil
    case myapiv1.PhaseReady:
        return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
    case myapiv1.PhaseFailed:
        return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
    default:
        return ctrl.Result{}, nil
    }
}
```

## ┌─ Resource Lifecycle Management [R076-R085] ─┐
### │ Proper Creation, Update, and Deletion      │
└─────────────────────────────────────────────────┘

**Rule R076**: Resources MUST implement proper owner references
**Rule R077**: Dependent resources MUST be managed through controller patterns
**Rule R078**: Cleanup MUST be guaranteed through finalizers
**Rule R079**: Updates MUST handle resource version conflicts
**Rule R080**: Creation MUST be idempotent with proper conflict resolution

### Owner Reference Pattern

```go
func (r *MyResourceReconciler) createDependentResource(ctx context.Context, owner *myapiv1.MyResource) error {
    dependent := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("%s-config", owner.Name),
            Namespace: owner.Namespace,
            Labels: map[string]string{
                "app.kubernetes.io/name":       owner.Name,
                "app.kubernetes.io/managed-by": "myresource-controller",
            },
        },
        Data: map[string]string{
            "config.yaml": owner.Spec.Configuration,
        },
    }
    
    // R076: Set owner reference for garbage collection
    if err := controllerutil.SetControllerReference(owner, dependent, r.Scheme); err != nil {
        return fmt.Errorf("failed to set controller reference: %w", err)
    }
    
    // R080: Idempotent creation with conflict resolution
    if err := r.Create(ctx, dependent); err != nil {
        if errors.IsAlreadyExists(err) {
            // Resource exists, check if update is needed
            var existing corev1.ConfigMap
            if err := r.Get(ctx, client.ObjectKeyFromObject(dependent), &existing); err != nil {
                return err
            }
            
            // Update if content differs
            if !reflect.DeepEqual(existing.Data, dependent.Data) {
                existing.Data = dependent.Data
                return r.Update(ctx, &existing)
            }
            return nil
        }
        return err
    }
    
    return nil
}
```

### Finalizer Cleanup Pattern

```go
func (r *MyResourceReconciler) reconcileDelete(ctx context.Context, resource *myapiv1.MyResource) (ctrl.Result, error) {
    log := r.Log.WithValues("myresource", resource.Name, "namespace", resource.Namespace)
    
    if !controllerutil.ContainsFinalizer(resource, myapiv1.MyResourceFinalizer) {
        // No finalizer, nothing to clean up
        return ctrl.Result{}, nil
    }
    
    // R078: Perform cleanup operations
    if err := r.cleanupExternalResources(ctx, resource); err != nil {
        log.Error(err, "Failed to cleanup external resources")
        // Don't remove finalizer if cleanup failed
        return ctrl.Result{RequeueAfter: time.Minute}, err
    }
    
    // Cleanup successful, remove finalizer
    controllerutil.RemoveFinalizer(resource, myapiv1.MyResourceFinalizer)
    if err := r.Update(ctx, resource); err != nil {
        return ctrl.Result{}, err
    }
    
    log.Info("Successfully cleaned up resource")
    return ctrl.Result{}, nil
}

func (r *MyResourceReconciler) cleanupExternalResources(ctx context.Context, resource *myapiv1.MyResource) error {
    // Clean up external resources that don't have owner references
    // Examples: external API calls, cloud resources, etc.
    
    // Delete external service
    if resource.Status.ExternalServiceID != "" {
        if err := r.deleteExternalService(resource.Status.ExternalServiceID); err != nil {
            return fmt.Errorf("failed to delete external service: %w", err)
        }
    }
    
    // Clean up any other external dependencies
    return nil
}
```

## ┌─ Event Handling and Watching [R086-R095] ─┐
### │ Efficient Resource Watching and Processing │
└─────────────────────────────────────────────────┘

**Rule R086**: Controllers MUST use proper predicates to filter events
**Rule R087**: Watch predicates MUST reduce unnecessary reconciliations
**Rule R088**: Controllers MUST handle watch errors gracefully
**Rule R089**: Event handlers MUST be non-blocking
**Rule R090**: Cross-namespace watching MUST be explicitly controlled

### Predicate Filtering Pattern

```go
func (r *MyResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    // R086: Use predicates to filter events
    return ctrl.NewControllerManagedBy(mgr).
        For(&myapiv1.MyResource{}).
        Owns(&corev1.ConfigMap{}).
        Owns(&appsv1.Deployment{}).
        WithOptions(controller.Options{
            MaxConcurrentReconciles: 3,
        }).
        WithEventFilter(r.buildPredicates()).
        Complete(r)
}

func (r *MyResourceReconciler) buildPredicates() predicate.Predicate {
    return predicate.Funcs{
        // R087: Reduce unnecessary reconciliations
        UpdateFunc: func(e event.UpdateEvent) bool {
            oldResource := e.ObjectOld.(*myapiv1.MyResource)
            newResource := e.ObjectNew.(*myapiv1.MyResource)
            
            // Only reconcile on meaningful changes
            return oldResource.Generation != newResource.Generation ||
                   !reflect.DeepEqual(oldResource.Spec, newResource.Spec) ||
                   oldResource.DeletionTimestamp != newResource.DeletionTimestamp
        },
        CreateFunc: func(e event.CreateEvent) bool {
            return true  // Always reconcile new resources
        },
        DeleteFunc: func(e event.DeleteEvent) bool {
            return false  // Deletion handled by finalizers
        },
        GenericFunc: func(e event.GenericEvent) bool {
            return false  // No generic events needed
        },
    }
}
```

### Multi-Resource Watching Pattern

```go
func (r *MyResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    controller := ctrl.NewControllerManagedBy(mgr).
        For(&myapiv1.MyResource{}).
        Owns(&corev1.ConfigMap{}).
        Owns(&appsv1.Deployment{})
    
    // R090: Watch related resources with proper mapping
    err := controller.Watches(
        &source.Kind{Type: &corev1.Secret{}},
        handler.EnqueueRequestsFromMapFunc(r.mapSecretToMyResource),
        builder.WithPredicates(r.secretPredicate()),
    ).Complete(r)
    
    return err
}

func (r *MyResourceReconciler) mapSecretToMyResource(obj client.Object) []reconcile.Request {
    secret := obj.(*corev1.Secret)
    
    // Find MyResource instances that reference this secret
    var resources myapiv1.MyResourceList
    if err := r.List(context.Background(), &resources, client.InNamespace(secret.Namespace)); err != nil {
        return nil
    }
    
    var requests []reconcile.Request
    for _, resource := range resources.Items {
        if resource.Spec.SecretRef != nil && resource.Spec.SecretRef.Name == secret.Name {
            requests = append(requests, reconcile.Request{
                NamespacedName: types.NamespacedName{
                    Name:      resource.Name,
                    Namespace: resource.Namespace,
                },
            })
        }
    }
    
    return requests
}
```

## ┌─ Error Handling and Retry Strategies [R096-R105] ─┐
### │ Robust Error Management                           │
└─────────────────────────────────────────────────────┘

**Rule R096**: Errors MUST be categorized as recoverable vs. non-recoverable
**Rule R097**: Retry backoff MUST use exponential strategy with jitter
**Rule R098**: Persistent errors MUST be surfaced through status conditions
**Rule R099**: Temporary failures MUST not block other reconciliations
**Rule R100**: Error metrics MUST be exposed for monitoring

### Error Classification and Retry Pattern

```go
type ReconcileError struct {
    Err         error
    Recoverable bool
    RetryAfter  time.Duration
}

func (r *MyResourceReconciler) reconcileNormal(ctx context.Context, resource *myapiv1.MyResource) (ctrl.Result, error) {
    // Attempt reconciliation with proper error handling
    if err := r.ensureDeployment(ctx, resource); err != nil {
        reconcileErr := r.classifyError(err)
        
        // R098: Update status with error information
        r.updateStatusWithError(ctx, resource, reconcileErr)
        
        if reconcileErr.Recoverable {
            // R097: Exponential backoff for recoverable errors
            backoff := r.calculateBackoff(resource.Status.RetryCount)
            
            // Update retry count
            resource.Status.RetryCount++
            r.Status().Update(ctx, resource)
            
            return ctrl.Result{RequeueAfter: backoff}, reconcileErr.Err
        } else {
            // Non-recoverable error - don't retry automatically
            return ctrl.Result{}, reconcileErr.Err
        }
    }
    
    // Success - reset retry count
    if resource.Status.RetryCount > 0 {
        resource.Status.RetryCount = 0
        r.Status().Update(ctx, resource)
    }
    
    return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

func (r *MyResourceReconciler) classifyError(err error) ReconcileError {
    // R096: Categorize errors
    switch {
    case errors.IsConflict(err):
        // Resource version conflicts are recoverable
        return ReconcileError{
            Err:         err,
            Recoverable: true,
            RetryAfter:  time.Second * 5,
        }
    case errors.IsServerTimeout(err), errors.IsTimeout(err):
        // Timeouts are usually recoverable
        return ReconcileError{
            Err:         err,
            Recoverable: true,
            RetryAfter:  time.Second * 30,
        }
    case errors.IsInternalError(err), errors.IsServiceUnavailable(err):
        // Server errors are recoverable
        return ReconcileError{
            Err:         err,
            Recoverable: true,
            RetryAfter:  time.Minute,
        }
    case errors.IsForbidden(err), errors.IsUnauthorized(err):
        // Permission errors are typically not recoverable without manual intervention
        return ReconcileError{
            Err:         err,
            Recoverable: false,
        }
    case errors.IsInvalid(err), errors.IsBadRequest(err):
        // Invalid specs are not recoverable without user changes
        return ReconcileError{
            Err:         err,
            Recoverable: false,
        }
    default:
        // Unknown errors - assume recoverable but with longer backoff
        return ReconcileError{
            Err:         err,
            Recoverable: true,
            RetryAfter:  time.Minute * 2,
        }
    }
}

func (r *MyResourceReconciler) calculateBackoff(retryCount int) time.Duration {
    // R097: Exponential backoff with jitter
    base := time.Second * 5
    backoff := time.Duration(math.Pow(2, float64(retryCount))) * base
    
    // Cap maximum backoff
    if backoff > time.Minute*10 {
        backoff = time.Minute * 10
    }
    
    // Add jitter to prevent thundering herd
    jitter := time.Duration(rand.Intn(int(backoff.Nanoseconds()/4))) * time.Nanosecond
    return backoff + jitter
}
```

## ┌─ Testing Patterns for Controllers [R106-R115] ─┐
### │ Comprehensive Controller Testing               │
└─────────────────────────────────────────────────────┘

**Rule R106**: Tests MUST use envtest for integration testing
**Rule R107**: Unit tests MUST mock external dependencies
**Rule R108**: Tests MUST verify both happy path and error scenarios
**Rule R109**: Tests MUST validate status updates and conditions
**Rule R110**: Tests MUST ensure proper cleanup and finalizer handling

### EnvTest Integration Pattern

```go
var _ = Describe("MyResource Controller", func() {
    var (
        ctx       context.Context
        cancel    context.CancelFunc
        k8sClient client.Client
        testEnv   *envtest.Environment
        reconciler *MyResourceReconciler
    )

    BeforeEach(func() {
        ctx, cancel = context.WithCancel(context.Background())
        
        // R106: Use envtest for integration testing
        testEnv = &envtest.Environment{
            CRDDirectoryPaths: []string{
                filepath.Join("..", "config", "crd", "bases"),
            },
        }

        cfg, err := testEnv.Start()
        Expect(err).NotTo(HaveOccurred())

        scheme := runtime.NewScheme()
        err = myapiv1.AddToScheme(scheme)
        Expect(err).NotTo(HaveOccurred())
        
        k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
        Expect(err).NotTo(HaveOccurred())

        // Setup reconciler
        reconciler = &MyResourceReconciler{
            Client: k8sClient,
            Scheme: scheme,
            Log:    ctrl.Log.WithName("test"),
        }
    })

    AfterEach(func() {
        cancel()
        err := testEnv.Stop()
        Expect(err).NotTo(HaveOccurred())
    })

    Context("When creating a MyResource", func() {
        It("Should create dependent resources", func() {
            // Create test resource
            resource := &myapiv1.MyResource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "test-resource",
                    Namespace: "default",
                },
                Spec: myapiv1.MyResourceSpec{
                    Configuration: "test-config",
                },
            }

            err := k8sClient.Create(ctx, resource)
            Expect(err).NotTo(HaveOccurred())

            // Trigger reconciliation
            _, err = reconciler.Reconcile(ctx, ctrl.Request{
                NamespacedName: types.NamespacedName{
                    Name:      resource.Name,
                    Namespace: resource.Namespace,
                },
            })
            Expect(err).NotTo(HaveOccurred())

            // Verify dependent ConfigMap was created
            var configMap corev1.ConfigMap
            err = k8sClient.Get(ctx, types.NamespacedName{
                Name:      fmt.Sprintf("%s-config", resource.Name),
                Namespace: resource.Namespace,
            }, &configMap)
            Expect(err).NotTo(HaveOccurred())
            Expect(configMap.Data["config.yaml"]).To(Equal("test-config"))

            // R109: Verify status updates
            var updatedResource myapiv1.MyResource
            err = k8sClient.Get(ctx, client.ObjectKeyFromObject(resource), &updatedResource)
            Expect(err).NotTo(HaveOccurred())
            Expect(updatedResource.Status.Phase).To(Equal(myapiv1.PhaseReady))
        })

        // R108: Test error scenarios
        It("Should handle invalid configuration", func() {
            resource := &myapiv1.MyResource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "invalid-resource",
                    Namespace: "default",
                },
                Spec: myapiv1.MyResourceSpec{
                    Configuration: "", // Invalid empty config
                },
            }

            err := k8sClient.Create(ctx, resource)
            Expect(err).NotTo(HaveOccurred())

            _, err = reconciler.Reconcile(ctx, ctrl.Request{
                NamespacedName: client.ObjectKeyFromObject(resource),
            })
            Expect(err).To(HaveOccurred())

            // Verify error is reflected in status
            var updatedResource myapiv1.MyResource
            err = k8sClient.Get(ctx, client.ObjectKeyFromObject(resource), &updatedResource)
            Expect(err).NotTo(HaveOccurred())
            Expect(updatedResource.Status.Phase).To(Equal(myapiv1.PhaseFailed))
        })

        // R110: Test finalizer handling
        It("Should clean up properly when deleted", func() {
            // Setup test will be here
        })
    })
})
```

## ┌─ Performance and Scaling Patterns [R116-R125] ─┐
### │ Controller Optimization Strategies              │
└───────────────────────────────────────────────────┘

**Rule R116**: Controllers MUST implement proper caching strategies
**Rule R117**: Reconciliation MUST be limited by MaxConcurrentReconciles
**Rule R118**: Large resource lists MUST use pagination
**Rule R119**: Expensive operations MUST be rate limited
**Rule R120**: Memory usage MUST be monitored and bounded

### Efficient Resource Listing Pattern

```go
func (r *MyResourceReconciler) listRelatedResources(ctx context.Context, resource *myapiv1.MyResource) ([]corev1.ConfigMap, error) {
    var configMaps corev1.ConfigMapList
    
    // R118: Use pagination for large lists
    listOptions := []client.ListOption{
        client.InNamespace(resource.Namespace),
        client.MatchingLabels(map[string]string{
            "app.kubernetes.io/managed-by": "myresource-controller",
            "app.kubernetes.io/instance":   resource.Name,
        }),
    }
    
    // Implement pagination for large result sets
    var allConfigMaps []corev1.ConfigMap
    continueToken := ""
    
    for {
        opts := append(listOptions, client.Continue(continueToken), client.Limit(100))
        
        if err := r.List(ctx, &configMaps, opts...); err != nil {
            return nil, err
        }
        
        allConfigMaps = append(allConfigMaps, configMaps.Items...)
        
        if configMaps.Continue == "" {
            break
        }
        continueToken = configMaps.Continue
    }
    
    return allConfigMaps, nil
}
```

## ┌─ Pattern Detection Queries [R126-R130] ─┐
### │ Automated Pattern Validation            │
└───────────────────────────────────────────────┘

### Detect Missing Error Handling
```bash
# Query: Find reconcile functions without proper error handling
grep -r "func.*Reconcile" --include="*.go" -A 20 | grep -L "if err" | head -10
```

### Detect Missing Finalizers
```bash
# Query: Find controllers that don't implement finalizers
grep -r "controllerutil.AddFinalizer\|controllerutil.RemoveFinalizer" --include="*.go" -L | head -10
```

### Detect Improper Status Updates
```bash
# Query: Find status updates not using Status() subresource
grep -r "\.Update.*Status" --include="*.go" | grep -v "\.Status().Update" | head -10
```

## ┌─ Validation Criteria [R131-R135] ─┐
### │ Controller Implementation Quality │
└─────────────────────────────────────────┘

**Validation R131**: Controller MUST implement proper RBAC annotations
**Validation R132**: Reconcile function MUST be idempotent
**Validation R133**: Error handling MUST provide useful diagnostics
**Validation R134**: Status conditions MUST follow Kubernetes conventions
**Validation R135**: Tests MUST cover all major reconciliation paths

This comprehensive Kubernetes patterns module provides the essential controller patterns, error handling strategies, and testing approaches needed for robust Kubernetes operator development.