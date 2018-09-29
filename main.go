package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yabou/HpPoints/endpoints"
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
	var err error
	var dataToAdd endpoints.AddPoints
	if r.Method == "GET" {
		for _, v := range r.Cookies() {
			err := bcrypt.CompareHashAndPassword([]byte("$2a$10$TyKWWjifpB6HUhBEOm7Aq.m9Fnex/rYbHOArvbft.io07TAZ9.Ace"), []byte(v.Value))
			if err == nil {
				housesData, err := HouseData()
				if err != nil {
					panic(err)
				}
				t, err := template.ParseFiles("tmplt/adminPanel.html")
				if err != nil {
					panic("err")
				}
				t.Execute(w, housesData)
			}
		}
	} else {
		r.ParseForm()
		dataToAdd.HouseName = r.FormValue("houses")
		dataToAdd.HousePoint, err = strconv.Atoi(r.FormValue("points"))
		if err != nil {
			panic(err)
		}
		house, err := endpoints.AddPoint(dataToAdd, db)
		if err != nil {
			panic(err)
		}
		t, err := template.ParseFiles("tmplt/pointsAdd.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, house)
	}
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

func AdminConnexion(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var err error

	if err != nil {
		panic(err)
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("tmplt/adminForm.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		err := VerifiedCredential(r.FormValue("email"), r.FormValue("password"))
		if err != nil {
			t, err = template.ParseFiles("tmplt/connexionError.html")
			if err != nil {
				panic(err)
			}
			err = t.Execute(w, nil)
			if err != nil {
				panic(err)
			}
		} else {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "username", Value: r.FormValue("email"), Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/modifyPoint", 302)
		}
	}
}

func init() {
	var err error
	dbPath := os.Getenv("DBPATH")
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
