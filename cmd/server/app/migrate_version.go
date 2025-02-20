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

// Package app provides the entrypoint for the minder migrations
package app

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"       // nolint
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stacklok/minder/internal/config"
	"github.com/stacklok/minder/internal/logger"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get the db version",
	Long:  `Command to get the database version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.ReadConfigFromViper(viper.GetViper())
		if err != nil {
			return fmt.Errorf("unable to read config: %w", err)
		}

		ctx := logger.FromFlags(cfg.LoggingConfig).WithContext(context.Background())

		// Database configuration
		dbConn, connString, err := cfg.Database.GetDBConnection(ctx)
		if err != nil {
			return fmt.Errorf("unable to connect to database: %w", err)
		}
		defer dbConn.Close()

		configPath := getMigrateConfigPath()
		m, err := migrate.New(configPath, connString)
		if err != nil {
			fmt.Printf("Error while creating migration instance (%s): %v\n", configPath, err)
			os.Exit(1)
		}

		version, dirty, err := m.Version()
		if err != nil {
			fmt.Printf("Error while getting migration version: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Version=%v dirty=%v\n", version, dirty)
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(versionCmd)
}
