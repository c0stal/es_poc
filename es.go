package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/olivere/elastic"
	uuid "github.com/satori/go.uuid"
	"github.com/sha1sum/aws_signing_client"
)

var awsClient *http.Client
var client *elastic.Client
const (
	index = "report"
)

// ElasticIndex ...
type ElasticIndex struct {
	FieldStr string `json:"field str"`
}

func init() {
	creds := credentials.NewStaticCredentials("access_key", "secret_key", "")
	signer := v4.NewSigner(creds)
	var err error
	awsClient, err = aws_signing_client.New(signer, nil, "es", "us-east-1")
	if err != nil {
		log.Println(err)
	}

}
func main() {

	u1 := uuid.Must(uuid.NewV4(), nil)
	fmt.Printf("UUIDv4: %s\n", u1)
	val:=computeHmac256(u1.String(), "")
	fmt.Println("final: ", val)

	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()

	// Check if the Elasticsearch index already exists
	exist, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Fatalf("IndexExists() ERROR: %v", err)

		// Index some documents if the index exists
	} else if exist {

		// Instantiate new Elasticsearch documents from the ElasticIndex struct
		newDoc := ElasticIndex{
			FieldStr: val,
		}
		indexResult, err :=client.Index().
		Index(index).
		Id("1").
		BodyJson(newDoc).
		Do(ctx)
	}
}
func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getClient() *elastic.Client {
	var err error
	client, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(awsClient),
		elastic.SetURL(""),
	)

	fmt.Println("elastic.NewClient() ERROR: %v", err)
	return client
}
