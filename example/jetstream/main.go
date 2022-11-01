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
