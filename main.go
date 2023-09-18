package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Product struct {
	Name   string `json:"product"`
	Price  int    `json:"price"`
	Rating int    `json:"rating"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Испоьзование: ./app имя файла")
		return
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()
	var products []Product
	products = readData(*file, filename)

	maxPriceProduct := Product{}
	maxRatingProduct := Product{}
	for _, product := range products {
		if product.Price > maxPriceProduct.Price {
			maxPriceProduct = product
		}
		if product.Rating > maxRatingProduct.Rating {
			maxRatingProduct = product
		}
	}

	fmt.Println("Самый дорогой продукт:", maxPriceProduct)
	fmt.Println("Продукт с самым высоким рейтингом:", maxRatingProduct)
}

func readData(file os.File, filename string) []Product {
	switch {
	case strings.HasSuffix(filename, ".csv"):
		return readCSV(file)
	case strings.HasSuffix(filename, ".json"):
		return readJSON(file)
	default:
		fmt.Println("Неподдерживаемый формат файла")
		return []Product{}
	}
}

func readCSV(file os.File) []Product {
	csvReader := csv.NewReader(&file)
	csvReader.FieldsPerRecord = -1
	products := make([]Product, 0)
	firstRow := true // Переменная для отслеживания первой строки

	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}

		if firstRow {
			firstRow = false
			continue
		}

		if len(record) != 3 {
			continue
		}

		name := record[0]
		price := record[1]
		rating := record[2]
		product := Product{
			Name:   name,
			Price:  parseInt(price),
			Rating: parseInt(rating),
		}
		products = append(products, product)
	}

	return products
}

func readJSON(file os.File) []Product {
	decoder := json.NewDecoder(&file)
	products := make([]Product, 0)
	_, err := decoder.Token()
	if err != nil {
		fmt.Println("Ошибка при чтении начала JSON-массива:", err)
		return []Product{}
	}

	for decoder.More() {
		var product Product
		if err := decoder.Decode(&product); err != nil {
			fmt.Println("Ошибка при декодировании JSON:", err)
			break
		}
		products = append(products, product)

	}

	_, err = decoder.Token()
	if err != nil {
		fmt.Println("Ошибка при чтении конца JSON-массива:", err)
	}

	return products
}

func parseInt(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в число:", err)
		return 0
	}
	return res
}
