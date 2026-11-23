package repository

import (
	"cmp"
	"context"
	"database/sql"
	"slices"
	"strings"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/samber/lo"
	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

// MigrateApplicationGroups puts every application without a group into one.
// Applications whose websites are nested, directly or through each other, share a group that users see,
// named after the application at the shortest path and owned by all owners of the applications.
// Every other application gets an implicit group with its own owners.
// Owners are read from application_owners, the owners before groups.
//
// It runs in one transaction, and applications that already have a group are left as they are,
// so running it again does nothing but migrate applications created without a group since.
func MigrateApplicationGroups(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return oops.Wrapf(err, "starting transaction")
	}
	defer tx.Rollback()

	mApps, err := models.Applications(
		models.ApplicationWhere.GroupID.IsNull(),
		qm.Load(models.ApplicationRels.Websites),
		qm.Load(models.ApplicationRels.Users),
		qm.OrderBy(models.ApplicationColumns.CreatedAt+", "+models.ApplicationColumns.ID),
		qm.For("UPDATE"),
	).All(ctx, tx)
	if err != nil {
		return oops.Wrapf(err, "getting applications without group")
	}

	apps := make([]*domain.Application, 0, len(mApps))
	owners := make(map[string][]string, len(mApps))
	for _, ma := range mApps {
		apps = append(apps, &domain.Application{
			ID:       ma.ID,
			Name:     ma.Name,
			Websites: ds.Map(ma.R.Websites, repoconvert.ToDomainWebsite),
		})
		owners[ma.ID] = ds.Map(ma.R.Users, func(u *models.User) string { return u.ID })
	}

	for _, set := range domain.GroupByNestedWebsites(apps) {
		group := groupForMigration(set, owners)
		if err = insertApplicationGroup(ctx, tx, group); err != nil {
			return err
		}
		_, err = models.Applications(
			models.ApplicationWhere.ID.IN(ds.Map(set, func(app *domain.Application) string { return app.ID })),
		).UpdateAll(ctx, tx, models.M{models.ApplicationColumns.GroupID: null.StringFrom(group.ID)})
		if err != nil {
			return oops.Wrapf(err, "setting application group")
		}
	}

	if err = tx.Commit(); err != nil {
		return oops.Wrapf(err, "committing transaction")
	}
	return nil
}

func groupForMigration(apps []*domain.Application, owners map[string][]string) *domain.ApplicationGroup {
	ownerIDs := lo.Uniq(lo.FlatMap(apps, func(app *domain.Application, _ int) []string { return owners[app.ID] }))
	if len(apps) == 1 {
		return domain.NewImplicitApplicationGroup(ownerIDs)
	}
	parent := slices.MinFunc(apps, func(a, b *domain.Application) int {
		return cmp.Or(cmp.Compare(shortestPathDepth(a), shortestPathDepth(b)), strings.Compare(a.Name, b.Name))
	})
	return domain.NewApplicationGroup(parent.Name, ownerIDs)
}

func shortestPathDepth(app *domain.Application) int {
	return lo.Min(ds.Map(app.Websites, func(w *domain.Website) int {
		return strings.Count(strings.TrimSuffix(w.PathPrefix, "/"), "/")
	}))
}
