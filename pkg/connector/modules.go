package connector

import (
	"context"

	"github.com/conductorone/baton-litmos/pkg/litmos"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
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
		userResourceType,
		module.Id,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *moduleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}
	modules, nextPageToken, err := o.client.ListModules(ctx, pToken, parentResourceID.Resource)
	if err != nil {
		return nil, nextPageToken, nil, err
	}

	resources := make([]*v2.Resource, 0, len(modules))
	for _, module := range modules {
		resource, err := moduleResource(ctx, &module, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nextPageToken, nil, nil
}

func (o *moduleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *moduleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newModuleBuilder(client litmos.Client) *moduleBuilder {
	return &moduleBuilder{
		client: client,
	}
}
