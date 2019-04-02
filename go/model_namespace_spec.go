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

// The definitions for a namespace
type NamespaceSpec struct {

	// A list of service accounts for this namespace
	ServiceAccounts []map[string]string `json:"service_accounts,omitempty"`
}
