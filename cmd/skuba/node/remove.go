/*
 * Copyright (c) 2019 SUSE LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package node

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/SUSE/skuba/internal/pkg/skuba/kubernetes"
	node "github.com/SUSE/skuba/pkg/skuba/actions/node/remove"
)

type removeOptions struct {
	drainTimeout time.Duration
}

// NewRemoveCmd creates a new `skuba node remove` cobra command
func NewRemoveCmd() *cobra.Command {
	removeOptions := removeOptions{}
	cmd := &cobra.Command{
		Use:   "remove <node-name>",
		Short: "Removes a node from the cluster",
		Run: func(cmd *cobra.Command, nodenames []string) {
			client, err := kubernetes.GetAdminClientSet()
			if err != nil {
				fmt.Printf("Unable to remove a node: %s\n", err)
				os.Exit(1)
			}

			if err := node.Remove(client, nodenames[0], removeOptions.drainTimeout); err != nil {
				fmt.Printf("Unable to remove a node %s: %s\n", nodenames[0], err)
				os.Exit(1)
			}
		},
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().DurationVar(&removeOptions.drainTimeout, "drain-timeout", 0, `Time to wait for the node to drain, before proceeding with node removal.
The time can be specified using abbreviations for units: e.g. 1h15m15s (Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h").
Will wait indefinitely by default.`)

	return cmd
}
