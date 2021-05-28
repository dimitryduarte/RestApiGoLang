package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"./model"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Testing go-swagger generation
//
// The purpose of this API is to provide a register of companies for the purpose of studying GoLang
//
//     Schemes: http
//     Host: localhost:8000
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Dimitry <dimitry.brito@q2pay.com.br>
//
//     Consumes:
//     - text/plain
//
//     Produces:
//     - text/plain
//
// swagger:meta

var err error
var company []model.Company
var taxes []model.Taxes
var dsn = "test_user:123456@tcp(127.0.0.1:3306)/webapi"
var client *redis.Client

//Endpoints

func GetCompany(resp http.ResponseWriter, req *http.Request) {
	var empresa []model.Company
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	result := db.Find(&empresa)

	if result.RowsAffected == 0 {
		json.NewEncoder(resp).Encode("{errorMessage: Não foram encontradas empresas para o Id informado}")
		return
	} else {
		json.NewEncoder(resp).Encode(empresa)
		return
	}
}

func GetTaxes(resp http.ResponseWriter, req *http.Request) {
	var taxes []model.Taxes
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	result := db.Find(&taxes)

	if result.RowsAffected == 0 {
		json.NewEncoder(resp).Encode("{errorMessage: Não foram encontradas taxas para o Id informado}")
		return
	} else {
		json.NewEncoder(resp).Encode(taxes)
		return
	}
}

func GetCompanyId(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var empresa model.Company
	var taxas model.Taxes
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.Find(&empresa, params["id"])
	db.Find(&taxas, empresa.TaxesId)
	empresa.Taxes = taxas
	json.NewEncoder(resp).Encode(empresa)
}

func GetTaxesId(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var taxas model.Taxes
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.Find(&taxas, params["id"])
	json.NewEncoder(resp).Encode(taxas)
}

func CreateCompany(resp http.ResponseWriter, req *http.Request) {
	var newcompany model.Company

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	reqBody, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(reqBody, &newcompany)

	db.Create(&newcompany)
	json.NewEncoder(resp).Encode(newcompany.IdCompany)
}

func CreateTaxes(resp http.ResponseWriter, req *http.Request) {
	var newtaxes model.Taxes

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	reqBody, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(reqBody, &newtaxes)

	db.Create(&newtaxes)

}

func UpdateCompany(rep http.ResponseWriter, req *http.Request) {
	var newcompany model.Company
	var empresa model.Company
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	reqBody, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(reqBody, &newcompany)

	db.Find(&empresa, &newcompany.IdCompany)

	empresa = newcompany

	db.Save(&empresa)
}

func UpdateTaxes(rep http.ResponseWriter, req *http.Request) {
	var newtax model.Taxes
	var taxes model.Taxes
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	reqBody, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(reqBody, &newtax)

	db.Find(&taxes, &newtax.TaxesId)

	taxes = newtax

	db.Save(&taxes)
}

func DeleteCompany(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.Delete(&model.Company{}, params["id"])

}

func DeleteTax(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.Delete(&model.Taxes{}, params["id"])

}

// função principal
func main() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(
		&model.Taxes{},
		&model.Company{},
		&model.Users{},
	)
	//Seed
	//db.Model(&User{}).Save((&User{username: "admin", password: "123456", status: "A"}))

	router := mux.NewRouter()

	//Declaração dos Endpoints
	//GET
	router.HandleFunc("/company", GetCompany).Methods("GET")
	router.HandleFunc("/taxes", GetTaxes).Methods("GET")
	router.HandleFunc("/company/{id}", GetCompanyId).Methods("GET")
	router.HandleFunc("/taxes/{id}", GetTaxesId).Methods("GET")

	//POST
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/company", CreateCompany).Methods("POST")
	router.HandleFunc("/taxes", CreateTaxes).Methods("POST")

	//PUT
	router.HandleFunc("/company", UpdateCompany).Methods("PUT")
	router.HandleFunc("/taxes", UpdateTaxes).Methods("PUT")

	//DELETE
	router.HandleFunc("/company/{id}", DeleteCompany).Methods("DELETE")
	router.HandleFunc("/taxes/{id}", DeleteTax).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func Login(resp http.ResponseWriter, req *http.Request) {
	var login model.Logins
	var user model.Users

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	reqBody, _ := ioutil.ReadAll(req.Body)

	//_ = json.NewDecoder(req.Body).Decode(&login)
	json.Unmarshal(reqBody, &login)

	log.Print(string(reqBody))

	// result := db.Where("username = ? AND password = ?", &login.username, &login.password).First(&user)
	result := db.Where("username = ? AND password = ?", "admin", "123456").First(&user)
	if result.RowsAffected == 0 {
		json.NewEncoder(resp).Encode("{errorMessage: Usuário e ou Senha incorretos!}")
		return
	} else {
		token, err := CreateToken(user.IdUser)
		if err != nil {
			json.NewEncoder(resp).Encode(err)
			return
		}

		saveErr := CreateAuth(user.IdUser, token)
		if saveErr != nil {
			json.NewEncoder(resp).Encode(saveErr)
			return
		}

		tokens := map[string]string{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		}

		json.NewEncoder(resp).Encode(tokens)
	}
}

func CreateToken(userid uint64) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(userid uint64, td *model.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// func CreateToken(userid uint64) (string, error) {
// 	var err error
// 	//Creating Access Token
// 	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
// 	atClaims := jwt.MapClaims{}
// 	atClaims["authorized"] = true
// 	atClaims["user_id"] = userid
// 	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
// 	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
// 	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }
