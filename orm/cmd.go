package orm

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type commander interface {
	Parse([]string)
	Run()
}

var (
	commands = make(map[string]commander)
)

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

func RunCommand() {
	if len(os.Args) < 2 || os.Args[1] != "orm" {
		return
	}

	BootStrap()

	args := argString(os.Args[2:])
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

type commandSyncDb struct {
	al      *alias
	force   bool
	verbose bool
}

func (d *commandSyncDb) Parse(args []string) {
	var name string

	flagSet := flag.NewFlagSet("orm command: syncdb", flag.ExitOnError)
	flagSet.StringVar(&name, "db", "default", "DataBase alias name")
	flagSet.BoolVar(&d.force, "force", false, "drop tables before create")
	flagSet.BoolVar(&d.verbose, "v", false, "verbose info")
	flagSet.Parse(args)

	d.al = getDbAlias(name)
}

func (d *commandSyncDb) Run() {
	var drops []string
	if d.force {
		drops = getDbDropSql(d.al)
	}

	db := d.al.DB

	if d.force {
		for i, mi := range modelCache.allOrdered() {
			query := drops[i]
			_, err := db.Exec(query)
			result := ""
			if err != nil {
				result = err.Error()
			}
			fmt.Printf("drop table `%s` %s\n", mi.table, result)
			if d.verbose {
				fmt.Printf("    %s\n\n", query)
			}
		}
	}

	sqls, indexes := getDbCreateSql(d.al)

	for i, mi := range modelCache.allOrdered() {
		fmt.Printf("create table `%s` \n", mi.table)

		queries := []string{sqls[i]}
		queries = append(queries, indexes[mi.table]...)

		for _, query := range queries {
			_, err := db.Exec(query)
			if d.verbose {
				query = "    " + strings.Join(strings.Split(query, "\n"), "\n    ")
				fmt.Println(query)
			}
			if err != nil {
				fmt.Printf("    %s\n", err.Error())
			}
		}
		if d.verbose {
			fmt.Println("")
		}
	}
}

type commandSqlAll struct {
	al *alias
}

func (d *commandSqlAll) Parse(args []string) {
	var name string

	flagSet := flag.NewFlagSet("orm command: sqlall", flag.ExitOnError)
	flagSet.StringVar(&name, "db", "default", "DataBase alias name")
	flagSet.Parse(args)

	d.al = getDbAlias(name)
}

func (d *commandSqlAll) Run() {
	sqls, indexes := getDbCreateSql(d.al)
	var all []string
	for i, mi := range modelCache.allOrdered() {
		queries := []string{sqls[i]}
		queries = append(queries, indexes[mi.table]...)
		sql := strings.Join(queries, "\n")
		all = append(all, sql)
	}
	fmt.Println(strings.Join(all, "\n\n"))
}

func init() {
	commands["syncdb"] = new(commandSyncDb)
	commands["sqlall"] = new(commandSqlAll)
}
