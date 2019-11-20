package main

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Tips :
type Tips struct {
	Tips []Tip `json:"Tips"`
}

// Tip :
type Tip struct {
	Class string `json:"class"`
	Text  string `json:"text"`
}

// Credenciales : tipo para  oAuth
type Credenciales struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

// readJSON: Hace el parse a un JSON
func readJSON(file string) Tips {
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Println(err)
	}
	// Se cierra al acabar la función
	defer jsonFile.Close()
	byteSlice, _ := ioutil.ReadAll(jsonFile)
	var tips Tips
	json.Unmarshal(byteSlice, &tips)
	return tips

}

func getClient(cred *Credenciales) (*twitter.Client, error) {
	// Creamos configuración de autentificación de oauth1
	config := oauth1.NewConfig(cred.ConsumerKey, cred.ConsumerSecret)
	// Creamos el token de acceso
	token := oauth1.NewToken(cred.AccessToken, cred.AccessSecret)

	//Creamos el cliente http
	httpClient := config.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	user, _, err := twitterClient.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's Account:\n%+v\n", user)
	return twitterClient, nil
}

//Difference : Función que saca los elementos diferentes de dos slices de tweets
func Difference(a, b []twitter.Tweet) (diff []twitter.Tweet) {
	m := make(map[int64]bool)

	for _, item := range b {
		m[item.ID] = true
	}

	for _, item := range a {
		if _, ok := m[item.ID]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
func main() {
	fmt.Println("El bot del vish 0.0.1")
	creds := Credenciales{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
	}

	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error intentando acceder a twitter.")
		log.Println(err)
	}
	// Inicialización de los tips.
	tips := readJSON("./resources/tips.json")
	par := &twitter.SearchTweetParams{Query: "#DimeBjarne", ResultType: "recent"}
	search, _, err := client.Search.Tweets(par)
	tweets := search.Statuses
	for {
		search, _, err := client.Search.Tweets(par)
		newTweets := search.Statuses
		if err != nil {
			log.Println("Error al buscar")
			log.Println(err)
		}
		diff := Difference(newTweets, tweets)
		if len(diff) > 0 {
			tweets = newTweets
		}
		for _, item := range diff {
			if strings.Contains(item.Text, "tip") || strings.Contains(item.Text, "consejo") {
				parametros := &twitter.StatusUpdateParams{InReplyToStatusID: item.ID}
				_, _, err := client.Statuses.Update("@"+item.User.ScreenName+" "+tips.Tips[rand.Int()%len(tips.Tips)].Text, parametros)
				if err != nil {
					log.Println(err)
				}
			} else {
				parametros := &twitter.StatusUpdateParams{InReplyToStatusID: item.ID}
				_, _, err := client.Statuses.Update("@"+item.User.ScreenName+" Bjarne Strouptrup te vigila ;-)", parametros)
				if err != nil {
					log.Println(err)
				}
			}
			log.Println("Nuevo tweet")
			time.Sleep(1000000000)
		}
		time.Sleep(5000000000)

	}
}
