package dl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CmdDL represents the download command.
var CmdDL = &cobra.Command{
	Use:   "dl",
	Short: "Download product catalog from MFG website (https://dropship.myfashiongrosir.com/shop/produk)",
	Long:  "Download product catalog from MFG website (https://dropship.myfashiongrosir.com/shop/produk). Example: mfg dl",
	Run:   run,
}

func init() {
	// CmdDL.Flags().BoolVarP(&verbose, "verbose", "v", true, "verbose mode")
}

func run(cmd *cobra.Command, args []string) {
	Start()
}

// Start download by web scraping then save in folder ./result
func Start() {
	var productUrlMap = make(map[string]int)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0"),
	)

	//* [Product List]
	// On every a element which has href attribute call callback
	c.OnHTML("body > div.container > div.row > div.col-lg-9 > div", func(e *colly.HTMLElement) {
		// get product url
		e.ForEach("div > a[href]", func(_ int, h *colly.HTMLElement) {
			link := h.Attr("href")

			if link[:51] == "https://dropship.myfashiongrosir.com/produk/detail/" {
				// save link
				productUrlMap[link]++
			}

			if link[:54] == "https://dropship.myfashiongrosir.com/shop/produk?page=" {
				c.Visit(e.Request.AbsoluteURL(link))
			}
		})

		// scrape next page
		e.ForEach("nav > ul > li > a[href]", func(i int, h *colly.HTMLElement) {
			rel := h.Attr("rel")
			if rel == "next" {
				link := h.Attr("href")
				c.Visit(link)
				fmt.Println("next page", i, rel, link)
			}
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on ...
	c.Visit("https://dropship.myfashiongrosir.com/shop/produk")

	//* [Product Detail]
	products := []Product{}
	c.OnHTML("body > div.container.d-none.d-lg-block.d-xl-block.d-xxl-block", func(h *colly.HTMLElement) {
		var p Product

		p.Name = strings.TrimSpace(h.DOM.Find("p.text-question.text-judul-produk").Text())
		p.Url = findURL(p.Name, productUrlMap)
		p.BuyPrice = strings.TrimSpace(h.DOM.Find("p.text-price").Text())
		p.SellPrice = strings.TrimSpace(strings.ReplaceAll(h.DOM.Find("p.text-mobile").Text(), "Potensi Harga Jual: ", ""))
		p.ImageUrl = h.ChildAttr("div:nth-child(5) > div > div.position-relative > img", "src")
		p.Description = strings.TrimSpace(h.DOM.Find("#detail-deskripsi textarea").Text())

		// product colors
		pattern := "warnaDiClick\\([0-9]+,'([a-zA-Z]+)'\\)"
		onclicks := h.ChildAttrs("div:nth-child(5) > div > div.div-text-detail > div button", "onclick")
		for _, v := range onclicks {
			reg := regexp.MustCompile(pattern)
			submatch := reg.FindAllStringSubmatch(v, -1)

			if len(submatch) > 0 {
				if len(submatch[0]) > 1 {
					color := submatch[0][1]
					p.Colors = append(p.Colors, color)
				}
			}
		}

		// foto lainnya
		h.ForEach("div:nth-child(5) > div > div.foto-lainnya > div ", func(_ int, x *colly.HTMLElement) {
			srcs := x.ChildAttrs("img", "src")
			p.ImageUrls = srcs
		})

		// product detail
		h.ForEach("#detail-produk", func(i int, h *colly.HTMLElement) {

			_html_, _ := h.DOM.Html()
			spans := strings.Split(_html_, "<br/>")

			bm := bluemonday.StripTagsPolicy()
			for _, v := range spans {
				v := strings.TrimSpace(bm.Sanitize(v))
				if strings.HasPrefix(v, "Material: ") {
					p.Material = strings.ReplaceAll(v, "Material: ", "")
				} else if strings.HasPrefix(v, "Style: ") {
					p.Style = strings.ReplaceAll(v, "Style: ", "")
				} else if strings.HasPrefix(v, "Berat: ") {
					p.Weight = strings.ReplaceAll(v, "Berat: ", "")
				}
			}
		})

		products = append(products, p)
	})

	// scrape product detail page
	for url := range productUrlMap {
		c.Visit(url)
	}

	logrus.Infof("Scraping Complete, got %d products", len(productUrlMap))

	err := GenerateTokpedExcelFile(products)
	if err != nil {
		logrus.Errorf("failed generate tokopedia. err=%s", err)
	}
	err = GenerateShopeeExcelFile(products)
	if err != nil {
		logrus.Errorf("failed generate shopee. err=%s", err)
	}
}

func findURL(title string, urls map[string]int) string {
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))

	for k, _ := range urls {
		if strings.Contains(k, slug) {
			return k
		}
	}

	return ""
}
