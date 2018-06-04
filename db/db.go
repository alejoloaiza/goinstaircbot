package db

import (
	"fmt"
	"goinstabot/config"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbpostgre *gorm.DB
var err error

type FollowingUser struct {
	gorm.Model
	UserId string
}

func DBConnectPostgres() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", config.Localconfig.DBHost, config.Localconfig.DBPort, config.Localconfig.DBUser, config.Localconfig.DBPass, config.Localconfig.DBName)
	dbpostgre, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//fmt.Println(">>>>>>>>>>>>>>>>> Successfully connected to Database <<<<<<<<<<<<<<<<<")
}
func DBInsertPostgres_Following(Username string) error {

	dbpostgre.Create(&FollowingUser{UserId: Username})
	if dbpostgre.Error != nil {
		fmt.Println(dbpostgre.Error)
	}
	return err
}
func DBDeletePostgres_Following(Username string) error {
	var userToDelete FollowingUser

	dbpostgre.First(&userToDelete, "UserId = ?", Username)

	dbpostgre.Delete(&userToDelete)
	if dbpostgre.Error != nil {
		fmt.Println(dbpostgre.Error)
	}
	return err
}
func DBSelectPostgres_Following() []string {
	var userToDelete FollowingUser

	var users []string
	rows, err := dbpostgre.Find(&userToDelete).Select("UserId").Rows()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var username string

		err = rows.Scan(&username)
		users = append(users, username)
	}
	return users
}
