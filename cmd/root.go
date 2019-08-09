/*

Copyright © 2019 Shohi Wang <oshohi@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/shohi/cube/pkg/config"
	"github.com/shohi/cube/pkg/kube"
)

var conf = config.Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cube",
	Short: "kubectl config manipulation tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return kube.Dispatch(conf)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	setupFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Printf("run kube error, err: %v\n", err)
		os.Exit(1)
	}
}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command) {
	flagSet := cmd.Flags()

	flagSet.StringVar(&conf.RemoteUser, "remote-user", "core", "remote user, e.g. root.")
	flagSet.StringVar(&conf.RemoteIP, "remote-ip", "", "remote master private ip, e.g. 172.17.31.1.")

	flagSet.IntVar(&conf.LocalPort, "local-port", 0, "local forwarding port, e.g. 7001.")
	flagSet.StringVar(&conf.SSHVia, "ssh-via", "", "ssh jump server, e.g. user@jump. If not set, SSH_VIA env will be used. ")
	flagSet.StringVar(&conf.NameSuffix, "name-suffix", "", "cluster name suffix, e.g. dev.")

	flagSet.BoolVar(&conf.DryRun, "dry-run", false, "dry-run mode. validate config and then exit.")
	flagSet.BoolVar(&conf.Purge, "purge", false, "remove configuration.")
	flagSet.BoolVar(&conf.Force, "force", false, "merge configuration forcedly. Only take effect when cluster name is unique")
	flagSet.BoolVar(&conf.PrintSSHForwarding, "print-ssh-forwarding", false, "print ssh forwarding command and exit.")

	cmd.MarkFlagRequired("remote-ip")
}
