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

func testApplications(t *testing.T) {
	t.Parallel()

	query := Applications()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testApplicationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
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

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testApplicationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Applications().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testApplicationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ApplicationSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testApplicationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ApplicationExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Application exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ApplicationExists to return true, but got false.")
	}
}

func testApplicationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	applicationFound, err := FindApplication(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if applicationFound == nil {
		t.Error("want a record, got nil")
	}
}

func testApplicationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Applications().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testApplicationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Applications().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testApplicationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	applicationOne := &Application{}
	applicationTwo := &Application{}
	if err = randomize.Struct(seed, applicationOne, applicationDBTypes, false, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}
	if err = randomize.Struct(seed, applicationTwo, applicationDBTypes, false, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = applicationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = applicationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Applications().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testApplicationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	applicationOne := &Application{}
	applicationTwo := &Application{}
	if err = randomize.Struct(seed, applicationOne, applicationDBTypes, false, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}
	if err = randomize.Struct(seed, applicationTwo, applicationDBTypes, false, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = applicationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = applicationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func applicationBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func applicationAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Application) error {
	*o = Application{}
	return nil
}

func testApplicationsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Application{}
	o := &Application{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, applicationDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Application object: %s", err)
	}

	AddApplicationHook(boil.BeforeInsertHook, applicationBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	applicationBeforeInsertHooks = []ApplicationHook{}

	AddApplicationHook(boil.AfterInsertHook, applicationAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	applicationAfterInsertHooks = []ApplicationHook{}

	AddApplicationHook(boil.AfterSelectHook, applicationAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	applicationAfterSelectHooks = []ApplicationHook{}

	AddApplicationHook(boil.BeforeUpdateHook, applicationBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	applicationBeforeUpdateHooks = []ApplicationHook{}

	AddApplicationHook(boil.AfterUpdateHook, applicationAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	applicationAfterUpdateHooks = []ApplicationHook{}

	AddApplicationHook(boil.BeforeDeleteHook, applicationBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	applicationBeforeDeleteHooks = []ApplicationHook{}

	AddApplicationHook(boil.AfterDeleteHook, applicationAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	applicationAfterDeleteHooks = []ApplicationHook{}

	AddApplicationHook(boil.BeforeUpsertHook, applicationBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	applicationBeforeUpsertHooks = []ApplicationHook{}

	AddApplicationHook(boil.AfterUpsertHook, applicationAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	applicationAfterUpsertHooks = []ApplicationHook{}
}

func testApplicationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testApplicationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(applicationColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testApplicationToManyBranches(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Application
	var b, c Branch

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, branchDBTypes, false, branchColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, branchDBTypes, false, branchColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.ApplicationID = a.ID
	c.ApplicationID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Branches().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ApplicationID == b.ApplicationID {
			bFound = true
		}
		if v.ApplicationID == c.ApplicationID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ApplicationSlice{&a}
	if err = a.L.LoadBranches(ctx, tx, false, (*[]*Application)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Branches); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Branches = nil
	if err = a.L.LoadBranches(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Branches); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testApplicationToManyUsers(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Application
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into `owners` (`app_id`, `user_id`) values (?, ?)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into `owners` (`app_id`, `user_id`) values (?, ?)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	check, err := a.Users().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ID == b.ID {
			bFound = true
		}
		if v.ID == c.ID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ApplicationSlice{&a}
	if err = a.L.LoadUsers(ctx, tx, false, (*[]*Application)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Users); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Users = nil
	if err = a.L.LoadUsers(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Users); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testApplicationToManyAddOpBranches(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Application
	var b, c, d, e Branch

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, applicationDBTypes, false, strmangle.SetComplement(applicationPrimaryKeyColumns, applicationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Branch{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Branch{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddBranches(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ApplicationID {
			t.Error("foreign key was wrong value", a.ID, first.ApplicationID)
		}
		if a.ID != second.ApplicationID {
			t.Error("foreign key was wrong value", a.ID, second.ApplicationID)
		}

		if first.R.Application != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Application != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Branches[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Branches[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Branches().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testApplicationToManyAddOpUsers(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Application
	var b, c, d, e User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, applicationDBTypes, false, strmangle.SetComplement(applicationPrimaryKeyColumns, applicationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*User{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*User{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddUsers(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.AppApplications[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.AppApplications[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Users[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Users[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Users().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testApplicationToManySetOpUsers(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Application
	var b, c, d, e User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, applicationDBTypes, false, strmangle.SetComplement(applicationPrimaryKeyColumns, applicationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*User{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
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

	err = a.SetUsers(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Users().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetUsers(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Users().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	// The following checks cannot be implemented since we have no handle
	// to these when we call Set(). Leaving them here as wishful thinking
	// and to let people know there's dragons.
	//
	// if len(b.R.AppApplications) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	// if len(c.R.AppApplications) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	if d.R.AppApplications[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.AppApplications[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Users[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Users[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testApplicationToManyRemoveOpUsers(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Application
	var b, c, d, e User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, applicationDBTypes, false, strmangle.SetComplement(applicationPrimaryKeyColumns, applicationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*User{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddUsers(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Users().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveUsers(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Users().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.AppApplications) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.AppApplications) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.AppApplications[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.AppApplications[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Users) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Users[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Users[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testApplicationToOneRepositoryUsingRepository(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Application
	var foreign Repository

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, applicationDBTypes, false, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, repositoryDBTypes, false, repositoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Repository struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.RepositoryID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Repository().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ApplicationSlice{&local}
	if err = local.L.LoadRepository(ctx, tx, false, (*[]*Application)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Repository == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Repository = nil
	if err = local.L.LoadRepository(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Repository == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testApplicationToOneSetOpRepositoryUsingRepository(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Application
	var b, c Repository

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, applicationDBTypes, false, strmangle.SetComplement(applicationPrimaryKeyColumns, applicationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, repositoryDBTypes, false, strmangle.SetComplement(repositoryPrimaryKeyColumns, repositoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, repositoryDBTypes, false, strmangle.SetComplement(repositoryPrimaryKeyColumns, repositoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Repository{&b, &c} {
		err = a.SetRepository(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Repository != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Applications[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.RepositoryID != x.ID {
			t.Error("foreign key was wrong value", a.RepositoryID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.RepositoryID))
		reflect.Indirect(reflect.ValueOf(&a.RepositoryID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.RepositoryID != x.ID {
			t.Error("foreign key was wrong value", a.RepositoryID, x.ID)
		}
	}
}

func testApplicationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
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

func testApplicationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ApplicationSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testApplicationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Applications().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	applicationDBTypes = map[string]string{`ID`: `varchar`, `RepositoryID`: `varchar`, `CreatedAt`: `datetime`, `UpdatedAt`: `datetime`, `DeletedAt`: `datetime`}
	_                  = bytes.MinRead
)

func testApplicationsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(applicationPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(applicationAllColumns) == len(applicationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testApplicationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(applicationAllColumns) == len(applicationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Application{}
	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, applicationDBTypes, true, applicationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(applicationAllColumns, applicationPrimaryKeyColumns) {
		fields = applicationAllColumns
	} else {
		fields = strmangle.SetComplement(
			applicationAllColumns,
			applicationPrimaryKeyColumns,
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

	slice := ApplicationSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testApplicationsUpsert(t *testing.T) {
	t.Parallel()

	if len(applicationAllColumns) == len(applicationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLApplicationUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Application{}
	if err = randomize.Struct(seed, &o, applicationDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Application: %s", err)
	}

	count, err := Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, applicationDBTypes, false, applicationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Application struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Application: %s", err)
	}

	count, err = Applications().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
