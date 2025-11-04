package v1alpha1

// IngressDef defines ingress spec and annotations
type IngressDef struct {
	ClassName   *string           `json:"className,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Host        string            `json:"host"`
}
