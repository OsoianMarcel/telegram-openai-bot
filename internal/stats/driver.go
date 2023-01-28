package stats

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// The stats driver structure.
type Driver struct {
	Stats     *Stats
	fileMutex sync.RWMutex
	filePath  string
}

// Creates a new stats driver instance.
func New(filePath string) *Driver {
	if filePath != "" {
		filePath = filepath.Clean(filePath)
	}

	return &Driver{
		Stats: &Stats{
			Users: make(map[string]User),
		},
		filePath: filePath,
	}
}

// Returns true if the stats file path has been set.
func (d *Driver) IsFileSet() bool {
	return d.filePath != ""
}

// Returns the stats file path.
func (d *Driver) GetFilePath() string {
	return d.filePath
}

// Loads the stats from the file.
func (d *Driver) LoadFromFile() error {
	if d.filePath == "" {
		return ErrNoFilePath
	}

	d.fileMutex.RLock()
	defer d.fileMutex.RUnlock()

	b, err := os.ReadFile(d.filePath)

	// If the file does not exists then stop here, otherwise return the error.
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	d.Stats.mutex.Lock()
	defer d.Stats.mutex.Unlock()

	return json.Unmarshal(b, &d.Stats)
}

// Saves the stats to the file.
func (d *Driver) WriteToFile() error {
	if d.filePath == "" {
		return ErrNoFilePath
	}

	d.Stats.mutex.RLock()
	defer d.Stats.mutex.RUnlock()
	b, err := json.MarshalIndent(d.Stats, "", "\t")
	if err != nil {
		return err
	}

	d.fileMutex.Lock()
	defer d.fileMutex.Unlock()
	if err := os.WriteFile(d.filePath, b, 0644); err != nil {
		return err
	}

	return nil
}
