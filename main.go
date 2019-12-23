package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
)

var Users = make(map[string]*User)
var Users2 = make([]*User, 0)

type (
	HomegateUsers struct {
		Id         int    `json:"id"`
		UserId     string `json:"user_id"`
		HomegateId string `json:"homegate_id"`
	}
	TokenDevices struct {
		TokenDevice string `json:"token_devices"
`
	}

	User struct {
		Email  string `json:"email" xml:"email" form:"email" query:"email"`
		Tokens []string
	}
	Device struct {
		Email string `json:"email" xml:"email" form:"email" query:"email"`
		Token string `json:"token"`
	}
	Data struct {
		Title        string `json:"title"`
		Body         string `json:"body"`
		Notification messaging.Notification
	}
	Notification struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	Payload struct {
		RegistrationIds []string     `json:"registration_ids"`
		Data            Data         `json:"data"`
		Notification    Notification `json:"notification"`
	}
)

func main() {
	//username := os.Getenv("db_user")
	//password := os.Getenv("db_pass")
	//dbName := os.Getenv("db_name")
	//dbHost := os.Getenv("db_host")
	//log.Print("valueeeeLL",dbHost,username,password)
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", "localhost", "postgres", "auth", "postgres") //Build connection string
	dbCon, err := gorm.Open("postgres", dbUri)
	if err != nil {
		panic(err)
	}
	defer dbCon.Close()

	database := dbCon.DB()
	err = database.Ping()
	if err != nil {
		panic(err)
	}
	log.Print("connected to DB")

	e := echo.New()

	e.POST("/user/register", RegisterUser)

	http.HandleFunc("device/register", RegisterDevice)
	e.POST("/device/push", Push(dbCon))

	e.Logger.Fatal(e.Start(":8081"))

}

func RegisterUser(c echo.Context) error {
	param := new(User)
	err := c.Bind(param)
	if err != nil {
		log.Printf(err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"status": err.Error()})
	}

	user := Users[param.Email]
	log.Print("user", user)
	if user == nil {
		newUser := &User{Email: param.Email}
		Users[param.Email] = newUser
		return c.JSON(200, echo.Map{
			"status": "dang ki thanh cong",
			"email":  param.Email,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "user da ton tai",
		"email":  param.Email,
	})
}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {

	//param := new(Device)
	//if err := context.Bind(param) ; err!= nil {
	//	return context.JSON(http.StatusBadRequest, echo.Map{"status": err.Error()})
	//}
	//user := Users[param.Email]
	//user.Tokens  = append(user.Tokens , param.Token)
	//
	//return context.JSON(200, echo.Map{
	//	"status" : "Your mobile has been registered",
	//	"email" : param.Email,
	//})
}

func Push(dbGorm *gorm.DB) echo.HandlerFunc {
	var arrTokenDevice = make([]string, 0)

	return func(c echo.Context) error {
		//var hgUser HomegateUsers
		homegateId := c.Request().Header.Get("homegate_id")
		//res1 := dbGorm.Debug().First(&hgUser, "homegate_id = ?", hgIdGetRequest)
		//if res1.Error != nil {
		//	log.Print("error")
		//}
		//resValue := res1.Value
		//b, err := json.Marshal(resValue)
		//if err != nil {
		//	log.Print(err.Error())
		//}
		//err = json.Unmarshal(b, &hgUser)

		getTokenDev, err := dbGorm.Debug().
			Table("users").
			Select("users.token_devices").
			Joins("left join homegate_users on users.id = homegate_users.user_id ").
			Where("homegate_id = ?", homegateId).
			Rows()
		if err != nil {
			log.Print(err.Error())
		}
		var tokenArr = make([]*TokenDevices, 0)

		defer getTokenDev.Close()
		for getTokenDev.Next() {
			t := new(TokenDevices)
			err = getTokenDev.Scan(&t.TokenDevice)
			if err != nil {
				log.Print(err)
			}
			tokenArr = append(tokenArr, t)
		}
		log.Print("tokenArr:::", tokenArr)
		//var tokenDev TokenDevices

		dbGorm.LogMode(true)
		arrTokenDevice = append(arrTokenDevice, "fZQiZf6T7HQ:APA91bG_4fyMRNzGN6g3-rzfZEr77C47KySEeGFz_wGOkDW7NwkGbcjtT1h2nh88aXwqnd7TYykuIBtoV29MQtoz4b-PdI0Vu1eeweeB2M_9CIX9eszU72-IZv3v_9Bs9T_KPZ-VN1HG")

		payload := Payload{
			RegistrationIds: arrTokenDevice,
			Data: Data{
				Title: "Humidity sensor",
				Body:  "a text",
				Notification: messaging.Notification{
					Title:    "sd",
					Body:     "sd",
					ImageURL: "https://img.lovepik.com/element/40035/8356.png_860.png",
				},
			},
			Notification: Notification{
				Title: "Humidity sensor",
				Body:  "High temperature . Should I decrease the air conditional",
				//Icon :
			},
		}
		bPayload, _ := json.Marshal(payload)
		go pushNotify(bPayload)
		return c.JSON(http.StatusOK, "sd")
	}
}

func Push1(c echo.Context) (err error) {
	var arrTokenDevice = make([]string, 0)
	arrTokenDevice = append(arrTokenDevice, "fZQiZf6T7HQ:APA91bG_4fyMRNzGN6g3-rzfZEr77C47KySEeGFz_wGOkDW7NwkGbcjtT1h2nh88aXwqnd7TYykuIBtoV29MQtoz4b-PdI0Vu1eeweeB2M_9CIX9eszU72-IZv3v_9Bs9T_KPZ-VN1HG")

	payload := Payload{
		RegistrationIds: arrTokenDevice,
		Data: Data{
			Title: "Humidity sensor",
			Body:  "a text",
			Notification: messaging.Notification{
				Title:    "sd",
				Body:     "sd",
				ImageURL: "https://img.lovepik.com/element/40035/8356.png_860.png",
			},
		},
		Notification: Notification{
			Title: "Humidity sensor",
			Body:  "High temperature . Should I decrease the air conditional",
			//Icon :
		},
	}
	bPayload, _ := json.Marshal(payload)
	go pushNotify(bPayload)

	return c.JSON(http.StatusOK, echo.Map{"status": "Dang push"})
}

func pushNotify(payload []byte) {
	req, _ := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key=AAAANbra0SM:APA91bEoUGX77jMMl60L297UJYAKa8nKk1ImI5eTsz6FyFiTU7FhFNGQadKbwUUpDOKfZqUHGJ0oI2G-ane5jux3IEvQ4j2e6EJkWWwQvUeKQARd8o-RwHPnzVm0m4NvaZMkHbZOKIqP")

	cl := http.Client{}
	response, _ := cl.Do(req)

	b, _ := ioutil.ReadAll(response.Body)
	responseString := string(b)
	log.Print(responseString)

	// Response is a message ID string.

	defer response.Body.Close()

}

func sendToToken(app *firebase.App) {
	// [START send_to_token_golang]
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "dCHGQ1nAdIc:APA91bFwBpS6bYtrX2gYTAD-atljc9jipasjg5XkN6z7CWPThc90mopFDA2tmH8kN_KnXqpRhMxAiIHao7GUaiK0MXDTtsxyT_JhbADSXTSbiETCakkwhW5pJOU5mo_wmHyJs_O_WK0q"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	log.Println("Successfully sent message:", response)
	// [END send_to_token_golang]
}
