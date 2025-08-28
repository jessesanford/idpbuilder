package extractor

import "fmt"

// ClusterNotFoundError indicates the expected Kind cluster was not found
type ClusterNotFoundError struct {
	Expected string
	Actual   string
}

func (e *ClusterNotFoundError) Error() string {
	if e.Actual != "" {
		return fmt.Sprintf("wrong Kind cluster: expected '%s' but connected to '%s'",
			e.Expected, e.Actual)
	}
	return fmt.Sprintf("Kind cluster '%s' not found. Please ensure the cluster is running.",
		e.Expected)
}

// PodNotFoundError indicates the Gitea pod was not found
type PodNotFoundError struct {
	AppName    string
	Namespaces []string
}

func (e *PodNotFoundError) Error() string {
	return fmt.Sprintf("no running %s pod found in namespaces: %v. "+
		"Please ensure %s is installed and running in the Kind cluster.",
		e.AppName, e.Namespaces, e.AppName)
}

// CertificateReadError indicates failure to read certificate from pod
type CertificateReadError struct {
	PodName   string
	Namespace string
	Path      string
	Err       error
}

func (e *CertificateReadError) Error() string {
	return fmt.Sprintf("failed to read certificate from %s/%s at path %s: %v",
		e.Namespace, e.PodName, e.Path, e.Err)
}

func (e *CertificateReadError) Unwrap() error {
	return e.Err
}

// CertificateParseError indicates failure to parse the certificate
type CertificateParseError struct {
	Reason string
}

func (e *CertificateParseError) Error() string {
	return fmt.Sprintf("failed to parse certificate: %s", e.Reason)
}