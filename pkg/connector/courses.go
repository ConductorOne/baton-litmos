package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-litmos/pkg/litmos"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const assignedEntitlement = "assigned"
const completedEntitlement = "completed"
const inProgressEntitlement = "in_progress"

type courseBuilder struct {
	client litmos.Client
}

func (o *courseBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return courseResourceType
}

func courseResource(ctx context.Context, course *litmos.Course, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"Id":                        course.Id,
		"Code":                      course.Code,
		"Name":                      course.Name,
		"Active":                    course.Active,
		"ForSale":                   course.ForSale,
		"OriginalId":                course.OriginalId,
		"Description":               course.Description,
		"EcommerceShortDescription": course.EcommerceShortDescription,
		"EcommerceLongDescription":  course.EcommerceLongDescription,
		"CourseCodeForBulkImport":   course.CourseCodeForBulkImport,
		"Price":                     course.Price,
		"AccessTillDate":            course.AccessTillDate,
		"AccessTillDays":            course.AccessTillDays,
		"CourseTeamLibrary":         course.CourseTeamLibrary,
		"CreatedBy":                 course.CreatedBy,
		"SeqId":                     course.SeqId,
	}

	groupTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	resource, err := rs.NewGroupResource(
		course.Name,
		courseResourceType,
		course.Id,
		groupTraitOptions,
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
	var rv []*v2.Entitlement
	assignedOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(fmt.Sprintf("Course %s %s", resource.DisplayName, assignedEntitlement)),
		entitlement.WithDescription(fmt.Sprintf("Assigned course %s in Litmos", resource.DisplayName)),
	}
	completedOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(fmt.Sprintf("Course %s %s", resource.DisplayName, completedEntitlement)),
		entitlement.WithDescription(fmt.Sprintf("Completed course %s in Litmos", resource.DisplayName)),
	}
	inProgressOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(fmt.Sprintf("Course %s %s", resource.DisplayName, inProgressEntitlement)),
		entitlement.WithDescription(fmt.Sprintf("In progress course %s in Litmos", resource.DisplayName)),
	}

	entitlements := []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			assignedEntitlement,
			assignedOptions...,
		),
		entitlement.NewAssignmentEntitlement(
			resource,
			completedEntitlement,
			completedOptions...,
		),
		entitlement.NewAssignmentEntitlement(
			resource,
			inProgressEntitlement,
			inProgressOptions...,
		),
	}
	rv = append(rv, entitlements...)
	return rv, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *courseBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, nextPageToken, err := o.client.ListCourseUsers(ctx, pToken, resource.Id.Resource)
	if err != nil {
		return nil, nextPageToken, nil, err
	}

	rv := make([]*v2.Grant, 0, len(users))
	for _, user := range users {
		rID, err := rs.NewResourceID(userResourceType, user.Id)
		if err != nil {
			return rv, nextPageToken, nil, err
		}

		grants := []*v2.Grant{grant.NewGrant(
			resource,
			assignedEntitlement,
			rID,
		)}
		if user.Completed {
			grants = append(grants, grant.NewGrant(
				resource,
				completedEntitlement,
				rID,
			))
		} else {
			grants = append(grants, grant.NewGrant(
				resource,
				inProgressEntitlement,
				rID,
			))
		}

		rv = append(rv, grants...)
	}

	return rv, nextPageToken, nil, nil
}

func newCourseBuilder(client litmos.Client) *courseBuilder {
	return &courseBuilder{
		client: client,
	}
}
