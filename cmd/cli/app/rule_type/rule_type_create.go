//
// Copyright 2023 Stacklok, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rule_type

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stacklok/minder/internal/util"
	minderv1 "github.com/stacklok/minder/pkg/api/protobuf/go/minder/v1"
)

// RuleType_createCmd represents the profile create command
var RuleType_createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a rule type within a minder control plane",
	Long: `The minder rule type create subcommand lets you create new rule types for a project
within a minder control plane.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			fmt.Fprintf(os.Stderr, "Error binding flags: %s\n", err)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := cmd.Flags().GetStringArray("file")
		if err != nil {
			return fmt.Errorf("error getting file flag: %w", err)
		}

		if err := validateFilesArg(files); err != nil {
			return fmt.Errorf("error validating file arg: %w", err)
		}

		conn, err := util.GrpcForCommand(cmd, viper.GetViper())
		if err != nil {
			return fmt.Errorf("error getting grpc connection: %w", err)
		}
		defer conn.Close()

		client := minderv1.NewProfileServiceClient(conn)

		expfiles, err := util.ExpandFileArgs(files)
		if err != nil {
			return fmt.Errorf("error expanding file args: %w", err)
		}

		table := initializeTable(cmd)

		ctx, cancel := util.GetAppContext()
		defer cancel()

		createFunc := func(fileName string, rt *minderv1.RuleType) (*minderv1.RuleType, error) {
			resprt, err := client.CreateRuleType(ctx, &minderv1.CreateRuleTypeRequest{
				RuleType: rt,
			})
			if err != nil {
				return nil, fmt.Errorf("error creating rule type from %s: %w", fileName, err)
			}

			return resprt.RuleType, nil
		}

		for _, f := range expfiles {
			if shouldSkipFile(f) {
				continue
			}

			if err := execOnOneRuleType(table, f, os.Stdin, createFunc); err != nil {
				return fmt.Errorf("error creating rule type %s: %w", f, err)
			}
		}

		table.Render()

		return nil
	},
}

func init() {
	ruleTypeCmd.AddCommand(RuleType_createCmd)
	RuleType_createCmd.Flags().StringArrayP("file", "f", []string{},
		"Path to the YAML defining the rule type (or - for stdin). Can be specified multiple times. Can be a directory.")
}
