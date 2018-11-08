package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	pixela "github.com/gainings/pixela-go-client"
)

// Point is left & top positions of bounding box in the Rekognition result
type Point struct {
	Left float64
	Top  float64
}

// !! fixed number from experiment (maybe require to change your env) !!
var assumedDatePoint = Point{Left: 0.393, Top: 0.111}
var assumedQuantityPoint = Point{Left: 0.268, Top: 0.282}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, s3Event events.S3Event) error {
	// for each s3 object
	for _, record := range s3Event.Records {
		// extract s3 object info
		bucket, key := getS3ObjectFromRecord(record)
		fmt.Printf("[%s] Bucket = %s, Key = %s \n", record.EventSource, bucket, key)

		// execute text detection of Rekognition
		res, rekerr := exeRekognitionDetectText(bucket, key)
		if rekerr != nil {
			fmt.Println("Error")
			fmt.Println(rekerr.Error())
		}

		// extract date & quantity from the above result
		date, quantity := getValueFromRekognitionResult(res.TextDetections)
		fmt.Printf("data: %s, quantity: %s\n", date, quantity)

		// record pixel
		perr := recordPixel(date, quantity)
		fmt.Println(perr)
	}

	return nil
}

func getS3ObjectFromRecord(record events.S3EventRecord) (string, string) {
	s := record.S3
	bucket := s.Bucket.Name
	key := s.Object.Key

	key = strings.Replace(key, "+", " ", -1)
	key = strings.Replace(key, "%3A", ":", -1)
	key = strings.Replace(key, "%2C", ",", -1)

	return bucket, key
}

func exeRekognitionDetectText(bucket, key string) (*rekognition.DetectTextOutput, error) {
	// create Rekognition client
	sess := session.Must(session.NewSession())
	rc := rekognition.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

	// set params
	params := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(key),
			},
		},
	}
	fmt.Printf("params: %s", params)

	// execute DetectText
	return rc.DetectText(params)
}

func getValueFromRekognitionResult(results []*rekognition.TextDetection) (string, string) {
	dateHypot, quantityHypot := math.MaxFloat64, math.MaxFloat64
	date, quantity := "", ""

	// for each detected text
	for _, td := range results {
		left, top := *td.Geometry.BoundingBox.Left, *td.Geometry.BoundingBox.Top

		// calc hypot with assumed date pos & update value
		tmpDHypot := math.Hypot(math.Abs(left-assumedDatePoint.Left), math.Abs(top-assumedDatePoint.Top))
		if tmpDHypot < dateHypot {
			// if td is most-likely-result (nearest to the assumed point), keep the result (with removing "/")
			dateHypot, date = tmpDHypot, strings.Replace(*td.DetectedText, "/", "", -1)
		}

		// calc hypot with assumed quantity pos & update value
		tmpQHypot := math.Hypot(math.Abs(left-assumedQuantityPoint.Left), math.Abs(top-assumedQuantityPoint.Top))
		if tmpQHypot < quantityHypot {
			// if td is most-likely-result (nearest to the assumed point), keep the result
			quantityHypot, quantity = tmpQHypot, *td.DetectedText
		}
	}

	return date, quantity
}

func recordPixel(date, quantity string) error {
	user := os.Getenv("PIXELA_USER")
	token := os.Getenv("PIXELA_TOKEN")
	graph := os.Getenv("PIXELA_GRAPH")
	c := pixela.NewClient(user, token)

	// try to record
	err := c.RegisterPixel(graph, date, quantity)
	if err == nil {
		fmt.Println("recorded")
		return err
	}

	// if fail, try to update
	err = c.UpdatePixelQuantity(graph, date, quantity)
	if err == nil {
		fmt.Println("updated")
	}

	return err
}

func main() {
	lambda.Start(Handler)
}
