package internal

import (
	"github.com/cSploit/daemon/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	Db          *gorm.DB
	models      []interface{}
	join_tables []string
)

func OpenDb() error {
	return openDb(false)
}

func openDb(drop_tables bool) error {
	var dd *gorm.DB
	var err error

	if dd, err = gorm.Open(config.Conf.Db.Dialect, config.Conf.Db.Args...); err != nil {
		return err
	}

	if drop_tables {
		dd.DropTableIfExists(models...)
		for _, table_name := range join_tables {
			dd.DropTableIfExists(table_name)
		}
	}

	dd = dd.Debug().AutoMigrate(models...)

	if dd.Error == nil {
		Db = dd
	}

	return dd.Error
}

func ClearDb() {
	for _, model := range models {
		Db.Delete(model)
	}
	for _, table_name := range join_tables {
		Db.Exec("DELETE FROM " + table_name)
	}
}

func OpenDbForTests() {
	if Db == nil {
		if err := openDb(true); err != nil {
			panic(err)
		}
	}
	ClearDb()
}

func RegisterModels(model ...interface{}) {
	models = append(models, model...)
}

func RegisterJoinTables(table_name ...string) {
	join_tables = append(join_tables, table_name...)
}

//TODO: UpdateCallback { e -> wsConns.each { c -> c.write("entity changed:" + e) } } [Event system]
