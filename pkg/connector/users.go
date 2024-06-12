package connector

import (
	"context"

	"github.com/conductorone/baton-litmos/pkg/litmos"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client litmos.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func userResource(ctx context.Context, user *litmos.User, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"user_id":      user.Id,
		"brand":        user.Brand,
		"access_level": user.AccessLevel,
	}

	status := rs.WithStatus(v2.UserTrait_Status_STATUS_ENABLED)
	if !user.Active {
		status = rs.WithStatus(v2.UserTrait_Status_STATUS_DISABLED)
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithUserLogin(user.UserName),
		rs.WithEmail(user.Email, true),
		status,
	}

	resource, err := rs.NewUserResource(
		user.UserName,
		userResourceType,
		user.Id,
		userTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	users, err := o.client.ListUsers(ctx, pToken)
	if err != nil {
		return nil, "", nil, err
	}

	resources := make([]*v2.Resource, 0, len(users))
	for _, user := range users {
		resource, err := userResource(ctx, &user, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, resource)
	}
	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client litmos.Client) *userBuilder {
	return &userBuilder{
		client: client,
	}
}
