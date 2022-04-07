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

// accountCmd represents the account command.
// Example:
//		scriptcli query account --address=0x98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5
var accountCmd = &cobra.Command{
	Use:     "account",
	Short:   "Get account status",
	Long:    `Get account status.`,
	Example: `scriptcli query account --address=0x98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5`,
	Run:     doAccountCmd,
}

func doAccountCmd(cmd *cobra.Command, args []string) {
	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	res, err := client.Call("script.GetAccount", rpc.GetAccountArgs{
		Address: addressFlag, Preview: previewFlag})
	if err != nil {
		utils.Error("Failed to get account details: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to get account details: %v\n", res.Error)
	}
	json, err := json.MarshalIndent(res.Result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {
	accountCmd.Flags().StringVar(&addressFlag, "address", "", "Address of the account")
	accountCmd.Flags().BoolVar(&previewFlag, "preview", false, "Preview account balance from the screened view")
	accountCmd.MarkFlagRequired("address")
}
