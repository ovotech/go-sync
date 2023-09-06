package membership

import (
	"context"
	"testing"

	"github.com/hashicorp/go-tfe"
	gosync "github.com/ovotech/go-sync"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
}

func TestMembership_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	memberships := newMockIOrganizationMemberships(t)

	adapter := New(&tfe.Client{}, "org")
	adapter.organizationMemberships = memberships

	memberships.EXPECT().List(ctx, "org", &tfe.OrganizationMembershipListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1},
	}).Return(&tfe.OrganizationMembershipList{
		Pagination: &tfe.Pagination{
			CurrentPage: 1,
			NextPage:    2,
			TotalPages:  2,
		},
		Items: []*tfe.OrganizationMembership{
			{Email: "foo@email"},
			{Email: "bar@email"},
		},
	}, nil)

	memberships.EXPECT().List(ctx, "org", &tfe.OrganizationMembershipListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 2},
	}).Return(&tfe.OrganizationMembershipList{
		Pagination: &tfe.Pagination{
			CurrentPage: 2,
			NextPage:    2,
			TotalPages:  2,
		},
		Items: []*tfe.OrganizationMembership{
			{Email: "baz@email"},
			{Email: "quz@email"},
		},
	}, nil)

	things, err := adapter.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, things, []string{
		"foo@email",
		"bar@email",
		"baz@email",
		"quz@email",
	})
}

func TestMembership_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	memberships := newMockIOrganizationMemberships(t)

	adapter := New(&tfe.Client{}, "org")
	adapter.organizationMemberships = memberships

	memberships.EXPECT().Create(ctx, "org", tfe.OrganizationMembershipCreateOptions{
		Email: tfe.String("foo@email"),
		Type:  "organization-memberships",
	}).Return(&tfe.OrganizationMembership{
		ID:     "foo-id",
		Email:  "foo@email",
		Status: tfe.OrganizationMembershipInvited,
	}, nil)

	err := adapter.Add(ctx, []string{"foo@email"})
	assert.NoError(t, err)
}

func TestMembership_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	memberships := newMockIOrganizationMemberships(t)

	adapter := New(&tfe.Client{}, "org")
	adapter.organizationMemberships = memberships

	memberships.EXPECT().List(ctx, "org", &tfe.OrganizationMembershipListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1},
		Emails:      []string{"foo@email"},
	}).Return(&tfe.OrganizationMembershipList{
		Pagination: &tfe.Pagination{
			CurrentPage: 1,
			NextPage:    1,
			TotalPages:  1,
		},
		Items: []*tfe.OrganizationMembership{
			{Email: "foo@email", ID: "foo-id"},
		},
	}, nil)

	memberships.EXPECT().Delete(ctx, "foo-id").Return(nil)

	err := adapter.Remove(ctx, []string{"foo@email"})

	assert.NoError(t, err)
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Token:        "token",
			Organisation: "org",
		})

		assert.NoError(t, err)
		assert.IsType(t, &Membership{}, adapter)
	})

	t.Run("missing token", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			Organisation: "org",
		})

		assert.ErrorIs(t, err, gosync.ErrMissingConfig)
		assert.ErrorContains(t, err, Token)
	})

	t.Run("missing organisation", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			Token: "token",
		})

		assert.ErrorIs(t, err, gosync.ErrMissingConfig)
		assert.ErrorContains(t, err, Organisation)
	})
}
