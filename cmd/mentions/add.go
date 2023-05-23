package mentions

import (
	"bytes"
	"fmt"
	"os"

	"github.com/babypouch/hop-cli/cmd/products"
	"github.com/babypouch/hop-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	mentionsCmd.AddCommand(addCmd)
}

type PublisherRequest struct {
	Data PublisherRequestAttributes `json:"data"`
}

type PublisherRequestAttributes struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PublisherAttributes struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PublisherData struct {
	Id         int                 `json:"id"`
	Attributes PublisherAttributes `json:"attributes"`
}

type PublisherResponse struct {
	Data []PublisherData `json:"data"`
}

type NewPublisherResponse struct {
	Data PublisherData `json:"data"`
}

type MentionRequestData struct {
	Title              string `json:"title"`
	URL                string `json:"url"`
	Product            int    `json:"product"`
	Publisher          int    `json:"publisher"`
	MentionPublishedAt string `json:"mention_published_at"`
}

type MentionAttributes struct {
	Title              string `json:"title"`
	URL                string `json:"url"`
	MentionPublishedAt string `json:"mention_published_at"`
}

type MentionData struct {
	Id         int               `json:"id"`
	Attributes MentionAttributes `json:"attributes"`
}

type MentionResponse struct {
	Data PublisherData `json:"data"`
}

type MentionRequest struct {
	Data MentionRequestData `json:"data"`
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add mentions to a product by id",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Adding mentions...")
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
			"q":         product.Data.Attributes.Name + " reviews",
			"domain":    "google.com",
			"lang":      "en",
			"device":    "desktop",
			"serp_type": "web",
			"loc_id":    "9041160",
			"loc":       "Toronto Pearson International Airport,Ontario,Canada",
		})

		serpResult := serpRes.Result().(*utils.SerpLiveResponse)
		selectResults := serpResult.GetOrganicResults()
		selectPromptItems := serpResult.GetOrganicSelectItems()

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U0001F336 {{ . | cyan }})",
			Inactive: "  {{ . | cyan }} ",
			Selected: "\U0001F336 {{ . | red | cyan }}",
		}

		selectedIndex := -1
		var selectedItems []utils.SerpSearchResultsOrganic

		for selectedIndex != 0 {
			selectPrompt := promptui.Select{
				Label:     "Select mention to add",
				Templates: templates,
				Items:     selectPromptItems,
				Size:      30,
			}
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			selectedIndex, _, err = selectPrompt.Run()
			selectedResult := selectResults[selectedIndex]
			selectedItems = append(selectedItems, selectedResult)
			fmt.Println(selectedItems)
		}

		for _, selectedResult := range selectedItems {
			fmt.Println("Adding selected items...")
			microlink := utils.GetMicrolinkClient()
			microRes := microlink.GetMetaData(selectedResult.Link)
			if microRes.IsError() {
				fmt.Println("There was an error getting metadata for " + selectedResult.Link)
				fmt.Println(microRes)
				os.Exit(1)
			}

			metadata := microRes.Result().(*utils.MicrolinkResponse)
			fmt.Println(metadata)
			getPublisherRes, _ := restyClient.R().
				SetResult(&PublisherResponse{}).
				SetQueryParam("filters[name][$eq]", metadata.Data.Publisher).
				Get(hostURL + "/api/publishers")
			pubRes := getPublisherRes.Result().(*PublisherResponse)
			if len(pubRes.Data) == 0 {
				fmt.Println("Publisher " + metadata.Data.Publisher + " not found.")
				publisherRequestData := PublisherRequestAttributes{
					Name: metadata.Data.Publisher,
					URL:  metadata.Data.URL,
				}

				publisherRes, _ := restyClient.R().
					SetResult(&NewPublisherResponse{}).
					SetBody(PublisherRequest{
						Data: publisherRequestData,
					}).
					Post(hostURL + "/api/publishers")

				if publisherRes.IsError() {
					fmt.Println("There was an error adding publisher " + metadata.Data.Publisher)
					fmt.Println(publisherRes)
					os.Exit(1)
				}
				fmt.Println(publisherRes)
				newPublisher := publisherRes.Result().(*NewPublisherResponse)
				mentionRequestData := MentionRequestData{
					Title:              metadata.Data.Title,
					Publisher:          newPublisher.Data.Id,
					URL:                metadata.Data.URL,
					Product:            product.Data.Id,
					MentionPublishedAt: metadata.Data.Date,
				}
				mentionsRes, _ := restyClient.R().
					SetResult(&MentionResponse{}).
					SetBody(MentionRequest{
						Data: mentionRequestData,
					}).
					Post(hostURL + "/api/mentions")
				fmt.Println(mentionsRes)

				fmt.Println("Uploading logo...")
				imageName, err := utils.BuildFileNameFromURL(metadata.Data.Logo.URL)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				imageRes, err := restyClient.R().Get(metadata.Data.Logo.URL)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				uploadRes, _ := restyClient.R().
					SetFileReader("files", imageName, bytes.NewReader(imageRes.Body())).
					SetFormData(map[string]string{
						"refId": fmt.Sprint(newPublisher.Data.Id),
						"ref":   "api::publisher.publisher",
						"field": "logo",
					}).
					Post(hostURL + "/api/upload")
				if uploadRes.IsError() {
					fmt.Println(uploadRes)
					os.Exit(1)
				}
				fmt.Println("Logo with url " + metadata.Data.Logo.URL + " uploaded.")
			} else {
				mentionRequestData := MentionRequestData{
					Title:              metadata.Data.Title,
					Publisher:          pubRes.Data[0].Id,
					URL:                metadata.Data.URL,
					Product:            product.Data.Id,
					MentionPublishedAt: metadata.Data.Date,
				}
				mentionsRes, _ := restyClient.R().
					SetBody(MentionRequest{
						Data: mentionRequestData,
					}).
					Post(hostURL + "/api/mentions")
				fmt.Println(mentionsRes)
			}
		}

	},
}
