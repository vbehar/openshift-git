package openshift

import (
	// install all required openshift resources
	_ "github.com/openshift/origin/pkg/api/install"

	"github.com/openshift/origin/pkg/cmd/util/clientcmd"

	"github.com/spf13/pflag"
)

// init the factory early, to bind the flags
var (
	Flags   = pflag.NewFlagSet("openshift-git", pflag.ContinueOnError)
	Factory = clientcmd.New(Flags)
)
