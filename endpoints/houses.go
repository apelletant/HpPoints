package endpoints

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type AddPoints struct {
	HouseName  string
	HousePoint int
}

type House struct {
	HouseID    int    `gorm:"primary_key;unique"`
	HouseName  string `gorm:"type:text"`
	HousePoint int    `gorm:"type:int"`
}

func GetHouses(db *gorm.DB) ([]House, error) {
	var house []House
	err := db.Find(&house).Error
	if err != nil {
		return house, errors.New("No house find")
	}
	return house, nil
}

func GetHouse(houseName string, db *gorm.DB) (House, error) {
	var house House
	err := db.Where("house_name = ?", houseName).Find(&house).Error
	if err != nil {
		return house, errors.New("unable to find this Houses")
	}
	return house, nil
}

func updateHouse(house *House, points int, db *gorm.DB) error {
	house.HousePoint += points
	err := db.Model(&house).Updates(house).Error
	if err != nil {
		return err
	}
	return nil
}

func AddPoint(toAdd AddPoints, db *gorm.DB) (House, error) {
	house, err := GetHouse(toAdd.HouseName, db)
	if err != nil {
		return house, err
	}
	err = updateHouse(&house, toAdd.HousePoint, db)
	if err != nil {
		return house, err
	}
	return house, nil
}
