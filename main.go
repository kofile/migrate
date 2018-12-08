// package main
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/pressly/goose"
)

type config struct {
	URL      string `env:"DB_URL"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Database string `env:"DB_DATABASE"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Options  string `env:"DB_OPTIONS"`
}

var (
	flags  = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir    = flags.String("dir", ".", "directory with migration files")
	useEnv = flags.Bool("env", false, "use local .env file")
)

func main() {
	cfg := config{}

	if *useEnv {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	err := env.Parse(&cfg)

	var connectionStr *url.URL

	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	if cfg.URL != "" {
		connectionStr, err = url.Parse(cfg.URL)

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	} else {
		connectionStr = &url.URL{}
		connectionStr.Host = makeHost(cfg.Host, cfg.Port)

		if cfg.Username != "" || cfg.Password != "" {
			if cfg.Password == "" {
				connectionStr.User = url.User(cfg.Username)
			} else {
				connectionStr.User = url.UserPassword(cfg.Username, cfg.Password)
			}
		}

		connectionStr.Path = cfg.Database
		connectionStr.RawQuery = cfg.Options
		connectionStr.Scheme = "postgres"
	}

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) > 1 && args[0] == "create" {
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("migrate run: %v", err)
		}
		return
	}

	if len(args) < 1 {
		flags.Usage()
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	log.Printf("Connecting to %s", connectionStr.String())
	db, err := sql.Open("postgres", connectionStr.String())

	if err != nil {
		log.Fatal(err)
	}

	arguments := []string{}

	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("migrate run: %v", err)
	}
}

func makeHost(host string, port string) (hostname string) {
	if port != "" {
		return fmt.Sprintf("%s:%s", host, port)
	}

	return host
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
	fmt.Println("v1.2.0")
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Drivers:
    postgres
Examples:
    migrate status
Available ENV vars:
    DB_URL
    DB_HOST
    DB_PORT
    DB_DATABASE
    DB_USERNAME
    DB_PASSWORD
    DB_OPTIONS    eg. sslmode=disable
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
`
)
