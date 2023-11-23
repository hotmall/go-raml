package main

import (
	"flag"
	"log"

	"github.com/hotmall/go-raml/docs/tutorial/go/client/goramldir"
)

var (
	appID     = flag.String("app_id", "", "application ID")
	appSecret = flag.String("app_secret", "", "application secret")
)

func main() {
	flag.Parse()
	if *appID == "" || *appSecret == "" {
		log.Fatalf("please specify itsyou.online application ID & API Key")
	}

	gr := goramldir.Newgoramldir()

	// create itsyou.online JWT token
	jwtToken, err := gr.GetOauth2AccessToken(*appID, *appSecret, []string{}, []string{})
	if err != nil {
		log.Fatalf("failed to create itsyou.online JWT token:%v", err)
	}

	log.Printf("got JWT token = %v\n", jwtToken)
	// set goramldir authorization header to use JWT token
	gr.AuthHeader = "Bearer " + jwtToken

	// calling GET /users
	users, resp, err := gr.Users.UsersGet(nil, nil)
	if err != nil {
		log.Fatalf("failed to GET /users. err = %v", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("GET /users failed. http status code = %v", resp.StatusCode)
	}

	log.Printf("users = %v\n", users)

}
