/**
 * Package busybox
 * @Author iFurySt <ifuryst@gmail.com>
 * @Date 2024/4/22
 */

package main

import (
	"fmt"
	"github.com/ifuryst/busybox/exit"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "busybox",
		Short: "Busybox is a simple task runner for Kubernetes and Docker test purposes",
		Long: `Busybox is a simple task runner for Kubernetes and Docker test purposes. 
Like to use the busybox simulates the exit behavior of a container in a Kubernetes Pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, Busybox here!")
		},
	}

	rootCmd.AddCommand(exit.NewCmdExit())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
