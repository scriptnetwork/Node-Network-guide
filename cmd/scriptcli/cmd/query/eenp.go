package query

import (
	"encoding/json"
	"fmt"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/common"
	"github.com/scripttoken/script/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// eenpCmd represents the eenp command.
// Example:
//
//	scriptcli query eenp --height=10
var eenpCmd = &cobra.Command{
	Use:     "eenp",
	Short:   "Get elite edge node pool",
	Example: `scriptcli query eenp --height=10`,
	Run:     doEenpCmd,
}

func doEenpCmd(cmd *cobra.Command, args []string) {
	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	height := heightFlag
	res, err := client.Call("script.GetEenpByHeight", rpc.GetEenpByHeightArgs{Height: common.JSONUint64(height)})
	if err != nil {
		utils.Error("Failed to get elite edge node pool: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to get elite edge node pool: %v\n", res.Error)
	}
	json, err := json.MarshalIndent(res.Result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%s\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {
	eenpCmd.Flags().Uint64Var(&heightFlag, "height", uint64(0), "height of the block")
	eenpCmd.MarkFlagRequired("height")
}
