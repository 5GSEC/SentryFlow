package speculator

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"

	"github.com/5gsec/sentryflow/speculator/pkg/apispec"
	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

var logger = util.GetLogger()

type Repository struct {
	Speculators      map[uint]*apispec.Speculator
	speculatorConfig apispec.Config
	lock             *sync.RWMutex
}

func NewRepository(config apispec.Config) *Repository {
	return &Repository{
		Speculators:      map[uint]*apispec.Speculator{},
		speculatorConfig: config,
		lock:             &sync.RWMutex{},
	}
}

func (r *Repository) Get(speculatorID uint) *apispec.Speculator {
	r.lock.RLock()
	defer r.lock.RUnlock()

	speculator, exists := r.Speculators[speculatorID]
	if !exists {
		r.Speculators[speculatorID] = apispec.CreateSpeculator(r.speculatorConfig)
		return r.Speculators[speculatorID]
	}

	return speculator
}

func (r *Repository) EncodeState(filePath string) error {
	const perm = 400
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.FileMode(perm))
	if err != nil {
		return fmt.Errorf("failed to open state file: %v", err)
	}
	defer closeFile(file)

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(r)
	if err != nil {
		return fmt.Errorf("failed to encode state: %v", err)
	}

	return nil
}

func DecodeState(filePath string, config apispec.Config) (*Repository, error) {
	r := Repository{}

	const perm = 400
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.FileMode(perm))
	if err != nil {
		return nil, fmt.Errorf("failed to open file (%v): %v", filePath, err)
	}
	defer closeFile(file)

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state: %v", err)
	}

	r.speculatorConfig = config
	r.lock = &sync.RWMutex{}
	logger.Info("Speculator state was decoded")

	return &r, nil
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		logger.Errorf("failed to close file: %v", err)
	}
}
