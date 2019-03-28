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

// The resource definition for a namespace in the cluster
type Namespace struct {

	// A globally unique human readible resource name
	Name string `json:"name"`

	// A spec for the namespace
	NamespaceSpec *NamespaceSpec `json:"spec"`

	// A cryptographic signature used to verify the request payload
	Signature string `json:"signature,omitempty"`
}
