package mentions

import (
	"fmt"
	"os"

	"github.com/babypouch/hop-cli/cmd/products"
	"github.com/babypouch/hop-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	mentionsCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get search results for product by id",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Configure list...")
		hostURL := viper.GetString("Host")
		validate := func(input string) error {
			return nil
		}
		prompt := promptui.Prompt{
			Label:    "Get product with id: ",
			Validate: validate,
		}

		restyClient := utils.GetRestyClient()

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			fmt.Println(err)
			os.Exit(1)
		}
		getRes, _ := restyClient.R().
			SetResult(&products.ProductResponse{}).
			Get(hostURL + "/api/products/" + result)
		if getRes.IsError() {
			fmt.Println("product with id " + result + " not found.")
			os.Exit(1)
		}
		fmt.Println("Product with id " + result + " found.")
		product := getRes.Result().(*products.ProductResponse)
		serp := utils.GetSerpClient()
		serpRes := serp.GetLive(map[string]string{
			"q":         product.Data.Attributes.Name + " review",
			"domain":    "google.com",
			"lang":      "en",
			"device":    "desktop",
			"serp_type": "web",
			"loc_id":    "9041160",
			"loc":       "Toronto Pearson International Airport,Ontario,Canada",
		})

		serpResult := serpRes.Result().(*utils.SerpLiveResponse)
		fmt.Println(serpResult.GetOrganic())
		selectItems := serpResult.GetOrganic()

		templates := &promptui.SelectTemplates{
			Label:    "{{ .Title }} - {{.Link}}",
			Active:   "\U0001F336 {{ .Title | cyan }} ({{ .Link | red }})",
			Inactive: "  {{ .Title | cyan }} ({{ .Link | red }})",
			Selected: "\U0001F336 {{ .Title | red | cyan }}",
			Details: `
	--------- Pepper ----------
	{{ "Title:" | faint }}	{{ .Title }}
	{{ "Link:" | faint }}	{{ .Link }}`,
		}

		selectPrompt := promptui.Select{
			Label:     "Select mention to add",
			Templates: templates,
			Items:     selectItems,
			Size:      20,
		}

		selectedIndex, _, err := selectPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		microlink := utils.GetMicrolinkClient()
		selectedResult := selectItems[selectedIndex]
		microRes := microlink.GetMetaData(selectedResult.Link)
		if microRes.IsError() {
			fmt.Println("There was an error getting metadata for " + selectedResult.Link)
			os.Exit(1)
		}
		fmt.Println(microRes)
	},
}
