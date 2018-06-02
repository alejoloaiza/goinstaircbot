package db

import (
	"database/sql"
	"fmt"
	"goinstabot/config"

	_ "github.com/lib/pq"
)

var dbpostgre *sql.DB
var err error

func DBConnectPostgres(configpath string) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Localconfig.DBHost, config.Localconfig.DBPort, config.Localconfig.DBUser, config.Localconfig.DBPass, config.Localconfig.DBName)
	dbpostgre, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//fmt.Println(">>>>>>>>>>>>>>>>> Successfully connected to Database <<<<<<<<<<<<<<<<<")
}

/*
func DBInsertPostgres(a *assets.Asset) {

	point := fmt.Sprintf(`'POINT( %.6f %.6f )'`, a.Lat, a.Lon)

	sqlStatement := `
		INSERT INTO parallel.webscrapingresults
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, postgis.ST_GeomFromText( ` + point + ` )  );`

	_, err := dbpostgre.Exec(sqlStatement, a.Business, a.Code, a.Type, a.Agency, a.Location, a.City, a.Area, a.Price, a.Numrooms, a.Numbaths, a.Parking, a.Status, a.Link)
	//fmt.Println(a.Business, a.Code, a.Type, a.Agency, a.Location, a.City, a.Area, a.Price, a.Numrooms, a.Numbaths, a.Status, a.Link)
	if err != nil {
		fmt.Println(err)
	}
}
*/
