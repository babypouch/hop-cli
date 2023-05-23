package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"net/url"
	"os"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/babypouch/hop-cli/cmd"
	"github.com/babypouch/hop-cli/utils"
	"github.com/go-resty/resty/v2"
	"github.com/gosimple/slug"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cmd.RootCmd.AddCommand(productCmd)
	productCmd.AddCommand(createCmd)
}

var productCmd = &cobra.Command{
	Use:   "products",
	Short: "Product related subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like create")
	},
}

var createCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a list of products to strapi",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostURL := viper.GetString("Host")
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
		var product *ProductResponse

		for i := 0; i < len(productInputs.Data); i++ {
			newSlug := slug.Make(productInputs.Data[i].Name)
			var brandRes *resty.Response
			newProductAttributes := &ProductRequestAttributes{
				Name:         productInputs.Data[i].Name,
				Slug:         newSlug,
				Amount:       productInputs.Data[i].Amount,
				OriginalURL:  productInputs.Data[i].OriginalURL,
				Pid:          productInputs.Data[i].Pid,
				Description:  productInputs.Data[i].Description,
				Thumbnail:    productInputs.Data[i].Thumbnail,
				PrimaryImage: productInputs.Data[i].PrimaryImage,
				Media:        productInputs.Data[i].Media,
			}
			brandParams := url.Values{}
			brandParams.Add("filters[name][$eq]", productInputs.Data[i].Brand)
			brandRes, _ = restyClient.R().
				SetResult(&BrandsResponse{}).
				Get(hostURL + "/api/brands?" + brandParams.Encode())
			brand := brandRes.Result().(*BrandsResponse)
			if len(brand.Data) == 0 {
				fmt.Println("Brand " + productInputs.Data[i].Brand + " does not exist or can't be found.")
				newBrandSlug := slug.Make(productInputs.Data[i].Brand)
				postBrandRes, _ := restyClient.R().
					SetResult(&BrandResponse{}).
					SetBody(BrandRequestBody{
						Data: BrandRequestAttributes{
							Name: productInputs.Data[i].Brand,
							Slug: newBrandSlug,
						},
					}).
					Post(hostURL + "/api/brands")
				fmt.Println(postBrandRes)
				newBrand := postBrandRes.Result().(*BrandResponse)
				newProductAttributes.Brand = newBrand.Data.Id
				fmt.Println("Brand " + productInputs.Data[i].Brand + " created.")
			} else {
				fmt.Println("Brand found: ", productInputs.Data[i].Brand)
				fmt.Println(brandRes)
				newProductAttributes.Brand = brand.Data[0].Id
			}
			collectionParams := url.Values{}
			collectionParams.Add("filters[name][$eq]", productInputs.Data[i].Collections)

			collectionRes, _ := restyClient.R().
				SetResult(&CollectionsResponse{}).
				Get(hostURL + "/api/collections?" + collectionParams.Encode())
			collection := collectionRes.Result().(*CollectionsResponse)

			if len(collection.Data) == 0 {
				fmt.Println("Collection " + productInputs.Data[i].Collections + " does not exist or can't be found.")
				newCollectionSlug := slug.Make(productInputs.Data[i].Collections)
				postCollectionRes, _ := restyClient.R().
					SetResult(&CollectionResponse{}).
					SetBody(CollectionRequestBody{
						Data: CollectionRequestAttributes{
							Name: productInputs.Data[i].Collections,
							Slug: newCollectionSlug,
						},
					}).
					Post(hostURL + "/api/collections")
				newCollection := postCollectionRes.Result().(*CollectionResponse)
				fmt.Println(postCollectionRes)
				newProductAttributes.Collections = []int{newCollection.Data.Id}
				fmt.Println("Collection " + productInputs.Data[i].Collections + " created.")
			} else {
				newProductAttributes.Collections = []int{collection.Data[0].Id}
			}

			getProductRes, _ := restyClient.R().
				SetResult(&ProductResponse{}).
				Get(hostURL + "/api/products/find-by-slug/" + newSlug)
			if getProductRes.IsSuccess() {
				fmt.Println("Product with slug " + newSlug + " already exists.")
				validate := func(input string) error {
					return nil
				}

				prompt := promptui.Prompt{
					Label: "Do you want to overwrite this product? (y/n)",

					Validate: validate,
				}

				result, err := prompt.Run()

				if err != nil {
					fmt.Println("Failed to update product with slug: ", newSlug)
					fmt.Printf("Prompt failed %v\n", err)
					return
				}

				if result == "y" {
					getProduct := getProductRes.Result().(*ProductResponse)

					productRes, _ := restyClient.R().
						SetResult(&ProductResponse{}).
						SetBody(ProductRequest{
							Data: newProductAttributes,
						}).
						Put(fmt.Sprintf("%s/api/products/%v", hostURL, getProduct.Data.Id))
					if productRes.IsSuccess() {
						fmt.Println("Successfully updated product with slug: ", newSlug)
						product = productRes.Result().(*ProductResponse)
					} else {
						fmt.Println("Failed to update product with slug: ", newSlug)
						fmt.Println(productRes)
						continue
					}
				} else {
					fmt.Println("Skip product with slug: ", newSlug)
					continue
				}
			} else {
				productRes, _ := restyClient.R().
					SetResult(&ProductResponse{}).
					SetBody(ProductRequest{
						Data: newProductAttributes,
					}).
					Post(hostURL + "/api/products")
				if productRes.IsError() {
					fmt.Println("There was an issue creating ", productInputs.Data[i].Name)
					fmt.Println(productRes)
					continue
				}
				fmt.Println("Creating product with slug: " + newSlug)
				product = productRes.Result().(*ProductResponse)
			}

			// Just skip the whole update block if there are no media objects
			if len(productInputs.Data[i].Media) == 0 {
				continue
			}

			uploadRequest := restyClient.R()
			var imagesToUpload []*ImageData

			if len(productInputs.Data[i].Thumbnail) > 0 {
				imageName, err := utils.BuildFileNameFromURL(productInputs.Data[i].Thumbnail)
				if err != nil {
					fmt.Println(err)
					continue
				}
				imageRes, err := restyClient.R().Get(productInputs.Data[i].Thumbnail)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, _, err = image.Decode(bytes.NewReader(imageRes.Body()))
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("Uploading thumbnail to strapi...")
				thumbnailUploadRes, _ := restyClient.R().
					SetFileReader("files", imageName, bytes.NewReader(imageRes.Body())).
					SetFormData(map[string]string{
						"refId": fmt.Sprint(product.Data.Id),
						"ref":   "api::product.product",
						"field": "main_thumbnail",
					}).
					Post(hostURL + "/api/upload")

				if thumbnailUploadRes.IsError() {
					fmt.Println(thumbnailUploadRes)
					return
				}

				fmt.Println("Thumbnail successfully uploaded to ", product.Data.Attributes.Name)
			}

			if len(productInputs.Data[i].PrimaryImage) > 0 {
				imageName, err := utils.BuildFileNameFromURL(productInputs.Data[i].PrimaryImage)
				if err != nil {
					fmt.Println(err)
					continue
				}
				imageRes, err := restyClient.R().Get(productInputs.Data[i].PrimaryImage)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, _, err = image.Decode(bytes.NewReader(imageRes.Body()))
				if err != nil {
					fmt.Println(err)
					continue
				}
				imagesToUpload = append(imagesToUpload, &ImageData{
					Name: imageName,
					Body: imageRes.Body(),
				})

			}

			for _, fileURL := range productInputs.Data[i].Media {
				imageName, err := utils.BuildFileNameFromURL(fileURL)
				if err != nil {
					fmt.Println(err)
					continue
				}
				imageRes, err := restyClient.R().Get(fileURL)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, _, err = image.Decode(bytes.NewReader(imageRes.Body()))
				if err != nil {
					fmt.Println(err)
					continue
				}
				imagesToUpload = append(imagesToUpload, &ImageData{
					Name: imageName,
					Body: imageRes.Body(),
				})
			}

			// Skip the rest of the code if there are no images to upload
			if len(imagesToUpload) == 0 {
				continue
			}

			for _, imageData := range imagesToUpload {
				uploadRequest = uploadRequest.SetFileReader("files", imageData.Name, bytes.NewReader(imageData.Body))
			}

			uploadRes, _ := uploadRequest.
				SetFormData(map[string]string{
					"refId": fmt.Sprint(product.Data.Id),
					"ref":   "api::product.product",
					"field": "product_images",
				}).
				Post(hostURL + "/api/upload")
			if uploadRes.IsError() {
				fmt.Println(uploadRes)
				return
			}
			fmt.Println(len(imagesToUpload), " successfully uploaded to ", product.Data.Attributes.Name)
		}
	},
}
