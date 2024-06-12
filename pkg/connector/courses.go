package connector

import (
	"context"

	"github.com/conductorone/baton-litmos/pkg/litmos"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type courseBuilder struct {
	client litmos.Client
}

func (o *courseBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return courseResourceType
}

func courseResource(ctx context.Context, course *litmos.Course, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	resource, err := rs.NewResource(
		course.Name,
		courseResourceType,
		course.Id,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *courseBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	courses, nextPageToken, err := o.client.ListCourses(ctx, pToken)
	if err != nil {
		return nil, nextPageToken, nil, err
	}

	resources := make([]*v2.Resource, 0, len(courses))
	for _, course := range courses {
		resource, err := courseResource(ctx, &course, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nextPageToken, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *courseBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *courseBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newCourseBuilder(client litmos.Client) *courseBuilder {
	return &courseBuilder{
		client: client,
	}
}
