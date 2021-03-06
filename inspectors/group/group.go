package group

import (
	"fmt"
	"github.com/proofpoint/kapprover/inspectors"
	certificates "k8s.io/api/certificates/v1beta1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	inspectors.Register("group", &group{"system:kubelet-bootstrap"})
}

// Group is an Inspector that verifies the CSR was submitted
// by a user in the configured group.
type group struct {
	requiredGroup string
}

func (g *group) Configure(config string) (inspectors.Inspector, error) {
	if config != "" {
		return &group{requiredGroup: config}, nil
	}
	return g, nil
}

func (g *group) Inspect(client kubernetes.Interface, request *certificates.CertificateSigningRequest) (string, error) {
	isRequiredGroup := false
	for _, group := range request.Spec.Groups {
		if group == g.requiredGroup {
			isRequiredGroup = true
			break
		}
	}
	if !isRequiredGroup {
		return fmt.Sprintf("Requesting user is not in the %s group", g.requiredGroup), nil
	}

	return "", nil
}
