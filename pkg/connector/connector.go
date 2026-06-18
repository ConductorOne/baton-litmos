package connector

import (
	"context"
	"io"

	cfg "github.com/conductorone/baton-litmos/pkg/config"
	"github.com/conductorone/baton-litmos/pkg/litmos"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	mapset "github.com/deckarep/golang-set/v2"
)

type LitmosConnector struct {
	client        litmos.Client
	limitCourses  mapset.Set[string]
	enableModules bool
}

// ResourceSyncers returns a ResourceSyncerV2 for each resource type that should be synced from the upstream service.
func (d *LitmosConnector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncerV2 {
	rv := []connectorbuilder.ResourceSyncerV2{
		newUserBuilder(d.client),
		newTeamBuilder(d.client),
		newCourseBuilder(d.client, d.limitCourses, d.enableModules),
	}
	if d.enableModules {
		rv = append(rv, newModuleBuilder(d.client))
	}
	return rv
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *LitmosConnector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *LitmosConnector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Litmos Baton Connector",
		Description: "A Baton connector for Litmos",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *LitmosConnector) Validate(ctx context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, apiKey, source string, limitCourses []string) (*LitmosConnector, error) {
	cli, err := litmos.NewClient(ctx, apiKey, source)
	if err != nil {
		return nil, err
	}
	lc := &LitmosConnector{
		client: *cli,
	}
	if len(limitCourses) > 0 {
		lc.limitCourses = mapset.NewSet(limitCourses...)
	}
	return lc, nil
}

// NewLambdaConnector satisfies cli.NewConnector for use with config.RunConnector.
func NewLambdaConnector(ctx context.Context, ac *cfg.Litmos, _ *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	cb, err := New(ctx, ac.ApiKey, ac.Source, ac.LimitedCourses)
	if err != nil {
		return nil, nil, err
	}
	return cb, nil, nil
}
