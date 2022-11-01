package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	StreamName     = "REVIEWS"
	StreamSubjects = "REVIEWS.*"

	SubjectNameReviewCreated  = "REVIEWS.rateGiven"
	SubjectNameReviewAnswered = "REVIEWS.rateAnswer"
)

type Review struct {
	Id      string `json:"_id"`
	Author  string `json:"author"`
	Store   string `json:"store"`
	Text    string `json:"text"`
	Rating  int    `json:"rating"`
	Created string `json:"created"`
}

func main() {
	log.Println("Starting...")

	js, err := JetStreamInit()
	if err != nil {
		log.Println(err)
		return
	}

	// Let's assume that publisher and consumer are services running on different servers.
	// So run publisher and consumer asynchronously to see how it works
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		publishReviews(js)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeReviews(js)
	}()

	wg.Wait()

	log.Println("Exit...")
}

func JetStreamInit() (nats.JetStreamContext, error) {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	// Create a stream if it does not exist
	err = CreateStream(js)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(StreamName)

	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func publishReviews(js nats.JetStreamContext) {
	reviews, err := getReviews()
	if err != nil {
		log.Println(err)
		return
	}

	for _, oneReview := range reviews {

		// create random message intervals to slow down
		r := rand.Intn(1500)
		time.Sleep(time.Duration(r) * time.Millisecond)

		reviewString, err := json.Marshal(oneReview)
		if err != nil {
			log.Println(err)
			continue
		}

		// publish to REVIEWS.rateGiven subject
		_, err = js.Publish(SubjectNameReviewCreated, reviewString)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Publisher  =>  Message: %s\n", oneReview.Text)
		}
	}
}

func getReviews() ([]Review, error) {
	rawReviews, _ := ioutil.ReadFile("./reviews.json")
	var reviewsObj []Review
	err := json.Unmarshal(rawReviews, &reviewsObj)

	return reviewsObj, err
}

func consumeReviews(js nats.JetStreamContext) {
	_, err := js.Subscribe(SubjectNameReviewCreated, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			log.Println("Unable to Ack", err)
			return
		}

		var review Review
		err = json.Unmarshal(m.Data, &review)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Consumer  =>  Subject: %s  -  ID: %s  -  Author: %s  -  Rating: %d\n", m.Subject, review.Id, review.Author, review.Rating)

		// send answer via JetStream using another subject if you need
		// js.Publish(SubjectNameReviewAnswered, []byte(review.Id))
	})

	if err != nil {
		log.Println("Subscribe failed")
		return
	}
}
