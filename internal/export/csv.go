package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"

	"github.com/athosbes/PeritiaGo/internal/models"
)

// ToCSV writes a slice of structs to a CSV file. Struct must have 'csv' tags.
func ToCSV(path string, data interface{}) (models.Evidence, error) {
	file, err := os.Create(path)
	if err != nil {
		return models.Evidence{}, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return models.Evidence{}, nil
	}

	if v.Len() == 0 {
		return models.Evidence{}, nil
	}

	elemType := v.Index(0).Type()
	var headers []string

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("csv")
		if tag != "" {
			headers = append(headers, tag)
		}
	}

	if err := writer.Write(headers); err != nil {
		return models.Evidence{}, err
	}

	for i := 0; i < v.Len(); i++ {
		elemVal := v.Index(i)
		var row []string
		for j := 0; j < elemVal.NumField(); j++ {
			field := elemType.Field(j)
			if field.Tag.Get("csv") != "" {
				val := elemVal.Field(j)
				// Format simple string presentation for CSV
				row = append(row, fmt.Sprintf("%v", val.Interface())) 
			}
		}
		if err := writer.Write(row); err != nil {
			return models.Evidence{}, err
		}
	}

	return models.Evidence{
		FileName: path,
		Path:     path,
	}, nil
}
