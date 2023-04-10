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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ApplicationConfig is an object representing the database table.
type ApplicationConfig struct { // アプリケーションID
	ApplicationID string `boil:"application_id" json:"application_id" toml:"application_id" yaml:"application_id"`
	// MariaDBを使用するか
	UseMariadb bool `boil:"use_mariadb" json:"use_mariadb" toml:"use_mariadb" yaml:"use_mariadb"`
	// MongoDBを使用するか
	UseMongodb bool `boil:"use_mongodb" json:"use_mongodb" toml:"use_mongodb" yaml:"use_mongodb"`
	// ビルドタイプ
	BuildType string `boil:"build_type" json:"build_type" toml:"build_type" yaml:"build_type"`
	// ベースイメージの名前
	BaseImage string `boil:"base_image" json:"base_image" toml:"base_image" yaml:"base_image"`
	// ビルドコマンド(shell)
	BuildCMD string `boil:"build_cmd" json:"build_cmd" toml:"build_cmd" yaml:"build_cmd"`
	// コンテナのエントリポイント(shell)
	EntrypointCMD string `boil:"entrypoint_cmd" json:"entrypoint_cmd" toml:"entrypoint_cmd" yaml:"entrypoint_cmd"`
	// 静的成果物のパス
	ArtifactPath string `boil:"artifact_path" json:"artifact_path" toml:"artifact_path" yaml:"artifact_path"`
	// Dockerfile名
	DockerfileName string `boil:"dockerfile_name" json:"dockerfile_name" toml:"dockerfile_name" yaml:"dockerfile_name"`
	// Entrypointの上書き(args)
	EntrypointOverride string `boil:"entrypoint_override" json:"entrypoint_override" toml:"entrypoint_override" yaml:"entrypoint_override"`
	// Commandの上書き(args)
	CommandOverride string `boil:"command_override" json:"command_override" toml:"command_override" yaml:"command_override"`

	R *applicationConfigR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L applicationConfigL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ApplicationConfigColumns = struct {
	ApplicationID      string
	UseMariadb         string
	UseMongodb         string
	BuildType          string
	BaseImage          string
	BuildCMD           string
	EntrypointCMD      string
	ArtifactPath       string
	DockerfileName     string
	EntrypointOverride string
	CommandOverride    string
}{
	ApplicationID:      "application_id",
	UseMariadb:         "use_mariadb",
	UseMongodb:         "use_mongodb",
	BuildType:          "build_type",
	BaseImage:          "base_image",
	BuildCMD:           "build_cmd",
	EntrypointCMD:      "entrypoint_cmd",
	ArtifactPath:       "artifact_path",
	DockerfileName:     "dockerfile_name",
	EntrypointOverride: "entrypoint_override",
	CommandOverride:    "command_override",
}

var ApplicationConfigTableColumns = struct {
	ApplicationID      string
	UseMariadb         string
	UseMongodb         string
	BuildType          string
	BaseImage          string
	BuildCMD           string
	EntrypointCMD      string
	ArtifactPath       string
	DockerfileName     string
	EntrypointOverride string
	CommandOverride    string
}{
	ApplicationID:      "application_config.application_id",
	UseMariadb:         "application_config.use_mariadb",
	UseMongodb:         "application_config.use_mongodb",
	BuildType:          "application_config.build_type",
	BaseImage:          "application_config.base_image",
	BuildCMD:           "application_config.build_cmd",
	EntrypointCMD:      "application_config.entrypoint_cmd",
	ArtifactPath:       "application_config.artifact_path",
	DockerfileName:     "application_config.dockerfile_name",
	EntrypointOverride: "application_config.entrypoint_override",
	CommandOverride:    "application_config.command_override",
}

// Generated where

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

var ApplicationConfigWhere = struct {
	ApplicationID      whereHelperstring
	UseMariadb         whereHelperbool
	UseMongodb         whereHelperbool
	BuildType          whereHelperstring
	BaseImage          whereHelperstring
	BuildCMD           whereHelperstring
	EntrypointCMD      whereHelperstring
	ArtifactPath       whereHelperstring
	DockerfileName     whereHelperstring
	EntrypointOverride whereHelperstring
	CommandOverride    whereHelperstring
}{
	ApplicationID:      whereHelperstring{field: "`application_config`.`application_id`"},
	UseMariadb:         whereHelperbool{field: "`application_config`.`use_mariadb`"},
	UseMongodb:         whereHelperbool{field: "`application_config`.`use_mongodb`"},
	BuildType:          whereHelperstring{field: "`application_config`.`build_type`"},
	BaseImage:          whereHelperstring{field: "`application_config`.`base_image`"},
	BuildCMD:           whereHelperstring{field: "`application_config`.`build_cmd`"},
	EntrypointCMD:      whereHelperstring{field: "`application_config`.`entrypoint_cmd`"},
	ArtifactPath:       whereHelperstring{field: "`application_config`.`artifact_path`"},
	DockerfileName:     whereHelperstring{field: "`application_config`.`dockerfile_name`"},
	EntrypointOverride: whereHelperstring{field: "`application_config`.`entrypoint_override`"},
	CommandOverride:    whereHelperstring{field: "`application_config`.`command_override`"},
}

// ApplicationConfigRels is where relationship names are stored.
var ApplicationConfigRels = struct {
	Application string
}{
	Application: "Application",
}

// applicationConfigR is where relationships are stored.
type applicationConfigR struct {
	Application *Application `boil:"Application" json:"Application" toml:"Application" yaml:"Application"`
}

// NewStruct creates a new relationship struct
func (*applicationConfigR) NewStruct() *applicationConfigR {
	return &applicationConfigR{}
}

func (r *applicationConfigR) GetApplication() *Application {
	if r == nil {
		return nil
	}
	return r.Application
}

// applicationConfigL is where Load methods for each relationship are stored.
type applicationConfigL struct{}

var (
	applicationConfigAllColumns            = []string{"application_id", "use_mariadb", "use_mongodb", "build_type", "base_image", "build_cmd", "entrypoint_cmd", "artifact_path", "dockerfile_name", "entrypoint_override", "command_override"}
	applicationConfigColumnsWithoutDefault = []string{"application_id", "use_mariadb", "use_mongodb", "build_type", "base_image", "build_cmd", "entrypoint_cmd", "artifact_path", "dockerfile_name", "entrypoint_override", "command_override"}
	applicationConfigColumnsWithDefault    = []string{}
	applicationConfigPrimaryKeyColumns     = []string{"application_id"}
	applicationConfigGeneratedColumns      = []string{}
)

type (
	// ApplicationConfigSlice is an alias for a slice of pointers to ApplicationConfig.
	// This should almost always be used instead of []ApplicationConfig.
	ApplicationConfigSlice []*ApplicationConfig
	// ApplicationConfigHook is the signature for custom ApplicationConfig hook methods
	ApplicationConfigHook func(context.Context, boil.ContextExecutor, *ApplicationConfig) error

	applicationConfigQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	applicationConfigType                 = reflect.TypeOf(&ApplicationConfig{})
	applicationConfigMapping              = queries.MakeStructMapping(applicationConfigType)
	applicationConfigPrimaryKeyMapping, _ = queries.BindMapping(applicationConfigType, applicationConfigMapping, applicationConfigPrimaryKeyColumns)
	applicationConfigInsertCacheMut       sync.RWMutex
	applicationConfigInsertCache          = make(map[string]insertCache)
	applicationConfigUpdateCacheMut       sync.RWMutex
	applicationConfigUpdateCache          = make(map[string]updateCache)
	applicationConfigUpsertCacheMut       sync.RWMutex
	applicationConfigUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var applicationConfigAfterSelectHooks []ApplicationConfigHook

var applicationConfigBeforeInsertHooks []ApplicationConfigHook
var applicationConfigAfterInsertHooks []ApplicationConfigHook

var applicationConfigBeforeUpdateHooks []ApplicationConfigHook
var applicationConfigAfterUpdateHooks []ApplicationConfigHook

var applicationConfigBeforeDeleteHooks []ApplicationConfigHook
var applicationConfigAfterDeleteHooks []ApplicationConfigHook

var applicationConfigBeforeUpsertHooks []ApplicationConfigHook
var applicationConfigAfterUpsertHooks []ApplicationConfigHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ApplicationConfig) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ApplicationConfig) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ApplicationConfig) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ApplicationConfig) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ApplicationConfig) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ApplicationConfig) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ApplicationConfig) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ApplicationConfig) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ApplicationConfig) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range applicationConfigAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddApplicationConfigHook registers your hook function for all future operations.
func AddApplicationConfigHook(hookPoint boil.HookPoint, applicationConfigHook ApplicationConfigHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		applicationConfigAfterSelectHooks = append(applicationConfigAfterSelectHooks, applicationConfigHook)
	case boil.BeforeInsertHook:
		applicationConfigBeforeInsertHooks = append(applicationConfigBeforeInsertHooks, applicationConfigHook)
	case boil.AfterInsertHook:
		applicationConfigAfterInsertHooks = append(applicationConfigAfterInsertHooks, applicationConfigHook)
	case boil.BeforeUpdateHook:
		applicationConfigBeforeUpdateHooks = append(applicationConfigBeforeUpdateHooks, applicationConfigHook)
	case boil.AfterUpdateHook:
		applicationConfigAfterUpdateHooks = append(applicationConfigAfterUpdateHooks, applicationConfigHook)
	case boil.BeforeDeleteHook:
		applicationConfigBeforeDeleteHooks = append(applicationConfigBeforeDeleteHooks, applicationConfigHook)
	case boil.AfterDeleteHook:
		applicationConfigAfterDeleteHooks = append(applicationConfigAfterDeleteHooks, applicationConfigHook)
	case boil.BeforeUpsertHook:
		applicationConfigBeforeUpsertHooks = append(applicationConfigBeforeUpsertHooks, applicationConfigHook)
	case boil.AfterUpsertHook:
		applicationConfigAfterUpsertHooks = append(applicationConfigAfterUpsertHooks, applicationConfigHook)
	}
}

// One returns a single applicationConfig record from the query.
func (q applicationConfigQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ApplicationConfig, error) {
	o := &ApplicationConfig{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for application_config")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ApplicationConfig records from the query.
func (q applicationConfigQuery) All(ctx context.Context, exec boil.ContextExecutor) (ApplicationConfigSlice, error) {
	var o []*ApplicationConfig

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ApplicationConfig slice")
	}

	if len(applicationConfigAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ApplicationConfig records in the query.
func (q applicationConfigQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count application_config rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q applicationConfigQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if application_config exists")
	}

	return count > 0, nil
}

// Application pointed to by the foreign key.
func (o *ApplicationConfig) Application(mods ...qm.QueryMod) applicationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`id` = ?", o.ApplicationID),
	}

	queryMods = append(queryMods, mods...)

	return Applications(queryMods...)
}

// LoadApplication allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (applicationConfigL) LoadApplication(ctx context.Context, e boil.ContextExecutor, singular bool, maybeApplicationConfig interface{}, mods queries.Applicator) error {
	var slice []*ApplicationConfig
	var object *ApplicationConfig

	if singular {
		var ok bool
		object, ok = maybeApplicationConfig.(*ApplicationConfig)
		if !ok {
			object = new(ApplicationConfig)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeApplicationConfig)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeApplicationConfig))
			}
		}
	} else {
		s, ok := maybeApplicationConfig.(*[]*ApplicationConfig)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeApplicationConfig)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeApplicationConfig))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &applicationConfigR{}
		}
		args = append(args, object.ApplicationID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &applicationConfigR{}
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
		foreign.R.ApplicationConfig = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ApplicationID == foreign.ID {
				local.R.Application = foreign
				if foreign.R == nil {
					foreign.R = &applicationR{}
				}
				foreign.R.ApplicationConfig = local
				break
			}
		}
	}

	return nil
}

// SetApplication of the applicationConfig to the related item.
// Sets o.R.Application to related.
// Adds o to related.R.ApplicationConfig.
func (o *ApplicationConfig) SetApplication(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Application) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `application_config` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"application_id"}),
		strmangle.WhereClause("`", "`", 0, applicationConfigPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ApplicationID}

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
		o.R = &applicationConfigR{
			Application: related,
		}
	} else {
		o.R.Application = related
	}

	if related.R == nil {
		related.R = &applicationR{
			ApplicationConfig: o,
		}
	} else {
		related.R.ApplicationConfig = o
	}

	return nil
}

// ApplicationConfigs retrieves all the records using an executor.
func ApplicationConfigs(mods ...qm.QueryMod) applicationConfigQuery {
	mods = append(mods, qm.From("`application_config`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`application_config`.*"})
	}

	return applicationConfigQuery{q}
}

// FindApplicationConfig retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindApplicationConfig(ctx context.Context, exec boil.ContextExecutor, applicationID string, selectCols ...string) (*ApplicationConfig, error) {
	applicationConfigObj := &ApplicationConfig{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `application_config` where `application_id`=?", sel,
	)

	q := queries.Raw(query, applicationID)

	err := q.Bind(ctx, exec, applicationConfigObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from application_config")
	}

	if err = applicationConfigObj.doAfterSelectHooks(ctx, exec); err != nil {
		return applicationConfigObj, err
	}

	return applicationConfigObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ApplicationConfig) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no application_config provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(applicationConfigColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	applicationConfigInsertCacheMut.RLock()
	cache, cached := applicationConfigInsertCache[key]
	applicationConfigInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			applicationConfigAllColumns,
			applicationConfigColumnsWithDefault,
			applicationConfigColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(applicationConfigType, applicationConfigMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(applicationConfigType, applicationConfigMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `application_config` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `application_config` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `application_config` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, applicationConfigPrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into application_config")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ApplicationID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for application_config")
	}

CacheNoHooks:
	if !cached {
		applicationConfigInsertCacheMut.Lock()
		applicationConfigInsertCache[key] = cache
		applicationConfigInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ApplicationConfig.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ApplicationConfig) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	applicationConfigUpdateCacheMut.RLock()
	cache, cached := applicationConfigUpdateCache[key]
	applicationConfigUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			applicationConfigAllColumns,
			applicationConfigPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update application_config, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `application_config` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, applicationConfigPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(applicationConfigType, applicationConfigMapping, append(wl, applicationConfigPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update application_config row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for application_config")
	}

	if !cached {
		applicationConfigUpdateCacheMut.Lock()
		applicationConfigUpdateCache[key] = cache
		applicationConfigUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q applicationConfigQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for application_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for application_config")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ApplicationConfigSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), applicationConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `application_config` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, applicationConfigPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in applicationConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all applicationConfig")
	}
	return rowsAff, nil
}

var mySQLApplicationConfigUniqueColumns = []string{
	"application_id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ApplicationConfig) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no application_config provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(applicationConfigColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLApplicationConfigUniqueColumns, o)

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

	applicationConfigUpsertCacheMut.RLock()
	cache, cached := applicationConfigUpsertCache[key]
	applicationConfigUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			applicationConfigAllColumns,
			applicationConfigColumnsWithDefault,
			applicationConfigColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			applicationConfigAllColumns,
			applicationConfigPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("models: unable to upsert application_config, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`application_config`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `application_config` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(applicationConfigType, applicationConfigMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(applicationConfigType, applicationConfigMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for application_config")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(applicationConfigType, applicationConfigMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for application_config")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for application_config")
	}

CacheNoHooks:
	if !cached {
		applicationConfigUpsertCacheMut.Lock()
		applicationConfigUpsertCache[key] = cache
		applicationConfigUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single ApplicationConfig record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ApplicationConfig) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ApplicationConfig provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), applicationConfigPrimaryKeyMapping)
	sql := "DELETE FROM `application_config` WHERE `application_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from application_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for application_config")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q applicationConfigQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no applicationConfigQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from application_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for application_config")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ApplicationConfigSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(applicationConfigBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), applicationConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `application_config` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, applicationConfigPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from applicationConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for application_config")
	}

	if len(applicationConfigAfterDeleteHooks) != 0 {
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
func (o *ApplicationConfig) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindApplicationConfig(ctx, exec, o.ApplicationID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ApplicationConfigSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ApplicationConfigSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), applicationConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `application_config`.* FROM `application_config` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, applicationConfigPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ApplicationConfigSlice")
	}

	*o = slice

	return nil
}

// ApplicationConfigExists checks if the ApplicationConfig row exists.
func ApplicationConfigExists(ctx context.Context, exec boil.ContextExecutor, applicationID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `application_config` where `application_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, applicationID)
	}
	row := exec.QueryRowContext(ctx, sql, applicationID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if application_config exists")
	}

	return exists, nil
}

// Exists checks if the ApplicationConfig row exists.
func (o *ApplicationConfig) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ApplicationConfigExists(ctx, exec, o.ApplicationID)
}
