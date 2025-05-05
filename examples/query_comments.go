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
	o.AddQueryComment("Batch user creation") // Use instance method
	o.AddQueryComment("Initial data load")   // Use instance method
	for _, user := range users {
		// Note: Comments added before the loop apply to each Insert within the loop
		id, err := o.Insert(&user)
		if err != nil {
			fmt.Printf("Failed to insert user: %v\n", err)
			continue
		}
		fmt.Printf("Inserted user with ID: %d\n", id)
	}
	o.ClearQueryComments() // Use instance method

	fmt.Println("\n=== Querying users with comments ===")
	var fetchedUsers []User
	o.AddQueryComment("User listing query")  // Use instance method
	o.AddQueryComment("Backend API request") // Use instance method
	_, err := o.QueryTable("user").All(&fetchedUsers)
	if err != nil {
		fmt.Printf("Failed to query users: %v\n", err)
		return
	}
	for _, user := range fetchedUsers {
		fmt.Printf("User: %s (Email: %s)\n", user.Name, user.Email)
	}
	o.ClearQueryComments() // Use instance method

	fmt.Println("\n=== Updating user with comments ===")
	o.AddQueryComment("Profile update")        // Use instance method
	o.AddQueryComment("User requested change") // Use instance method
	alice := User{Id: 1, Profile: "Senior Developer"}
	if num, err := o.Update(&alice, "Profile"); err == nil {
		fmt.Printf("Updated %d user(s)\n", num)
	}
	o.ClearQueryComments() // Use instance method

	fmt.Println("\n=== Deleting user with comments ===")
	o.AddQueryComment("User deletion")           // Use instance method
	o.AddQueryComment("Account closure request") // Use instance method
	charlie := User{Id: 3}
	if num, err := o.Delete(&charlie); err == nil {
		fmt.Printf("Deleted %d user(s)\n", num)
	}
	o.ClearQueryComments() // Use instance method

	fmt.Println("\n=== Complex query with comments ===")
	var developers []User
	o.AddQueryComment("Developer search")      // Use instance method
	o.AddQueryComment("HR department request") // Use instance method
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
	o2 := orm.NewOrm()
	_, err = o2.QueryTable("user").
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
