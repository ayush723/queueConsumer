package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"consumer/db"
	"consumer/models"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
)

type request struct {
	Url string `json:"url"`
}

func main() {
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()
	channelAmqp, _ := amqpConnection.Channel()
	defer channelAmqp.Close()

	forever := make(chan bool)
	msgs, err := channelAmqp.Consume(os.Getenv("RABBITMQ_QUEUE"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var req request
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("RSS URL:", req.Url)
			entries, err := GetFeedEntries(req.Url)
			if err != nil {
				log.Fatal(err)
			}
			var ctx context.Context
			collection := db.Client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
			for _, entry := range entries[2:] {
				collection.InsertOne(ctx, bson.M{
					"title":     entry.Title,
					"thumbnail": entry.Thumbnail.URL,
					"url":       entry.Link.Href,
				})
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func GetFeedEntries(url string) ([]models.Entry, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/83.0.4103.116 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	byteValue, _ := ioutil.ReadAll(resp.Body)
	var feed models.Feed
	xml.Unmarshal(byteValue, &feed)
	return feed.Entries, nil

}
