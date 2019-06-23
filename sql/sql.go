package sql

import (
	db "database/sql"
	"fmt"

	"github.com/aeridya/core"
	"github.com/aeridya/module"
)

// DBType specifies the type of Database available and being used.
type DBType int

const (
	//MYSQL is a constant to define that MySQL is being used
	MYSQL DBType = iota
	//PGSQL is a constant to define that Postgresql is being used
	PGSQL
	//SQLITE is a constant to define that sqlite3 is being used
	SQLITE
)

// String converts the DBType to a string
func (d DBType) String() string {
	switch d {
	case MYSQL:
		return "mysql"
	case PGSQL:
		return "postgres"
	case SQLITE:
		return "sqlite3"
	default:
		return "unknown"
	}
}

var (
	// DB is the database reference created by "database/sql"
	DB *db.DB

	//sqlinfo Instance information
	info *sqlinfo
)

func init() {
	info = &sqlinfo{}
}

// sqlinfo contains all of the needed information to connect to a sql database.
type sqlinfo struct {
	module.Module
	Username string
	Password string
	Name     string
	Host     string
	Port     string
	Options  string
	Driver   DBType
}

func Username(user string) module.Option {
	return func() {
		info.Username = user
	}
}

func Password(pass string) module.Option {
	return func() {
		info.Password = pass
	}
}

func Name(name string) module.Option {
	return func() {
		info.Name = name
	}
}

func Hostname(host string) module.Option {
	return func() {
		info.Host = host
	}
}

func Port(port string) module.Option {
	return func() {
		info.Port = port
	}
}

func Options(options string) module.Option {
	return func() {
		info.Options = options
	}
}

func Driver(driver DBType) module.Option {
	return func() {
		info.Driver = driver
	}
}

// Connect accepts options as-needed and connects to the SQL instance depending
// on the driver being used.
// Returns an error
func Connect(options ...module.Option) error {
	info.ParseOpts(options)
	switch info.Driver {
	case MYSQL:
		if err := connectMYSQL(); err != nil {
			return err
		}
	case PGSQL:
		if err := connectPSQL(); err != nil {
			return err
		}
	case SQLITE:
		if err := connectSQLITE(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown Driver Type: %s", info.Driver.String())
	}
	return DB.Ping()
}

// GetInfo returns the SQL information currently in use
func GetInfo() *sqlinfo {
	return info
}

// ReadConfig is a helper function that retrieves the Username, Password,
// Database, Host, and Port from the module Configuration file and sets
// the info
func ReadConfig() (err error) {
	if info.Name, err = core.Config.GetString("database", "dbname"); err != nil {
		return fmt.Errorf("Unable to get DBName from config: %s", err)
	}
	if info.Username, err = core.Config.GetString("database", "dbuser"); err != nil {
		return fmt.Errorf("Unable to get DBUser from database config: %s", err)
	}
	if info.Password, err = core.Config.GetString("database", "dbpass"); err != nil {
		return fmt.Errorf("Unable to get DBPass from database config: %s", err)
	}
	if info.Host, err = core.Config.GetString("database", "dbhost"); err != nil {
		return fmt.Errorf("Unable to get DBHost from database config: %s", err)
	}
	if info.Port, err = core.Config.GetString("database", "dbport"); err != nil {
		return fmt.Errorf("Unable to get DBPort from database config: %s", err)
	}
	return nil
}

// connects to the Postgresql database
func connectPSQL() error {
	var err error
	DB, err = db.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", info.Username, info.Password, info.Host, info.Port, info.Name, info.Options))
	if err != nil {
		return fmt.Errorf("Unable to connect to Postgres Database: %s", err)
	}
	return nil
}

// connects to the MySQL database
func connectMYSQL() error {
	var err error
	DB, err = db.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", info.Username, info.Password, info.Host, info.Port, info.Name, info.Options))
	if err != nil {
		return fmt.Errorf("Unable to connect to MySQL Database: %s", err)
	}
	return nil
}

// connects to the SQLITE database
func connectSQLITE() error {
	return fmt.Errorf("SQLITE not implemented yet...")
}
