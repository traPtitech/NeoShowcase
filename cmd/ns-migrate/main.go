package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/traPtitech/neoshowcase/migrations"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
	"net/http"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

var (
	ms         = &migrate.HttpFileSystemMigrationSource{FileSystem: http.FS(migrations.FS)}
	dbHost     string
	dbPort     int
	dbName     string
	dbUser     string
	dbPassword string
	dryrun     bool
)

var rootCommand = &cobra.Command{
	Use:              "ns-migrate",
	Short:            "NeoShowcase DB Migration Tool",
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cliutil.PrintVersion,
}

func upCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Migrates the database to the most recent version available",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := connectDB()
			if err != nil {
				return err
			}
			defer db.Close()

			if dryrun {
				migrations, _, err := migrate.PlanMigration(db, "mysql", ms, migrate.Up, 0)
				if err != nil {
					return err
				}
				for _, m := range migrations {
					fmt.Printf("-- Do migration %s\n", m.Id)
					for _, s := range m.Up {
						fmt.Println(s)
					}
				}
			} else {
				n, err := migrate.Exec(db, "mysql", ms, migrate.Up)
				if err != nil {
					return err
				}
				fmt.Printf("Do %d migrations\n", n)
			}
			return nil
		},
	}
	return cmd
}

func downCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Undo a database migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := connectDB()
			if err != nil {
				return err
			}
			defer db.Close()

			if dryrun {
				migrations, _, err := migrate.PlanMigration(db, "mysql", ms, migrate.Down, 1)
				if err != nil {
					return err
				}
				for _, m := range migrations {
					fmt.Printf("-- Undo migration %s\n", m.Id)
					for _, s := range m.Down {
						fmt.Println(s)
					}
				}
			} else {
				n, err := migrate.ExecMax(db, "mysql", ms, migrate.Down, 1)
				if err != nil {
					return err
				}
				fmt.Printf("Undo %d migrations\n", n)
			}
			return nil
		},
	}
	return cmd
}

func listMigrationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show the list of migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			migrations, err := ms.FindMigrations()
			if err != nil {
				return err
			}
			for _, m := range migrations {
				fmt.Println(m.Id)
			}
			return nil
		},
	}
	return cmd
}

func connectDB() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.ParseTime = true
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%d", dbHost, dbPort)
	cfg.DBName = dbName
	cfg.User = dbUser
	cfg.Passwd = dbPassword
	return sql.Open("mysql", cfg.FormatDSN())
}

func init() {
	rootCommand.AddCommand(
		upCommand(),
		downCommand(),
		listMigrationsCommand(),
	)
	flags := rootCommand.PersistentFlags()
	cliutil.SetupDebugFlag(flags)
	cliutil.SetupLogLevelFlag(flags)

	flags.StringVarP(&dbHost, "dbhost", "H", cliutil.GetEnvOrDefault("DB_HOST", "localhost"), "database host name")
	flags.IntVarP(&dbPort, "dbport", "P", cliutil.GetIntEnvOrDefault("DB_PORT", 3306), "database port number")
	flags.StringVarP(&dbName, "dbname", "n", cliutil.GetEnvOrDefault("DB_NAME", "neoshowcase"), "database name")
	flags.StringVarP(&dbUser, "dbuser", "u", cliutil.GetEnvOrDefault("DB_USER", "root"), "database user name")
	flags.StringVarP(&dbPassword, "dbpassword", "p", cliutil.GetEnvOrDefault("DB_PASS", "password"), "database user password")
	flags.BoolVarP(&dryrun, "dryrun", "d", false, "dry run")
	migrate.SetTable("migrations")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
