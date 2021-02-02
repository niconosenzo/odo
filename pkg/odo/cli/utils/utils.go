package utils

import (
	"fmt"

	odoutil "github.com/openshift/odo/pkg/odo/util"

	"github.com/openshift/odo/pkg/odo/cli/component"
	"github.com/spf13/cobra"
)

const RecommendedCommandName = "utils"

// NewCmdUtils implements the utils odo command
func NewCmdUtils(name, fullName string) *cobra.Command {
	terminalCmd := NewCmdTerminal(terminalCommandName, odoutil.GetFullName(fullName, terminalCommandName))
	convertCmd := NewCmdConvert(convertCommandName, odoutil.GetFullName(fullName, convertCommandName))
	utilsCmd := &cobra.Command{
		Use:   name,
		Short: "Utilities for terminal commands and modifying odo configurations",
		Long:  "Utilities for terminal commands and modifying odo configurations",
		Example: fmt.Sprintf("%s\n",
			terminalCmd.Example),
	}

	utilsCmd.Annotations = map[string]string{"command": "utility"}
	utilsCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)

	utilsCmd.AddCommand(terminalCmd)
	utilsCmd.AddCommand(convertCmd)
	return utilsCmd
}

// GetComponentContext gets the component context
func GetComponentContext(po *component.PushOptions) string {
	return po.ComponentContext
}
