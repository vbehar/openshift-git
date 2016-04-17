package openshift

import (
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"

	"k8s.io/kubernetes/pkg/client/restclient"
	kclientcmd "k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	kclientcmdapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

var (
	Factory = newFactory(Flags)
)

// newFactory returns an OpenShift's Factory
// It first tries to use the config that is made available when we are running in a cluster
// and then fallback to a standard factory (using the default config files)
func newFactory(flags *pflag.FlagSet) *clientcmd.Factory {
	factory, err := getFactoryFromCluster()
	if err != nil {
		glog.V(1).Infof("Seems like we are not running in an OpenShift environment (%s), falling back to building a std factory...", err)
		factory = clientcmd.New(flags)
	}

	return factory
}

// getFactoryFromCluster returns an OpenShift's Factory
// using the config that is made available when we are running in a cluster
// (using environment variables and token secret file)
// or an error if those are not available (meaning we are not running in a cluster)
func getFactoryFromCluster() (*clientcmd.Factory, error) {
	clusterConfig, err := restclient.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// keep only what we need to initialize a factory
	overrides := &kclientcmd.ConfigOverrides{
		ClusterInfo: kclientcmdapi.Cluster{
			Server: clusterConfig.Host,
		},
		AuthInfo: kclientcmdapi.AuthInfo{
			Token: clusterConfig.BearerToken,
		},
		Context: kclientcmdapi.Context{},
	}

	if len(clusterConfig.TLSClientConfig.CAFile) > 0 {
		// FIXME "x509: cannot validate certificate for x.x.x.x because it doesn't contain any IP SANs"
		// overrides.ClusterInfo.CertificateAuthority = clusterConfig.TLSClientConfig.CAFile
		overrides.ClusterInfo.InsecureSkipTLSVerify = true
	} else {
		overrides.ClusterInfo.InsecureSkipTLSVerify = true
	}

	config := kclientcmd.NewDefaultClientConfig(*kclientcmdapi.NewConfig(), overrides)

	factory := clientcmd.NewFactory(config)
	return factory, nil
}
