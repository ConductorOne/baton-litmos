package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-litmos/pkg/litmos"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
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
		"CodeForBulkImport": team.TeamCodeForBulkImport,
		"Id":                team.Id,
		"Name":              team.Name,
		"ParentTeamId":      team.ParentTeamId,
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

func (o *teamBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, opts rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	pToken := &opts.PageToken
	teams, nextPageToken, err := o.client.ListTeams(ctx, pToken)
	if err != nil {
		return nil, nil, err
	}

	resources := make([]*v2.Resource, 0, len(teams))
	for _, team := range teams {
		resource, err := teamResource(ctx, &team, parentResourceID)
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

func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
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
	return rv, nil, nil
}

func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, opts rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	pToken := &opts.PageToken
	users, nextPageToken, err := o.client.ListTeamUsers(ctx, pToken, resource.Id.Resource)
	if err != nil {
		return nil, nil, err
	}

	rv := make([]*v2.Grant, 0, len(users))
	for _, user := range users {
		u, err := userResource(ctx, &user, nil)
		if err != nil {
			return nil, nil, err
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

	if nextPageToken == "" {
		return rv, nil, nil
	}
	return rv, &rs.SyncOpResults{NextPageToken: nextPageToken}, nil
}

func newTeamBuilder(client litmos.Client) *teamBuilder {
	return &teamBuilder{
		client: client,
	}
}
