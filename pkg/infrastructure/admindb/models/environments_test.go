// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testEnvironments(t *testing.T) {
	t.Parallel()

	query := Environments()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testEnvironmentsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEnvironmentsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Environments().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEnvironmentsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := EnvironmentSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEnvironmentsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := EnvironmentExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Environment exists: %s", err)
	}
	if !e {
		t.Errorf("Expected EnvironmentExists to return true, but got false.")
	}
}

func testEnvironmentsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	environmentFound, err := FindEnvironment(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if environmentFound == nil {
		t.Error("want a record, got nil")
	}
}

func testEnvironmentsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Environments().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testEnvironmentsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Environments().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testEnvironmentsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	environmentOne := &Environment{}
	environmentTwo := &Environment{}
	if err = randomize.Struct(seed, environmentOne, environmentDBTypes, false, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}
	if err = randomize.Struct(seed, environmentTwo, environmentDBTypes, false, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = environmentOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = environmentTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Environments().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testEnvironmentsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	environmentOne := &Environment{}
	environmentTwo := &Environment{}
	if err = randomize.Struct(seed, environmentOne, environmentDBTypes, false, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}
	if err = randomize.Struct(seed, environmentTwo, environmentDBTypes, false, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = environmentOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = environmentTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func environmentBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func environmentAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Environment) error {
	*o = Environment{}
	return nil
}

func testEnvironmentsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Environment{}
	o := &Environment{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, environmentDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Environment object: %s", err)
	}

	AddEnvironmentHook(boil.BeforeInsertHook, environmentBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	environmentBeforeInsertHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.AfterInsertHook, environmentAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	environmentAfterInsertHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.AfterSelectHook, environmentAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	environmentAfterSelectHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.BeforeUpdateHook, environmentBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	environmentBeforeUpdateHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.AfterUpdateHook, environmentAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	environmentAfterUpdateHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.BeforeDeleteHook, environmentBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	environmentBeforeDeleteHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.AfterDeleteHook, environmentAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	environmentAfterDeleteHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.BeforeUpsertHook, environmentBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	environmentBeforeUpsertHooks = []EnvironmentHook{}

	AddEnvironmentHook(boil.AfterUpsertHook, environmentAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	environmentAfterUpsertHooks = []EnvironmentHook{}
}

func testEnvironmentsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testEnvironmentsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(environmentColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testEnvironmentOneToOneWebsiteUsingWebsite(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var foreign Website
	var local Environment

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &foreign, websiteDBTypes, true, websiteColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Website struct: %s", err)
	}
	if err := randomize.Struct(seed, &local, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreign.EnvironmentID = local.ID
	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Website().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.EnvironmentID != foreign.EnvironmentID {
		t.Errorf("want: %v, got %v", foreign.EnvironmentID, check.EnvironmentID)
	}

	slice := EnvironmentSlice{&local}
	if err = local.L.LoadWebsite(ctx, tx, false, (*[]*Environment)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Website == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Website = nil
	if err = local.L.LoadWebsite(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Website == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testEnvironmentOneToOneSetOpWebsiteUsingWebsite(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b, c Website

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, false, strmangle.SetComplement(environmentPrimaryKeyColumns, environmentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, websiteDBTypes, false, strmangle.SetComplement(websitePrimaryKeyColumns, websiteColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, websiteDBTypes, false, strmangle.SetComplement(websitePrimaryKeyColumns, websiteColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Website{&b, &c} {
		err = a.SetWebsite(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Website != x {
			t.Error("relationship struct not set to correct value")
		}
		if x.R.Environment != &a {
			t.Error("failed to append to foreign relationship struct")
		}

		if a.ID != x.EnvironmentID {
			t.Error("foreign key was wrong value", a.ID)
		}

		zero := reflect.Zero(reflect.TypeOf(x.EnvironmentID))
		reflect.Indirect(reflect.ValueOf(&x.EnvironmentID)).Set(zero)

		if err = x.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ID != x.EnvironmentID {
			t.Error("foreign key was wrong value", a.ID, x.EnvironmentID)
		}

		if _, err = x.Delete(ctx, tx); err != nil {
			t.Fatal("failed to delete x", err)
		}
	}
}

func testEnvironmentToManyBuildLogs(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b, c BuildLog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, buildLogDBTypes, false, buildLogColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, buildLogDBTypes, false, buildLogColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.EnvironmentID, a.ID)
	queries.Assign(&c.EnvironmentID, a.ID)
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.BuildLogs().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.EnvironmentID, b.EnvironmentID) {
			bFound = true
		}
		if queries.Equal(v.EnvironmentID, c.EnvironmentID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := EnvironmentSlice{&a}
	if err = a.L.LoadBuildLogs(ctx, tx, false, (*[]*Environment)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.BuildLogs); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.BuildLogs = nil
	if err = a.L.LoadBuildLogs(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.BuildLogs); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testEnvironmentToManyAddOpBuildLogs(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b, c, d, e BuildLog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, false, strmangle.SetComplement(environmentPrimaryKeyColumns, environmentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*BuildLog{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, buildLogDBTypes, false, strmangle.SetComplement(buildLogPrimaryKeyColumns, buildLogColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*BuildLog{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddBuildLogs(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if !queries.Equal(a.ID, first.EnvironmentID) {
			t.Error("foreign key was wrong value", a.ID, first.EnvironmentID)
		}
		if !queries.Equal(a.ID, second.EnvironmentID) {
			t.Error("foreign key was wrong value", a.ID, second.EnvironmentID)
		}

		if first.R.Environment != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Environment != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.BuildLogs[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.BuildLogs[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.BuildLogs().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testEnvironmentToManySetOpBuildLogs(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b, c, d, e BuildLog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, false, strmangle.SetComplement(environmentPrimaryKeyColumns, environmentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*BuildLog{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, buildLogDBTypes, false, strmangle.SetComplement(buildLogPrimaryKeyColumns, buildLogColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetBuildLogs(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.BuildLogs().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetBuildLogs(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.BuildLogs().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.EnvironmentID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.EnvironmentID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.EnvironmentID) {
		t.Error("foreign key was wrong value", a.ID, d.EnvironmentID)
	}
	if !queries.Equal(a.ID, e.EnvironmentID) {
		t.Error("foreign key was wrong value", a.ID, e.EnvironmentID)
	}

	if b.R.Environment != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Environment != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Environment != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Environment != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.BuildLogs[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.BuildLogs[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testEnvironmentToManyRemoveOpBuildLogs(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b, c, d, e BuildLog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, false, strmangle.SetComplement(environmentPrimaryKeyColumns, environmentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*BuildLog{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, buildLogDBTypes, false, strmangle.SetComplement(buildLogPrimaryKeyColumns, buildLogColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddBuildLogs(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.BuildLogs().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveBuildLogs(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.BuildLogs().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.EnvironmentID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.EnvironmentID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Environment != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Environment != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Environment != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Environment != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.BuildLogs) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.BuildLogs[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.BuildLogs[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testEnvironmentToOneApplicationUsingApplication(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Environment
	var foreign Application

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, environmentDBTypes, false, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, applicationDBTypes, false, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.ApplicationID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Application().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := EnvironmentSlice{&local}
	if err = local.L.LoadApplication(ctx, tx, false, (*[]*Environment)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Application == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Application = nil
	if err = local.L.LoadApplication(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Application == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testEnvironmentToOneBuildLogUsingBuild(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Environment
	var foreign BuildLog

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, buildLogDBTypes, false, buildLogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BuildLog struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.BuildID, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Build().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := EnvironmentSlice{&local}
	if err = local.L.LoadBuild(ctx, tx, false, (*[]*Environment)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Build == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Build = nil
	if err = local.L.LoadBuild(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Build == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testEnvironmentToOneSetOpApplicationUsingApplication(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b, c Application

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, false, strmangle.SetComplement(environmentPrimaryKeyColumns, environmentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, applicationDBTypes, false, strmangle.SetComplement(applicationPrimaryKeyColumns, applicationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, applicationDBTypes, false, strmangle.SetComplement(applicationPrimaryKeyColumns, applicationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Application{&b, &c} {
		err = a.SetApplication(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Application != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Environments[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ApplicationID != x.ID {
			t.Error("foreign key was wrong value", a.ApplicationID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ApplicationID))
		reflect.Indirect(reflect.ValueOf(&a.ApplicationID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ApplicationID != x.ID {
			t.Error("foreign key was wrong value", a.ApplicationID, x.ID)
		}
	}
}
func testEnvironmentToOneSetOpBuildLogUsingBuild(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b, c BuildLog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, false, strmangle.SetComplement(environmentPrimaryKeyColumns, environmentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, buildLogDBTypes, false, strmangle.SetComplement(buildLogPrimaryKeyColumns, buildLogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, buildLogDBTypes, false, strmangle.SetComplement(buildLogPrimaryKeyColumns, buildLogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*BuildLog{&b, &c} {
		err = a.SetBuild(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Build != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.BuildEnvironments[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.BuildID, x.ID) {
			t.Error("foreign key was wrong value", a.BuildID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.BuildID))
		reflect.Indirect(reflect.ValueOf(&a.BuildID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.BuildID, x.ID) {
			t.Error("foreign key was wrong value", a.BuildID, x.ID)
		}
	}
}

func testEnvironmentToOneRemoveOpBuildLogUsingBuild(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Environment
	var b BuildLog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, environmentDBTypes, false, strmangle.SetComplement(environmentPrimaryKeyColumns, environmentColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, buildLogDBTypes, false, strmangle.SetComplement(buildLogPrimaryKeyColumns, buildLogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetBuild(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveBuild(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Build().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Build != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.BuildID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.BuildEnvironments) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testEnvironmentsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testEnvironmentsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := EnvironmentSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testEnvironmentsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Environments().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	environmentDBTypes = map[string]string{`ID`: `varchar`, `ApplicationID`: `varchar`, `BranchName`: `varchar`, `BuildType`: `enum('image','static')`, `CreatedAt`: `datetime`, `UpdatedAt`: `datetime`, `BuildID`: `varchar`}
	_                  = bytes.MinRead
)

func testEnvironmentsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(environmentPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(environmentAllColumns) == len(environmentPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testEnvironmentsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(environmentAllColumns) == len(environmentPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Environment{}
	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, environmentDBTypes, true, environmentPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(environmentAllColumns, environmentPrimaryKeyColumns) {
		fields = environmentAllColumns
	} else {
		fields = strmangle.SetComplement(
			environmentAllColumns,
			environmentPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := EnvironmentSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testEnvironmentsUpsert(t *testing.T) {
	t.Parallel()

	if len(environmentAllColumns) == len(environmentPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLEnvironmentUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Environment{}
	if err = randomize.Struct(seed, &o, environmentDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Environment: %s", err)
	}

	count, err := Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, environmentDBTypes, false, environmentPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Environment struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Environment: %s", err)
	}

	count, err = Environments().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
