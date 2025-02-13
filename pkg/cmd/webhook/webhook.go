package webhook

import (
	"os"

	admissionserver "github.com/openshift/generic-admission-server/pkg/cmd/server"
	"github.com/spf13/cobra"
	genericapiserver "k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"open-cluster-management.io/work/pkg/webhook"
)

func NewAdmissionHook() *cobra.Command {
	o := admissionserver.NewAdmissionServerOptions(os.Stdout, os.Stderr, &webhook.ManifestWorkAdmissionHook{})

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Start Managed Cluster Admission Server",
		RunE: func(c *cobra.Command, args []string) error {
			stopCh := genericapiserver.SetupSignalHandler()

			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunAdmissionServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	featureGate := utilfeature.DefaultMutableFeatureGate
	featureGate.AddFlag(flags)
	o.RecommendedOptions.FeatureGate = featureGate

	o.RecommendedOptions.AddFlags(cmd.Flags())

	return cmd
}
