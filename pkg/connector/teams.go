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

const memberEntitlement = "member"

type teamBuilder struct {
	client litmos.Client
}

func (o *teamBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return teamResourceType
}

func teamResource(ctx context.Context, team *litmos.Team, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"code_for_bulk_import": team.TeamCodeForBulkImport,
		"parent_team_id":       team.ParentTeamId,
	}

	groupTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	resource, err := rs.NewGroupResource(
		team.Name,
		teamResourceType,
		team.Id,
		groupTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *teamBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	teams, nextPageToken, err := o.client.ListTeams(ctx, pToken)
	if err != nil {
		return nil, nextPageToken, nil, err
	}

	resources := make([]*v2.Resource, 0, len(teams))
	for _, team := range teams {
		resource, err := teamResource(ctx, &team, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nextPageToken, nil, nil
}

func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	assignmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(fmt.Sprintf("Team %s %s", resource.DisplayName, memberEntitlement)),
		entitlement.WithDescription(fmt.Sprintf("Member of team %s in Litmos", resource.DisplayName)),
	}

	rv = append(rv, entitlement.NewAssignmentEntitlement(
		resource,
		memberEntitlement,
		assignmentOptions...,
	))
	return rv, "", nil, nil
}

func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {

	users, nextPageToken, err := o.client.ListTeamUsers(ctx, pToken, resource.Id.Resource)
	if err != nil {
		return nil, nextPageToken, nil, err
	}

	rv := make([]*v2.Grant, 0, len(users))
	for _, user := range users {
		u, err := userResource(ctx, &user, nil)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(
			rv,
			grant.NewGrant(
				resource,
				memberEntitlement,
				u.Id,
			),
		)
	}

	return rv, nextPageToken, nil, nil
}

func newTeamBuilder(client litmos.Client) *teamBuilder {
	return &teamBuilder{
		client: client,
	}
}
