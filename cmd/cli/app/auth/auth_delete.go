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

package auth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stacklok/minder/internal/auth"
	"github.com/stacklok/minder/internal/util"
	"github.com/stacklok/minder/internal/util/cli"
	pb "github.com/stacklok/minder/pkg/api/protobuf/go/minder/v1"
)

// auth_deleteCmd represents the account deletion command
var auth_deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Permanently delete account",
	Long:  `Permanently delete account. All associated user data will be permanently removed.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			fmt.Fprintf(os.Stderr, "Error binding flags: %s\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := util.GetAppContext()
		defer cancel()

		conn, err := util.GrpcForCommand(cmd, viper.GetViper())
		util.ExitNicelyOnError(err, "Error getting grpc connection")
		defer conn.Close()
		client := pb.NewUserServiceClient(conn)

		// Ensure the user already exists in the local database
		_, _, err = userRegistered(ctx, client)
		util.ExitNicelyOnError(err, "Error fetching user")

		// Get user details - name, email from the jwt token
		userDetails, err := auth.GetUserDetails(ctx, cmd, viper.GetViper())
		util.ExitNicelyOnError(err, "Error fetching user details")

		// Confirm user wants to delete their account
		yes := cli.PrintYesNoPrompt(cmd,
			fmt.Sprintf(
				"You are about to permanently delete your account. \n\nName: %s\nEmail: %s",
				userDetails.Name,
				userDetails.Email,
			),
			"Are you sure?",
			"Delete account operation cancelled.",
			false)
		if !yes {
			return
		}

		_, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{})
		util.ExitNicelyOnError(err, "Error registering user")

		// This step is added to avoid confusing the users by seeing their credentials locally, however it is not
		// directly related to user deletion because the token will expire after 5 minutes and cannot be refreshed
		err = util.RemoveCredentials()
		if err != nil {
			cli.PrintCmd(cmd, cli.WarningBanner.Render("Failed to remove locally stored credentials."))
		}
		cli.PrintCmd(cmd, cli.SuccessBanner.Render("Successfully deleted account. It may take up to 48 hours for "+
			"all data to be removed."))
	},
}

func init() {
	AuthCmd.AddCommand(auth_deleteCmd)
}
