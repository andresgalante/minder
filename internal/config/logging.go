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

package config

// LoggingConfig is the configuration for the logging package
type LoggingConfig struct {
	Level   string `mapstructure:"level" default:"debug"`
	Format  string `mapstructure:"format" default:"json"`
	LogFile string `mapstructure:"logFile" default:""`
}
