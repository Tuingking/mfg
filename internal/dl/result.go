package dl

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	tokpedTmpl = "internal/dl/result/tokopedia.xlsx"
	shopeeTmpl = "internal/dl/result/shopee.xlsx"

	templatePath = "internal/dl/template/"
	resultPath   = "internal/dl/result/"
)

func GenerateTokpedExcelFile(products []Product) error {
	err := copyTemplate("tokopedia")
	if err != nil {
		return errors.Wrap(err, "copyTemplate tokopedia")
	}

	f, err := excelize.OpenFile(tokpedTmpl)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	sheetName := "ISI Template Impor Produk"
	sheetIndex := f.GetSheetIndex(sheetName)
	f.SetActiveSheet(sheetIndex)

	offset := 4
	for i, v := range products {
		// B - Nama Produk
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+offset), v.Name)
		// C - Deskripsi Produk
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+offset), v.Description)
		// D - Katergori Kode
		// f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+offset), v.BuyPrice)
		// E - Berat
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+offset), cleanWeight(v.Weight))
		// F - Minimum Pemesanan
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+offset), 1)
		// I - Kondisi
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", i+offset), "Baru")
		// J - Gambar 1
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", i+offset), v.ImageUrl)
		// Gambar 2-5 (K-N)
		col := 'K'
		for _, img := range v.ImageUrls {
			if col == 'O' {
				break
			}
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(col), i+offset), img)
			col++
		}

		f.SetCellValue(sheetName, fmt.Sprintf("R%d", i+offset), v.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("S%d", i+offset), "Aktif")
		f.SetCellValue(sheetName, fmt.Sprintf("T%d", i+offset), 10)
		f.SetCellValue(sheetName, fmt.Sprintf("U%d", i+offset), cleanPrice(v.BuyPrice))
		f.SetCellValue(sheetName, fmt.Sprintf("V%d", i+offset), "optional")
	}

	err = f.Save()
	if err != nil {
		logrus.Errorf("failed save file. err=%+v\n", err)
		return errors.Wrap(err, "save file")
	}

	logrus.Info("result saved in result/tokopedia.xlsx ✅")

	return nil
}

func copyTemplate(templateName string) error {
	var (
		src  = fmt.Sprintf("%s/%s.xlsx", templatePath, templateName)
		dest = fmt.Sprintf("%s/%s.xlsx", resultPath, templateName)
	)
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func cleanPrice(v string) string {
	m := regexp.MustCompile(`[^\d]`)
	return m.ReplaceAllString(v, "")
}

func cleanWeight(v string) string {
	v = strings.ReplaceAll(v, "Gram", "")
	v = strings.ReplaceAll(v, "gr", "")
	return strings.TrimSpace(v)
}

func GenerateShopeeExcelFile(products []Product) error {
	err := copyTemplate("shopee")
	if err != nil {
		return errors.Wrap(err, "copyTemplate shopee")
	}

	f, err := excelize.OpenFile(shopeeTmpl)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	sheetName := "Template"
	sheetIndex := f.GetSheetIndex(sheetName)
	f.SetActiveSheet(sheetIndex)

	offset := 6
	for i, v := range products {
		// B - Nama Produk
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+offset), v.Name)
		// C - Deskripsi Produk
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+offset), v.Description)
		// D - SKU Induk (optional)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+offset), v.Name)
		// E - Produk Berbahaya
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+offset), "No (ID)")
		// L - Price
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", i+offset), cleanPrice(v.BuyPrice))
		// M - Stock
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", i+offset), 10)
		// N - Foto Sampul
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", i+offset), v.ImageUrl)
		// [O-W] - Gambar 1-8
		col := 'O'
		for _, img := range v.ImageUrls {
			if col == 'X' {
				break
			}
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(col), i+offset), img)
			col++
		}
		// X - Weight
		f.SetCellValue(sheetName, fmt.Sprintf("X%d", i+offset), cleanWeight(v.Weight))
		f.SetCellValue(sheetName, fmt.Sprintf("AA%d", i+offset), "Aktif")
		f.SetCellValue(sheetName, fmt.Sprintf("AB%d", i+offset), "Nonaktif")
		f.SetCellValue(sheetName, fmt.Sprintf("AC%d", i+offset), "Nonaktif")
	}

	err = f.Save()
	if err != nil {
		logrus.Errorf("failed save file. err=%+v\n", err)
		return errors.Wrap(err, "save file")
	}

	logrus.Info("result saved in result/shopee.xlsx ✅")

	return nil
}
