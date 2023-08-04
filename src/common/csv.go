package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Implementation inspired from - https://stackoverflow.com/questions/24999079/reading-csv-file-in-go

type CSVReader struct {
	fileReader *os.File
	reader     *csv.Reader
	Headers    map[string]int
	Row        []string
}

func (s *CSVReader) Rows() bool {
	record, err := s.reader.Read()
	s.Row = record
	return err == nil
}

func (s *CSVReader) Text(header string) string {
	return s.Row[s.Headers[header]]
}

func (s *CSVReader) Bool(header string) bool {
	value, err := strconv.ParseBool(s.Text(header))
	if err != nil {
		return false
	}
	return value
}

func (s *CSVReader) Int(header string) int64 {
	value, err := strconv.ParseInt(s.Text(header), 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func (s *CSVReader) Close() error {
	return s.fileReader.Close()
}

func ReadCSVFile(filePath string) (*CSVReader, error) {
	fileReader, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file '%s' : %s", filePath, err)
	}
	reader := csv.NewReader(fileReader)
	records, err2 := reader.Read()
	if err2 != nil {
		return nil, fmt.Errorf("failed reading file '%s' : %s", filePath, err)
	}
	headers := map[string]int{}
	for i, key := range records {
		headers[key] = i
	}
	output := CSVReader{fileReader: fileReader, reader: reader, Headers: headers}
	return &output, nil
}
