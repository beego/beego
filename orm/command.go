package orm

import (
	"flag"
	"fmt"
	"os"
)

func printHelp() {

}

func getSqlAll() (sql string) {
	for _, mi := range modelCache.allOrdered() {
		_ = mi
	}
	return
}

func runCommand() {
	if len(os.Args) < 2 || os.Args[1] != "orm" {
		return
	}

	_ = flag.NewFlagSet("orm command", flag.ExitOnError)

	args := argString(os.Args[2:])
	cmd := args.Get(0)

	switch cmd {
	case "syncdb":
	case "sqlall":
		sql := getSqlAll()
		fmt.Println(sql)
	default:
		if cmd != "" {
			fmt.Printf("unknown command %s", cmd)
		} else {
			printHelp()
		}

		os.Exit(2)
	}
}
