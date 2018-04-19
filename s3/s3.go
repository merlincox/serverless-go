package s3

import (

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "bytes"
    "io"
    //"fmt"
)

type Reader struct {
    service *s3.S3
}

func (reader *Reader) ReadBytes(bucket string, key string) ([]byte, error) {

     results, err := reader.service.GetObject(&s3.GetObjectInput{
         Bucket: aws.String(bucket),
         Key: aws.String(key),
     })

    if err != nil {
        return nil, err
    }

    defer results.Body.Close()

    buf := bytes.NewBuffer(nil)

    if _, err := io.Copy(buf, results.Body); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

func NewReader(options session.Options) Reader {

    s3Session := session.Must(session.NewSessionWithOptions(options))

    reader := Reader{
        service: s3.New(s3Session),
    }

    return reader
}

