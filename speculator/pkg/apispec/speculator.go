package apispec

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type SpecKey string

type Config struct {
	OperationGeneratorConfig OperationGeneratorConfig
}

type Speculator struct {
	Specs map[SpecKey]*Spec `json:"specs,omitempty"`

	config Config
}

// nolint:gochecknoinits
func init() {
	gob.Register(json.RawMessage{})
}

func CreateSpeculator(config Config) *Speculator {
	logger.Info("Creating Speculator")
	logger.Debugf("Speculator Config %+v", config)
	return &Speculator{
		Specs:  make(map[SpecKey]*Spec),
		config: config,
	}
}

func GetSpecKey(host, port string) SpecKey {
	return SpecKey(host + ":" + port)
}

func GetHostAndPortFromSpecKey(key SpecKey) (host, port string, err error) {
	const hostAndPortLen = 2
	hostAndPort := strings.Split(string(key), ":")
	if len(hostAndPort) != hostAndPortLen {
		return "", "", fmt.Errorf("invalid key: %v", key)
	}
	host = hostAndPort[0]
	if len(host) == 0 {
		return "", "", fmt.Errorf("no host for key: %v", key)
	}
	port = hostAndPort[1]
	if len(port) == 0 {
		return "", "", fmt.Errorf("no port for key: %v", key)
	}
	return host, port, nil
}

func (s *Speculator) SuggestedReview(specKey SpecKey) (*SuggestedSpecReview, error) {
	spec, exists := s.Specs[specKey]
	if !exists {
		return nil, fmt.Errorf("spec doesn't exist for key %v", specKey)
	}

	return spec.CreateSuggestedReview(), nil
}

type AddressInfo struct {
	IP   string
	Port string
}

func GetAddressInfoFromAddress(address string) (*AddressInfo, error) {
	const addrLen = 2
	addr := strings.Split(address, ":")
	if len(addr) != addrLen {
		return nil, fmt.Errorf("invalid address: %v", addr)
	}

	return &AddressInfo{
		IP:   addr[0],
		Port: addr[1],
	}, nil
}

func (s *Speculator) InitSpec(host, port string) error {
	specKey := GetSpecKey(host, port)
	if _, exists := s.Specs[specKey]; exists {
		return fmt.Errorf("spec was already initialized using host and port: %s:%s", host, port)
	}
	s.Specs[specKey] = CreateDefaultSpec(host, port, s.config.OperationGeneratorConfig)
	return nil
}

func (s *Speculator) LearnTelemetry(telemetry *Telemetry) error {
	destInfo, err := GetAddressInfoFromAddress(telemetry.DestinationAddress)
	if err != nil {
		return fmt.Errorf("failed get destination info: %v", err)
	}
	specKey := GetSpecKey(telemetry.Request.Host, destInfo.Port)
	if _, exists := s.Specs[specKey]; !exists {
		s.Specs[specKey] = CreateDefaultSpec(telemetry.Request.Host, destInfo.Port, s.config.OperationGeneratorConfig)
	}
	spec := s.Specs[specKey]
	if err := spec.LearnTelemetry(telemetry); err != nil {
		return fmt.Errorf("failed to insert telemetry: %v. %v", telemetry, err)
	}

	return nil
}

func (s *Speculator) GetPathID(specKey SpecKey, path string, specSource SpecSource) (string, error) {
	spec, exists := s.Specs[specKey]
	if !exists {
		return "", fmt.Errorf("no spec for key %v", specKey)
	}

	pathID, err := spec.GetPathID(path, specSource)
	if err != nil {
		return "", fmt.Errorf("failed to get path id. specKey=%v, specSource=%v: %v", specKey, specSource, err)
	}

	return pathID, nil
}

func (s *Speculator) DiffTelemetry(telemetry *Telemetry, diffSource SpecSource) (*APIDiff, error) {
	destInfo, err := GetAddressInfoFromAddress(telemetry.DestinationAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination info: %v", err)
	}
	specKey := GetSpecKey(telemetry.Request.Host, destInfo.Port)
	spec, exists := s.Specs[specKey]
	if !exists {
		return nil, fmt.Errorf("no spec for key %v", specKey)
	}

	apiDiff, err := spec.DiffTelemetry(telemetry, diffSource)
	if err != nil {
		return nil, fmt.Errorf("failed to run DiffTelemetry: %v", err)
	}

	return apiDiff, nil
}

func (s *Speculator) HasApprovedSpec(key SpecKey) bool {
	spec, exists := s.Specs[key]
	if !exists {
		return false
	}

	return spec.HasApprovedSpec()
}

func (s *Speculator) GetApprovedSpecVersion(key SpecKey) OASVersion {
	spec, exists := s.Specs[key]
	if !exists {
		return Unknown
	}

	return spec.ApprovedSpec.GetSpecVersion()
}

func (s *Speculator) LoadProvidedSpec(key SpecKey, providedSpec []byte, pathToPathID map[string]string) error {
	spec, exists := s.Specs[key]
	if !exists {
		return fmt.Errorf("no spec found with key: %v", key)
	}

	if err := spec.LoadProvidedSpec(providedSpec, pathToPathID); err != nil {
		return fmt.Errorf("failed to load provided spec: %w", err)
	}

	return nil
}

func (s *Speculator) UnsetProvidedSpec(key SpecKey) error {
	spec, exists := s.Specs[key]
	if !exists {
		return fmt.Errorf("no spec found with key: %v", key)
	}
	spec.UnsetProvidedSpec()
	return nil
}

func (s *Speculator) UnsetApprovedSpec(key SpecKey) error {
	spec, exists := s.Specs[key]
	if !exists {
		return fmt.Errorf("no spec found with key: %v", key)
	}
	spec.UnsetApprovedSpec()
	return nil
}

func (s *Speculator) HasProvidedSpec(key SpecKey) bool {
	spec, exists := s.Specs[key]
	if !exists {
		return false
	}

	return spec.HasProvidedSpec()
}

func (s *Speculator) GetProvidedSpecVersion(key SpecKey) OASVersion {
	spec, exists := s.Specs[key]
	if !exists {
		return Unknown
	}

	return spec.ProvidedSpec.GetSpecVersion()
}

func (s *Speculator) DumpSpecs() {
	logger.Infof("Generating Open API Specs...\n")
	for specKey, spec := range s.Specs {
		approvedYaml, err := spec.GenerateOASYaml(OASv3)
		if err != nil {
			logger.Errorf("failed to generate OAS yaml for %v.: %v", specKey, err)
			continue
		}
		logger.Infof("Spec for %s:\n%s\n\n", specKey, approvedYaml)
	}
}

func (s *Speculator) ApplyApprovedReview(specKey SpecKey, approvedReview *ApprovedSpecReview, version OASVersion) error {
	if err := s.Specs[specKey].ApplyApprovedReview(approvedReview, version); err != nil {
		return fmt.Errorf("failed to apply approved review for spec: %v. %w", specKey, err)
	}
	return nil
}

func (s *Speculator) EncodeState(filePath string) error {
	file, err := openFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open state file: %v", err)
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(s)
	if err != nil {
		return fmt.Errorf("failed to encode state: %v", err)
	}
	closeSpeculatorStateFile(file)

	return nil
}

func DecodeSpeculatorState(filePath string, config Config) (*Speculator, error) {
	r := &Speculator{}
	file, err := openFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file (%v): %v", filePath, err)
	}
	defer closeSpeculatorStateFile(file)

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state: %v", err)
	}

	r.config = config

	logger.Info("Speculator state was decoded")
	logger.Debugf("Speculator Config %+v", config)

	return r, nil
}

func openFile(filePath string) (*os.File, error) {
	const perm = 400
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.FileMode(perm))
	if err != nil {
		return nil, fmt.Errorf("failed to open file (%v) for writing: %v", filePath, err)
	}

	return file, nil
}

func closeSpeculatorStateFile(f *os.File) {
	if err := f.Close(); err != nil {
		logger.Errorf("Failed to close file: %v", err)
	}
}
