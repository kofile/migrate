// package main
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version     = "v1.3.2"
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND`

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

func main() {
	var err error

	flag.StringP("dir", "d", ".", "directory with migration files")
	flag.StringP("env-file", "e", "", "use .env file")
	flag.BoolP("help", "h", false, "")
	flag.BoolP("version", "v", false, "")

	flag.Usage = usage

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
	viper.SetConfigName("migraterc")
	viper.AddConfigPath("$HOME/.config/migrate")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	args := flag.Args()

	if len(args) >= 1 && args[0] == "help" {
		flag.Usage()
		return
	}

	if viper.GetBool("version") {
		fmt.Println(version)
		return
	}

	if len(args) < 1 {
		flag.Usage()
		return
	}

	envFile := viper.GetString("env-file")

	if envFile != "" {
		err := godotenv.Load(envFile)

		if err != nil {
			log.Fatalf("Error loading %s file", envFile)
		}
	}

	viper.BindEnv("env-file", "ENV_FILE")
	viper.BindEnv("database.url", "DB_URL")

	var connectionStr *url.URL

	if viper.IsSet("database.url") {
		connectionStr, err = url.Parse(viper.GetString("database.url"))

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	} else {
		connectionStr = &url.URL{}
		errors := make([]string, 0)

		if !viper.IsSet("database.host") {
			errors = append(errors, "database.host is required!")
		}

		if !viper.IsSet("database.port") {
			errors = append(errors, "database.port is required!")
		}

		if !viper.IsSet("database.user") {
			errors = append(errors, "database.user is required!")
		}

		if !viper.IsSet("database.pass") {
			errors = append(errors, "database.pass is required!")
		}

		if !viper.IsSet("database.name") {
			errors = append(errors, "database.name is required!")
		}

		if len(errors) != 0 {
			for _, message := range errors {
				fmt.Println(message)
			}
			os.Exit(1)
		}

		dbHost := viper.GetString("database.host")
		dbPort := viper.GetString("database.port")
		dbUser := viper.GetString("database.user")
		dbPass := viper.GetString("database.pass")
		dbName := viper.GetString("database.name")
		dbOpts := ""

		if viper.IsSet("database.options") {
			dbOpts = viper.GetString("database.options")
		}

		connectionStr.Host = makeHost(dbHost, dbPort)

		if dbUser != "" || dbPass != "" {
			if dbPass == "" {
				connectionStr.User = url.User(dbUser)
			} else {
				connectionStr.User = url.UserPassword(dbUser, dbPass)
			}
		}

		connectionStr.Path = dbName
		connectionStr.RawQuery = dbOpts
		connectionStr.Scheme = "postgres"
	}

	dir := ""

	if viper.IsSet("migrations.directory") {
		dir = viper.GetString("migrations.directory")
	}

	if len(args) > 1 && args[0] == "create" {
		if err := goose.Run("create", nil, dir, args[1:]...); err != nil {
			log.Fatalf("migrate run: %v", err)
		}
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

	if err := goose.Run(command, db, dir, arguments...); err != nil {
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
	fmt.Println(fmt.Sprintf("migrate %s", version))
	fmt.Println(usagePrefix)
	flag.PrintDefaults()
	fmt.Println(usageCommands)
}
