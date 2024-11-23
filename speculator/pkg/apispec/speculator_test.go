package apispec

import (
	"os"
	"testing"

	"github.com/gofrs/uuid"
)

func TestGetHostAndPortFromSpecKey(t *testing.T) {
	type args struct {
		key SpecKey
	}
	tests := []struct {
		name     string
		args     args
		wantHost string
		wantPort string
		wantErr  bool
	}{
		{
			name: "invalid key",
			args: args{
				key: "invalid",
			},
			wantHost: "",
			wantPort: "",
			wantErr:  true,
		},
		{
			name: "invalid:key:invalid",
			args: args{
				key: "invalid",
			},
			wantHost: "",
			wantPort: "",
			wantErr:  true,
		},
		{
			name: "invalid key - no host",
			args: args{
				key: ":8080",
			},
			wantHost: "",
			wantPort: "",
			wantErr:  true,
		},
		{
			name: "invalid key - no port",
			args: args{
				key: "host:",
			},
			wantHost: "",
			wantPort: "",
			wantErr:  true,
		},
		{
			name: "valid key",
			args: args{
				key: "host:8080",
			},
			wantHost: "host",
			wantPort: "8080",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPort, err := GetHostAndPortFromSpecKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostAndPortFromSpecKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHost != tt.wantHost {
				t.Errorf("GetHostAndPortFromSpecKey() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotPort != tt.wantPort {
				t.Errorf("GetHostAndPortFromSpecKey() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func TestDecodeState(t *testing.T) {
	testSpec := GetSpecKey("host", "port")
	uid, _ := uuid.NewV4()
	uidStr := uid.String()
	testStatePath := "/tmp/" + uidStr + "state.gob"
	defer func() {
		_ = os.Remove(testStatePath)
	}()

	speculatorConfig := Config{
		OperationGeneratorConfig: OperationGeneratorConfig{
			ResponseHeadersToIgnore: []string{"before"},
		},
	}
	speculator := CreateSpeculator(speculatorConfig)
	speculator.Specs[testSpec] = CreateDefaultSpec("host", "port", speculator.config.OperationGeneratorConfig)

	if err := speculator.EncodeState(testStatePath); err != nil {
		t.Errorf("EncodeSpeculatorState() error = %v", err)
		return
	}

	newSpeculatorConfig := Config{
		OperationGeneratorConfig: OperationGeneratorConfig{
			ResponseHeadersToIgnore: []string{"after"},
		},
	}
	got, err := DecodeSpeculatorState(testStatePath, newSpeculatorConfig)
	if err != nil {
		t.Errorf("DecodeSpeculatorState() error = %v", err)
		return
	}

	// OpGenerator on the decoded state should hold the previous OperationGeneratorConfig
	responseHeadersToIgnore := got.Specs[testSpec].OpGenerator.ResponseHeadersToIgnore
	if _, ok := responseHeadersToIgnore["before"]; !ok {
		t.Errorf("ResponseHeadersToIgnore not as expected = %+v", responseHeadersToIgnore)
		return
	}
}
