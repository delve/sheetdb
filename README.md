# sheetdb

[![CircleCI](https://circleci.com/gh/takuoki/sheetdb/tree/master.svg?style=shield&circle-token=9bbc178fd927c6b27f6d726ffd66e3d5deb06fcc)](https://circleci.com/gh/takuoki/sheetdb/tree/master)
[![GoDoc](https://godoc.org/github.com/takuoki/sheetdb?status.svg)](https://godoc.org/github.com/takuoki/sheetdb)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

A golang package for using Google spreadsheets as a database instead of the actual database management system.

**!!! Caution !!!**

* Currently we are not measuring performance. It is intended for use in small applications where performance is not an issue.
* The Google Sheets API has a [usage limit](https://developers.google.com/sheets/api/limits). Do not use this package for applications that require access beyond this usage limit.

---

<!-- vscode-markdown-toc -->
* [Features](#Features)
* [Installation](#Installation)
* [Requirement](#Requirement)
* [Usage](#Usage)
	* [1. Write models](#Writemodels)
	* [2. Generate codes from models](#Generatecodesfrommodels)
	* [3. Set up Google spreadsheet](#SetupGooglespreadsheet)
	* [4. Initialize models](#Initializemodels)
		* [Set up nortificaltion](#Setupnortificaltion)
		* [Create new client](#Createnewclient)
		* [Load sheet data](#Loadsheetdata)
	* [5. Use CRUD functions](#UseCRUDfunctions)
		* [Read (Get/Select)](#ReadGetSelect)
		* [Create (Add/Insert)](#CreateAddInsert)
		* [Update](#Update)
		* [Delete](#Delete)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

## <a name='Features'></a>Features

* Load sheet data into cache
* Apply cache update information to sheet asynchronously
* Exclusive control when updating cache and inserting a row to a sheet
* Automatic generation of CRUD functions based on model (structure definition)
* Automatic numbering and initial value setting of ID
* Unique and non-null constraints
* Cascade delete child data when deleting parent data
* Notification mechanism on asynchronous update error

The following features are not included.

* SQL query
* Transaction control (commit and rollback)
* Read Lock for Update

## <a name='Installation'></a>Installation

```bash
go get github.com/takuoki/sheetdb
```

## <a name='Requirement'></a>Requirement

This package uses Google OAuth2.0. So before executing tool, you have to prepare credentials.json.
See [Go Quickstart](https://developers.google.com/sheets/api/quickstart/go), or [blog post (Japanese)](https://medium.com/veltra-engineering/how-to-use-google-sheets-api-with-golang-9e50ee9e0abc) for the details.

## <a name='Usage'></a>Usage

### <a name='Writemodels'></a>1. Write models

Write the structure of the model as follows.
Please refer to [the sheetdb-modeler tool documentation](tools/sheetdb-modeler) for details.

```go
//go:generate sheetdb-modeler -type=User -children=Foo,FooChild,Bar -initial=10001

// User is a struct of user.
type User struct {
  UserID   int           `json:"user_id" db:"primarykey"`
  Name     string        `json:"name"`
  Email    string        `json:"email" db:"unique"`
  Sex      Sex           `json:"sex"`
  Birthday *sheetdb.Date `json:"birthday"`
}

//go:generate sheetdb-modeler -type=Foo -parent=User -children=FooChild

// Foo is a struct of foo which is a child of user.
type Foo struct {
  UserID int     `json:"user_id" db:"primarykey"`
  FooID  int     `json:"foo_id" db:"primarykey"`
  Value  float32 `json:"value"`
  Note   string  `json:"note" db:"allowempty"`
}

//go:generate sheetdb-modeler -type=FooChild -parent=Foo

// FooChild is a struct of foo child.
type FooChild struct {
  UserID  int    `json:"user_id" db:"primarykey"`
  FooID   int    `json:"foo_id" db:"primarykey"`
  ChildID int    `json:"child_id" db:"primarykey"`
  Value   string `json:"value" db:"unique"`
}

//go:generate sheetdb-modeler -type=Bar -parent=User

// Bar is a struct of bar which is a child of user.
type Bar struct {
  UserID   int              `json:"user_id" db:"primarykey"`
  Datetime sheetdb.Datetime `json:"datetime" db:"primarykey"`
  Value    float32          `json:"value"`
  Note     string           `json:"note" db:"allowempty"`
}
```

### <a name='Generatecodesfrommodels'></a>2. Generate codes from models

You can generate in bulk with the `go generate` command by putting `//go:generate sheetdb-modeler` comments in the code of the target package.
Please refer to [the sheetdb-modeler tool documentation](tools/sheetdb-modeler#Howtogeneratemodels) for details.

```bash
go generate ./sample
```

### <a name='SetupGooglespreadsheet'></a>3. Set up Google spreadsheet

Prepare a spreadsheet according to the header comments of each generated file.

```go
// Code generated by "sheetdb-modeler"; DO NOT EDIT.
// Create a Spreadsheet (sheet name: "users") as data storage.
// The spreadsheet header is as follows:
//   user_id | name | email | sex | birthday | updated_at | deleted_at
// Please copy and paste this header on the first line of the sheet.
```

### <a name='Initializemodels'></a>4. Initialize models

#### <a name='Setupnortificaltion'></a>Set up nortificaltion

By setting Logger, you can customize log output freely.
If you want to send an alert to Slack, you can use `SlackLogger` in this package.
Configuration of Logger is package scope. Please note that it is not client scope.

```go
sheetdb.SetLogger(sheetdb.NewSlackLogger(
  "project name",
  "service name",
  "icon_emoji",
  "https://hooks.slack.com/services/zzz/zzz/zzzzz",
  sheetdb.LevelError,
))
```

#### <a name='Createnewclient'></a>Create new client

Create a new client using the `New` function.
Set the created client in the package global variable. The name of the variable is the name specified with the `-client` option (the default is `dbClient`).

```go
var dbClient *sheetdb.Client

// Initialize initializes this package.
func Initialize(ctx context.Context) error {
  client, err := sheetdb.New(
    ctx,
    `{"installed":{"client_id":"..."}`, // Google API credentials
    `{"access_token":"..."`,            // Google API token
    "xxxxx",                            // Google spreadsheet ID
  )
  // ...
  dbClient = client
  return nil
}
```

#### <a name='Loadsheetdata'></a>Load sheet data

```go
err := client.LoadData(ctx)
```

### <a name='UseCRUDfunctions'></a>5. Use CRUD functions

The functions in this section are generated automatically by [sheetdb-modeler](tools/sheetdb-modeler).

#### <a name='ReadGetSelect'></a>Read (Get/Select)

`GetModelName` function returns an instance of that model by the primary key(s).
If it can not be found, this function returns `*sheetdb.NotFoundError`.

```go
user, err := GetUser(userID)
foo, err := GetFoo(userID, fooID)
fooChild, err := GetFooChild(userID, fooID, childID)
bar, err := GetBar(userID, datetime)
```

If the model is a child model of another model, `GetModelName` method is also added to the parent model.

```go
foo, err := user.GetFoo(fooID)
fooChild, err := foo.GetFooChild(childID)
bar, err := user.GetBar(datetime)
```

For fields defined as unique, `GetModelNameByFieldName` function is also generated.
For child models, this method is also added to the parent model.

```go
user, err := GetUserByEmail(email)
fooChild, err := GetFooChildByValue(userID, fooID, value)

fooChild, err := foo.GetFooChildByValue(value)
```

`GetModelNames` function returns all instances of that model.
If the model is a child model of another model, this function returns all instances that parent model has.
For child models, this method is also added to the parent model.

```go
users, err := GetUsers()
foos, err := GetFoos(userID)
fooChildren, err := GetFooChildren(userID, fooID)
bars, err := GetBars(userID)

foos, err := user.GetFoos()
fooChildren, err := foo.GetFooChildren()
bars, err := user.GetBars()
```

You can get filtered results using `ModelFilter` option.

```go
users, err := GetUsers(sample.UserFilter(func(user *sample.User) bool {
  return user.Sex == sample.Male
}))
```

You can also get sorted results using `ModelSort` option.
If the sort option is not specified, the order of results is random.

```go
users, err := GetUsers(sample.UserSort(func(users []*sample.User) {
  sort.Slice(users, func(i, j int) bool {
    return users[i].UserID < users[j].UserID
  })
}))
```

#### <a name='CreateAddInsert'></a>Create (Add/Insert)

`AddModel` adds a new instance.
If the primary key matches [the automatic numbering rule](tools/sheetdb-modeler#AutonumberingofID), it will be automatically numbered.
For child models, `AddModel` method is added to the parent model.

```go
user, err := AddUser(name, email, sex, birthday)
foo, err := AddFoo(userID, value, note)
fooChild, err := AddFooChild(userID, fooID, value)
bar, err := AddBar(userID, datetime, value, note)

foo, err := user.AddFoo(value, note)
fooChild, err := foo.AddFooChild(value)
bar, err := user.AddBar(datetime, value, note)
```

#### <a name='Update'></a>Update

`UpdateModel` updates an existing instance.
You can not change the primary key(s).
For child models, `UpdateModel` method is added to the parent model.

```go
user, err := UpdateUser(userID, name, email, sex, birthday)
foo, err := UpdateFoo(userID, fooID, value, note)
fooChild, err := UpdateFooChild(userID, fooID, childID, value)
bar, err := UpdateBar(userID, datetime, value, note)

foo, err := user.UpdateFoo(fooID, value, note)
fooChild, err := foo.UpdateFooChild(childID, value, note)
bar, err := user.UpdateBar(datetime, value, note)
```

#### <a name='Delete'></a>Delete

`DeleteModel` deletes an existing instance.
If the model has child models, this function deletes the instances of the child model in cascade.
For child models, `DeleteModel` method is added to the parent model.

```go
err = DeleteUser(userID)
err = DeleteFoo(userID, fooID)
err = DeleteFooChild(userID, fooID, childID)
err = DeleteBar(userID, datetime)

err = user.DeleteFoo(fooID)
err = foo.DeleteFooChild(childID)
err = user.DeleteBar(datetime)
```
