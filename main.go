package main

//export les variables boloss

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var (
	db *gorm.DB
)

type House struct {
	HouseID    int    `gorm:"primary_key;unique"`
	HouseName  string `gorm:"type:text"`
	HousePoint int    `gorm:"type:int"`
}

type Admin struct {
	AdminID       int    `gorm:"primary_key;unique"`
	AdminEmail    string `gorm:"type:text"`
	AdminPassword string `gorm:"type:text"`
}

func HouseData() ([]House, error) {
	var houses []House
	err := db.Find(&houses).Error
	return houses, err
}

func ModifyPoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println(r.Form)
	}

	fmt.Println(r.Form["houses"])
	fmt.Println(r.Form["points"])
	house := r.FormValue("ravenclaw")
	fmt.Println(house)
}

//VerifiedCredential used to see if a user exist and if his login credential are correct
func VerifiedCredential(email string, pass string) error {
	var admin Admin
	err := db.Where("admin_email = ?", email).First(&admin).Error
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}

//Index  salut
func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmplt/index.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

//AdminConnexion salut
func AdminConnexion(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var housesData []House
	if r.Method == "GET" {
		t, _ := template.ParseFiles("tmplt/adminForm.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		err := VerifiedCredential(r.FormValue("email"), r.FormValue("password"))
		if err != nil {
			t, err = template.ParseFiles("tmplt/connexionError.html")
		} else {
			housesData, err = HouseData()
			if err != nil {
				panic(err)
			}
			t, err = template.ParseFiles("tmplt/adminPanel.html")
		}
		if err != nil {
			panic(err)
		}
		if err := t.Execute(w, housesData); err != nil {
			panic(err)
		}
	}
}

func init() {
	var err error
	dbPath := os.Getenv("DBPATH")
	fmt.Println(dbPath)
	db, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	defer db.Close()

	http.HandleFunc("/", Index)
	http.HandleFunc("/admin", AdminConnexion)
	http.HandleFunc("/modifyPoint", ModifyPoint)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
