package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"io/ioutil"
	"os"

	"github.com/kwngo/hop-cli/models"
	"github.com/kwngo/hop-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(mentionsCmd)
	mentionsCmd.AddCommand(uploadCmd)
}

var mentionsCmd = &cobra.Command{
	Use:   "mentions",
	Short: "Set up mentions",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

type FuzzySearchItem struct {
	Id         int               `json:"id"`
	Attributes ProductAttributes `json:"attributes"`
}

type FuzzySearchItems struct {
	Item FuzzySearchItem `json:"item"`
}

type FuzzySearchResponse struct {
	Data []FuzzySearchItems `json:"data"`
}

func (f FuzzySearchResponse) NameList() []string {
	var list []string
	for _, item := range f.Data {
		list = append(list, item.Item.Attributes.Name)
	}
	return list
}

func (f FuzzySearchResponse) ProductList() []ProductAttributes {
	var list []ProductAttributes
	for _, item := range f.Data {
		list = append(list, item.Item.Attributes)
	}
	return list
}
func (f FuzzySearchResponse) FindByName(name string) FuzzySearchItem {
	for _, item := range f.Data {
		if item.Item.Attributes.Name == name {
			return item.Item
		}
	}
	return FuzzySearchItem{}
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload mentions to strapi",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uploading reviews...")
		hostURL := viper.GetString("Host")
		jsonFile, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var mentionInputs models.MentionInputs
		json.Unmarshal(byteValue, &mentionInputs)

		restyClient := utils.GetRestyClient()
		for _, input := range mentionInputs.Data {
			params := url.Values{}
			params.Add("_search", input.ProductName)
			fuzzySearchRes, _ := restyClient.R().
				SetResult(&FuzzySearchResponse{}).
				Get(hostURL + "/api/products/fuzzy-search?" + params.Encode())
			if fuzzySearchRes.IsError() {
				fmt.Println("Issue running fuzzy search: " + input.ProductName)
				continue
			}

			fuzzySearchResponse := fuzzySearchRes.Result().(*FuzzySearchResponse)
			if len(fuzzySearchResponse.Data) == 0 {
				fmt.Println("No results returned for fuzzy search on: " + input.ProductName)
				continue
			}
			prompt := promptui.Select{
				Label: "Select a result",
				Items: fuzzySearchResponse.NameList(),
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			matchedResult := fuzzySearchResponse.FindByName(result)

			fmt.Printf("You choose %q\n", matchedResult.Attributes.Name)
			fmt.Printf("You choose %v\n", matchedResult.Id)
		}
	},
}
