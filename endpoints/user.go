package endpoints

type Admin struct {
	AdminID       int    `gorm:"primary_key;unique"`
	AdminEmail    string `gorm:"type:text"`
	AdminPassword string `gorm:"type:text"`
}
