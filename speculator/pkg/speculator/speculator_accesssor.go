package speculator

import (
	"github.com/5gsec/sentryflow/speculator/pkg/apispec"
)

type SpeculatorsAccessor interface {
	DiffTelemetry(speculatorID uint, telemetry *apispec.Telemetry, diffSource apispec.SpecSource) (*apispec.APIDiff, error)
	HasApprovedSpec(speculatorID uint, specKey apispec.SpecKey) bool
	HasProvidedSpec(speculatorID uint, specKey apispec.SpecKey) bool
	GetProvidedSpecVersion(speculatorID uint, specKey apispec.SpecKey) apispec.OASVersion
	GetApprovedSpecVersion(speculatorID uint, specKey apispec.SpecKey) apispec.OASVersion
}

func NewSpeculatorAccessor(speculators *Repository) SpeculatorsAccessor {
	return &Impl{speculators: speculators}
}

type Impl struct {
	speculators *Repository
}

func (s *Impl) DiffTelemetry(speculatorID uint, telemetry *apispec.Telemetry, diffSource apispec.SpecSource) (*apispec.APIDiff, error) {
	//nolint: wrapcheck
	return s.speculators.Get(speculatorID).DiffTelemetry(telemetry, diffSource)
}

func (s *Impl) HasApprovedSpec(speculatorID uint, specKey apispec.SpecKey) bool {
	return s.speculators.Get(speculatorID).HasApprovedSpec(specKey)
}

func (s *Impl) HasProvidedSpec(speculatorID uint, specKey apispec.SpecKey) bool {
	return s.speculators.Get(speculatorID).HasProvidedSpec(specKey)
}

func (s *Impl) GetProvidedSpecVersion(speculatorID uint, specKey apispec.SpecKey) apispec.OASVersion {
	return s.speculators.Get(speculatorID).GetProvidedSpecVersion(specKey)
}

func (s *Impl) GetApprovedSpecVersion(speculatorID uint, specKey apispec.SpecKey) apispec.OASVersion {
	return s.speculators.Get(speculatorID).GetApprovedSpecVersion(specKey)
}
