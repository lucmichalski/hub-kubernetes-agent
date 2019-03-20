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

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clientset *kubernetes.Clientset

func getClient(server, token, caCert string) (*kubernetes.Clientset, error) {
	decodedCert, err := base64.StdEncoding.DecodeString(caCert)

	if err != nil {
		fmt.Println("decode error:", err)
		return nil, err
	}

	config := &rest.Config{
		Host:            server,
		BearerToken:     token,
		TLSClientConfig: rest.TLSClientConfig{CAData: decodedCert},
	}

	client, err := kubernetes.NewForConfig(config)
	_ = client

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func NamespacesList(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})

	json := simplejson.New()

	if errors.IsNotFound(err) {
		log.Printf("Listing namespaces\n")
		json.Set("status", "not found")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing namespaces: %v\n",
			statusError.ErrStatus.Message)
		json.Set("status", "error")
	} else if err != nil {
		json.Set("status", "error")
		panic(err.Error())
	} else {
		var namespaceList []string

		for _, namespace := range namespaces.Items {
			namespaceList = append(namespaceList, namespace.Name)
		}
		json.Set("namespaces", namespaceList)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}

func NamespacesNameGet(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	vars := mux.Vars(r)
	name := vars["name"]
	namespace, err := clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	json := simplejson.New()

	if err != nil {
		log.Println(err)
		json.Set("status", "error")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting namespace %s: %v\n",
			name, statusError.ErrStatus.Message)
		json.Set("status", "error")
	} else {
		json.Set("name", name)
		log.Printf("Found namespace: %s\n", name)
	}
	_ = namespace

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	payload, err := json.MarshalJSON()

	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}

func NamespacesNameDelete(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	vars := mux.Vars(r)
	name := vars["name"]

	log.Printf("Deleting namespace: %v", name)

	json := simplejson.New()

	deletePolicy := metav1.DeletePropagationForeground

	if err := clientset.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); errors.IsNotFound(err) {
		log.Printf("Namespace %s not found\n", name)
		json.Set("status", "not found")
	} else if err != nil {
		json.Set("status", "error")
		panic(err.Error())
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting namespace %s: %v\n",
			name, statusError.ErrStatus.Message)
		json.Set("status", "error")
	} else {
		json.Set("status", "deleted")
		log.Printf("Deleted namespace %s\n", name)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}

func NamespacesNamePut(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	vars := mux.Vars(r)
	name := vars["name"]
	namespace, err := clientset.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
	json := simplejson.New()

	if err != nil {
		log.Println(err)
		json.Set("status", "error")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error creating namespace %s: %v\n",
			name, statusError.ErrStatus.Message)
		json.Set("status", "error")
	} else {
		json.Set("name", name)
		log.Printf("Created namespace: %s\n", name)
	}
	_ = namespace

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	payload, err := json.MarshalJSON()

	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}

func ServiceAccountsNamespaceGet(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	vars := mux.Vars(r)
	namespace := vars["namespace"]

	if namespace == "" {
		namespace = "default"
	}

	serviceAccounts, err := clientset.CoreV1().ServiceAccounts(namespace).List(metav1.ListOptions{})
	_ = serviceAccounts

	json := simplejson.New()

	if errors.IsNotFound(err) {
		log.Printf("service accounts \n")
		json.Set("status", "not found")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing service accounts for namespace %s: %v\n",
			namespace, statusError.ErrStatus.Message)
		json.Set("status", "error")
	} else if err != nil {
		json.Set("status", "error")
		panic(err.Error())
	} else {
		var serviceAccountList []string
		for _, sa := range serviceAccounts.Items {
			serviceAccountList = append(serviceAccountList, sa.Name)
		}
		json.Set("service-accounts", serviceAccountList)
		log.Printf("Listed service accounts for namespace %s\n", namespace)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}

func ServiceAccountsNamespaceNameDelete(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	vars := mux.Vars(r)
	name := vars["name"]
	namespace := vars["namespace"]

	if namespace == "" {
		namespace = "default"
	}
	json := simplejson.New()

	deletePolicy := metav1.DeletePropagationForeground
	if err := clientset.CoreV1().ServiceAccounts(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); errors.IsNotFound(err) {
		log.Printf("service account %s not found\n", name)
		json.Set("status", "not found")
	} else if err != nil {
		json.Set("status", "error")
		panic(err.Error())
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error deleting service account %s: %v\n",
			name, statusError.ErrStatus.Message)
		json.Set("status", "error")
	} else {
		json.Set("status", "deleted")
		log.Printf("Deleted service account %s\n", name)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}

func ServiceAccountsNamespaceNameGet(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	vars := mux.Vars(r)
	name := vars["name"]
	namespace := vars["namespace"]

	if namespace == "" {
		namespace = "default"
	}

	serviceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
	_ = serviceAccount

	json := simplejson.New()

	if errors.IsNotFound(err) {
		log.Printf("Service account %s not found\n", name)
		json.Set("status", "not found")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting service account %s: %v\n",
			name, statusError.ErrStatus.Message)
		json.Set("status", "error")
	} else if err != nil {
		json.Set("status", "error")
		panic(err.Error())
	} else {
		json.Set("name", name)
		log.Printf("Found service account %s\n", name)
		secret, err := clientset.CoreV1().Secrets(namespace).Get(serviceAccount.Secrets[0].Name, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		json.Set("token", string(secret.Data["token"]))
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}

func ServiceAccountsNamespaceNamePut(w http.ResponseWriter, r *http.Request) {
	clientset, err := getClient(r.Header.Get("X-Kube-API-URL"), r.Header.Get("X-Kube-Token"), r.Header.Get("X-Kube-CA"))

	vars := mux.Vars(r)
	name := vars["name"]
	namespace := vars["namespace"]

	serviceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Create(&corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
	if err != nil {
		log.Println(err)
	}
	_ = serviceAccount
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	log.Printf("Creating service account: %v", name)
	json := simplejson.New()
	json.Set("name", name)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	w.Write(payload)
}
