// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/osiloke/jumler/scraper"
	"github.com/spf13/cobra"
)

// categoryCmd represents the category command
var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Collects category listing",
	Long:  `Collects category listing.`,
	Run: func(cmd *cobra.Command, args []string) {
		pages, _ := cmd.Flags().GetInt("pages")
		category, _ := cmd.Flags().GetString("category")
		baseURL, _ := cmd.Flags().GetString("baseURL")
		filename := fmt.Sprintf("%s_%d.jsonlines", category, pages)
		apiURL, _ := cmd.Flags().GetString("dostow-endpoint")
		apiKey, _ := cmd.Flags().GetString("dostow-key")
		storeName, _ := cmd.Flags().GetString("store")

		os.Remove(filename)
		output, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		var res chan map[string]interface{}

		res, err = scraper.ScrapeCategory(baseURL, category, pages)
		if err != nil {
			return
		}
		if len(apiURL) > 0 && len(apiKey) > 0 {
			scraper.DostowWriter(apiURL, apiKey, storeName, category, res)
		} else {
			scraper.StreamItemResultsToIO(output, res)
		}
	},
}

func init() {
	rootCmd.AddCommand(categoryCmd)

	categoryCmd.Flags().StringP("dostow-endpoint", "e", "http://localhost:4995/v1/", "Dostow api endpoint")
	categoryCmd.Flags().StringP("dostow-key", "k", "21ca3b50-3c4c-497b-90e7-35e71adc53eb", "Dostow api key")
	categoryCmd.Flags().StringP("baseURL", "u", "https://www.jumia.com.ng/", "Url to jackbian forums")
	categoryCmd.Flags().StringP("category", "s", "groceries", "Category to scrape")
	categoryCmd.Flags().StringP("store", "r", "product", "store name")
	categoryCmd.Flags().IntP("pages", "z", 1, "Number of pages")
}
