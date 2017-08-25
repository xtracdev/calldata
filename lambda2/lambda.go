package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bitly/go-simplejson"
	"github.com/xtracdev/calldata"
	"github.com/xtracdev/calldata/apitimings"
	"io"
	"log"
	"os"
	"strings"
)

const (
	callRecordStreamBucketEnv = "CALL_RECORD_BUCKET_NAME"
)

var (
	callRecordBucketName string
)


func processCallRecord(s3svc *s3.S3, key, bucket string, lineBuffer *bytes.Buffer) {
	records := string(lineBuffer.Bytes())

	input := &s3.PutObjectInput{
		Body: aws.ReadSeekCloser(strings.NewReader(records)),
		Bucket:aws.String(bucket),
		Key:aws.String(key),
	}

	_, err := s3svc.PutObject(input)

	if err != nil {
		fmt.Println("Error putting record to bucket", err.Error())
	}
}

func processBody(s3svc *s3.S3, key, bucket string, body io.Reader) error {

	var reader io.Reader = body

	gzipReader, err := gzip.NewReader(body)
	if err == nil {
		fmt.Println("assuming gzip encoding")
		reader = gzipReader
		defer gzipReader.Close()
	} else {
		fmt.Println("Error creating gzip reader... read as uncompressed", err.Error())
	}

	sr, err := sumoreader.NewSumoReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	var lineBuffer bytes.Buffer
	emptyBuffer := true
	for sr.Scan() {
		line := sr.Text()
		if strings.Contains(line, "{") {
			at, err := apitimings.NewAPITimingRec(line)
			if err != nil {
				fmt.Println("Parse error, not a timing record... skipping")
				continue
			}

			isCallRec, err := at.IsCallRecord()
			if  err != nil || isCallRec== false {
				fmt.Println("Not a timing record... skipping")
				continue
			}

			cr, err := at.CallRecord()
			if err != nil {
				fmt.Println("Not a call record... skipping")
				continue
			}

			if emptyBuffer {
				emptyBuffer = false
				fmt.Fprintf(&lineBuffer, "%s\n",at.Header())
			}

			fmt.Printf("call record:\n%s\n", cr)
			fmt.Fprintf(&lineBuffer, "%s\n", cr)

		}
	}

	if lineBuffer.Len() > 0 {
		fmt.Println("Writing bytes to s3")
		processCallRecord(s3svc, key, bucket, &lineBuffer)
	} else {
		fmt.Println("No bytes to write")
	}

	if err := sr.Err(); err != nil {
		return err
	}

	return nil
}

func main() {

	callRecordBucketName = os.Getenv(callRecordStreamBucketEnv)
	if callRecordBucketName == "" {
		log.Fatalf("No value in the environment for %s", callRecordStreamBucketEnv)
	}

	fmt.Printf("Go binary called with args %v\n", os.Args)
	buf := bytes.NewBuffer([]byte(os.Args[1]))
	js, err := simplejson.NewFromReader(buf)
	if err != nil {
		log.Fatal(err.Error)
	}

	records := js.Get("Records")

	arr, err := records.Array()
	if err != nil {
		log.Fatal(err)
	}

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	s3svc := s3.New(sess)

	for i := 0; i < len(arr); i++ {
		s3Rec := records.GetIndex(i).Get("s3")

		key := s3Rec.Get("object").Get("key")
		bucket := s3Rec.Get("bucket")
		arn := bucket.Get("arn")
		bucketName := bucket.Get("name")

		fmt.Printf("process %s in bucket %s (%s)\n", key.MustString(), bucketName.MustString(), arn.MustString())

		params := &s3.GetObjectInput{
			Bucket: aws.String(bucketName.MustString()),
			Key:    aws.String(key.MustString()),
		}

		resp, err := s3svc.GetObject(params)
		if err != nil {
			fmt.Printf("Error on GetObject: %s\n", err.Error())
			continue
		}

		if resp.Body == nil {
			fmt.Println("Nil body - nothing to read.")
			continue
		}

		defer resp.Body.Close()

		err = processBody(s3svc, key.MustString(), callRecordBucketName, resp.Body)
		if err != nil {
			fmt.Printf("Error processing body: %s\n", err.Error())
		}

	}

}
