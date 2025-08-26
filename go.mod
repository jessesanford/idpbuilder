module github.com/idpbuilder/idpbuilder-oci-mgmt/split-003

go 1.21

require (
	github.com/containers/buildah v1.30.0
	github.com/containers/storage v1.48.0
	github.com/containers/image/v5 v5.26.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
)

// Note: In a real implementation, we would have dependencies on other splits:
// github.com/idpbuilder/idpbuilder-oci-mgmt/split-001 v0.0.0
// github.com/idpbuilder/idpbuilder-oci-mgmt/split-002 v0.0.0
// 
// But since they're in different branches, we'll define compatible interfaces locally