/**
 * Package busybox
 * @Author iFurySt <ifuryst@gmail.com>
 * @Date 2024/4/22
 */

package main

import (
	"fmt"
	"github.com/ifuryst/busybox/exit"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests",
	},
	[]string{"path"},
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
}

func Execute() {
	var addr string
	var rootCmd = &cobra.Command{
		Use:   "busybox",
		Short: "Busybox is a simple task runner for Kubernetes and Docker test purposes",
		Long: `Busybox is a simple task runner for Kubernetes and Docker test purposes. 
Like to use the busybox simulates the exit behavior of a container in a Kubernetes Pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, Busybox here!")

			fmt.Println("Busybox is running on", addr)

			gin.SetMode(gin.ReleaseMode)
			r := gin.Default()
			r.Use(prometheusMiddleware())
			r.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "pong"})
			})
			r.GET("/metrics", func(c *gin.Context) {
				promhttp.Handler().ServeHTTP(c.Writer, c.Request)
			})

			_ = r.Run(addr)
		},
	}

	rootCmd.AddCommand(exit.NewCmdExit())

	rootCmd.Flags().StringVar(&addr, "addr", "127.0.0.1:8888", "The address to listen on for HTTP requests.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call the next handler
		c.Next()

		// After the request has been handled, increment the counter
		httpRequestsTotal.With(prometheus.Labels{"path": c.Request.URL.Path}).Inc()
	}
}
