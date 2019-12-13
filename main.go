package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

)

var Users = make(map[string]*User)
var Users2 = make([]*User, 0)

type (
	User struct {
		Email string `json:"email" xml:"email" form:"email" query:"email"`
		Tokens []string
	}
	Device struct {
		Email string `json:"email" xml:"email" form:"email" query:"email"`
		Token string `json:"token"`
	}
	Data struct {
		 Feature string `json:"feature"`
		 Body string `json:"body"`
	}

	Payload struct {
		RegistrationIds []string `json:"registration_ids"`
		Data Data `json:"data"`
	}
)

func main()  {
	e := echo.New()
	e.POST("/user/register", RegisterUser)

	http.HandleFunc("device/register", RegisterDevice)
    e.POST("/device/push", Push)

	e.Logger.Fatal(e.Start(":4001"))


}

func RegisterUser(c echo.Context) error {
	param := new(User)
	err := c.Bind(param)
    if err != nil {
    	log.Printf(err.Error())
    	return c.JSON(http.StatusBadRequest , echo.Map{"status" : err.Error()})
	}

	user := Users[param.Email]
	log.Print("user", user)
 	if user == nil {
          newUser := &User{Email: param.Email}
          Users[param.Email] = newUser
	      return c.JSON(200, echo.Map{
	      	"status" : "dang ki thanh cong",
	      	"email"  : param.Email,
		   })
	}

	return c.JSON(http.StatusOK , echo.Map{
		"status" : "user da ton tai",
		"email" : param.Email,
	})
}

func RegisterDevice(w http.ResponseWriter, r *http.Request)  {


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

func Push(c echo.Context) error  {
	param := new(User)
	err:= c.Bind(param)
	if err != nil {
		log.Print(err)
	}

	//user  := Users[param.Email]
    payload := Payload{
		RegistrationIds: []string{"e-F294Izgqs:APA91bFGaCiRpfHog24IL7Vw1FmtoF4KYtFvYHIGJr0L6CfnMFns8NBZ9Zf8WgzLId_KmDgjhembvOdeo0-rlhklLqnFK1fi1YYMRTB4ywCwF-HJ7s1hmx9bJzsb6IkpqS4iogmgV_Or"},
		Data:            Data{
			Feature: "Hello Nam",
			Body:    "Nice to meet u," ,
		},
	}
	bPayload , _  := json.Marshal(payload)
	go pushNotify(bPayload)

	return c.JSON(http.StatusOK, echo.Map{"status": "Dang push"})
}

func pushNotify(payload []byte) {
	  req,  _ := http.NewRequest("POST","https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(payload))

	  req.Header.Set("Content-Type","application/json")
	  req.Header.Set("Authorization","key=AAAAmlwBCtY:APA91bEn7bT8YgMCY9vR0iaHAJmonIA7_xownDfpur5GajXpahdtMzX6h0SWdKGg-eJmZYjQyV1H55Td1cyop6N7puYxHRDz0d27cCxaIH9oNDFNxbhYHrDl-Ab56Ak36olfKMuPI6dX")

	  cl := http.Client{}
	  response , _ :=  cl.Do(req)

	  b, _ := ioutil.ReadAll(response.Body)
	  responseString := string(b)
	  log.Print(responseString)

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
	registrationToken := "eCtSd7Ymvo4:APA91bEGEbqbSvQKVUvw_xO-8aQV6OaT8xOCESVgyouaB4GWZSvUkX3718Ll4yKApqrA0ZYv-Inu91oB2nzPefe7YAoTpl6DL2PN_1JjXrZLccmQzI1dW_pTbqAU7stTDlI1V-WAh6gN"

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