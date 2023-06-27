// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/client/orm/internal/utils"

	"github.com/beego/beego/v2/client/orm/internal/models"
)

type commander interface {
	Parse([]string)
	Run() error
}

var commands = make(map[string]commander)

// print help.
func printHelp(errs ...string) {
	content := `orm command usage:

    syncdb     - auto create tables
    sqlall     - print sql of create tables
    help       - print this help
`

	if len(errs) > 0 {
		fmt.Println(errs[0])
	}
	fmt.Println(content)
	os.Exit(2)
}

// RunCommand listens for orm command and runs if command arguments have been passed.
func RunCommand() {
	if len(os.Args) < 2 || os.Args[1] != "orm" {
		return
	}

	BootStrap()

	args := utils.ArgString(os.Args[2:])
	name := args.Get(0)

	if name == "help" {
		printHelp()
	}

	if cmd, ok := commands[name]; ok {
		cmd.Parse(os.Args[3:])
		cmd.Run()
		os.Exit(0)
	} else {
		if name == "" {
			printHelp()
		} else {
			printHelp(fmt.Sprintf("unknown command %s", name))
		}
	}
}

// sync database struct command interface.
type commandSyncDb struct {
	al        *alias
	force     bool
	verbose   bool
	noInfo    bool
	rtOnError bool
}

// Parse the orm command line arguments.
func (d *commandSyncDb) Parse(args []string) {
	var name string

	flagSet := flag.NewFlagSet("orm command: syncdb", flag.ExitOnError)
	flagSet.StringVar(&name, "db", "default", "DataBase alias name")
	flagSet.BoolVar(&d.force, "force", false, "drop tables before create")
	flagSet.BoolVar(&d.verbose, "v", false, "verbose info")
	flagSet.Parse(args)

	d.al = getDbAlias(name)
}

// Run orm line command.
func (d *commandSyncDb) Run() error {
	var drops []string
	var err error
	if d.force {
		drops, err = defaultModelCache.getDbDropSQL(d.al)
		if err != nil {
			return err
		}
	}

	db := d.al.DB

	if d.force && len(drops) > 0 {
		for i, mi := range defaultModelCache.allOrdered() {
			query := drops[i]
			if !d.noInfo {
				fmt.Printf("drop table `%s`\n", mi.Table)
			}
			_, err := db.Exec(query)
			if d.verbose {
				fmt.Printf("    %s\n\n", query)
			}
			if err != nil {
				if d.rtOnError {
					return err
				}
				fmt.Printf("    %s\n", err.Error())
			}
		}
	}

	createQueries, indexes, err := defaultModelCache.getDbCreateSQL(d.al)
	if err != nil {
		return err
	}

	tables, err := d.al.DbBaser.GetTables(db)
	if err != nil {
		if d.rtOnError {
			return err
		}
		fmt.Printf("    %s\n", err.Error())
	}

	ctx := context.Background()
	for i, mi := range defaultModelCache.allOrdered() {

		if !models.IsApplicableTableForDB(mi.AddrField, d.al.Name) {
			fmt.Printf("table `%s` is not applicable to database '%s'\n", mi.Table, d.al.Name)
			continue
		}

		if tables[mi.Table] {
			if !d.noInfo {
				fmt.Printf("table `%s` already exists, skip\n", mi.Table)
			}

			var fields []*models.FieldInfo
			columns, err := d.al.DbBaser.GetColumns(ctx, db, mi.Table)
			if err != nil {
				if d.rtOnError {
					return err
				}
				fmt.Printf("    %s\n", err.Error())
			}

			for _, fi := range mi.Fields.FieldsDB {
				if _, ok := columns[fi.Column]; !ok {
					fields = append(fields, fi)
				}
			}

			for _, fi := range fields {
				query := getColumnAddQuery(d.al, fi)

				if !d.noInfo {
					fmt.Printf("add column `%s` for table `%s`\n", fi.FullName, mi.Table)
				}

				_, err := db.Exec(query)
				if d.verbose {
					fmt.Printf("    %s\n", query)
				}
				if err != nil {
					if d.rtOnError {
						return err
					}
					fmt.Printf("    %s\n", err.Error())
				}
			}

			for _, idx := range indexes[mi.Table] {
				if !d.al.DbBaser.IndexExists(ctx, db, idx.Table, idx.Name) {
					if !d.noInfo {
						fmt.Printf("create index `%s` for table `%s`\n", idx.Name, idx.Table)
					}

					query := idx.SQL
					_, err := db.Exec(query)
					if d.verbose {
						fmt.Printf("    %s\n", query)
					}
					if err != nil {
						if d.rtOnError {
							return err
						}
						fmt.Printf("    %s\n", err.Error())
					}
				}
			}

			continue
		}

		if !d.noInfo {
			fmt.Printf("create table `%s` \n", mi.Table)
		}

		queries := []string{createQueries[i]}
		for _, idx := range indexes[mi.Table] {
			queries = append(queries, idx.SQL)
		}

		for _, query := range queries {
			_, err := db.Exec(query)
			if d.verbose {
				query = "    " + strings.Join(strings.Split(query, "\n"), "\n    ")
				fmt.Println(query)
			}
			if err != nil {
				if d.rtOnError {
					return err
				}
				fmt.Printf("    %s\n", err.Error())
			}
		}
		if d.verbose {
			fmt.Println("")
		}
	}

	return nil
}

// database creation commander interface implement.
type commandSQLAll struct {
	al *alias
}

// Parse orm command line arguments.
func (d *commandSQLAll) Parse(args []string) {
	var name string

	flagSet := flag.NewFlagSet("orm command: sqlall", flag.ExitOnError)
	flagSet.StringVar(&name, "db", "default", "DataBase alias name")
	flagSet.Parse(args)

	d.al = getDbAlias(name)
}

// Run orm line command.
func (d *commandSQLAll) Run() error {
	createQueries, indexes, err := defaultModelCache.getDbCreateSQL(d.al)
	if err != nil {
		return err
	}
	var all []string
	for i, mi := range defaultModelCache.allOrdered() {
		queries := []string{createQueries[i]}
		for _, idx := range indexes[mi.Table] {
			queries = append(queries, idx.SQL)
		}
		sql := strings.Join(queries, "\n")
		all = append(all, sql)
	}
	fmt.Println(strings.Join(all, "\n\n"))

	return nil
}

func init() {
	commands["syncdb"] = new(commandSyncDb)
	commands["sqlall"] = new(commandSQLAll)
}

// RunSyncdb run syncdb command line.
// name: Table's alias name (default is "default")
// force: Run the next sql command even if the current gave an error
// verbose: Print all information, useful for debugging
func RunSyncdb(name string, force bool, verbose bool) error {
	BootStrap()

	al := getDbAlias(name)
	cmd := new(commandSyncDb)
	cmd.al = al
	cmd.force = force
	cmd.noInfo = !verbose
	cmd.verbose = verbose
	cmd.rtOnError = true
	return cmd.Run()
}
