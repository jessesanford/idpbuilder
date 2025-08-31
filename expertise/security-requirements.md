# Security Requirements and Validation Patterns Expertise Module

## ┌─ RBAC and Authorization Patterns [R191-R200] ─┐
### │ Kubernetes Role-Based Access Control         │
└─────────────────────────────────────────────────────┘

**Rule R191**: All controllers MUST define minimal required RBAC permissions
**Rule R192**: Service accounts MUST follow principle of least privilege
**Rule R193**: Cross-namespace access MUST be explicitly authorized
**Rule R194**: RBAC permissions MUST be validated in tests
**Rule R195**: Cluster-wide permissions MUST be justified and documented

### RBAC Permission Definition Pattern

```go
// R191: Minimal required RBAC permissions with kubebuilder annotations
// +kubebuilder:rbac:groups=mygroup.example.com,resources=myresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mygroup.example.com,resources=myresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mygroup.example.com,resources=myresources/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *MyResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Controller implementation with minimal permissions
    return ctrl.Result{}, nil
}

// Generate RBAC manifests with proper scoping
//go:generate controller-gen rbac:roleName=manager-role paths="./..."

// R192: Service account with minimal privileges
apiVersion: v1
kind: ServiceAccount
metadata:
  name: myresource-controller
  namespace: system
automountServiceAccountToken: true  # Only if needed
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: myresource-manager-role
rules:
# Minimal permissions for MyResource management
- apiGroups:
  - mygroup.example.com
  resources:
  - myresources
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - mygroup.example.com
  resources:
  - myresources/finalizers
  verbs:
  - update
- apiGroups:
  - mygroup.example.com
  resources:
  - myresources/status
  verbs:
  - get
  - patch
  - update
# Dependent resource permissions
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
# R193: Explicit cross-namespace permissions (if needed)
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
  resourceNames:  # Restrict to specific secrets only
  - webhook-cert
  - ca-bundle
```

### RBAC Validation Testing Pattern

```go
// R194: RBAC validation in tests
func TestRBACPermissions(t *testing.T) {
    ctx := context.Background()
    
    testCases := []struct {
        name           string
        serviceAccount string
        namespace      string
        resource       client.Object
        verb           string
        shouldAllow    bool
    }{
        {
            name:           "controller_should_create_myresource",
            serviceAccount: "myresource-controller",
            namespace:      "default",
            resource:       &myapiv1.MyResource{},
            verb:           "create",
            shouldAllow:    true,
        },
        {
            name:           "controller_should_not_create_secrets",
            serviceAccount: "myresource-controller",
            namespace:      "default",
            resource:       &corev1.Secret{},
            verb:           "create",
            shouldAllow:    false,
        },
        {
            name:           "controller_should_not_access_other_namespaces",
            serviceAccount: "myresource-controller",
            namespace:      "kube-system",
            resource:       &myapiv1.MyResource{},
            verb:           "get",
            shouldAllow:    false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create SubjectAccessReview
            sar := &authorizationv1.SubjectAccessReview{
                Spec: authorizationv1.SubjectAccessReviewSpec{
                    User: fmt.Sprintf("system:serviceaccount:%s:%s", "system", tc.serviceAccount),
                    ResourceAttributes: &authorizationv1.ResourceAttributes{
                        Namespace: tc.namespace,
                        Verb:      tc.verb,
                        Group:     tc.resource.GetObjectKind().GroupVersionKind().Group,
                        Version:   tc.resource.GetObjectKind().GroupVersionKind().Version,
                        Resource:  strings.ToLower(tc.resource.GetObjectKind().GroupVersionKind().Kind) + "s",
                    },
                },
            }

            err := k8sClient.Create(ctx, sar)
            require.NoError(t, err)

            if tc.shouldAllow {
                assert.True(t, sar.Status.Allowed, 
                    "Expected permission to be allowed for %s", tc.name)
            } else {
                assert.False(t, sar.Status.Allowed, 
                    "Expected permission to be denied for %s", tc.name)
            }
        })
    }
}
```

## ┌─ Secret and Credential Management [R201-R210] ─┐
### │ Secure Handling of Sensitive Data             │
└─────────────────────────────────────────────────────┘

**Rule R201**: Secrets MUST never be logged or exposed in status
**Rule R202**: Credential rotation MUST be supported
**Rule R203**: Secret access MUST be audited
**Rule R204**: Secrets MUST be encrypted at rest
**Rule R205**: Default credentials MUST be changed

### Secure Secret Handling Pattern

```go
// R201: Secure secret handling without exposure
func (r *MyResourceReconciler) reconcileSecrets(ctx context.Context, resource *myapiv1.MyResource) error {
    log := r.Log.WithValues("myresource", resource.Name, "namespace", resource.Namespace)
    
    // Fetch secret securely
    var secret corev1.Secret
    secretRef := resource.Spec.SecretRef
    if secretRef == nil {
        return fmt.Errorf("secretRef is required")
    }
    
    err := r.Get(ctx, types.NamespacedName{
        Name:      secretRef.Name,
        Namespace: resource.Namespace, // Secrets must be in same namespace
    }, &secret)
    if err != nil {
        // R201: Never log secret content or detailed errors that might expose secrets
        log.Error(err, "Failed to retrieve secret", "secretName", secretRef.Name)
        return fmt.Errorf("failed to retrieve secret: %w", err)
    }
    
    // Validate required secret keys
    requiredKeys := []string{"username", "password", "api-token"}
    for _, key := range requiredKeys {
        if _, exists := secret.Data[key]; !exists {
            return fmt.Errorf("required secret key missing: %s", key)
        }
    }
    
    // R203: Audit secret access
    r.recordSecretAccess(ctx, resource, secretRef.Name)
    
    // Use secret data without exposing it
    return r.configureWithSecret(ctx, resource, secret.Data)
}

func (r *MyResourceReconciler) configureWithSecret(ctx context.Context, resource *myapiv1.MyResource, secretData map[string][]byte) error {
    // Process secret data securely
    // Never log or store secret values
    
    // Example: Create deployment with secret mounted as volume
    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      resource.Name + "-deployment",
            Namespace: resource.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Template: corev1.PodTemplateSpec{
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{{
                        Name:  "app",
                        Image: resource.Spec.Image,
                        // Mount secret as volume, not environment variable
                        VolumeMounts: []corev1.VolumeMount{{
                            Name:      "credentials",
                            MountPath: "/etc/credentials",
                            ReadOnly:  true,
                        }},
                    }},
                    Volumes: []corev1.Volume{{
                        Name: "credentials",
                        VolumeSource: corev1.VolumeSource{
                            Secret: &corev1.SecretVolumeSource{
                                SecretName: resource.Spec.SecretRef.Name,
                                // R204: Use defaultMode to restrict permissions
                                DefaultMode: int32Ptr(0400), // Read-only for owner
                            },
                        },
                    }},
                },
            },
        },
    }
    
    return r.Create(ctx, deployment)
}

func int32Ptr(i int32) *int32 { return &i }

// R203: Audit secret access
func (r *MyResourceReconciler) recordSecretAccess(ctx context.Context, resource *myapiv1.MyResource, secretName string) {
    event := &corev1.Event{
        ObjectMeta: metav1.ObjectMeta{
            GenerateName: "secret-access-",
            Namespace:    resource.Namespace,
        },
        InvolvedObject: corev1.ObjectReference{
            APIVersion: resource.APIVersion,
            Kind:       resource.Kind,
            Name:       resource.Name,
            Namespace:  resource.Namespace,
            UID:        resource.UID,
        },
        Type:    corev1.EventTypeNormal,
        Reason:  "SecretAccessed",
        Message: fmt.Sprintf("Accessed secret %s for configuration", secretName),
        Source: corev1.EventSource{
            Component: "myresource-controller",
        },
        FirstTimestamp: metav1.Now(),
        LastTimestamp:  metav1.Now(),
        Count:          1,
    }
    
    // Create audit event
    r.Create(ctx, event)
}
```

### Secret Rotation Support Pattern

```go
// R202: Support credential rotation
func (r *MyResourceReconciler) handleSecretRotation(ctx context.Context, resource *myapiv1.MyResource) (ctrl.Result, error) {
    // Watch for secret changes
    var secret corev1.Secret
    err := r.Get(ctx, types.NamespacedName{
        Name:      resource.Spec.SecretRef.Name,
        Namespace: resource.Namespace,
    }, &secret)
    if err != nil {
        return ctrl.Result{}, err
    }
    
    // Check if secret has been updated since last reconciliation
    lastSecretVersion := resource.Status.LastSecretResourceVersion
    if secret.ResourceVersion != lastSecretVersion {
        log := r.Log.WithValues("myresource", resource.Name)
        log.Info("Secret rotation detected, updating configuration")
        
        // Reconfigure with new secret
        if err := r.configureWithSecret(ctx, resource, secret.Data); err != nil {
            return ctrl.Result{}, err
        }
        
        // Update status with new secret version
        resource.Status.LastSecretResourceVersion = secret.ResourceVersion
        if err := r.Status().Update(ctx, resource); err != nil {
            return ctrl.Result{}, err
        }
        
        // Restart dependent workloads to pick up new secrets
        return r.restartWorkloads(ctx, resource)
    }
    
    return ctrl.Result{}, nil
}

func (r *MyResourceReconciler) restartWorkloads(ctx context.Context, resource *myapiv1.MyResource) (ctrl.Result, error) {
    // Rolling restart of deployment to pick up new secrets
    var deployment appsv1.Deployment
    err := r.Get(ctx, types.NamespacedName{
        Name:      resource.Name + "-deployment",
        Namespace: resource.Namespace,
    }, &deployment)
    if err != nil {
        return ctrl.Result{}, err
    }
    
    // Update deployment with restart annotation
    if deployment.Spec.Template.Annotations == nil {
        deployment.Spec.Template.Annotations = make(map[string]string)
    }
    deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
    
    err = r.Update(ctx, &deployment)
    if err != nil {
        return ctrl.Result{}, err
    }
    
    return ctrl.Result{RequeueAfter: time.Minute}, nil
}
```

## ┌─ Network Security and TLS [R211-R220] ─┐
### │ Secure Communication Patterns         │
└─────────────────────────────────────────────┘

**Rule R211**: All API communications MUST use TLS
**Rule R212**: Certificate validation MUST be enforced
**Rule R213**: Webhook configurations MUST use proper CA bundles
**Rule R214**: Service mesh security policies MUST be defined
**Rule R215**: Network policies MUST restrict pod-to-pod communication

### TLS Configuration Pattern

```go
// R211: TLS-secured webhook configuration
func (r *MyResourceReconciler) setupWebhook(ctx context.Context) error {
    // Generate or retrieve TLS certificate
    cert, err := r.getTLSCertificate(ctx)
    if err != nil {
        return fmt.Errorf("failed to get TLS certificate: %w", err)
    }
    
    // R213: Configure webhook with proper CA bundle
    webhook := &admissionregistrationv1.ValidatingAdmissionWebhook{
        ObjectMeta: metav1.ObjectMeta{
            Name: "myresource-validator",
        },
        Webhooks: []admissionregistrationv1.ValidatingWebhook{{
            Name: "validate.myresource.example.com",
            ClientConfig: admissionregistrationv1.WebhookClientConfig{
                Service: &admissionregistrationv1.ServiceReference{
                    Name:      "webhook-service",
                    Namespace: "system",
                    Path:      strPtr("/validate"),
                },
                CABundle: cert.CACert, // R213: Proper CA bundle
            },
            Rules: []admissionregistrationv1.RuleWithOperations{{
                Operations: []admissionregistrationv1.OperationType{
                    admissionregistrationv1.Create,
                    admissionregistrationv1.Update,
                },
                Rule: admissionregistrationv1.Rule{
                    APIGroups:   []string{"mygroup.example.com"},
                    APIVersions: []string{"v1"},
                    Resources:   []string{"myresources"},
                },
            }},
            AdmissionReviewVersions: []string{"v1", "v1beta1"},
            SideEffects:             sideEffectsPtr(admissionregistrationv1.SideEffectClassNone),
            FailurePolicy:           failurePolicyPtr(admissionregistrationv1.Fail),
        }},
    }
    
    return r.Create(ctx, webhook)
}

// R212: Certificate validation enforcement
func (r *MyResourceReconciler) createTLSClient(ctx context.Context, endpoint string, caCert []byte) (*http.Client, error) {
    // Create certificate pool with CA certificate
    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(caCert) {
        return nil, fmt.Errorf("failed to parse CA certificate")
    }
    
    // Configure TLS with certificate validation
    tlsConfig := &tls.Config{
        RootCAs:            caCertPool,
        InsecureSkipVerify: false, // R212: Never skip certificate validation
        MinVersion:         tls.VersionTLS12, // Enforce minimum TLS version
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        },
    }
    
    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
        // Additional security headers
        DisableKeepAlives: false,
        IdleConnTimeout:   30 * time.Second,
        MaxIdleConns:      10,
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }, nil
}
```

### Network Security Policy Pattern

```yaml
# R215: Network policies for pod-to-pod communication restriction
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: myresource-controller-netpol
  namespace: system
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: myresource-controller
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Allow webhook traffic from API server
  - from:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    - podSelector:
        matchLabels:
          component: kube-apiserver
    ports:
    - protocol: TCP
      port: 9443
  # Allow metrics scraping from monitoring
  - from:
    - namespaceSelector:
        matchLabels:
          name: monitoring
    ports:
    - protocol: TCP
      port: 8080
  egress:
  # Allow DNS resolution
  - to: []
    ports:
    - protocol: UDP
      port: 53
  # Allow API server communication
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: TCP
      port: 443
  # Allow webhook certificate access
  - to:
    - podSelector:
        matchLabels:
          app.kubernetes.io/name: cert-manager
    ports:
    - protocol: TCP
      port: 9443
---
# R214: Service mesh security policy (Istio example)
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: myresource-controller
  namespace: system
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: myresource-controller
  mtls:
    mode: STRICT
---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: myresource-controller
  namespace: system
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: myresource-controller
  rules:
  - from:
    - source:
        principals:
        - cluster.local/ns/kube-system/sa/kube-apiserver
  - to:
    - operation:
        ports: ["9443"]
        methods: ["POST"]
        paths: ["/validate", "/mutate"]
```

## ┌─ Pod Security and Runtime Security [R221-R230] ─┐
### │ Container and Runtime Security Patterns        │
└───────────────────────────────────────────────────────┘

**Rule R221**: Pods MUST run with non-root user
**Rule R222**: Security contexts MUST be restrictive
**Rule R223**: Container images MUST be scanned for vulnerabilities
**Rule R224**: Privileged containers MUST be avoided
**Rule R225**: Resource limits MUST be enforced

### Secure Pod Configuration Pattern

```go
// R221, R222: Secure pod configuration with restrictive security context
func (r *MyResourceReconciler) createSecureDeployment(ctx context.Context, resource *myapiv1.MyResource) error {
    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      resource.Name + "-deployment",
            Namespace: resource.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &resource.Spec.Replicas,
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app.kubernetes.io/name":     resource.Name,
                    "app.kubernetes.io/instance": resource.Name,
                },
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app.kubernetes.io/name":     resource.Name,
                        "app.kubernetes.io/instance": resource.Name,
                    },
                    Annotations: map[string]string{
                        // Security scanning annotations
                        "container.apparmor.security.beta.kubernetes.io/app": "runtime/default",
                        "seccomp.security.alpha.kubernetes.io/pod":           "runtime/default",
                    },
                },
                Spec: corev1.PodSpec{
                    // R222: Restrictive pod security context
                    SecurityContext: &corev1.PodSecurityContext{
                        RunAsNonRoot: boolPtr(true),          // R221: Non-root user
                        RunAsUser:    int64Ptr(10001),        // Specific non-root UID
                        RunAsGroup:   int64Ptr(10001),        // Specific non-root GID
                        FSGroup:      int64Ptr(10001),        // File system group
                        SeccompProfile: &corev1.SeccompProfile{
                            Type: corev1.SeccompProfileTypeRuntimeDefault,
                        },
                    },
                    Containers: []corev1.Container{{
                        Name:  "app",
                        Image: resource.Spec.Image, // R223: Image should be scanned
                        // R222: Restrictive container security context
                        SecurityContext: &corev1.SecurityContext{
                            AllowPrivilegeEscalation: boolPtr(false), // R224: No privilege escalation
                            Privileged:               boolPtr(false), // R224: Not privileged
                            RunAsNonRoot:             boolPtr(true),  // R221: Non-root
                            RunAsUser:                int64Ptr(10001),
                            RunAsGroup:               int64Ptr(10001),
                            ReadOnlyRootFilesystem:   boolPtr(true), // Read-only root FS
                            Capabilities: &corev1.Capabilities{
                                Drop: []corev1.Capability{"ALL"}, // Drop all capabilities
                                Add:  []corev1.Capability{},       // Add only necessary ones
                            },
                        },
                        // R225: Resource limits enforcement
                        Resources: corev1.ResourceRequirements{
                            Limits: corev1.ResourceList{
                                corev1.ResourceCPU:    resource1.MustParse("500m"),
                                corev1.ResourceMemory: resource1.MustParse("512Mi"),
                            },
                            Requests: corev1.ResourceList{
                                corev1.ResourceCPU:    resource1.MustParse("100m"),
                                corev1.ResourceMemory: resource1.MustParse("128Mi"),
                            },
                        },
                        // Liveness and readiness probes
                        LivenessProbe: &corev1.Probe{
                            ProbeHandler: corev1.ProbeHandler{
                                HTTPGet: &corev1.HTTPGetAction{
                                    Path: "/health",
                                    Port: intstr.FromInt(8080),
                                },
                            },
                            InitialDelaySeconds: 30,
                            PeriodSeconds:       10,
                        },
                        ReadinessProbe: &corev1.Probe{
                            ProbeHandler: corev1.ProbeHandler{
                                HTTPGet: &corev1.HTTPGetAction{
                                    Path: "/ready",
                                    Port: intstr.FromInt(8080),
                                },
                            },
                            InitialDelaySeconds: 5,
                            PeriodSeconds:       5,
                        },
                        // Temporary file system for read-only root FS
                        VolumeMounts: []corev1.VolumeMount{{
                            Name:      "tmp-volume",
                            MountPath: "/tmp",
                        }},
                    }},
                    Volumes: []corev1.Volume{{
                        Name: "tmp-volume",
                        VolumeSource: corev1.VolumeSource{
                            EmptyDir: &corev1.EmptyDirVolumeSource{
                                SizeLimit: resource1.NewQuantity(100*1024*1024, resource1.BinarySI), // 100MB
                            },
                        },
                    }},
                    // Additional security measures
                    ServiceAccountName:            resource.Name + "-sa",
                    AutomountServiceAccountToken:  boolPtr(false), // Only if not needed
                    HostNetwork:                   false,
                    HostPID:                       false,
                    HostIPC:                       false,
                },
            },
        },
    }
    
    return r.Create(ctx, deployment)
}

// Helper functions
func boolPtr(b bool) *bool       { return &b }
func int64Ptr(i int64) *int64    { return &i }
func strPtr(s string) *string    { return &s }
```

### Image Security Scanning Integration

```go
// R223: Container image vulnerability scanning
func (r *MyResourceReconciler) validateImageSecurity(ctx context.Context, imageRef string) error {
    // Integration with image scanning tools (e.g., Trivy, Clair, etc.)
    scanner := &ImageScanner{
        Client: r.scannerClient,
    }
    
    report, err := scanner.ScanImage(ctx, imageRef)
    if err != nil {
        return fmt.Errorf("failed to scan image %s: %w", imageRef, err)
    }
    
    // Check for critical vulnerabilities
    criticalCount := report.GetVulnerabilityCount(VulnerabilitySeverityCritical)
    if criticalCount > 0 {
        return fmt.Errorf("image %s contains %d critical vulnerabilities", imageRef, criticalCount)
    }
    
    // Check for high vulnerabilities beyond threshold
    highCount := report.GetVulnerabilityCount(VulnerabilitySeverityHigh)
    if highCount > 5 { // Configurable threshold
        return fmt.Errorf("image %s contains too many high vulnerabilities: %d", imageRef, highCount)
    }
    
    // Validate image signature if using cosign or similar
    if err := r.validateImageSignature(ctx, imageRef); err != nil {
        return fmt.Errorf("image signature validation failed: %w", err)
    }
    
    return nil
}

type ImageScanner struct {
    Client ScannerClient
}

type ScanReport interface {
    GetVulnerabilityCount(severity VulnerabilitySeverity) int
    GetVulnerabilities() []Vulnerability
}

type VulnerabilitySeverity string

const (
    VulnerabilitySeverityLow      VulnerabilitySeverity = "LOW"
    VulnerabilitySeverityMedium   VulnerabilitySeverity = "MEDIUM"
    VulnerabilitySeverityHigh     VulnerabilitySeverity = "HIGH"
    VulnerabilitySeverityCritical VulnerabilitySeverity = "CRITICAL"
)
```

## ┌─ Admission Control and Validation [R231-R240] ─┐
### │ Webhook-Based Security Enforcement            │
└─────────────────────────────────────────────────────┘

**Rule R231**: Validating webhooks MUST enforce security policies
**Rule R232**: Mutating webhooks MUST apply security defaults
**Rule R233**: Webhook failures MUST fail closed for security
**Rule R234**: Admission reviews MUST be logged for audit
**Rule R235**: Webhook timeouts MUST have reasonable defaults

### Validating Admission Webhook Pattern

```go
// R231: Security policy enforcement through validating webhooks
func (r *MyResourceWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) error {
    resource := obj.(*myapiv1.MyResource)
    
    var allErrors field.ErrorList
    
    // Validate security requirements
    if errs := r.validateSecurityContext(resource); len(errs) > 0 {
        allErrors = append(allErrors, errs...)
    }
    
    if errs := r.validateImageSecurity(ctx, resource); len(errs) > 0 {
        allErrors = append(allErrors, errs...)
    }
    
    if errs := r.validateNetworkPolicies(resource); len(errs) > 0 {
        allErrors = append(allErrors, errs...)
    }
    
    if errs := r.validateRBACRequirements(resource); len(errs) > 0 {
        allErrors = append(allErrors, errs...)
    }
    
    // R234: Log admission decision for audit
    r.logAdmissionDecision(ctx, "CREATE", resource, len(allErrors) == 0)
    
    if len(allErrors) > 0 {
        return apierrors.NewInvalid(
            schema.GroupKind{Group: myapiv1.GroupVersion.Group, Kind: "MyResource"},
            resource.Name,
            allErrors,
        )
    }
    
    return nil
}

func (r *MyResourceWebhook) validateSecurityContext(resource *myapiv1.MyResource) field.ErrorList {
    var allErrors field.ErrorList
    specPath := field.NewPath("spec")
    
    // Validate that runAsNonRoot is explicitly set
    if resource.Spec.SecurityContext == nil {
        allErrors = append(allErrors, 
            field.Required(specPath.Child("securityContext"), 
                "securityContext is required for security policy compliance"))
        return allErrors
    }
    
    secCtxPath := specPath.Child("securityContext")
    
    // Enforce non-root user
    if resource.Spec.SecurityContext.RunAsNonRoot == nil || !*resource.Spec.SecurityContext.RunAsNonRoot {
        allErrors = append(allErrors,
            field.Invalid(secCtxPath.Child("runAsNonRoot"), 
                resource.Spec.SecurityContext.RunAsNonRoot,
                "runAsNonRoot must be true for security compliance"))
    }
    
    // Enforce non-zero UID
    if resource.Spec.SecurityContext.RunAsUser != nil && *resource.Spec.SecurityContext.RunAsUser == 0 {
        allErrors = append(allErrors,
            field.Invalid(secCtxPath.Child("runAsUser"),
                resource.Spec.SecurityContext.RunAsUser,
                "runAsUser cannot be 0 (root)"))
    }
    
    // Validate that privileged mode is disabled
    if resource.Spec.SecurityContext.Privileged != nil && *resource.Spec.SecurityContext.Privileged {
        allErrors = append(allErrors,
            field.Invalid(secCtxPath.Child("privileged"),
                resource.Spec.SecurityContext.Privileged,
                "privileged containers are not allowed"))
    }
    
    return allErrors
}

func (r *MyResourceWebhook) validateImageSecurity(ctx context.Context, resource *myapiv1.MyResource) field.ErrorList {
    var allErrors field.ErrorList
    specPath := field.NewPath("spec")
    
    // Validate image registry is allowed
    allowedRegistries := []string{
        "gcr.io/my-org/",
        "registry.company.com/",
        "docker.io/library/", // Only official library images
    }
    
    imageAllowed := false
    for _, allowedRegistry := range allowedRegistries {
        if strings.HasPrefix(resource.Spec.Image, allowedRegistry) {
            imageAllowed = true
            break
        }
    }
    
    if !imageAllowed {
        allErrors = append(allErrors,
            field.Invalid(specPath.Child("image"),
                resource.Spec.Image,
                fmt.Sprintf("image must be from allowed registries: %v", allowedRegistries)))
    }
    
    // Validate image tag is not 'latest'
    if strings.HasSuffix(resource.Spec.Image, ":latest") {
        allErrors = append(allErrors,
            field.Invalid(specPath.Child("image"),
                resource.Spec.Image,
                "image tag 'latest' is not allowed for security and reproducibility"))
    }
    
    return allErrors
}

// R234: Audit logging for admission decisions
func (r *MyResourceWebhook) logAdmissionDecision(ctx context.Context, operation string, resource *myapiv1.MyResource, allowed bool) {
    logEntry := map[string]interface{}{
        "timestamp":  time.Now().UTC(),
        "operation":  operation,
        "resource":   fmt.Sprintf("%s/%s", resource.Namespace, resource.Name),
        "user":       r.getUserFromContext(ctx),
        "allowed":    allowed,
        "component":  "admission-webhook",
    }
    
    if !allowed {
        logEntry["reason"] = "security-policy-violation"
    }
    
    // Log to structured audit log
    r.auditLogger.Info("admission-decision", logEntry)
}
```

### Mutating Admission Webhook Pattern

```go
// R232: Apply security defaults through mutating webhooks
func (r *MyResourceWebhook) Default(ctx context.Context, obj runtime.Object) error {
    resource := obj.(*myapiv1.MyResource)
    
    r.logger.V(1).Info("Applying security defaults", "resource", resource.Name)
    
    // Apply default security context if not provided
    if resource.Spec.SecurityContext == nil {
        resource.Spec.SecurityContext = &myapiv1.SecurityContext{}
    }
    
    // Set secure defaults
    if resource.Spec.SecurityContext.RunAsNonRoot == nil {
        resource.Spec.SecurityContext.RunAsNonRoot = boolPtr(true)
    }
    
    if resource.Spec.SecurityContext.RunAsUser == nil {
        resource.Spec.SecurityContext.RunAsUser = int64Ptr(10001) // Default non-root user
    }
    
    if resource.Spec.SecurityContext.ReadOnlyRootFilesystem == nil {
        resource.Spec.SecurityContext.ReadOnlyRootFilesystem = boolPtr(true)
    }
    
    if resource.Spec.SecurityContext.AllowPrivilegeEscalation == nil {
        resource.Spec.SecurityContext.AllowPrivilegeEscalation = boolPtr(false)
    }
    
    // Apply default resource limits if not provided
    if resource.Spec.Resources == nil {
        resource.Spec.Resources = &myapiv1.ResourceRequirements{
            Limits: map[string]string{
                "cpu":    "500m",
                "memory": "512Mi",
            },
            Requests: map[string]string{
                "cpu":    "100m",
                "memory": "128Mi",
            },
        }
    }
    
    // Apply default network policy if not specified
    if resource.Spec.NetworkPolicy == nil {
        resource.Spec.NetworkPolicy = &myapiv1.NetworkPolicy{
            Enabled: true,
            AllowedIngress: []string{}, // Default to no ingress
            AllowedEgress:  []string{"dns", "kube-apiserver"}, // Minimal egress
        }
    }
    
    // R234: Log mutation for audit
    r.logMutation(ctx, resource)
    
    return nil
}

func (r *MyResourceWebhook) logMutation(ctx context.Context, resource *myapiv1.MyResource) {
    r.auditLogger.Info("resource-mutation",
        "timestamp", time.Now().UTC(),
        "resource", fmt.Sprintf("%s/%s", resource.Namespace, resource.Name),
        "user", r.getUserFromContext(ctx),
        "mutations", "security-defaults-applied",
        "component", "mutating-webhook",
    )
}
```

## ┌─ Pattern Detection Queries [R241-R245] ─┐
### │ Security Anti-Pattern Detection         │
└───────────────────────────────────────────────┘

### Detect Insecure Pod Configurations
```bash
# Query: Find pods running as root
grep -r "runAsUser.*0\|runAsNonRoot.*false" --include="*.yaml" --include="*.go" | head -10

# Query: Find privileged containers
grep -r "privileged.*true\|allowPrivilegeEscalation.*true" --include="*.yaml" --include="*.go" | head -10
```

### Detect Missing Security Contexts
```bash
# Query: Find deployments without security contexts
grep -r -L "securityContext" --include="*deployment*.yaml" --include="*deployment*.go" | head -10
```

### Detect Hardcoded Secrets
```bash
# Query: Find potential hardcoded secrets
grep -r -i "password\|secret\|token\|key.*=" --include="*.go" | grep -v "_test.go" | head -10
```

## ┌─ Validation Criteria [R246-R250] ─┐
### │ Security Implementation Quality │
└─────────────────────────────────────────┘

**Validation R246**: All workloads MUST have security contexts defined
**Validation R247**: No containers MUST run as root or privileged
**Validation R248**: All network traffic MUST be secured with TLS
**Validation R249**: RBAC permissions MUST follow least privilege
**Validation R250**: Security policies MUST be enforced via admission control

### Security Implementation Checklist

- [ ] All pods run with non-root security context
- [ ] No privileged containers or privilege escalation
- [ ] Resource limits and requests are defined
- [ ] Network policies restrict traffic appropriately
- [ ] TLS is used for all API communications
- [ ] Secrets are properly mounted and not logged
- [ ] RBAC follows principle of least privilege
- [ ] Admission webhooks enforce security policies
- [ ] Container images are scanned for vulnerabilities
- [ ] Security events are logged for audit trails

This comprehensive security expertise module provides patterns and requirements for implementing robust security controls across all aspects of Kubernetes applications, from RBAC and secret management to network security and runtime protection.