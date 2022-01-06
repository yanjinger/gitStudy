package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/bxcodec/faker/support/slice"
	"github.com/xuri/excelize/v2"
	"path/filepath"

	"golang.org/x/text/encoding"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const DestSuffix = ".xlsx"

var Suffix = ".prn"

type Sheet struct {
	Name      string
	ColumnNum int
	RowNum    int
	Rows      []*Row
}

type Row struct {
	ColumnNum int
	Cells     []*Cell
}

type Cell struct {
	Value  string
	VMerge int
}

func (c *Cell) Int64() (int64, error) {
	i, err := strconv.Atoi(c.Value)
	return int64(i), err
}

func (c *Cell) Float() (float64, error) {
	f, err := strconv.ParseFloat(c.Value, 64)
	return f, err
}

func ParseCSV(filename string) (*Sheet, error) {
	if strings.Contains(filename, "art_resource") {
		fmt.Println("")
	}
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("csv文件打开失败！")
	}

	reader := csv.NewReader(file)
	s := &Sheet{}
	rowNum := 0
	for {
		contents, err := reader.Read()
		if err == io.EOF {
			break
		}
		cells := make([]*Cell, len(contents))
		if len(cells) == 0 {
			continue
		}
		for i1, i2 := range contents {
			cells[i1] = &Cell{
				Value: i2,
			}
		}
		if strings.HasPrefix(cells[0].Value, "#") {
			continue
		}
		rowNum++
		if s.ColumnNum != 0 && s.ColumnNum != len(cells) {
			err = fmt.Errorf(fmt.Sprintf("%s line:%d cell:%d splitChar num error", filename, rowNum, len(cells)))
			return nil, err
		}
		s.ColumnNum = len(cells)
		row := &Row{
			ColumnNum: len(cells),
			Cells:     cells,
		}
		s.Rows = append(s.Rows, row)

	}
	s.RowNum = rowNum
	s.Name = filename

	//bb := []byte(s.Rows[0].Cells[0].Value)
	//fmt.Println(bb[:3])
	s.Rows[0].Cells[0].Value = "ID"
	return s, nil
}

func Excel2Csv(filename string) error {
	if strings.Contains(filename, "art_resource") {
		fmt.Println("")
	}
	dir := filepath.Dir(filename)
	file, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	sheetName := strings.ReplaceAll(filename, DestSuffix, Suffix)
	sheetName = filepath.Base(sheetName)
	sheetlist := file.GetSheetList()
	if len(sheetlist) != 1 {
		return fmt.Errorf("len(sheetList) != 1")
	}
	if !slice.Contains(sheetlist, sheetName) {
		return fmt.Errorf("not found sheet %s", sheetName)
	}
	cols, err := file.GetCols(sheetName)
	rows, err := file.GetRows(sheetName)
	rowNum := len(rows)
	colNum := len(cols)

	var allcontents [][]string
	for i := 0; i < rowNum; i++ {
		var contents []string
		for j := 0; j < colNum; j++ {
			columnName := GetColumnIndexName2(j + 1)
			axis := fmt.Sprintf("%s%d", columnName, i+1)
			content, err := file.GetCellValue(sheetName, axis)
			if err != nil {
				return err
			}
			contents = append(contents, content)
		}
		allcontents = append(allcontents, contents)
	}

	var bb []byte = []byte{239, 187, 191}
	bom := string(bb)
	allcontents[0][0] = bom + allcontents[0][0]

	newfilename := filepath.Join(dir, sheetName)
	csvFile, _ := os.Create(newfilename)
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	w.WriteAll(allcontents)

	fmt.Printf("excel2csv %s -> %s \n", filename, newfilename)
	return nil
}

func ParseTxt(filename string, splitChar string, decoder *encoding.Decoder) (*Sheet, error) {
	if strings.Contains(filename, "art_resource") {
		fmt.Println("")
	}
	allBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	allBytes, err = decoder.Bytes(allBytes)
	if err != nil {
		return nil, err
	}
	buffer1 := bytes.NewBuffer(allBytes)
	s := &Sheet{}
	rowNum := 0
	for {
		b_UTF16LE, err := buffer1.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if strings.HasPrefix(string(b_UTF16LE), "#") {
			continue
		}
		if len(b_UTF16LE) > 0 {
			rowNum++
			cells := parseLine(string(b_UTF16LE), splitChar)
			st := string(b_UTF16LE)
			_ = st
			if s.ColumnNum != 0 && s.ColumnNum != len(cells) {
				err = fmt.Errorf(fmt.Sprintf("%s line:%d cell:%d splitChar num error", filename, rowNum, len(cells)))
				return nil, err
			}
			s.ColumnNum = len(cells)
			row := &Row{
				ColumnNum: len(cells),
				Cells:     cells,
			}
			s.Rows = append(s.Rows, row)
		}
		if err == io.EOF {
			break
		}
	}
	s.RowNum = rowNum
	s.Name = filename
	return s, nil
}

func parseLine(content string, splitChar string) []*Cell {
	cells := []*Cell{}
	content = strings.TrimSuffix(content, "\r\n")
	content = strings.TrimSuffix(content, "\n")
	array := strings.Split(content, splitChar)
	for _, s := range array {
		if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
			s = strings.TrimPrefix(s, "\"")
			s = strings.TrimSuffix(s, "\"")
		}
		s = strings.TrimSpace(s)
		cell := &Cell{Value: s}
		cells = append(cells, cell)
	}
	return cells
}

func MergeSheet(main string, sheet1, sheet2 *Sheet) (s *Sheet, errlist []error) {
	// 分表合并之前的检查  PS： 为什么不在检查逻辑地方检查？ 因为那里的检查是基于和表后的数据的
	headRowNum := 4
	if sheet1.ColumnNum != sheet2.ColumnNum {
		err := fmt.Errorf("[%s][%d]!=[%s][%d] column is not equal", sheet1.Name, sheet1.ColumnNum, sheet2.Name, sheet2.ColumnNum)
		errlist = append(errlist, err)
		return
	}
	column := sheet1.ColumnNum
	for i := 0; i < headRowNum; i++ {
		for j := 0; j < column; j++ {
			value1 := sheet1.Rows[i].Cells[j].Value
			value2 := sheet2.Rows[i].Cells[j].Value
			if value1 != value2 {
				err := fmt.Errorf("[%s][%s]  row[%d]column[%d] [%s]!=[%s]  Header is not equal. row&column start with 0",
					sheet1.Name, sheet2.Name, i, j, value1, value2)
				errlist = append(errlist, err)
			}
		}
	}
	for i := headRowNum; i < sheet1.RowNum; i++ {
		for j := headRowNum; j < sheet2.RowNum; j++ {
			value1 := sheet1.Rows[i].Cells[0].Value
			value2 := sheet2.Rows[j].Cells[0].Value
			if value1 == value2 {
				err := fmt.Errorf("[%s][%s]  row[%d]row[%d] [%s]==[%s]  id repeat",
					sheet1.Name, sheet2.Name, i, j, value1, value2)
				errlist = append(errlist, err)
			}
		}
	}
	if len(errlist) > 0 {
		return
	}
	s = &Sheet{
		Name:      main,
		ColumnNum: column,
		RowNum:    sheet1.RowNum + sheet2.RowNum - headRowNum,
	}
	s.Rows = make([]*Row, 0)
	for i := 0; i < sheet1.RowNum; i++ {
		cells := make([]*Cell, 0)
		for j := 0; j < sheet1.ColumnNum; j++ {
			value := sheet1.Rows[i].Cells[j].Value
			newCell := &Cell{
				Value: value,
			}
			cells = append(cells, newCell)
		}
		row := &Row{
			ColumnNum: column,
			Cells:     cells,
		}
		s.Rows = append(s.Rows, row)
	}
	for i := headRowNum; i < sheet2.RowNum; i++ {
		cells := make([]*Cell, 0)
		for j := 0; j < sheet2.ColumnNum; j++ {
			value := sheet2.Rows[i].Cells[j].Value
			newCell := &Cell{
				Value: value,
			}
			cells = append(cells, newCell)
		}
		row := &Row{
			ColumnNum: column,
			Cells:     cells,
		}
		s.Rows = append(s.Rows, row)
	}

	return s, nil
}

func ScanEnumString(dir, enumPrefix string) (s *Sheet, err error) {

	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	enumList := []string{}
	for _, info := range list {
		if strings.HasPrefix(info.Name(), enumPrefix) && strings.HasSuffix(info.Name(), Suffix) {
			enumList = append(enumList, filepath.Join(dir, info.Name()))
		}
	}

	dataList := []*Sheet{}
	for _, name := range enumList {
		s, err := ParseCSV(name)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, s)
	}

	if len(dataList) == 0 {
		return nil, fmt.Errorf("not find enum_string.csv")
	}
	if len(dataList) == 1 {
		return dataList[0], nil
	}
	ret := dataList[0]
	var errlist []error
	for i := 1; i < len(dataList); i++ {
		ret, errlist = MergeSheet(enumPrefix, ret, dataList[i])
		if len(errlist) > 0 {
			return nil, errlist[0]
		}
	}
	return ret, nil
}

// 第0列是A
func GetColumnIndexName1(index int) string {
	return GetColumnIndexName2(index + 1)
}

// 第1列是A
func GetColumnIndexName2(index int) string {
	if index <= 0 {
		return "0"
	}
	str := ""
	var index1, indexTemp int
	indexTemp = index
	for true {
		indexTemp--
		index1 = indexTemp % 26
		indexTemp = indexTemp / 26
		str = string(rune(index1+65)) + str
		if indexTemp == 0 {
			break
		}
	}
	return str
}

func CopyFile(sour, des string) error {

	if f, _ := os.Stat(des); f != nil {
		os.Remove(des)
	}

	source, err := os.Open(sour)
	destination, err := os.Create(des)
	defer destination.Close()
	defer source.Close()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024*1024*200)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
