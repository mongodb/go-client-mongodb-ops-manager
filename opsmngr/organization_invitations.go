package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const invitationBasePath = orgsBasePath + "/%s/invites"

// Invitations gets all unaccepted invitations to the specified Ops Manager organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/invitations/organizations/get-all-invitations/
func (s *OrganizationsServiceOp) Invitations(ctx context.Context, orgID string, opts *atlas.InvitationOptions) ([]*atlas.Invitation, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(invitationBasePath, orgID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*atlas.Invitation
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Invitation gets details for one unaccepted invitation to the specified Ops Manager organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/invitations/organizations/get-one-invitation/
func (s *OrganizationsServiceOp) Invitation(ctx context.Context, orgID, invitationID string) (*atlas.Invitation, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	if invitationID == "" {
		return nil, nil, atlas.NewArgError("invitationID", "must be set")
	}

	basePath := fmt.Sprintf(invitationBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, invitationID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Invitation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// InviteUser invites one user to the Ops Manager organization that you specify.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/invitations/organizations/create-one-invitation/
func (s *OrganizationsServiceOp) InviteUser(ctx context.Context, invitation *atlas.Invitation) (*atlas.Invitation, *Response, error) {
	if invitation.OrgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(invitationBasePath, invitation.OrgID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, invitation)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Invitation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// UpdateInvitation updates one pending invitation to the Ops Manager organization that you specify.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/invitations/organizations/update-one-invitation/
func (s *OrganizationsServiceOp) UpdateInvitation(ctx context.Context, invitation *atlas.Invitation) (*atlas.Invitation, *Response, error) {
	if invitation.OrgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	return s.updateInvitation(ctx, invitation)
}

// UpdateInvitationByID updates one invitation to the Ops Manager organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/invitations/organizations/update-one-invitation-by-id/
func (s *OrganizationsServiceOp) UpdateInvitationByID(ctx context.Context, invitationID string, invitation *atlas.Invitation) (*atlas.Invitation, *Response, error) {
	if invitation.OrgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	if invitationID == "" {
		return nil, nil, atlas.NewArgError("invitationID", "must be set")
	}

	invitation.ID = invitationID

	return s.updateInvitation(ctx, invitation)
}

// DeleteInvitation deletes one unaccepted invitation to the specified Ops Manager organization. You can't delete an invitation that a user has accepted.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/invitations/organizations/delete-one-invitation/
func (s *OrganizationsServiceOp) DeleteInvitation(ctx context.Context, orgID, invitationID string) (*Response, error) {
	if orgID == "" {
		return nil, atlas.NewArgError("orgID", "must be set")
	}

	if invitationID == "" {
		return nil, atlas.NewArgError("invitationID", "must be set")
	}

	basePath := fmt.Sprintf(invitationBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, invitationID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

func (s *OrganizationsServiceOp) updateInvitation(ctx context.Context, invitation *atlas.Invitation) (*atlas.Invitation, *Response, error) {
	path := fmt.Sprintf(invitationBasePath, invitation.OrgID)

	if invitation.ID != "" {
		path = fmt.Sprintf("%s/%s", path, invitation.ID)
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, invitation)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Invitation)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
