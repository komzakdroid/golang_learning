package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
)

type UIService struct {
	cache      *cache.Cache
	schemaPath string
}

func NewUIService() *UIService {
	schemaPath := os.Getenv("SCHEMA_BASE_PATH")
	if schemaPath == "" {
		schemaPath = "./schemas"
	}
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &UIService{cache: c, schemaPath: schemaPath}
}

func (s *UIService) GetScreenSchema(screenName, version string) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("schema_%s_%s", screenName, version)
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(map[string]interface{}), nil
	}

	filePath := filepath.Join(s.schemaPath, version, fmt.Sprintf("%s.json", screenName))
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("schema not found: %w", err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	s.cache.Set(cacheKey, schema, cache.DefaultExpiration)
	return schema, nil
}

func (s *UIService) GetAvailableScreens(version string) ([]string, error) {
	dirPath := filepath.Join(s.schemaPath, version)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	screens := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			screenName := file.Name()[:len(file.Name())-5]
			screens = append(screens, screenName)
		}
	}
	return screens, nil
}

func (s *UIService) ClearCache() {
	s.cache.Flush()
}
