// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apps

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/google/kf/pkg/kf/apps"
	"github.com/google/kf/pkg/kf/commands/config"
	"github.com/google/kf/pkg/kf/commands/utils"
	"github.com/spf13/cobra"
)

// NewAppsCommand creates a apps command.
func NewAppsCommand(p *config.KfParams, appsClient apps.Client) *cobra.Command {
	var apps = &cobra.Command{
		Use:     "apps",
		Short:   "List pushed apps",
		Example: `  kf apps`,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.ValidateNamespace(p); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			fmt.Fprintf(cmd.OutOrStdout(), "Getting apps in space %s\n", p.Namespace)

			apps, err := appsClient.List(p.Namespace)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout())

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 8, 4, 1, ' ', tabwriter.StripEscape)
			fmt.Fprintln(w, "name\t\trequested state\t\tinstances\tmemory\tdisk\turls")
			for _, app := range apps {

				// Requested State
				requestedState := "started"
				if app.Spec.Instances.Stopped {
					requestedState = "stopped"
				} else if cond := app.Status.GetCondition("Ready"); cond != nil && cond.Status == "Pending" {
					requestedState = "starting"
				} else if !app.DeletionTimestamp.IsZero() {
					requestedState = "deleting"
				}

				// Instances
				var instances string
				switch {
				case app.Spec.Instances.Exactly != nil:
					instances = strconv.FormatInt(int64(*app.Spec.Instances.Exactly), 10)
				case app.Spec.Instances.Min == nil && app.Spec.Instances.Max == nil:
					instances = "?"
				case app.Spec.Instances.Min != nil && app.Spec.Instances.Max != nil:
					instances = fmt.Sprintf(
						"%d - %d",
						*app.Spec.Instances.Min,
						*app.Spec.Instances.Max,
					)
				case app.Spec.Instances.Max != nil:
					instances = fmt.Sprintf(
						"0 - %d",
						*app.Spec.Instances.Max,
					)
				case app.Spec.Instances.Min != nil:
					instances = fmt.Sprintf(
						"%d - ∞",
						*app.Spec.Instances.Min,
					)
				}

				// Memory
				var memory string
				if containers := app.Spec.Template.Spec.Containers; len(containers) > 0 {
					if mem, ok := containers[0].Resources.Requests["memory"]; ok {
						memory = mem.String()
					}
				}

				// Disk
				// TODO(#431): Persistent disks
				var disk string
				if containers := app.Spec.Template.Spec.Containers; len(containers) > 0 {
					if d, ok := containers[0].Resources.Requests["ephemeral-storage"]; ok {
						disk = d.String()
					}
				}

				// URL
				var urls []string
				for _, route := range app.Spec.Routes {
					var hostnamePrefix string
					if route.Hostname != "" {
						hostnamePrefix = route.Hostname + "."
					}
					urls = append(
						urls,
						fmt.Sprintf(
							"%s%s%s",
							hostnamePrefix,
							route.Domain,
							path.Join("/", route.Path),
						),
					)
				}

				if app.Name == "" {
					continue
				}

				fmt.Fprintf(w, "%s\t\t%s\t\t%s\t%s\t%s\t%s\n",
					app.Name,
					requestedState,
					instances,
					memory,
					disk,
					strings.Join(urls, ", "),
				)
			}

			w.Flush()

			return nil
		},
	}

	return apps
}
