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

	"github.com/shohi/cube/cmd/add"
	"github.com/shohi/cube/cmd/del"
	"github.com/shohi/cube/cmd/forward"
	"github.com/shohi/cube/cmd/history"
	"github.com/shohi/cube/cmd/list"
	"github.com/shohi/cube/cmd/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cube",
	Short: "kubectl config manipulation tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	rootCmd.AddCommand(history.New())
	rootCmd.AddCommand(list.New())
	rootCmd.AddCommand(version.New())
	rootCmd.AddCommand(add.New())
	rootCmd.AddCommand(del.New())
	rootCmd.AddCommand(forward.New())

	if err := rootCmd.Execute(); err != nil {
		log.Printf("run kube error, err: %v\n", err)
		os.Exit(1)
	}
}
