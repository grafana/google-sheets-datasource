package v1

import (
	"context"
	"io"

	"github.com/grafana/kindsys"
)

var _ kindsys.ResourceKind = &rawSheetsKind{}

// This implements a playlist directly in golang
type rawSheetsKind struct {
	migrator kindsys.ResourceMigrator
}

func NewRawSheetsKind() kindsys.ResourceKind {
	return &rawSheetsKind{
		// migrator: newMigrator(hooks),
	}
}

func (k *rawSheetsKind) GetMachineNames() kindsys.MachineNames {
	return kindsys.MachineNames{
		Plural:   "datasources",
		Singular: "datasource",
	}
}

func (k *rawSheetsKind) GetKindInfo() kindsys.KindInfo {
	return kindsys.KindInfo{
		Group:       "googlesheets.ext.grafana.com",
		Kind:        "Datasource",
		Description: "",
		Maturity:    kindsys.MaturityMerged,
	}
}

func (k *rawSheetsKind) CurrentVersion() string {
	return "v1-0" // would be nice pick a middle one
}

func (k *rawSheetsKind) GetVersions() []kindsys.VersionInfo {
	return []kindsys.VersionInfo{
		{
			Version:         "v0-0",
			SoftwareVersion: "v6.0",
		},
	}
}

// NOTE: this files are not used to do validation, but can be used generically to describe the kind
func (k *rawSheetsKind) GetJSONSchema(version string) (string, error) {
	panic("implement")
}

type ResourceV0 = kindsys.GenericResource[DatasourceSpec, kindsys.SimpleCustomMetadata, DatasourceStatus]

func (k *rawSheetsKind) Read(reader io.Reader, strict bool) (kindsys.Resource, error) {
	panic("implement")
}

func (k *rawSheetsKind) Migrate(ctx context.Context, obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
	panic("implement")
}
