package cmd

import (
	"bytes"
	"fmt"

	"github.com/kwngo/hop-cli/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(attachCmd)
}

var attachCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a image to strapi",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Testing upload")
		url := args[0]
		fileName, err := buildFileNameFromURL(url)
		if err != nil {
			fmt.Println(err)
		}
		restyClient := utils.GetRestyClient()

		res, err := restyClient.R().Get(url)

		if err != nil {
			fmt.Println(err)
		}
		newProductAttributes := &ProductRequestAttributes{
			Name: "product_test_t",
			Slug: "product-test-t",
		}

		productRes, err := restyClient.R().
			SetResult(&ProductResponse{}).
			SetBody(ProductRequest{
				Data: newProductAttributes,
			}).
			Post("http://localhost:1337/api/products")

		if err != nil {
			fmt.Println(err)
			return
		}
		product := productRes.Result().(*ProductResponse)
		fmt.Println(productRes)
		fmt.Println("product: ")
		fmt.Println(product)
		uploadRes, err := restyClient.R().
			SetFileReader("files", fileName, bytes.NewReader(res.Body())).
			SetFormData(map[string]string{
				"refId": fmt.Sprint(product.Data.Id),
				"ref":   "api::product.product",
				"field": "product_images",
			}).
			Post("http://localhost:1337/api/upload")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(uploadRes)
	},
}
