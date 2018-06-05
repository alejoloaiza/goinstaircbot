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
type BlockedUser struct {
	gorm.Model
	UserId string
}

func DBConnectPostgres() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", config.Localconfig.DBHost, config.Localconfig.DBPort, config.Localconfig.DBUser, config.Localconfig.DBPass, config.Localconfig.DBName)
	dbpostgre, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	dbpostgre.LogMode(false)
	dbpostgre.CreateTable(&BlockedUser{})
	dbpostgre.CreateTable(&FollowingUser{})

	//fmt.Println(">>>>>>>>>>>>>>>>> Successfully connected to Database <<<<<<<<<<<<<<<<<")
}
func DBClosePostgress() {
	dbpostgre.Close()
}
func DBInsertPostgres_Following(Username string) error {

	dbpostgre.Create(&FollowingUser{UserId: Username})
	if dbpostgre.Error != nil {
		fmt.Println(dbpostgre.Error)
	}
	return dbpostgre.Error
}
func DBInsertPostgres_Blocked(Username string) error {

	var count int
	dbpostgre.Model(&BlockedUser{}).Where("User_Id = ?", Username).Count(&count)
	if count == 0 {
		dbpostgre.Create(&BlockedUser{UserId: Username})
	}
	if dbpostgre.Error != nil {
		fmt.Println(dbpostgre.Error)
	}
	return dbpostgre.Error
}
func DBDeletePostgres_Following(Username string) error {
	var userToDelete FollowingUser

	dbpostgre.First(&userToDelete, "User_Id = ?", Username)

	dbpostgre.Delete(&userToDelete)
	if dbpostgre.Error != nil {
		fmt.Println(dbpostgre.Error)
	}
	return dbpostgre.Error
}
func DBSelectPostgres_Following() []string {

	var users []string

	// TODO: To check why is not working with Model instead of Table
	//dbpostgre.Model(&FollowingUser{}).Pluck("User_Id", &users)

	dbpostgre.Table("following_users").Pluck("User_Id", &users)

	return users
}
func DBSelectPostgres_Blocked() []string {

	var users []string

	dbpostgre.Table("blocked_users").Pluck("User_Id", &users)

	return users
}
