// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Website is an object representing the database table.
type Website struct { // サイトID
	ID string `boil:"id" json:"id" toml:"id" yaml:"id"`
	// サイトURLのFQDN
	FQDN string `boil:"fqdn" json:"fqdn" toml:"fqdn" yaml:"fqdn"`
	// HTTPポート番号
	HTTPPort int `boil:"http_port" json:"http_port" toml:"http_port" yaml:"http_port"`
	// 作成日時
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	// 更新日時
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	// ブランチID
	BranchID string `boil:"branch_id" json:"branch_id" toml:"branch_id" yaml:"branch_id"`

	R *websiteR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L websiteL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var WebsiteColumns = struct {
	ID        string
	FQDN      string
	HTTPPort  string
	CreatedAt string
	UpdatedAt string
	BranchID  string
}{
	ID:        "id",
	FQDN:      "fqdn",
	HTTPPort:  "http_port",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	BranchID:  "branch_id",
}

var WebsiteTableColumns = struct {
	ID        string
	FQDN      string
	HTTPPort  string
	CreatedAt string
	UpdatedAt string
	BranchID  string
}{
	ID:        "websites.id",
	FQDN:      "websites.fqdn",
	HTTPPort:  "websites.http_port",
	CreatedAt: "websites.created_at",
	UpdatedAt: "websites.updated_at",
	BranchID:  "websites.branch_id",
}

// Generated where

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var WebsiteWhere = struct {
	ID        whereHelperstring
	FQDN      whereHelperstring
	HTTPPort  whereHelperint
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
	BranchID  whereHelperstring
}{
	ID:        whereHelperstring{field: "`websites`.`id`"},
	FQDN:      whereHelperstring{field: "`websites`.`fqdn`"},
	HTTPPort:  whereHelperint{field: "`websites`.`http_port`"},
	CreatedAt: whereHelpertime_Time{field: "`websites`.`created_at`"},
	UpdatedAt: whereHelpertime_Time{field: "`websites`.`updated_at`"},
	BranchID:  whereHelperstring{field: "`websites`.`branch_id`"},
}

// WebsiteRels is where relationship names are stored.
var WebsiteRels = struct {
	Branch string
}{
	Branch: "Branch",
}

// websiteR is where relationships are stored.
type websiteR struct {
	Branch *Branch `boil:"Branch" json:"Branch" toml:"Branch" yaml:"Branch"`
}

// NewStruct creates a new relationship struct
func (*websiteR) NewStruct() *websiteR {
	return &websiteR{}
}

func (r *websiteR) GetBranch() *Branch {
	if r == nil {
		return nil
	}
	return r.Branch
}

// websiteL is where Load methods for each relationship are stored.
type websiteL struct{}

var (
	websiteAllColumns            = []string{"id", "fqdn", "http_port", "created_at", "updated_at", "branch_id"}
	websiteColumnsWithoutDefault = []string{"id", "fqdn", "created_at", "updated_at", "branch_id"}
	websiteColumnsWithDefault    = []string{"http_port"}
	websitePrimaryKeyColumns     = []string{"id"}
	websiteGeneratedColumns      = []string{}
)

type (
	// WebsiteSlice is an alias for a slice of pointers to Website.
	// This should almost always be used instead of []Website.
	WebsiteSlice []*Website
	// WebsiteHook is the signature for custom Website hook methods
	WebsiteHook func(context.Context, boil.ContextExecutor, *Website) error

	websiteQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	websiteType                 = reflect.TypeOf(&Website{})
	websiteMapping              = queries.MakeStructMapping(websiteType)
	websitePrimaryKeyMapping, _ = queries.BindMapping(websiteType, websiteMapping, websitePrimaryKeyColumns)
	websiteInsertCacheMut       sync.RWMutex
	websiteInsertCache          = make(map[string]insertCache)
	websiteUpdateCacheMut       sync.RWMutex
	websiteUpdateCache          = make(map[string]updateCache)
	websiteUpsertCacheMut       sync.RWMutex
	websiteUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var websiteAfterSelectHooks []WebsiteHook

var websiteBeforeInsertHooks []WebsiteHook
var websiteAfterInsertHooks []WebsiteHook

var websiteBeforeUpdateHooks []WebsiteHook
var websiteAfterUpdateHooks []WebsiteHook

var websiteBeforeDeleteHooks []WebsiteHook
var websiteAfterDeleteHooks []WebsiteHook

var websiteBeforeUpsertHooks []WebsiteHook
var websiteAfterUpsertHooks []WebsiteHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Website) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Website) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Website) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Website) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Website) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Website) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Website) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Website) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Website) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range websiteAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddWebsiteHook registers your hook function for all future operations.
func AddWebsiteHook(hookPoint boil.HookPoint, websiteHook WebsiteHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		websiteAfterSelectHooks = append(websiteAfterSelectHooks, websiteHook)
	case boil.BeforeInsertHook:
		websiteBeforeInsertHooks = append(websiteBeforeInsertHooks, websiteHook)
	case boil.AfterInsertHook:
		websiteAfterInsertHooks = append(websiteAfterInsertHooks, websiteHook)
	case boil.BeforeUpdateHook:
		websiteBeforeUpdateHooks = append(websiteBeforeUpdateHooks, websiteHook)
	case boil.AfterUpdateHook:
		websiteAfterUpdateHooks = append(websiteAfterUpdateHooks, websiteHook)
	case boil.BeforeDeleteHook:
		websiteBeforeDeleteHooks = append(websiteBeforeDeleteHooks, websiteHook)
	case boil.AfterDeleteHook:
		websiteAfterDeleteHooks = append(websiteAfterDeleteHooks, websiteHook)
	case boil.BeforeUpsertHook:
		websiteBeforeUpsertHooks = append(websiteBeforeUpsertHooks, websiteHook)
	case boil.AfterUpsertHook:
		websiteAfterUpsertHooks = append(websiteAfterUpsertHooks, websiteHook)
	}
}

// One returns a single website record from the query.
func (q websiteQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Website, error) {
	o := &Website{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for websites")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Website records from the query.
func (q websiteQuery) All(ctx context.Context, exec boil.ContextExecutor) (WebsiteSlice, error) {
	var o []*Website

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Website slice")
	}

	if len(websiteAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Website records in the query.
func (q websiteQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count websites rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q websiteQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if websites exists")
	}

	return count > 0, nil
}

// Branch pointed to by the foreign key.
func (o *Website) Branch(mods ...qm.QueryMod) branchQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`id` = ?", o.BranchID),
	}

	queryMods = append(queryMods, mods...)

	return Branches(queryMods...)
}

// LoadBranch allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (websiteL) LoadBranch(ctx context.Context, e boil.ContextExecutor, singular bool, maybeWebsite interface{}, mods queries.Applicator) error {
	var slice []*Website
	var object *Website

	if singular {
		object = maybeWebsite.(*Website)
	} else {
		slice = *maybeWebsite.(*[]*Website)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &websiteR{}
		}
		args = append(args, object.BranchID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &websiteR{}
			}

			for _, a := range args {
				if a == obj.BranchID {
					continue Outer
				}
			}

			args = append(args, obj.BranchID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`branches`),
		qm.WhereIn(`branches.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Branch")
	}

	var resultSlice []*Branch
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Branch")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for branches")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for branches")
	}

	if len(websiteAfterSelectHooks) != 0 {
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
		object.R.Branch = foreign
		if foreign.R == nil {
			foreign.R = &branchR{}
		}
		foreign.R.Website = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.BranchID == foreign.ID {
				local.R.Branch = foreign
				if foreign.R == nil {
					foreign.R = &branchR{}
				}
				foreign.R.Website = local
				break
			}
		}
	}

	return nil
}

// SetBranch of the website to the related item.
// Sets o.R.Branch to related.
// Adds o to related.R.Website.
func (o *Website) SetBranch(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Branch) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `websites` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"branch_id"}),
		strmangle.WhereClause("`", "`", 0, websitePrimaryKeyColumns),
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

	o.BranchID = related.ID
	if o.R == nil {
		o.R = &websiteR{
			Branch: related,
		}
	} else {
		o.R.Branch = related
	}

	if related.R == nil {
		related.R = &branchR{
			Website: o,
		}
	} else {
		related.R.Website = o
	}

	return nil
}

// Websites retrieves all the records using an executor.
func Websites(mods ...qm.QueryMod) websiteQuery {
	mods = append(mods, qm.From("`websites`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`websites`.*"})
	}

	return websiteQuery{q}
}

// FindWebsite retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindWebsite(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Website, error) {
	websiteObj := &Website{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `websites` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, websiteObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from websites")
	}

	if err = websiteObj.doAfterSelectHooks(ctx, exec); err != nil {
		return websiteObj, err
	}

	return websiteObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Website) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no websites provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(websiteColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	websiteInsertCacheMut.RLock()
	cache, cached := websiteInsertCache[key]
	websiteInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			websiteAllColumns,
			websiteColumnsWithDefault,
			websiteColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(websiteType, websiteMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(websiteType, websiteMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `websites` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `websites` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `websites` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, websitePrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into websites")
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
		return errors.Wrap(err, "models: unable to populate default values for websites")
	}

CacheNoHooks:
	if !cached {
		websiteInsertCacheMut.Lock()
		websiteInsertCache[key] = cache
		websiteInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Website.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Website) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	websiteUpdateCacheMut.RLock()
	cache, cached := websiteUpdateCache[key]
	websiteUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			websiteAllColumns,
			websitePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update websites, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `websites` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, websitePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(websiteType, websiteMapping, append(wl, websitePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update websites row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for websites")
	}

	if !cached {
		websiteUpdateCacheMut.Lock()
		websiteUpdateCache[key] = cache
		websiteUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q websiteQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for websites")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for websites")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o WebsiteSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), websitePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `websites` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, websitePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in website slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all website")
	}
	return rowsAff, nil
}

var mySQLWebsiteUniqueColumns = []string{
	"id",
	"fqdn",
	"branch_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Website) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no websites provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(websiteColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLWebsiteUniqueColumns, o)

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

	websiteUpsertCacheMut.RLock()
	cache, cached := websiteUpsertCache[key]
	websiteUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			websiteAllColumns,
			websiteColumnsWithDefault,
			websiteColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			websiteAllColumns,
			websitePrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert websites, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`websites`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `websites` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(websiteType, websiteMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(websiteType, websiteMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for websites")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(websiteType, websiteMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for websites")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for websites")
	}

CacheNoHooks:
	if !cached {
		websiteUpsertCacheMut.Lock()
		websiteUpsertCache[key] = cache
		websiteUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Website record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Website) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Website provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), websitePrimaryKeyMapping)
	sql := "DELETE FROM `websites` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from websites")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for websites")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q websiteQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no websiteQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from websites")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for websites")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o WebsiteSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(websiteBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), websitePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `websites` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, websitePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from website slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for websites")
	}

	if len(websiteAfterDeleteHooks) != 0 {
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
func (o *Website) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindWebsite(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *WebsiteSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := WebsiteSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), websitePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `websites`.* FROM `websites` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, websitePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in WebsiteSlice")
	}

	*o = slice

	return nil
}

// WebsiteExists checks if the Website row exists.
func WebsiteExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `websites` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if websites exists")
	}

	return exists, nil
}
