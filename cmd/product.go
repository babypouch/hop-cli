package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"os"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gosimple/slug"
	"github.com/kwngo/hop-cli/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(productCmd)
	productCmd.AddCommand(createCmd)
}

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Product related subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like create")
	},
}

var createCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a list of products to strapi",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating products from file: ", args[0])
		jsonFile, err := os.Open(args[0])
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var productInputs ProductInputs
		json.Unmarshal(byteValue, &productInputs)

		restyClient := utils.GetRestyClient()

		for i := 0; i < len(productInputs.Data); i++ {
			newSlug := slug.Make(productInputs.Data[i].Name)

			newProductAttributes := &ProductRequestAttributes{
				Name:         productInputs.Data[i].Name,
				Slug:         newSlug,
				Amount:       productInputs.Data[i].Amount,
				OriginalURL:  productInputs.Data[i].OriginalURL,
				Description:  productInputs.Data[i].Description,
				Thumbnail:    productInputs.Data[i].Thumbnail,
				PrimaryImage: productInputs.Data[i].PrimaryImage,
				Media:        productInputs.Data[i].Media,
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
			fmt.Println(productRes)
			product := productRes.Result().(*ProductResponse)
			uploadRequest := restyClient.R()
			for j := 0; j < len(productInputs.Data[i].Media); j++ {
				fileName, err := buildFileNameFromURL(productInputs.Data[i].Media[j])
				if err != nil {
					fmt.Println(err)
					return
				}

				imageRes, err := restyClient.R().Get(productInputs.Data[i].Media[j])
				if err != nil {
					fmt.Println(err)
					return
				}
				_, _, err = image.Decode(bytes.NewReader(imageRes.Body()))
				if err != nil {
					continue
				}

				uploadRequest = uploadRequest.SetFileReader("files", fileName, bytes.NewReader(imageRes.Body()))
			}
			fmt.Println("woah!!!")
			fmt.Println(uploadRequest)

			uploadRes, err := uploadRequest.
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

		}
	},
}
