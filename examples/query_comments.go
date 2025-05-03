package main

import (
	"fmt"
	"os"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"
)

// User represents a user in the system
type User struct {
	Id      int    `orm:"auto"`
	Name    string `orm:"size(100)"`
	Email   string `orm:"size(100)"`
	Profile string `orm:"type(text);null"`
}

func init() {
	// Enable ORM debug mode to see SQL queries
	orm.Debug = true

	// Register the model with ORM
	orm.RegisterModel(new(User))

	// Register SQLite driver
	err := orm.RegisterDriver("sqlite3", orm.DRSqlite)
	if err != nil {
		fmt.Printf("Failed to register driver: %v\n", err)
		os.Exit(1)
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Printf("Failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	// Setup the database
	dataSource := "data/test.db"
	err = orm.RegisterDataBase("default", "sqlite3", dataSource)
	if err != nil {
		fmt.Printf("Failed to register database: %v\n", err)
		os.Exit(1)
	}

	// Create tables
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Printf("Failed to create tables: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	o := orm.NewOrm()

	// Create some test users
	users := []User{
		{Name: "Alice", Email: "alice@example.com", Profile: "Developer"},
		{Name: "Bob", Email: "bob@example.com", Profile: "Manager"},
		{Name: "Charlie", Email: "charlie@example.com", Profile: "Designer"},
	}

	fmt.Println("\n=== Inserting users with comments ===")
	orm.AddQueryComment("Batch user creation")
	orm.AddQueryComment("Initial data load")
	for _, user := range users {
		id, err := o.Insert(&user)
		if err != nil {
			fmt.Printf("Failed to insert user: %v\n", err)
			continue
		}
		fmt.Printf("Inserted user with ID: %d\n", id)
	}
	orm.ClearQueryComments()

	fmt.Println("\n=== Querying users with comments ===")
	var fetchedUsers []User
	orm.AddQueryComment("User listing query")
	orm.AddQueryComment("Backend API request")
	_, err := o.QueryTable("user").All(&fetchedUsers)
	if err != nil {
		fmt.Printf("Failed to query users: %v\n", err)
		return
	}
	for _, user := range fetchedUsers {
		fmt.Printf("User: %s (Email: %s)\n", user.Name, user.Email)
	}
	orm.ClearQueryComments()

	fmt.Println("\n=== Updating user with comments ===")
	orm.AddQueryComment("Profile update")
	orm.AddQueryComment("User requested change")
	alice := User{Id: 1, Profile: "Senior Developer"}
	if num, err := o.Update(&alice, "Profile"); err == nil {
		fmt.Printf("Updated %d user(s)\n", num)
	}
	orm.ClearQueryComments()

	fmt.Println("\n=== Deleting user with comments ===")
	orm.AddQueryComment("User deletion")
	orm.AddQueryComment("Account closure request")
	charlie := User{Id: 3}
	if num, err := o.Delete(&charlie); err == nil {
		fmt.Printf("Deleted %d user(s)\n", num)
	}
	orm.ClearQueryComments()

	fmt.Println("\n=== Complex query with comments ===")
	var developers []User
	orm.AddQueryComment("Developer search")
	orm.AddQueryComment("HR department request")
	_, err = o.QueryTable("user").
		Filter("Profile__icontains", "developer").
		OrderBy("Name").
		All(&developers)
	if err != nil {
		fmt.Printf("Failed to query developers: %v\n", err)
		return
	}
	for _, dev := range developers {
		fmt.Printf("Developer: %s (Profile: %s)\n", dev.Name, dev.Profile)
	}
}
