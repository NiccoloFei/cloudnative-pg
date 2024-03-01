/*
Copyright The CloudNativePG Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package status

import (
	"github.com/spf13/cobra"

	"github.com/cloudnative-pg/cloudnative-pg/internal/cmd/plugin"
	"github.com/cloudnative-pg/cloudnative-pg/internal/cmd/plugin/utils"
)

// NewCmd create the new "status" subcommand
func NewCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status [cluster]",
		Short: "Get the status of a PostgreSQL cluster",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return utils.CompleteClusters(cmd.Context(), plugin.Client, args, toComplete), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			clusterName := args[0]

			verbose, _ := cmd.Flags().GetBool("verbose")
			output, _ := cmd.Flags().GetString("output")

			return Status(ctx, clusterName, verbose, plugin.OutputFormat(output))
		},
	}

	statusCmd.Flags().BoolP(
		"verbose", "v", false, "Include PostgreSQL configuration, HBA rules, and full replication slots info")
	statusCmd.Flags().StringP(
		"output", "o", "text", "Output format. One of text|json")

	return statusCmd
}
