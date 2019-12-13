package main

import (
	"fmt"
  b64 "encoding/base64"
	// "flag"
	// "path/filepath"
	"log"

  apiv1 "k8s.io/api/core/v1"
  // "k8s.io/apimachinery/pkg/api/errors"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
	// "k8s.io/client-go/util/homedir"
	// "k8s.io/client-go/tools/clientcmd"

	"github.com/gorilla/mux"
	"net/http"
)

func getclientset()(*kubernetes.Clientset, error){

	/*
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	*/

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
        return clientset, err

}


func Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}


func createSecret(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	
	clientset, err := getclientset()
        if err!= nil {
	        panic(err)
	}

	namespace := params["namespace"]
	user := "postgres"
  password := "test"

	userEnc := b64.StdEncoding.EncodeToString([]byte(user))
	fmt.Println(user)
	fmt.Println(userEnc)

	passEnc := b64.StdEncoding.EncodeToString([]byte(password))
	fmt.Println(password)
	fmt.Println(passEnc)

	secret := apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "new-secret",
			Namespace:  namespace,
		},
		Data: map[string][]byte{
		  "ja-postgres-user": []byte(userEnc),
      "ja-postgres-login": []byte(passEnc),
		},
		Type: "Opaque",
	}

	secretOut, err := clientset.CoreV1().Secrets(namespace).Create(&secret)
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
	  fmt.Println("success?", secretOut)
	}
}


func deleteSecret(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	namespace := params["namespace"]
	name :=  params["name"]
	clientset, err := getclientset()
        if err!= nil {
	        panic(err)
	}
	err = clientset.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Secret deleted")
}


func main() {
	router := mux.NewRouter()

	router.HandleFunc("/hello-world", Hello)
	router.HandleFunc("/secrets/{namespace}", createSecret).Methods("POST")
	router.HandleFunc("/secrets/{namespace}/{name}", deleteSecret).Methods("DELETE")

	srv := &http.Server{
	Handler: router,
	Addr: "0.0.0.0:5000",

	}

	fmt.Println("running")

	log.Fatal(srv.ListenAndServe())

}

