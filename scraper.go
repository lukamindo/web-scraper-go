package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// defining a data structure to store the scraped data
type PokemonProduct struct {
	url   string
	image string
	name  string
	price string
}

func contains(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func main() {
	var pokemonProducts []PokemonProduct

	var pagesToScrape []string

	pageToScrape := "https://scrapeme.live/shop/page/1/"

	pagesDiscovered := []string{pageToScrape}

	i := 1
	limit := 5

	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		newPaginationLink := e.Attr("href")

		// if the page discovered is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}
	writer.Write(headers)

	c.OnScraped(func(response *colly.Response) {

		for _, pokemonProduct := range pokemonProducts {
			// converting a PokemonProduct to an array of strings
			record := []string{
				pokemonProduct.url,
				pokemonProduct.image,
				pokemonProduct.name,
				pokemonProduct.price,
			}

			writer.Write(record)
		}
		defer writer.Flush()

		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]
			// incrementing the iteration counter
			i++
			// visiting a new page
			c.Visit(pageToScrape)
		}
	})

	c.Visit(pageToScrape)

}
