package utils

import (
	"database/sql"
	"path"
	"runtime"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _, filename, _, _ = runtime.Caller(0)
var PROJECT_BASE = path.Join(path.Dir(filename), "..")

var zapSingleton *zap.Logger
var loadZapOnce sync.Once

func GetZapLogger(fields ...zapcore.Field) *zap.Logger {
	loadZapOnce.Do(func() {
		zapSingleton, _ = zap.NewProduction()
	})
	return zapSingleton.With(fields...)
}

var sqliteSingleton *sql.DB
var loadSqliteOnce sync.Once

func GetSqliteDB() *sql.DB {
	var err error
	dbPath := path.Join(PROJECT_BASE, "db.sqlite")
	loadSqliteOnce.Do(func() {
		sqliteSingleton, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			panic(err)
		}
	})
	return sqliteSingleton
}

var viperSingleton *viper.Viper
var loadViperOnce sync.Once

func GetDefaultViper() *viper.Viper {
	loadViperOnce.Do(func() {
		viperSingleton = viper.New()
		viperSingleton.AddConfigPath(path.Join(PROJECT_BASE, "config"))
		viperSingleton.SetConfigName("config")
		viperSingleton.SetConfigType("toml")
		if err := viperSingleton.ReadInConfig(); err != nil {
			panic(err)
		}
	})
	return viperSingleton
}
