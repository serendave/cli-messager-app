package apps

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func Slack() {
	addBotUrl := os.Getenv("ADD_SLACK_BOT_URL")
	fmt.Println("Add the bot to your workspace by going to:")
	fmt.Println(addBotUrl)
	server()
}

func NewWorkspace(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Success! You can now return to your console app.\n"))
	code := req.URL.Query().Get("code")
	clientId := os.Getenv("SLACK_CLIENT_ID")
	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	fmt.Println("The code: " + code)
	accessToken, err := getAccessToken(code, clientId, clientSecret)
	if err != nil {
		log.Fatal("Getting access token error: ", err)
	}
	client := slack.New(accessToken)
	users, err := client.GetUsers()
	if err != nil {
		log.Fatal("Error retrieving users: ", err)
	}

	for i, user := range users {
		if user.IsBot != true {
			fmt.Printf("%d: %s\n", i+1, user.Profile.RealNameNormalized)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	var selectedUser slack.User
	for {
		fmt.Print("Enter the number of the user you want to message to: ")
		char, _, err := reader.ReadRune()

		if err != nil {
			panic(err)
		}

		number := int(char - '0')

		if number >= len(users) {
			fmt.Println("Please, provide a valid number")
		} else {
			selectedUser = users[number-1]
			openConversationParams := slack.OpenConversationParameters{
				Users: []string{selectedUser.ID},
			}

			channel, _, _, err := client.OpenConversation(&openConversationParams)
			if err != nil {
				panic(err)
			}

			message, err := reader.ReadString('\n')

			for {
				fmt.Println("Type message you want to send to this user: ")
				message, err = reader.ReadString('\n')
				if message == "\n" {
					break
				}

				if err != nil {
					panic(err)
				}

				attachment := slack.Attachment{
					Pretext: message,
				}
				// PostMessage will send the message away.
				// First parameter is just the channelID, makes no sense to accept it
				_, _, err := client.PostMessage(
					channel.ID,
					// uncomment the item below to add a extra Header to the message, try it out :)
					//slack.MsgOptionText("New message from bot", false),
					slack.MsgOptionAttachments(attachment),
				)

				if err != nil {
					panic(err)
				} else {
					fmt.Println("Message sent successfully")
				}
			}
		}
	}
}

type SlackAuthResponse struct {
	Access_token string
	Scope        string
}

func getAccessToken(code string, clientId string, clientSecret string) (string, error) {
	slackAuthResponse := new(SlackAuthResponse)
	err := getJson("https://slack.com/api/oauth.v2.access?code="+code+"&client_id="+clientId+"&client_secret="+clientSecret, slackAuthResponse)
	if err != nil {
		return "", err
	}
	return slackAuthResponse.Access_token, nil
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func server() {
	http.HandleFunc("/", NewWorkspace)
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
