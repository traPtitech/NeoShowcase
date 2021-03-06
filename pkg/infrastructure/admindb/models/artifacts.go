// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// Artifact is an object representing the database table.
type Artifact struct {
	ID         string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	BuildLogID string    `boil:"build_log_id" json:"build_log_id" toml:"build_log_id" yaml:"build_log_id"`
	Size       int64     `boil:"size" json:"size" toml:"size" yaml:"size"`
	CreatedAt  time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt  null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *artifactR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L artifactL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ArtifactColumns = struct {
	ID         string
	BuildLogID string
	Size       string
	CreatedAt  string
	DeletedAt  string
}{
	ID:         "id",
	BuildLogID: "build_log_id",
	Size:       "size",
	CreatedAt:  "created_at",
	DeletedAt:  "deleted_at",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var ArtifactWhere = struct {
	ID         whereHelperstring
	BuildLogID whereHelperstring
	Size       whereHelperint64
	CreatedAt  whereHelpertime_Time
	DeletedAt  whereHelpernull_Time
}{
	ID:         whereHelperstring{field: "`artifacts`.`id`"},
	BuildLogID: whereHelperstring{field: "`artifacts`.`build_log_id`"},
	Size:       whereHelperint64{field: "`artifacts`.`size`"},
	CreatedAt:  whereHelpertime_Time{field: "`artifacts`.`created_at`"},
	DeletedAt:  whereHelpernull_Time{field: "`artifacts`.`deleted_at`"},
}

// ArtifactRels is where relationship names are stored.
var ArtifactRels = struct {
	BuildLog string
}{
	BuildLog: "BuildLog",
}

// artifactR is where relationships are stored.
type artifactR struct {
	BuildLog *BuildLog `boil:"BuildLog" json:"BuildLog" toml:"BuildLog" yaml:"BuildLog"`
}

// NewStruct creates a new relationship struct
func (*artifactR) NewStruct() *artifactR {
	return &artifactR{}
}

// artifactL is where Load methods for each relationship are stored.
type artifactL struct{}

var (
	artifactAllColumns            = []string{"id", "build_log_id", "size", "created_at", "deleted_at"}
	artifactColumnsWithoutDefault = []string{"id", "build_log_id", "size", "created_at", "deleted_at"}
	artifactColumnsWithDefault    = []string{}
	artifactPrimaryKeyColumns     = []string{"id"}
)

type (
	// ArtifactSlice is an alias for a slice of pointers to Artifact.
	// This should generally be used opposed to []Artifact.
	ArtifactSlice []*Artifact
	// ArtifactHook is the signature for custom Artifact hook methods
	ArtifactHook func(context.Context, boil.ContextExecutor, *Artifact) error

	artifactQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	artifactType                 = reflect.TypeOf(&Artifact{})
	artifactMapping              = queries.MakeStructMapping(artifactType)
	artifactPrimaryKeyMapping, _ = queries.BindMapping(artifactType, artifactMapping, artifactPrimaryKeyColumns)
	artifactInsertCacheMut       sync.RWMutex
	artifactInsertCache          = make(map[string]insertCache)
	artifactUpdateCacheMut       sync.RWMutex
	artifactUpdateCache          = make(map[string]updateCache)
	artifactUpsertCacheMut       sync.RWMutex
	artifactUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var artifactBeforeInsertHooks []ArtifactHook
var artifactBeforeUpdateHooks []ArtifactHook
var artifactBeforeDeleteHooks []ArtifactHook
var artifactBeforeUpsertHooks []ArtifactHook

var artifactAfterInsertHooks []ArtifactHook
var artifactAfterSelectHooks []ArtifactHook
var artifactAfterUpdateHooks []ArtifactHook
var artifactAfterDeleteHooks []ArtifactHook
var artifactAfterUpsertHooks []ArtifactHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Artifact) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Artifact) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Artifact) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Artifact) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Artifact) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Artifact) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Artifact) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Artifact) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Artifact) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range artifactAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddArtifactHook registers your hook function for all future operations.
func AddArtifactHook(hookPoint boil.HookPoint, artifactHook ArtifactHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		artifactBeforeInsertHooks = append(artifactBeforeInsertHooks, artifactHook)
	case boil.BeforeUpdateHook:
		artifactBeforeUpdateHooks = append(artifactBeforeUpdateHooks, artifactHook)
	case boil.BeforeDeleteHook:
		artifactBeforeDeleteHooks = append(artifactBeforeDeleteHooks, artifactHook)
	case boil.BeforeUpsertHook:
		artifactBeforeUpsertHooks = append(artifactBeforeUpsertHooks, artifactHook)
	case boil.AfterInsertHook:
		artifactAfterInsertHooks = append(artifactAfterInsertHooks, artifactHook)
	case boil.AfterSelectHook:
		artifactAfterSelectHooks = append(artifactAfterSelectHooks, artifactHook)
	case boil.AfterUpdateHook:
		artifactAfterUpdateHooks = append(artifactAfterUpdateHooks, artifactHook)
	case boil.AfterDeleteHook:
		artifactAfterDeleteHooks = append(artifactAfterDeleteHooks, artifactHook)
	case boil.AfterUpsertHook:
		artifactAfterUpsertHooks = append(artifactAfterUpsertHooks, artifactHook)
	}
}

// One returns a single artifact record from the query.
func (q artifactQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Artifact, error) {
	o := &Artifact{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for artifacts")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Artifact records from the query.
func (q artifactQuery) All(ctx context.Context, exec boil.ContextExecutor) (ArtifactSlice, error) {
	var o []*Artifact

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Artifact slice")
	}

	if len(artifactAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Artifact records in the query.
func (q artifactQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count artifacts rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q artifactQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if artifacts exists")
	}

	return count > 0, nil
}

// BuildLog pointed to by the foreign key.
func (o *Artifact) BuildLog(mods ...qm.QueryMod) buildLogQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`id` = ?", o.BuildLogID),
	}

	queryMods = append(queryMods, mods...)

	query := BuildLogs(queryMods...)
	queries.SetFrom(query.Query, "`build_logs`")

	return query
}

// LoadBuildLog allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (artifactL) LoadBuildLog(ctx context.Context, e boil.ContextExecutor, singular bool, maybeArtifact interface{}, mods queries.Applicator) error {
	var slice []*Artifact
	var object *Artifact

	if singular {
		object = maybeArtifact.(*Artifact)
	} else {
		slice = *maybeArtifact.(*[]*Artifact)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &artifactR{}
		}
		args = append(args, object.BuildLogID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &artifactR{}
			}

			for _, a := range args {
				if a == obj.BuildLogID {
					continue Outer
				}
			}

			args = append(args, obj.BuildLogID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`build_logs`),
		qm.WhereIn(`build_logs.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load BuildLog")
	}

	var resultSlice []*BuildLog
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice BuildLog")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for build_logs")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for build_logs")
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
		object.R.BuildLog = foreign
		if foreign.R == nil {
			foreign.R = &buildLogR{}
		}
		foreign.R.Artifact = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.BuildLogID == foreign.ID {
				local.R.BuildLog = foreign
				if foreign.R == nil {
					foreign.R = &buildLogR{}
				}
				foreign.R.Artifact = local
				break
			}
		}
	}

	return nil
}

// SetBuildLog of the artifact to the related item.
// Sets o.R.BuildLog to related.
// Adds o to related.R.Artifact.
func (o *Artifact) SetBuildLog(ctx context.Context, exec boil.ContextExecutor, insert bool, related *BuildLog) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `artifacts` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"build_log_id"}),
		strmangle.WhereClause("`", "`", 0, artifactPrimaryKeyColumns),
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

	o.BuildLogID = related.ID
	if o.R == nil {
		o.R = &artifactR{
			BuildLog: related,
		}
	} else {
		o.R.BuildLog = related
	}

	if related.R == nil {
		related.R = &buildLogR{
			Artifact: o,
		}
	} else {
		related.R.Artifact = o
	}

	return nil
}

// Artifacts retrieves all the records using an executor.
func Artifacts(mods ...qm.QueryMod) artifactQuery {
	mods = append(mods, qm.From("`artifacts`"))
	return artifactQuery{NewQuery(mods...)}
}

// FindArtifact retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindArtifact(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Artifact, error) {
	artifactObj := &Artifact{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `artifacts` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, artifactObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from artifacts")
	}

	return artifactObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Artifact) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no artifacts provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(artifactColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	artifactInsertCacheMut.RLock()
	cache, cached := artifactInsertCache[key]
	artifactInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			artifactAllColumns,
			artifactColumnsWithDefault,
			artifactColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(artifactType, artifactMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(artifactType, artifactMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `artifacts` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `artifacts` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `artifacts` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, artifactPrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into artifacts")
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
		return errors.Wrap(err, "models: unable to populate default values for artifacts")
	}

CacheNoHooks:
	if !cached {
		artifactInsertCacheMut.Lock()
		artifactInsertCache[key] = cache
		artifactInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Artifact.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Artifact) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	artifactUpdateCacheMut.RLock()
	cache, cached := artifactUpdateCache[key]
	artifactUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			artifactAllColumns,
			artifactPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update artifacts, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `artifacts` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, artifactPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(artifactType, artifactMapping, append(wl, artifactPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update artifacts row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for artifacts")
	}

	if !cached {
		artifactUpdateCacheMut.Lock()
		artifactUpdateCache[key] = cache
		artifactUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q artifactQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for artifacts")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for artifacts")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ArtifactSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), artifactPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `artifacts` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, artifactPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in artifact slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all artifact")
	}
	return rowsAff, nil
}

var mySQLArtifactUniqueColumns = []string{
	"id",
	"build_log_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Artifact) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no artifacts provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(artifactColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLArtifactUniqueColumns, o)

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

	artifactUpsertCacheMut.RLock()
	cache, cached := artifactUpsertCache[key]
	artifactUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			artifactAllColumns,
			artifactColumnsWithDefault,
			artifactColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			artifactAllColumns,
			artifactPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert artifacts, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`artifacts`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `artifacts` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(artifactType, artifactMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(artifactType, artifactMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for artifacts")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(artifactType, artifactMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for artifacts")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for artifacts")
	}

CacheNoHooks:
	if !cached {
		artifactUpsertCacheMut.Lock()
		artifactUpsertCache[key] = cache
		artifactUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Artifact record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Artifact) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Artifact provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), artifactPrimaryKeyMapping)
	sql := "DELETE FROM `artifacts` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from artifacts")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for artifacts")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q artifactQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no artifactQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from artifacts")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for artifacts")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ArtifactSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(artifactBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), artifactPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `artifacts` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, artifactPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from artifact slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for artifacts")
	}

	if len(artifactAfterDeleteHooks) != 0 {
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
func (o *Artifact) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindArtifact(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ArtifactSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ArtifactSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), artifactPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `artifacts`.* FROM `artifacts` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, artifactPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ArtifactSlice")
	}

	*o = slice

	return nil
}

// ArtifactExists checks if the Artifact row exists.
func ArtifactExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `artifacts` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if artifacts exists")
	}

	return exists, nil
}
