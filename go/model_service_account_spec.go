/*
 * hub-kubernetes-agent
 *
 * an agent used to provision and configure Kubernetes resources
 *
 * API version: v1beta
 * Contact: support@appvia.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

// The resource specification of a service account
type ServiceAccountSpec struct {
	// The name of the service account
	Name string `json:"name"`
	// The token associated with the service account
	Token string `json:"token,omitempty"`
	// The namespace this service account is in
	Namespace string `json:"namespace,omitempty"`
}
