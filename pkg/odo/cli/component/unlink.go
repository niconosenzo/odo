package component

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"

	appCmd "github.com/openshift/odo/pkg/odo/cli/application"
	projectCmd "github.com/openshift/odo/pkg/odo/cli/project"
	"github.com/openshift/odo/pkg/odo/util/completion"
	svc "github.com/openshift/odo/pkg/service"

	"github.com/openshift/odo/pkg/odo/util"
	ktemplates "k8s.io/kubectl/pkg/util/templates"

	"github.com/spf13/cobra"
)

// UnlinkRecommendedCommandName is the recommended unlink command name
const UnlinkRecommendedCommandName = "unlink"

var (
	unlinkExample = ktemplates.Examples(`# Unlink the 'my-postgresql' service from the current component 
%[1]s my-postgresql

# Unlink the 'my-postgresql' service  from the 'nodejs' component
%[1]s my-postgresql --component nodejs

# Unlink the 'backend' component from the current component (backend must have a single exposed port)
%[1]s backend

# Unlink the 'backend' service  from the 'nodejs' component
%[1]s backend --component nodejs

# Unlink the backend's 8080 port from the current component 
%[1]s backend --port 8080`)

	unlinkLongDesc = `Unlink component or service from a component. 
For this command to be successful, the service or component needs to have been linked prior to the invocation using 'odo link'`
)

// UnlinkOptions encapsulates the options for the odo link command
type UnlinkOptions struct {
	componentContext string
	*commonLinkOptions
}

// NewUnlinkOptions creates a new UnlinkOptions instance
func NewUnlinkOptions() *UnlinkOptions {
	options := UnlinkOptions{}
	options.commonLinkOptions = newCommonLinkOptions()
	options.commonLinkOptions.csvSupport, _ = svc.IsCSVSupported()
	return &options
}

// Complete completes UnlinkOptions after they've been created
func (o *UnlinkOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	err = o.complete(name, cmd, args)
	if err != nil {
		return err
	}

	o.csvSupport, err = o.Client.GetKubeClient().IsCSVSupported()
	if err != nil {
		return err
	}

	if o.csvSupport && o.Context.EnvSpecificInfo != nil {
		o.operation = o.KClient.UnlinkSecret
	} else {
		o.operation = o.Client.UnlinkSecret
	}
	return err
}

// Validate validates the UnlinkOptions based on completed values
func (o *UnlinkOptions) Validate() (err error) {
	return o.validate(false)
}

// Run contains the logic for the odo link command
func (o *UnlinkOptions) Run() (err error) {
	return o.run()
}

// NewCmdUnlink implements the link odo command
func NewCmdUnlink(name, fullName string) *cobra.Command {
	o := NewUnlinkOptions()

	unlinkCmd := &cobra.Command{
		Use:         fmt.Sprintf("%s <service> --component [component] OR %s <component> --component [component]", name, name),
		Short:       "Unlink component to a service or component",
		Long:        unlinkLongDesc,
		Example:     fmt.Sprintf(unlinkExample, fullName),
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"command": "component"},
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	unlinkCmd.PersistentFlags().StringVar(&o.port, "port", "", "Port of the backend to which to unlink")
	unlinkCmd.PersistentFlags().BoolVarP(&o.wait, "wait", "w", false, "If enabled the link will return only when the component is fully running after the link is deleted")

	unlinkCmd.SetUsageTemplate(util.CmdUsageTemplate)
	//Adding `--project` flag
	projectCmd.AddProjectFlag(unlinkCmd)
	//Adding `--application` flag
	appCmd.AddApplicationFlag(unlinkCmd)
	//Adding `--component` flag
	AddComponentFlag(unlinkCmd)
	// Adding context flag
	genericclioptions.AddContextFlag(unlinkCmd, &o.componentContext)

	completion.RegisterCommandHandler(unlinkCmd, completion.UnlinkCompletionHandler)

	return unlinkCmd
}
