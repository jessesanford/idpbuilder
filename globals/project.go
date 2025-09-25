package globals

import (
	"fmt"
	"os"
)

// ProjectName can be overridden via PROJECT_NAME environment variable
var ProjectName = getProjectName()

const (
	NginxNamespace  string = "ingress-nginx"
	ArgoCDNamespace string = "argocd"

	SelfSignedCertSecretName = "idpbuilder-cert"
	SelfSignedCertCMName     = "idpbuilder-cert"
	SelfSignedCertCMKeyName  = "ca.crt"
	DefaultSANWildcard       = "*.cnoe.localtest.me"
	DefaultHostName          = "cnoe.localtest.me"
)

func getProjectName() string {
	if name := os.Getenv("PROJECT_NAME"); name != "" {
		return name
	}
	return "idpbuilder" // Default only
}

func GetProjectNamespace(name string) string {
	return fmt.Sprintf("%s-%s", ProjectName, name)
}
