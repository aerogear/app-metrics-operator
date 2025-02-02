package utils

import (
	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"os"
	"strings"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

)

const (
	// AppNamespaceEnvVar is the constant for env variable AppNamespaceEnvVar
	// which is the namespace where the APP CR can applied.
	// The namespaces should be informed split by ";".
	AppNamespaceEnvVar = "APP_NAMESPACES"
	OperatorNamespaceForLocalEnv   = "app-metrics"

)

var log = logf.Log.WithName("app-metrics-operator.utils")

// GetAppNamespaces returns the namespace the operator should be watching for changes
func GetAppNamespaces() (string, error) {
	ns, found := os.LookupEnv(AppNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", AppNamespaceEnvVar)
	}
	return ns, nil
}

// IsValidAppNamespace return true when the namespace informed is declared in the ENV VAR APP_NAMESPACES
func IsValidAppNamespace(namespace string) (bool, error) {
	appNamespacesEnvVar, err := GetAppNamespaces()
	if err != nil {
		log.Error(err, "Unable to check if is app namespace %s is valid", namespace)

		// To skip when it is local env or a unit test
		_, err := k8sutil.GetOperatorNamespace()
		if err != nil {
			//Return true for the local env and for the unit tests
			if err == k8sutil.ErrNoNamespace { //
				log.Info("Allow to continue because it is outside of the cluster")
				return true, nil
			}
			return false, err
		}

		return false, err
	}
	for _, ns := range strings.Split(appNamespacesEnvVar, ";") {
		if ns == namespace {
			return true, nil
		}
	}
	err = fmt.Errorf("Invalid Namespace")
	return false, err
}

// IsValidOperatorNamespace return true when the namespace informed is declared in the ENV VAR APP_NAMESPACES
func IsValidOperatorNamespace(namespace string) (bool, error) {
	ns, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		//Return true for the local env and for the unit tests
		if err == k8sutil.ErrNoNamespace {
			log.Info("Allow to continue because it is outside of the cluster")
			return true, nil
		}
		return false, err
	}
	if ns == namespace {
		return true, nil
	}
	err = fmt.Errorf("Invalid Namespace")
	return false, err
}