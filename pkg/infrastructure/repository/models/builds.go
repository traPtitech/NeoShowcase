// Code generated by SQLBoiler 4.14.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Build is an object representing the database table.
type Build struct { // ビルドID
	ID string `boil:"id" json:"id" toml:"id" yaml:"id"`
	// コミットハッシュ
	Commit string `boil:"commit" json:"commit" toml:"commit" yaml:"commit"`
	// ビルドの状態
	Status string `boil:"status" json:"status" toml:"status" yaml:"status"`
	// ビルド追加日時
	QueuedAt time.Time `boil:"queued_at" json:"queued_at" toml:"queued_at" yaml:"queued_at"`
	// ビルド開始日時
	StartedAt null.Time `boil:"started_at" json:"started_at,omitempty" toml:"started_at" yaml:"started_at,omitempty"`
	// ビルド更新日時
	UpdatedAt null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	// ビルド終了日時
	FinishedAt null.Time `boil:"finished_at" json:"finished_at,omitempty" toml:"finished_at" yaml:"finished_at,omitempty"`
	// 再ビルド可能フラグ
	Retriable bool `boil:"retriable" json:"retriable" toml:"retriable" yaml:"retriable"`
	// アプリケーションID
	ApplicationID string `boil:"application_id" json:"application_id" toml:"application_id" yaml:"application_id"`

	R *buildR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L buildL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var BuildColumns = struct {
	ID            string
	Commit        string
	Status        string
	QueuedAt      string
	StartedAt     string
	UpdatedAt     string
	FinishedAt    string
	Retriable     string
	ApplicationID string
}{
	ID:            "id",
	Commit:        "commit",
	Status:        "status",
	QueuedAt:      "queued_at",
	StartedAt:     "started_at",
	UpdatedAt:     "updated_at",
	FinishedAt:    "finished_at",
	Retriable:     "retriable",
	ApplicationID: "application_id",
}

var BuildTableColumns = struct {
	ID            string
	Commit        string
	Status        string
	QueuedAt      string
	StartedAt     string
	UpdatedAt     string
	FinishedAt    string
	Retriable     string
	ApplicationID string
}{
	ID:            "builds.id",
	Commit:        "builds.commit",
	Status:        "builds.status",
	QueuedAt:      "builds.queued_at",
	StartedAt:     "builds.started_at",
	UpdatedAt:     "builds.updated_at",
	FinishedAt:    "builds.finished_at",
	Retriable:     "builds.retriable",
	ApplicationID: "builds.application_id",
}

// Generated where

var BuildWhere = struct {
	ID            whereHelperstring
	Commit        whereHelperstring
	Status        whereHelperstring
	QueuedAt      whereHelpertime_Time
	StartedAt     whereHelpernull_Time
	UpdatedAt     whereHelpernull_Time
	FinishedAt    whereHelpernull_Time
	Retriable     whereHelperbool
	ApplicationID whereHelperstring
}{
	ID:            whereHelperstring{field: "`builds`.`id`"},
	Commit:        whereHelperstring{field: "`builds`.`commit`"},
	Status:        whereHelperstring{field: "`builds`.`status`"},
	QueuedAt:      whereHelpertime_Time{field: "`builds`.`queued_at`"},
	StartedAt:     whereHelpernull_Time{field: "`builds`.`started_at`"},
	UpdatedAt:     whereHelpernull_Time{field: "`builds`.`updated_at`"},
	FinishedAt:    whereHelpernull_Time{field: "`builds`.`finished_at`"},
	Retriable:     whereHelperbool{field: "`builds`.`retriable`"},
	ApplicationID: whereHelperstring{field: "`builds`.`application_id`"},
}

// BuildRels is where relationship names are stored.
var BuildRels = struct {
	Application string
	Artifact    string
}{
	Application: "Application",
	Artifact:    "Artifact",
}

// buildR is where relationships are stored.
type buildR struct {
	Application *Application `boil:"Application" json:"Application" toml:"Application" yaml:"Application"`
	Artifact    *Artifact    `boil:"Artifact" json:"Artifact" toml:"Artifact" yaml:"Artifact"`
}

// NewStruct creates a new relationship struct
func (*buildR) NewStruct() *buildR {
	return &buildR{}
}

func (r *buildR) GetApplication() *Application {
	if r == nil {
		return nil
	}
	return r.Application
}

func (r *buildR) GetArtifact() *Artifact {
	if r == nil {
		return nil
	}
	return r.Artifact
}

// buildL is where Load methods for each relationship are stored.
type buildL struct{}

var (
	buildAllColumns            = []string{"id", "commit", "status", "queued_at", "started_at", "updated_at", "finished_at", "retriable", "application_id"}
	buildColumnsWithoutDefault = []string{"id", "commit", "status", "queued_at", "started_at", "updated_at", "finished_at", "retriable", "application_id"}
	buildColumnsWithDefault    = []string{}
	buildPrimaryKeyColumns     = []string{"id"}
	buildGeneratedColumns      = []string{}
)

type (
	// BuildSlice is an alias for a slice of pointers to Build.
	// This should almost always be used instead of []Build.
	BuildSlice []*Build
	// BuildHook is the signature for custom Build hook methods
	BuildHook func(context.Context, boil.ContextExecutor, *Build) error

	buildQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	buildType                 = reflect.TypeOf(&Build{})
	buildMapping              = queries.MakeStructMapping(buildType)
	buildPrimaryKeyMapping, _ = queries.BindMapping(buildType, buildMapping, buildPrimaryKeyColumns)
	buildInsertCacheMut       sync.RWMutex
	buildInsertCache          = make(map[string]insertCache)
	buildUpdateCacheMut       sync.RWMutex
	buildUpdateCache          = make(map[string]updateCache)
	buildUpsertCacheMut       sync.RWMutex
	buildUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var buildAfterSelectHooks []BuildHook

var buildBeforeInsertHooks []BuildHook
var buildAfterInsertHooks []BuildHook

var buildBeforeUpdateHooks []BuildHook
var buildAfterUpdateHooks []BuildHook

var buildBeforeDeleteHooks []BuildHook
var buildAfterDeleteHooks []BuildHook

var buildBeforeUpsertHooks []BuildHook
var buildAfterUpsertHooks []BuildHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Build) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Build) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Build) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Build) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Build) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Build) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Build) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Build) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Build) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range buildAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddBuildHook registers your hook function for all future operations.
func AddBuildHook(hookPoint boil.HookPoint, buildHook BuildHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		buildAfterSelectHooks = append(buildAfterSelectHooks, buildHook)
	case boil.BeforeInsertHook:
		buildBeforeInsertHooks = append(buildBeforeInsertHooks, buildHook)
	case boil.AfterInsertHook:
		buildAfterInsertHooks = append(buildAfterInsertHooks, buildHook)
	case boil.BeforeUpdateHook:
		buildBeforeUpdateHooks = append(buildBeforeUpdateHooks, buildHook)
	case boil.AfterUpdateHook:
		buildAfterUpdateHooks = append(buildAfterUpdateHooks, buildHook)
	case boil.BeforeDeleteHook:
		buildBeforeDeleteHooks = append(buildBeforeDeleteHooks, buildHook)
	case boil.AfterDeleteHook:
		buildAfterDeleteHooks = append(buildAfterDeleteHooks, buildHook)
	case boil.BeforeUpsertHook:
		buildBeforeUpsertHooks = append(buildBeforeUpsertHooks, buildHook)
	case boil.AfterUpsertHook:
		buildAfterUpsertHooks = append(buildAfterUpsertHooks, buildHook)
	}
}

// One returns a single build record from the query.
func (q buildQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Build, error) {
	o := &Build{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for builds")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Build records from the query.
func (q buildQuery) All(ctx context.Context, exec boil.ContextExecutor) (BuildSlice, error) {
	var o []*Build

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Build slice")
	}

	if len(buildAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Build records in the query.
func (q buildQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count builds rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q buildQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if builds exists")
	}

	return count > 0, nil
}

// Application pointed to by the foreign key.
func (o *Build) Application(mods ...qm.QueryMod) applicationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`id` = ?", o.ApplicationID),
	}

	queryMods = append(queryMods, mods...)

	return Applications(queryMods...)
}

// Artifact pointed to by the foreign key.
func (o *Build) Artifact(mods ...qm.QueryMod) artifactQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`build_id` = ?", o.ID),
	}

	queryMods = append(queryMods, mods...)

	return Artifacts(queryMods...)
}

// LoadApplication allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (buildL) LoadApplication(ctx context.Context, e boil.ContextExecutor, singular bool, maybeBuild interface{}, mods queries.Applicator) error {
	var slice []*Build
	var object *Build

	if singular {
		var ok bool
		object, ok = maybeBuild.(*Build)
		if !ok {
			object = new(Build)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeBuild)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeBuild))
			}
		}
	} else {
		s, ok := maybeBuild.(*[]*Build)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeBuild)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeBuild))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &buildR{}
		}
		args = append(args, object.ApplicationID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &buildR{}
			}

			for _, a := range args {
				if a == obj.ApplicationID {
					continue Outer
				}
			}

			args = append(args, obj.ApplicationID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`applications`),
		qm.WhereIn(`applications.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Application")
	}

	var resultSlice []*Application
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Application")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for applications")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for applications")
	}

	if len(applicationAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Application = foreign
		if foreign.R == nil {
			foreign.R = &applicationR{}
		}
		foreign.R.Builds = append(foreign.R.Builds, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ApplicationID == foreign.ID {
				local.R.Application = foreign
				if foreign.R == nil {
					foreign.R = &applicationR{}
				}
				foreign.R.Builds = append(foreign.R.Builds, local)
				break
			}
		}
	}

	return nil
}

// LoadArtifact allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-1 relationship.
func (buildL) LoadArtifact(ctx context.Context, e boil.ContextExecutor, singular bool, maybeBuild interface{}, mods queries.Applicator) error {
	var slice []*Build
	var object *Build

	if singular {
		var ok bool
		object, ok = maybeBuild.(*Build)
		if !ok {
			object = new(Build)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeBuild)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeBuild))
			}
		}
	} else {
		s, ok := maybeBuild.(*[]*Build)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeBuild)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeBuild))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &buildR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &buildR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`artifacts`),
		qm.WhereIn(`artifacts.build_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Artifact")
	}

	var resultSlice []*Artifact
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Artifact")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for artifacts")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for artifacts")
	}

	if len(artifactAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Artifact = foreign
		if foreign.R == nil {
			foreign.R = &artifactR{}
		}
		foreign.R.Build = object
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ID == foreign.BuildID {
				local.R.Artifact = foreign
				if foreign.R == nil {
					foreign.R = &artifactR{}
				}
				foreign.R.Build = local
				break
			}
		}
	}

	return nil
}

// SetApplication of the build to the related item.
// Sets o.R.Application to related.
// Adds o to related.R.Builds.
func (o *Build) SetApplication(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Application) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `builds` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"application_id"}),
		strmangle.WhereClause("`", "`", 0, buildPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ApplicationID = related.ID
	if o.R == nil {
		o.R = &buildR{
			Application: related,
		}
	} else {
		o.R.Application = related
	}

	if related.R == nil {
		related.R = &applicationR{
			Builds: BuildSlice{o},
		}
	} else {
		related.R.Builds = append(related.R.Builds, o)
	}

	return nil
}

// SetArtifact of the build to the related item.
// Sets o.R.Artifact to related.
// Adds o to related.R.Build.
func (o *Build) SetArtifact(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Artifact) error {
	var err error

	if insert {
		related.BuildID = o.ID

		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	} else {
		updateQuery := fmt.Sprintf(
			"UPDATE `artifacts` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, []string{"build_id"}),
			strmangle.WhereClause("`", "`", 0, artifactPrimaryKeyColumns),
		)
		values := []interface{}{o.ID, related.ID}

		if boil.IsDebug(ctx) {
			writer := boil.DebugWriterFrom(ctx)
			fmt.Fprintln(writer, updateQuery)
			fmt.Fprintln(writer, values)
		}
		if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
			return errors.Wrap(err, "failed to update foreign table")
		}

		related.BuildID = o.ID
	}

	if o.R == nil {
		o.R = &buildR{
			Artifact: related,
		}
	} else {
		o.R.Artifact = related
	}

	if related.R == nil {
		related.R = &artifactR{
			Build: o,
		}
	} else {
		related.R.Build = o
	}
	return nil
}

// Builds retrieves all the records using an executor.
func Builds(mods ...qm.QueryMod) buildQuery {
	mods = append(mods, qm.From("`builds`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`builds`.*"})
	}

	return buildQuery{q}
}

// FindBuild retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBuild(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Build, error) {
	buildObj := &Build{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `builds` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, buildObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from builds")
	}

	if err = buildObj.doAfterSelectHooks(ctx, exec); err != nil {
		return buildObj, err
	}

	return buildObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Build) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no builds provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(buildColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	buildInsertCacheMut.RLock()
	cache, cached := buildInsertCache[key]
	buildInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			buildAllColumns,
			buildColumnsWithDefault,
			buildColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(buildType, buildMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(buildType, buildMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `builds` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `builds` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `builds` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, buildPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into builds")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for builds")
	}

CacheNoHooks:
	if !cached {
		buildInsertCacheMut.Lock()
		buildInsertCache[key] = cache
		buildInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Build.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Build) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	buildUpdateCacheMut.RLock()
	cache, cached := buildUpdateCache[key]
	buildUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			buildAllColumns,
			buildPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update builds, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `builds` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, buildPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(buildType, buildMapping, append(wl, buildPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update builds row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for builds")
	}

	if !cached {
		buildUpdateCacheMut.Lock()
		buildUpdateCache[key] = cache
		buildUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q buildQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for builds")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for builds")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BuildSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), buildPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `builds` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, buildPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in build slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all build")
	}
	return rowsAff, nil
}

var mySQLBuildUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Build) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no builds provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(buildColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLBuildUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	buildUpsertCacheMut.RLock()
	cache, cached := buildUpsertCache[key]
	buildUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			buildAllColumns,
			buildColumnsWithDefault,
			buildColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			buildAllColumns,
			buildPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert builds, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`builds`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `builds` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(buildType, buildMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(buildType, buildMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for builds")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(buildType, buildMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for builds")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for builds")
	}

CacheNoHooks:
	if !cached {
		buildUpsertCacheMut.Lock()
		buildUpsertCache[key] = cache
		buildUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Build record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Build) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Build provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), buildPrimaryKeyMapping)
	sql := "DELETE FROM `builds` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from builds")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for builds")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q buildQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no buildQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from builds")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for builds")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BuildSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(buildBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), buildPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `builds` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, buildPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from build slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for builds")
	}

	if len(buildAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Build) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindBuild(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BuildSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := BuildSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), buildPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `builds`.* FROM `builds` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, buildPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in BuildSlice")
	}

	*o = slice

	return nil
}

// BuildExists checks if the Build row exists.
func BuildExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `builds` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if builds exists")
	}

	return exists, nil
}

// Exists checks if the Build row exists.
func (o *Build) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return BuildExists(ctx, exec, o.ID)
}