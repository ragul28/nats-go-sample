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
	NATS_URL = "localhost:4223"
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
}
func JetStreamInit() (nats.JetStreamContext, error) {
	// Connect to NATS
	nc, err := nats.Connect(NATS_URL)
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
