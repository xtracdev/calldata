package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bitly/go-simplejson"
	"github.com/xtracdev/calldata"
	"github.com/xtracdev/calldata/apitimings"
	"io"
	"log"
	"os"
	"strings"
)

func putRecord(firehoseSvc *firehose.Firehose, streamName, record string) {
	params := &firehose.PutRecordInput{
		DeliveryStreamName: aws.String(streamName),
		Record: &firehose.Record{
			Data: []byte(record + "\n"),
		},
	}

	_, err := firehoseSvc.PutRecord(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func processCallRecord(firehoseSvc *firehose.Firehose, callRecord string) {
	putRecord(firehoseSvc, "call-record-stream", callRecord)
}

func processSvcCallRecord(firehoseSvc *firehose.Firehose, svcCall string) {
	putRecord(firehoseSvc, "svc-call-stream", svcCall)
}

func processBody(fireHoseSvc *firehose.Firehose, body io.Reader) error {

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

	for sr.Scan() {
		line := sr.Text()
		if strings.Contains(line, "{") {
			//fmt.Println(sr.Text())
			at, err := apitimings.NewAPITimingRec(line)
			if err != nil {
				fmt.Println("Not a timing record... skipping")
				continue
			}
			cr, err := at.CallRecord()
			if err != nil {
				fmt.Println("Not a call record... skipping")
				continue
			}

			fmt.Printf("call record:\n%s\n", cr)
			processCallRecord(fireHoseSvc, cr)

			calls, err := at.ServiceCalls()
			if err != nil {
				fmt.Println("No service call records... skipping")
				continue
			}

			for _, c := range calls {
				fmt.Printf("service call:\n%s\n", c)
				processSvcCallRecord(fireHoseSvc, c)
			}
		}
	}

	if err := sr.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
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
	fireHoseSvc := firehose.New(sess)

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

		err = processBody(fireHoseSvc, resp.Body)
		if err != nil {
			fmt.Printf("Error processing body: %s\n", err.Error())
		}

	}

}
