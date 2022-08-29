package query

import (
	"encoding/json"
	"fmt"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// versionCmd represents the version command.
// Example:
//
//	scriptcli query version
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get the Script version",
	Example: `scriptcli query version`,
	Run: func(cmd *cobra.Command, args []string) {
		client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

		res, err := client.Call("script.GetVersion", rpc.GetVersionArgs{})
		if err != nil {
			utils.Error("Failed to get version: %v\n", err)
		}
		if res.Error != nil {
			utils.Error("Failed to get version: %v\n", res.Error)
		}
		json, err := json.MarshalIndent(res.Result, "", "    ")
		if err != nil {
			utils.Error("Failed to parse server response: %v\n%s\n", err, string(json))
		}
		fmt.Println(string(json))
	},
}
