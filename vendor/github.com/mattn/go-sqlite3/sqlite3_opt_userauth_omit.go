// Copyright (C) 2018 G.J.R. Timmer <gjr.timmer@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build !sqlite_userauth

package sqlite3

import (
	"C"
)

// Authenticate will perform an authentication of the provided username
// and password against the database.
//
// If a database contains the SQLITE_USER table, then the
// call to Authenticate must be invoked with an
// appropriate username and password prior to enable read and write
//access to the database.
//
// Return SQLITE_OK on success or SQLITE_ERROR if the username/password
// combination is incorrect or unknown.
//
// If the SQLITE_USER table is not present in the database file, then
// this interface is a harmless no-op returnning SQLITE_OK.
func (c *SQLiteConn) Authenticate(username, password string) error {
	// NOOP
	return nil
}

// authenticate provides the actual authentication to SQLite.
// This is not exported for usage in Go.
// It is however exported for usage within SQL by the user.
//
// Returns:
//	C.SQLITE_OK (0)
//	C.SQLITE_ERROR (1)
//  C.SQLITE_AUTH (23)
func (c *SQLiteConn) authenticate(username, password string) int {
	// NOOP
	return 0
}

// AuthUserAdd can be used (by an admin user only)
// to create a new user.  When called on a no-authentication-required
// database, this routine converts the database into an authentication-
// required database, automatically makes the added user an
// administrator, and logs in the current connection as that user.
// The AuthUserAdd only works for the "main" database, not
// for any ATTACH-ed databases. Any call to AuthUserAdd by a
// non-admin user results in an error.
func (c *SQLiteConn) AuthUserAdd(username, password string, admin bool) error {
	// NOOP
	return nil
}

// authUserAdd enables the User Authentication if not enabled.
// Otherwise it will add a user.
//
// When user authentication is already enabled then this function
// can only be called by an admin.
//
// This is not exported for usage in Go.
// It is however exported for usage within SQL by the user.
//
// Returns:
//	C.SQLITE_OK (0)
//	C.SQLITE_ERROR (1)
//  C.SQLITE_AUTH (23)
func (c *SQLiteConn) authUserAdd(username, password string, admin int) int {
	// NOOP
	return 0
}

// AuthUserChange can be used to change a users
// login credentials or admin privilege.  Any user can change their own
// login credentials.  Only an admin user can change another users login
// credentials or admin privilege setting.  No user may change their own
// admin privilege setting.
func (c *SQLiteConn) AuthUserChange(username, password string, admin bool) error {
	// NOOP
	return nil
}

// authUserChange allows to modify a user.
// Users can change their own password.
//
// Only admins can change passwords for other users
// and modify the admin flag.
//
// The admin flag of the current logged in user cannot be changed.
// THis ensures that their is always an admin.
//
// This is not exported for usage in Go.
// It is however exported for usage within SQL by the user.
//
// Returns:
//	C.SQLITE_OK (0)
//	C.SQLITE_ERROR (1)
//  C.SQLITE_AUTH (23)
func (c *SQLiteConn) authUserChange(username, password string, admin int) int {
	// NOOP
	return 0
}

// AuthUserDelete can be used (by an admin user only)
// to delete a user.  The currently logged-in user cannot be deleted,
// which guarantees that there is always an admin user and hence that
// the database cannot be converted into a no-authentication-required
// database.
func (c *SQLiteConn) AuthUserDelete(username string) error {
	// NOOP
	return nil
}

// authUserDelete can be used to delete a user.
//
// This function can only be executed by an admin.
//
// This is not exported for usage in Go.
// It is however exported for usage within SQL by the user.
//
// Returns:
//	C.SQLITE_OK (0)
//	C.SQLITE_ERROR (1)
//  C.SQLITE_AUTH (23)
func (c *SQLiteConn) authUserDelete(username string) int {
	// NOOP
	return 0
}

// AuthEnabled checks if the database is protected by user authentication
func (c *SQLiteConn) AuthEnabled() (exists bool) {
	// NOOP
	return false
}

// authEnabled perform the actual check for user authentication.
//
// This is not exported for usage in Go.
// It is however exported for usage within SQL by the user.
//
// Returns:
//	0 - Disabled
//  1 - Enabled
func (c *SQLiteConn) authEnabled() int {
	// NOOP
	return 0
}

// EOF
