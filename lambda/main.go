package main

import (
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "encoding/json"
    "os"
    "projects/serverless-go/s3"

    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
)

type S3Data struct {
    Name string `json:"name"`
    Value string `json:"value"`
}

type Body struct {
    Message string `json:"message"`
    Release string `json:"release"`
    Bucket  string `json:"bucket"`
    S3Data  S3Data `json:"s3data"`
}

type ErrorBody struct {
    Message string `json:"message"`
}

func getS3Data(s3Bucket, s3Filename, region string) (S3Data, error) {

    var s3data S3Data

    options := session.Options{
        Config: aws.Config{Region: aws.String(region)},
    }

    s3Reader := s3.NewReader(options)

    raw, err := s3Reader.ReadBytes(s3Bucket, s3Filename)

    if err != nil {
        return s3data, err
    }

    err = json.Unmarshal(raw, &s3data)

    return s3data, err
}


func makeBody(msg string) (Body, error) {

    var body Body

    s3bucket := os.Getenv("BUCKET")
    s3Filename := os.Getenv("FILENAME")
    region := os.Getenv("REGION")

    s3Data, err := getS3Data(s3bucket, s3Filename, region)

    if err != nil {
        return body, err
    }

    return Body{
        Message: msg,
        Release: os.Getenv("RELEASE"),
        Bucket: s3bucket,
        S3Data: s3Data,
    }, nil
}

func makeJsonBody(msg string) (string, int) {

    body, err := makeBody(msg)

    if err != nil {

        errorBody := ErrorBody{
            Message: err.Error(),
        }

        jsonErr, _ := json.Marshal(errorBody)

        return string(jsonErr), 500
    }

    jsonBody, _ := json.Marshal(body)

    return string(jsonBody), 200
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    body, code := makeJsonBody("Mock API is literally GO!")

    return events.APIGatewayProxyResponse{
        Body: body,
        StatusCode: code,
        Headers: map[string]string{"Access-Control-Allow-Origin": "*"},
    }, nil
}

func main() {
    lambda.Start(Handler)
}
