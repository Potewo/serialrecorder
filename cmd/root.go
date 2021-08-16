/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"github.com/Potewo/serialrecorder/file"
	"github.com/Potewo/serialrecorder/serial"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Baudrate   int
	DeviceName string
	Output   string
}

var cfgFile string
var config Config
var baudrate int
var deviceName string
var fileName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "serialrecorder",
	Short: "Show and record bytes through serial port.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("BaudRate: %#v\n", config.Baudrate)
		fmt.Printf("DeviceName: %#v\n", config.DeviceName)
		fmt.Printf("FileName %#v\n", config.Output)
		if config.DeviceName == "" {
			fmt.Fprintln(os.Stderr, "Device name is required")
			os.Exit(1)
		}
		if config.Output != "" {
			if err := file.Open(config.Output); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		if err := serial.Init(config.DeviceName, config.Baudrate, config.Output); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to initialize serial port.\n", err)
			os.Exit(1)
		}
		if config.Output != "" {
			for {
				s, err := serial.Read()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to read serial data:", err)
					os.Exit(1)
				}
				fmt.Printf("%s", s)
				if err := file.Append(string(s)); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		} else {
			for {
				s, err := serial.Read()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to read serial data:", err)
					os.Exit(1)
				}
				fmt.Printf("%s", s)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$HOME/.serialrecorder.yml", "config file (default is $HOME/.serialrecorder.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().IntVarP(&baudrate, "baudrate", "b", 9600, "set serial baudrate (required)")
	rootCmd.Flags().StringVarP(&deviceName, "devicename", "d", "", "set device name (required)")
	rootCmd.Flags().StringVarP(&fileName, "output", "o", "", "set filepath to save string")
	viper.BindPFlag("Baudrate", rootCmd.Flags().Lookup("baudrate"))
	viper.BindPFlag("DeviceName", rootCmd.Flags().Lookup("devicename"))
	viper.BindPFlag("Output", rootCmd.Flags().Lookup("output"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".serialrecorder" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".serialrecorder")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "Using config file:", err)
	}
}
