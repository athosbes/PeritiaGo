package export

import (
	"encoding/json"

	"github.com/athosbes/PeritiaGo/internal/filesystem"
	"github.com/athosbes/PeritiaGo/internal/models"
)

// ToJSON writes any data structure to a JSON file.
func ToJSON(path string, data interface{}) (models.Evidence, error) {
	file, err := filesystem.Create(path)
	if err != nil {
		return models.Evidence{}, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return models.Evidence{}, err
	}

	return models.Evidence{
		FileName: path,
		Path:     path,
	}, nil
}
