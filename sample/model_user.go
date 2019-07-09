// Code generated by "sheetdb-modeler"; DO NOT EDIT.
// Create a Spreadsheet (sheet name: "users") as data storage.
// The spreadsheet header is as follows:
//   user_id | name | email | sex | birthday | updated_at | deleted_at
// Please copy and paste this header on the first line of the sheet.

package sample

import (
	"strconv"
	"sync"
	"time"

	"github.com/takuoki/gsheets"
	"github.com/takuoki/sheetdb"
)

const (
	// Sheet definition
	_User_sheetName        = "users"
	_User_column_UserID    = 0 // A
	_User_column_Name      = 1 // B
	_User_column_Email     = 2 // C
	_User_column_Sex       = 3 // D
	_User_column_Birthday  = 4 // E
	_User_column_UpdatedAt = 5 // F
	_User_column_DeletedAt = 6 // G

	// Parent children relation for compile check
	_User_child_Foo           = 0
	_User_child_Bar           = 0
	_User_numOfChildren       = 3
	_User_numOfDirectChildren = 2
)

var (
	_User_mutex           = sync.RWMutex{}
	_User_cache           = map[int]*User{} // map[userID]*User
	_User_rowNoMap        = map[int]int{}   // map[userID]rowNo
	_User_maxRowNo        = 0
	_User_Email_uniqueMap = map[string]*User{}
)

func _() {
	// An "undeclared name" compiler error signifies that parent and child model set names are different.
	// Make sure that all models in parent-child relationship exist in the same model set (spreadsheet) and try again.
	_ = _Foo_modelSetName_default
	_ = _Bar_modelSetName_default
	// An "undeclared name" compiler error signifies that parent-children option conflicts between models.
	// Make sure that the parent-children options are correct for all relevant models and try again.
	_ = _Foo_parent_User
	_ = _Bar_parent_User
	// An "invalid array index" compiler error signifies that the children option is incorrect.
	// Make sure that all child models are specified, including not only the direct child model
	// but also the grandchild model, and try again.
	var x [1]struct{}
	_ = x[_User_numOfChildren-_User_numOfDirectChildren-_Foo_numOfChildren-_Bar_numOfChildren]
}

func init() {
	sheetdb.RegisterModel("default", "User", _User_sheetName, _User_load)
}

func _User_load(data *gsheets.Sheet) error {

	_User_mutex.Lock()
	defer _User_mutex.Unlock()

	_User_cache = map[int]*User{}
	_User_rowNoMap = map[int]int{}
	_User_maxRowNo = 0
	_User_Email_uniqueMap = map[string]*User{}

	for i, r := range data.Rows() {
		if i == 0 {
			continue
		}
		if r.Value(_User_column_DeletedAt) != "" {
			_User_maxRowNo++
			continue
		}
		if r.Value(_User_column_UserID) == "" {
			break
		}

		userID, err := _User_parseUserID(r.Value(_User_column_UserID))
		if err != nil {
			return err
		}
		name := r.Value(_User_column_Name)
		if err := _User_validateName(name); err != nil {
			return err
		}
		email := r.Value(_User_column_Email)
		if err := _User_validateEmail(email, nil); err != nil {
			return err
		}
		sex, err := _User_parseSex(r.Value(_User_column_Sex))
		if err != nil {
			return err
		}
		birthday, err := _User_parseBirthday(r.Value(_User_column_Birthday))
		if err != nil {
			return err
		}

		if _, ok := _User_cache[userID]; ok {
			return &sheetdb.DuplicationError{FieldName: "UserID"}
		}

		user := User{
			UserID:   userID,
			Name:     name,
			Email:    email,
			Sex:      sex,
			Birthday: birthday,
		}

		_User_maxRowNo++
		_User_cache[user.UserID] = &user
		_User_rowNoMap[user.UserID] = _User_maxRowNo
		_User_Email_uniqueMap[user.Email] = &user
	}

	return nil
}

// GetUser returns a user by UserID.
// If it can not be found, this function returns *sheetdb.NotFoundError.
func GetUser(userID int) (*User, error) {
	_User_mutex.RLock()
	defer _User_mutex.RUnlock()
	if v, ok := _User_cache[userID]; ok {
		return v, nil
	}
	return nil, &sheetdb.NotFoundError{Model: "User"}
}

// GetUserByEmail returns a user by Email.
// If it can not be found, this function returns *sheetdb.NotFoundError.
func GetUserByEmail(email string) (*User, error) {
	_User_mutex.RLock()
	defer _User_mutex.RUnlock()
	if v, ok := _User_Email_uniqueMap[email]; ok {
		return v, nil
	}
	return nil, &sheetdb.NotFoundError{Model: "User"}
}

// UserQuery is used for selecting users.
type UserQuery struct {
	filter func(user *User) bool
	sort   func(users []*User)
}

// UserQueryOption is an option to change the behavior of UserQuery.
type UserQueryOption func(query *UserQuery) *UserQuery

// UserFilter is an option to change the filtering behavior of UserQuery.
func UserFilter(filterFunc func(user *User) bool) func(query *UserQuery) *UserQuery {
	return func(query *UserQuery) *UserQuery {
		if query != nil {
			query.filter = filterFunc
		}
		return query
	}
}

// UserSort is an option to change the sorting behavior of UserQuery.
func UserSort(sortFunc func(users []*User)) func(query *UserQuery) *UserQuery {
	return func(query *UserQuery) *UserQuery {
		if query != nil {
			query.sort = sortFunc
		}
		return query
	}
}

// GetUsers returns all users.
// If any options are specified, the result according to the specified option is returned.
// If there are no user to return, this function returns an nil array.
// If the sort option is not specified, the order of users is random.
func GetUsers(opts ...UserQueryOption) ([]*User, error) {
	userQuery := &UserQuery{}
	for _, opt := range opts {
		userQuery = opt(userQuery)
	}
	_User_mutex.RLock()
	defer _User_mutex.RUnlock()
	var users []*User
	if userQuery.filter != nil {
		for _, v := range _User_cache {
			if userQuery.filter(v) {
				users = append(users, v)
			}
		}
	} else {
		for _, v := range _User_cache {
			users = append(users, v)
		}
	}
	if userQuery.sort != nil {
		userQuery.sort(users)
	}
	return users, nil
}

// AddUser adds new user.
// UserID is generated automatically.
// If any fields are invalid, this function returns error.
func AddUser(name string, email string, sex Sex, birthday *sheetdb.Date) (*User, error) {
	_User_mutex.Lock()
	defer _User_mutex.Unlock()
	if err := _User_validateName(name); err != nil {
		return nil, err
	}
	if err := _User_validateEmail(email, nil); err != nil {
		return nil, err
	}
	user := &User{
		UserID:   _User_maxRowNo + 10001,
		Name:     name,
		Email:    email,
		Sex:      sex,
		Birthday: birthday,
	}
	if err := user._asyncAdd(_User_maxRowNo + 1); err != nil {
		return nil, err
	}
	_User_maxRowNo++
	_User_cache[user.UserID] = user
	_User_rowNoMap[user.UserID] = _User_maxRowNo
	_User_Email_uniqueMap[user.Email] = user
	return user, nil
}

// UpdateUser updates user.
// If it can not be found, this function returns *sheetdb.NotFoundError.
// If any fields are invalid, this function returns error.
func UpdateUser(userID int, name string, email string, sex Sex, birthday *sheetdb.Date) (*User, error) {
	_User_mutex.Lock()
	defer _User_mutex.Unlock()
	user, ok := _User_cache[userID]
	if !ok {
		return nil, &sheetdb.NotFoundError{Model: "User"}
	}
	if err := _User_validateName(name); err != nil {
		return nil, err
	}
	if err := _User_validateEmail(email, &user.Email); err != nil {
		return nil, err
	}
	userCopy := *user
	userCopy.Name = name
	userCopy.Email = email
	userCopy.Sex = sex
	userCopy.Birthday = birthday
	if err := (&userCopy)._asyncUpdate(); err != nil {
		return nil, err
	}
	if userCopy.Email != user.Email {
		delete(_User_Email_uniqueMap, user.Email)
	}
	*user = userCopy
	_User_Email_uniqueMap[userCopy.Email] = &userCopy
	return user, nil
}

// DeleteUser deletes user and it's children foo, fooChild and bar.
// If it can not be found, this function returns *sheetdb.NotFoundError.
func DeleteUser(userID int) error {
	_User_mutex.Lock()
	defer _User_mutex.Unlock()
	_Foo_mutex.Lock()
	defer _Foo_mutex.Unlock()
	_FooChild_mutex.Lock()
	defer _FooChild_mutex.Unlock()
	_Bar_mutex.Lock()
	defer _Bar_mutex.Unlock()
	user, ok := _User_cache[userID]
	if !ok {
		return &sheetdb.NotFoundError{Model: "User"}
	}
	var foos []*Foo
	for _, v := range _Foo_cache[userID] {
		foos = append(foos, v)
	}
	var fooChildren []*FooChild
	for _, v := range _FooChild_cache[userID] {
		for _, v := range v {
			fooChildren = append(fooChildren, v)
		}
	}
	var bars []*Bar
	for _, v := range _Bar_cache[userID] {
		bars = append(bars, v)
	}
	if err := user._asyncDelete(foos, fooChildren, bars); err != nil {
		return err
	}
	delete(_User_cache, userID)
	delete(_User_Email_uniqueMap, user.Email)
	delete(_Foo_cache, userID)
	delete(_FooChild_cache, userID)
	for _, v := range fooChildren {
		delete(_FooChild_Value_uniqueMap, v.Value)
	}
	delete(_Bar_cache, userID)
	return nil
}

func _User_validateName(name string) error {
	if name == "" {
		return &sheetdb.EmptyStringError{FieldName: "Name"}
	}
	return nil
}

func _User_validateEmail(email string, oldEmail *string) error {
	if email == "" {
		return &sheetdb.EmptyStringError{FieldName: "Email"}
	}
	if oldEmail == nil || *oldEmail != email {
		if _, ok := _User_Email_uniqueMap[email]; ok {
			return &sheetdb.DuplicationError{FieldName: "Email"}
		}
	}
	return nil
}

func _User_parseUserID(userID string) (int, error) {
	v, err := strconv.Atoi(userID)
	if err != nil {
		return 0, &sheetdb.InvalidValueError{FieldName: "UserID", Err: err}
	}
	return v, nil
}

func _User_parseSex(sex string) (Sex, error) {
	v, err := NewSex(sex)
	if err != nil {
		return v, &sheetdb.InvalidValueError{FieldName: "Sex", Err: err}
	}
	return v, nil
}

func _User_parseBirthday(birthday string) (*sheetdb.Date, error) {
	var val *sheetdb.Date
	if birthday != "" {
		v, err := sheetdb.NewDate(birthday)
		if err != nil {
			return nil, &sheetdb.InvalidValueError{FieldName: "Birthday", Err: err}
		}
		val = &v
	}
	return val, nil
}

func (m *User) _asyncAdd(rowNo int) error {
	data := []gsheets.UpdateValue{
		{
			SheetName: _User_sheetName,
			RowNo:     rowNo,
			Values: []interface{}{
				m.UserID,
				m.Name,
				m.Email,
				m.Sex.String(),
				m.Birthday.String(),
				time.Now(),
				"",
			},
		},
	}
	return dbClient.AsyncUpdate(data)
}

func (m *User) _asyncUpdate() error {
	data := []gsheets.UpdateValue{
		{
			SheetName: _User_sheetName,
			RowNo:     _User_rowNoMap[m.UserID],
			Values: []interface{}{
				m.UserID,
				m.Name,
				m.Email,
				m.Sex.String(),
				m.Birthday.String(),
				time.Now(),
				"",
			},
		},
	}
	return dbClient.AsyncUpdate(data)
}

func (m *User) _asyncDelete(foos []*Foo, fooChildren []*FooChild, bars []*Bar) error {
	now := time.Now()
	data := []gsheets.UpdateValue{
		{
			SheetName: _User_sheetName,
			RowNo:     _User_rowNoMap[m.UserID],
			Values: []interface{}{
				m.UserID,
				m.Name,
				m.Email,
				m.Sex.String(),
				m.Birthday.String(),
				now,
				now,
			},
		},
	}
	for _, v := range foos {
		data = append(data, gsheets.UpdateValue{
			SheetName: _Foo_sheetName,
			RowNo:     _Foo_rowNoMap[v.UserID][v.FooID],
			Values: []interface{}{
				v.UserID,
				v.FooID,
				v.Value,
				v.Note,
				now,
				now,
			},
		})
	}
	for _, v := range fooChildren {
		data = append(data, gsheets.UpdateValue{
			SheetName: _FooChild_sheetName,
			RowNo:     _FooChild_rowNoMap[v.UserID][v.FooID][v.ChildID],
			Values: []interface{}{
				v.UserID,
				v.FooID,
				v.ChildID,
				v.Value,
				now,
				now,
			},
		})
	}
	for _, v := range bars {
		data = append(data, gsheets.UpdateValue{
			SheetName: _Bar_sheetName,
			RowNo:     _Bar_rowNoMap[v.UserID][v.Datetime],
			Values: []interface{}{
				v.UserID,
				v.Datetime.String(),
				v.Value,
				v.Note,
				now,
				now,
			},
		})
	}
	return dbClient.AsyncUpdate(data)
}
