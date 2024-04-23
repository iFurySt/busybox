/**
 * Package exit
 * @Author iFurySt <ifuryst@gmail.com>
 * @Date 2024/4/22
 */

package exit

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	isExiting atomic.Bool
)

func NewCmdExit() *cobra.Command {
	var sigtermDuration, sighupDuration, sigintDuration, sigquitDuration string

	cmd := &cobra.Command{
		Use:   "exit",
		Short: "Simulates the exit behavior of a container in a Kubernetes Pod.\n\nYou can specify the duration in ns, us, µs, ms, s, m, h. For example, 1s, 1m, 1h.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, Busybox Exit here! I'll exit until receive a signal...")

			durations := map[os.Signal]time.Duration{
				syscall.SIGTERM: mustParseDuration(sigtermDuration),
				syscall.SIGHUP:  mustParseDuration(sighupDuration),
				syscall.SIGINT:  mustParseDuration(sigintDuration),
				syscall.SIGQUIT: mustParseDuration(sigquitDuration),
			}

			handleSignal(durations)

			forever := make(chan struct{})
			<-forever
		},
	}

	cmd.Flags().StringVar(&sigtermDuration, "sigterm-duration", "1s", "Duration to wait before SIGTERM (15)")
	cmd.Flags().StringVar(&sighupDuration, "sighup-duration", "1s", "Duration to wait before SIGHUP (1)")
	cmd.Flags().StringVar(&sigintDuration, "sigint-duration", "1s", "Duration to wait before SIGINT (2)")
	cmd.Flags().StringVar(&sigquitDuration, "sigquit-duration", "1s", "Duration to wait before SIGQUIT (2)")

	return cmd
}

func handleSignal(durations map[os.Signal]time.Duration) {
	// 创建一个接收信号的通道
	sigs := make(chan os.Signal, 1)

	// 注册要接收的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// 使用goroutine来异步处理信号
	go func() {
		// 在这个循环中，我们会一直等待信号的到来
		for {
			sig := <-sigs
			duration := durations[sig]
			fmt.Println("-> Signal received:", sig)
			tryExit(duration)
		}
	}()
}

func tryExit(duration time.Duration) {
	go func() {
		if !isExiting.CompareAndSwap(false, true) {
			return
		}
		fmt.Println("Waiting for", duration.String(), "before exit")
		countDown(int(duration.Seconds()))
		fmt.Println("Exiting now...")
		os.Exit(0)
	}()
}
