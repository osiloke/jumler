package scraper

var categoryItemSchema = `
{
	"name": "Jumia Category",
	"id": "jumia-category",
	"type": "object",
	"css": [".products .sku.-gallery"],
	"properties": [
		{
			"css": ["a.link", "href"],
			"id": "link",
			"type": "string"
		},
		{
			"css": [".lazy.image", "data-src"],
			"id": "image",
			"type": "string"
		},
		{
			"css": ["h2.title span.name"],
			"id": "name",
			"type": "string"
		},
		{
			"css": ["h2.title span.brand"],
			"id": "brand",
			"type": "string"
		},
		{
			"css": ["button", "data-sku"],
			"id": "sku",
			"type": "string"
		}
	]
}
`
