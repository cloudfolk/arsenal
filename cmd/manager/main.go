package main

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"go.uber.org/dig"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/cloudfolk/arsenal/cmd"
	"github.com/cloudfolk/arsenal/pkg/app"
)

func init() {
	// Setup flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", appName))

	rootCmd.PersistentFlags().String("api.host", "", "the grpc server listening host")
	rootCmd.PersistentFlags().Int("api.port", 8088, "the grpc server listening port")

	rootCmd.PersistentFlags().String("metrics.host", "", "the metrics server listening host")
	rootCmd.PersistentFlags().Int32("metrics.port", 8383, "the metrics server listening port")
}

const (
	appName = "manager"
)

var cfgFile string

// rootCmd is the root command
var rootCmd = &cobra.Command{
	Use:               appName,
	Short:             appName + " runs arsenal-operator that manages asynchronuous tasks",
	Long:              appName + " runs arsenal-operator that manages asynchronuous tasks",
	PersistentPreRunE: cmd.InitViper,
	RunE: func(cmd *cobra.Command, args []string) error {
		container := dig.New()

		initializers := []interface{}{
			// actors
			app.NewManagerActor,
			app.NewSignalActor,

			// actors' dependencies
			NewRunGroup,
			KubernetesConfig,
			NewManager,
			func() logr.Logger {
				return logf.ZapLogger(true)
			},
		}

		for _, initFn := range initializers {
			if err := container.Provide(initFn); err != nil {
				return err
			}
		}

		// Invoke actors
		return container.Invoke(func(runGroup *run.Group, r app.ActorsResult) error {
			for _, actor := range r.Actors {
				runGroup.Add(actor.Run, actor.Interrupt)
			}

			// Run blocks until all the actors return. In the normal case, that’ll be when someone hits ctrl-C,
			// triggering the signal handler. If something breaks, its error will be propegated through. In all
			// cases, the first returned error triggers the interrupt function for all actors. And in this way,
			// we can reliably and coherently ensure that every goroutine that’s Added to the group is stopped,
			// when Run returns.
			return runGroup.Run()
		})
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Print(err)
		os.Exit(1)
	}
}
