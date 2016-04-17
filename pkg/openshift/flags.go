package openshift

import (
	"github.com/spf13/pflag"
)

var (
	Flags = pflag.NewFlagSet("openshift-git", pflag.ContinueOnError)
)
