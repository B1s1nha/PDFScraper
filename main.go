package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"github.com/gocolly/colly"
)

func main() {
	dataFiltrada := "24/04/2024"
	contador := 0

	c := colly.NewCollector()
	pdf := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitando página inicial: ", r.URL)
	})

	c.OnHTML("div.card.article-card.legal", func(e *colly.HTMLElement) {
		date := e.ChildText("div.article-card__summary_date time")

		if strings.Contains(date, dataFiltrada) {
			contador++
			link := e.ChildAttr("a", "href")
			fullLink := e.Request.AbsoluteURL(link)
			fmt.Printf("Seguindo o link: %s\n", fullLink)
			pdf.Visit(fullLink)
		}
	})

	pdf.OnHTML("a[href$='.pdf']", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href != "" {
			fullPDFLink := e.Request.AbsoluteURL(href)
			fmt.Println("Link do PDF encontrado: ", fullPDFLink)
			downloadPDF(fullPDFLink)
		}
	})

	url := "https://gauchazh.clicrbs.com.br/publicidade-legal/ultimas-noticias/"

	err := c.Visit(url)
	if err != nil {
		fmt.Println("Erro ao visitar o site:", err)
	}

	fmt.Printf("Número de ocorrências para a data %s: %d\n", dataFiltrada, contador)
}

func downloadPDF(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer o download do PDF: ", err)
		return
	}
	defer response.Body.Close()

	out, err := os.Create("balanco.pdf")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo PDF: ", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println("Erro ao salvar o PDF: ", err)
		return
	}

	fmt.Println("PDF baixado com sucesso!")
}

