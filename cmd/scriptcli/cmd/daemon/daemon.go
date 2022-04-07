package daemon

import (
	"context"
	"log"
	"sync"

	"github.com/spf13/cobra"
	"github.com/scripttoken/script/cmd/scriptcli/rpc"
)

// startDaemonCmd runs the scriptcli daemon
// Example:
//		scriptcli daemon start --port=16889
var startDaemonCmd = &cobra.Command{
	Use:     "start",
	Short:   "Run the thatacli daemon",
	Long:    `Run the thatacli daemon.`,
	Example: `scriptcli daemon start --port=16889`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := cmd.Flag("config").Value.String()
		server, err := rpc.NewScriptcliRPCServer(cfgPath, portFlag)
		if err != nil {
			log.Fatalf("Failed to run the ScriptCli Daemon: %v", err)
		}
		daemon := &ScriptCliDaemon{
			RPC: server,
		}
		daemon.Start(context.Background())
		daemon.Wait()
	},
}

func init() {
	startDaemonCmd.Flags().StringVar(&portFlag, "port", "16889", "Port to run the ScriptCli Daemon")
}

type ScriptCliDaemon struct {
	RPC *rpc.ScriptcliRPCServer

	// Life cycle
	wg      *sync.WaitGroup
	quit    chan struct{}
	ctx     context.Context
	cancel  context.CancelFunc
	stopped bool
}

func (d *ScriptCliDaemon) Start(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	d.ctx = c
	d.cancel = cancel

	if d.RPC != nil {
		d.RPC.Start(d.ctx)
	}
}

func (d *ScriptCliDaemon) Stop() {
	d.cancel()
}

func (d *ScriptCliDaemon) Wait() {
	if d.RPC != nil {
		d.RPC.Wait()
	}
}
