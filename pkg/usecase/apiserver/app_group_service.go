package apiserver

import (
	"context"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

// getExplicitApplicationGroup returns the group a user created.
// Implicit groups are reported as not found, since users do not see them.
func (s *Service) getExplicitApplicationGroup(ctx context.Context, id string) (*domain.ApplicationGroup, error) {
	group, err := s.groupRepo.GetApplicationGroup(ctx, id)
	if err != nil {
		return nil, err
	}
	if group.Implicit {
		return nil, domain.NewError(domain.ErrorTypeNotFound, "application group not found", nil)
	}
	return group, nil
}

func (s *Service) deleteImplicitApplicationGroup(ctx context.Context, id string) error {
	group, err := s.groupRepo.GetApplicationGroup(ctx, id)
	if err != nil {
		return err
	}
	if !group.Implicit {
		return nil
	}
	return s.groupRepo.DeleteApplicationGroup(ctx, id)
}

func (s *Service) CreateApplicationGroup(ctx context.Context, name string) (*domain.ApplicationGroup, error) {
	group := domain.NewApplicationGroup(name, []string{web.GetUser(ctx).ID})
	if err := group.Validate(); err != nil {
		return nil, domain.NewError(domain.ErrorTypeBadRequest, "invalid application group", err)
	}
	if err := s.groupRepo.CreateApplicationGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

type GetAppGroupScope int

const (
	GetAppGroupScopeMine GetAppGroupScope = iota
	GetAppGroupScopeAll
)

// GetApplicationGroups returns the groups users created. Implicit groups are not included.
func (s *Service) GetApplicationGroups(ctx context.Context, scope GetAppGroupScope) ([]*domain.ApplicationGroup, error) {
	cond := domain.GetApplicationGroupCondition{Implicit: optional.From(false)}
	switch scope {
	case GetAppGroupScopeMine:
		cond.UserID = optional.From(web.GetUser(ctx).ID)
	case GetAppGroupScopeAll:
	default:
		return nil, oops.New("unexpected scope type")
	}
	return s.groupRepo.GetApplicationGroups(ctx, cond)
}

func (s *Service) GetApplicationGroup(ctx context.Context, id string) (*domain.ApplicationGroup, error) {
	return s.getExplicitApplicationGroup(ctx, id)
}

func (s *Service) UpdateApplicationGroup(ctx context.Context, id string, args *domain.UpdateApplicationGroupArgs) error {
	group, err := s.getExplicitApplicationGroup(ctx, id)
	if err != nil {
		return err
	}
	if !group.IsOwner(web.GetUser(ctx)) {
		return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for this group", nil)
	}
	group.Apply(args)
	if err = group.Validate(); err != nil {
		return domain.NewError(domain.ErrorTypeBadRequest, "invalid application group", err)
	}
	return s.groupRepo.UpdateApplicationGroup(ctx, id, args)
}

// DeleteApplicationGroup deletes the group. Applications in the group have to be moved or deleted first.
func (s *Service) DeleteApplicationGroup(ctx context.Context, id string) error {
	group, err := s.getExplicitApplicationGroup(ctx, id)
	if err != nil {
		return err
	}
	if !group.IsOwner(web.GetUser(ctx)) {
		return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for this group", nil)
	}
	apps, err := s.appRepo.GetApplications(ctx, domain.GetApplicationCondition{GroupID: optional.From(id)})
	if err != nil {
		return err
	}
	if len(apps) > 0 {
		return domain.NewError(domain.ErrorTypeBadRequest, "all applications in the group must be moved or deleted first", nil)
	}
	return s.groupRepo.DeleteApplicationGroup(ctx, id)
}

// UpdateStandaloneApplicationOwners sets the owners of an application that is not in any group a user created.
// The owners of an application in such a group are the owners of the group.
func (s *Service) UpdateStandaloneApplicationOwners(ctx context.Context, appID string, ownerIDs []string) error {
	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return err
	}
	group, err := s.groupRepo.GetApplicationGroup(ctx, app.GroupID)
	if err != nil {
		return err
	}
	if !group.IsOwner(web.GetUser(ctx)) {
		return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for this application", nil)
	}
	if !group.Implicit {
		return domain.NewError(domain.ErrorTypeBadRequest, "owners of an application in a group are managed by the group", nil)
	}
	args := &domain.UpdateApplicationGroupArgs{OwnerIDs: optional.From(ownerIDs)}
	group.Apply(args)
	if err = group.Validate(); err != nil {
		return domain.NewError(domain.ErrorTypeBadRequest, "invalid owners", err)
	}
	return s.groupRepo.UpdateApplicationGroup(ctx, group.ID, args)
}

// MoveApplication moves the application into the group of groupID,
// or makes it standalone when groupID is not given.
// The user has to be an owner of both the current group and the destination.
// Websites of the application are checked again for conflicts in the destination.
func (s *Service) MoveApplication(ctx context.Context, appID string, groupID optional.Of[string]) error {
	user := web.GetUser(ctx)
	app, err := s.appRepo.GetApplication(ctx, appID)
	if err != nil {
		return err
	}
	src, err := s.groupRepo.GetApplicationGroup(ctx, app.GroupID)
	if err != nil {
		return err
	}
	if !src.IsOwner(user) {
		return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for this application", nil)
	}

	var dest *domain.ApplicationGroup
	createDest := false
	if groupID.Valid {
		dest, err = s.getExplicitApplicationGroup(ctx, groupID.V)
		if err != nil {
			return err
		}
		if !dest.IsOwner(user) {
			return domain.NewError(domain.ErrorTypeForbidden, "you do not have permission for the destination group", nil)
		}
	} else {
		if src.Implicit {
			return nil
		}
		dest = domain.NewImplicitApplicationGroup(lo.Uniq(src.OwnerIDs))
		createDest = true
	}
	if dest.ID == src.ID {
		return nil
	}

	app.GroupID = dest.ID
	if err = s.validateApp(ctx, app); err != nil {
		return err
	}

	// Each step leaves at most an implicit group without apps on failure, which no one sees.
	if createDest {
		if err = s.groupRepo.CreateApplicationGroup(ctx, dest); err != nil {
			return err
		}
	}
	if err = s.appRepo.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{GroupID: optional.From(dest.ID)}); err != nil {
		return oops.Wrapf(err, "updating application group")
	}
	if src.Implicit {
		if err = s.groupRepo.DeleteApplicationGroup(ctx, src.ID); err != nil {
			return err
		}
	}
	return nil
}
