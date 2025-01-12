package initPackage

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gorm.io/gorm"
)

var Db *gorm.DB

var InfluxDB *influxdb2.Client

func SetDatabases(database *gorm.DB, influxClient *influxdb2.Client) {
	Db = database

	InfluxDB = influxClient
}