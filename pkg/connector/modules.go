package connector

import (
	"context"

	"github.com/conductorone/baton-litmos/pkg/litmos"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type moduleBuilder struct {
	client litmos.Client
}

func (o *moduleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return moduleResourceType
}

func moduleResource(ctx context.Context, module *litmos.Module, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	resource, err := rs.NewResource(
		module.Name,
		moduleResourceType,
		module.Id,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *moduleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, opts rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	if parentResourceID == nil {
		return nil, nil, nil
	}
	pToken := &opts.PageToken
	modules, nextPageToken, err := o.client.ListModules(ctx, pToken, parentResourceID.Resource)
	if err != nil {
		return nil, nil, err
	}

	resources := make([]*v2.Resource, 0, len(modules))
	for _, module := range modules {
		resource, err := moduleResource(ctx, &module, parentResourceID)
		if err != nil {
			return nil, nil, err
		}
		resources = append(resources, resource)
	}

	if nextPageToken == "" {
		return resources, nil, nil
	}
	return resources, &rs.SyncOpResults{NextPageToken: nextPageToken}, nil
}

func (o *moduleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

func (o *moduleBuilder) Grants(ctx context.Context, resource *v2.Resource, opts rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

func newModuleBuilder(client litmos.Client) *moduleBuilder {
	return &moduleBuilder{
		client: client,
	}
}
