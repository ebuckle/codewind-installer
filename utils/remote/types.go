package remote

import "k8s.io/apimachinery/pkg/types"

// Codewind represents a Codewind instance: name, namespace, volume, serviceaccount, and pull secrets
type Codewind struct {
	PFEName            string
	PerformanceName    string
	GatekeeperName     string
	KeycloakName       string
	PFEImage           string
	PerformanceImage   string
	GatekeeperImage    string
	KeycloakImage      string
	Namespace          string
	WorkspaceID        string
	PVCName            string
	ServiceAccountName string
	PullSecret         string
	OwnerReferenceName string
	OwnerReferenceUID  types.UID
	Privileged         bool
	Ingress            string
	OnOpenShift        bool
}

// ServiceAccountPatch contains an array of imagePullSecrets that will be patched into a Kubernetes service account
type ServiceAccountPatch struct {
	ImagePullSecrets *[]ImagePullSecret `json:"imagePullSecrets"`
}

// ImagePullSecret represents a Kubernetes imagePullSecret
type ImagePullSecret struct {
	Name string `json:"name"`
}
