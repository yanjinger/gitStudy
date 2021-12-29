package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"github.com/yanjinger/gitStudy/wave/utils"
	"golang.org/x/text/encoding/unicode"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {

	viper.SetConfigFile("./config.toml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("don't find ./config.toml !")
	}
	startRow := viper.GetInt("startRow")
	sheetName := viper.GetString("sheetName")
	utils.Suffix = viper.GetString("Suffix")

	t1 := time.Now()

	dir, _ := os.Getwd()
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 寻找模板
	template := ""
	for _, f := range filelist {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".xlsx") {
			template = f.Name()
			break
		}
	}
	template = filepath.Join(dir, template)
	fmt.Println("Notice  template is ", template)

	excel, err := excelize.OpenFile(template)

	if err != nil {
		fmt.Println(err)
		return
	}

	sheetIndex := excel.GetSheetIndex(sheetName)
	if sheetIndex < 0 {
		fmt.Println("找不到", sheetName)
		return
	}
	allFile := []string{}
	allFile = Scan(dir, allFile)

	excelFile1, err := excelize.OpenFile(template)
	cols, err := excelFile1.GetCols(sheetName)
	rows, err := excelFile1.GetRows(sheetName)
	rowNum := len(rows)
	colNum := len(cols)

	fmt.Println("start process...")

	var wg sync.WaitGroup
	for _, info := range allFile {
		infoV := info
		wg.Add(1)
		go func() {
			fmt.Println("\t process ", infoV)
			newName := strings.Replace(infoV, utils.Suffix, utils.DestSuffix, -1)
			utils.CopyFile(template, newName)
			excelFile, err := excelize.OpenFile(newName)
			if err != nil {
				fmt.Println(err)
				return
			}

			txtSheet, err := utils.ParseTxt(infoV, "\t", unicode.UTF8.NewDecoder())
			if rowNum < txtSheet.RowNum {
				fmt.Println(" template row too less.  template row=", rowNum, "dataTable row=", txtSheet.RowNum)
				return
			}

			for i := startRow; i < txtSheet.RowNum; i++ {
				if len(txtSheet.Rows[i].Cells) != 5 {
					fmt.Println(" need 5 column")
					return
				}
				hz, err := txtSheet.Rows[i].Cells[0].Float()
				if err != nil {
					fmt.Println(err)
					return
				}
				hz = hz / 1e9
				txtSheet.Rows[i].Cells[0].Value = fmt.Sprintf("%.2f", hz)

				for j := 0; j < 5; j++ {
					columnName := utils.GetColumnIndexName1(j)
					axis := fmt.Sprintf("%s%d", columnName, i-startRow+3+1)
					excelFile.SetCellValue(sheetName, axis, txtSheet.Rows[i].Cells[j].Value)
				}
			}

			startCol := utils.GetColumnIndexName1(0)
			endCol := utils.GetColumnIndexName1(colNum - 1)
			sqref := fmt.Sprintf("%s%d:%s%d", startCol, txtSheet.RowNum, endCol, rowNum)
			err = excelFile.DeleteDataValidation(sheetName, sqref)

			for v := rowNum; v >= txtSheet.RowNum+3+1-startRow; v-- {
				err := excelFile.RemoveRow(sheetName, v)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if err != nil {
				fmt.Println(err)
				return
			}
			excelFile.SaveAs(newName)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("process end. time =", time.Now().Sub(t1).String())
}

func Scan(dir string, fl []string) []string {
	filelist, err := ioutil.ReadDir(dir)

	if err != nil {
		return fl
	}
	for _, f := range filelist {
		f1 := f
		if f.IsDir() {
			fl = Scan(filepath.Join(dir, f.Name()), fl)
		} else if strings.HasSuffix(f.Name(), utils.Suffix) {
			allName := filepath.Join(dir, f1.Name())
			fl = append(fl, allName)
		}
	}
	return fl
}
