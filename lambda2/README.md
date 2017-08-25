GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o lambda

zip lambda.zip index.js lambda

Build using the Makefile, e.g. make

Input:

Go binary called with args [./lambda {"Records":[{"eventVersion":"2.0","eventTime":"1970-01-01T00:00:00.000Z","requestParameters":{"sourceIPAddress":"127.0.0.1"},"s3":{"configurationId":"testConfigRule","object":{"eTag":"0123456789abcdef0123456789abcdef","sequencer":"0A1B2C3D4E5F678901","key":"HappyFace.jpg","size":1024},"bucket":{"arn":"arn:aws:s3:::mybucket","name":"sourcebucket","ownerIdentity":{"principalId":"EXAMPLE"}},"s3SchemaVersion":"1.0"},"responseElements":{"x-amz-id-2":"EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH","x-amz-request-id":"EXAMPLE123456789"},"awsRegion":"us-east-1","eventName":"ObjectCreated:Put","userIdentity":{"principalId":"EXAMPLE"},"eventSource":"aws:s3"}]}]

<pre>
{
  "Records":[
    {
      "eventVersion":"2.0",
      "eventTime":"1970-01-01T00:00:00.000Z",
      "requestParameters":{
        "sourceIPAddress":"127.0.0.1"
      },
      "s3":{
        "configurationId":"testConfigRule",
        "object":{
          "eTag":"0123456789abcdef0123456789abcdef",
          "sequencer":"0A1B2C3D4E5F678901",
          "key":"HappyFace.jpg",
          "size":1024
        },
        "bucket":{
          "arn":"arn:aws:s3:::mybucket",
          "name":"sourcebucket",
          "ownerIdentity":{
            "principalId":"EXAMPLE"
          }
        },
        "s3SchemaVersion":"1.0"
      },
      "responseElements":{
        "x-amz-id-2":"EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH",
        "x-amz-request-id":"EXAMPLE123456789"
      },
      "awsRegion":"us-east-1",
      "eventName":"ObjectCreated:Put",
      "userIdentity":{
        "principalId":"EXAMPLE"
      },
      "eventSource":"aws:s3"
    }
  ]
}
</pre>

For testing on the command line locally:

<pre>
./lambda '{"Records":[{"eventVersion":"2.0","eventTime":"1970-01-01T00:00:00.000Z","requestParameters":{"sourceIPAddress":"127.0.0.1"},"s3":{"configurationId":"testConfigRule","object":{"eTag":"0123456789abcdef0123456789abcdef","sequencer":"0A1B2C3D4E5F678901","key":"HappyFace.jpg","size":1024},"bucket":{"arn":"arn:aws:s3:::mybucket","name":"sourcebucket","ownerIdentity":{"principalId":"EXAMPLE"}},"s3SchemaVersion":"1.0"},"responseElements":{"x-amz-id-2":"EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH","x-amz-request-id":"EXAMPLE123456789"},"awsRegion":"us-east-1","eventName":"ObjectCreated:Put","userIdentity":{"principalId":"EXAMPLE"},"eventSource":"aws:s3"}]}'
</pre>

Dependency:

go get github.com/bitly/go-simplejson


Notes to my future self

To run without errors for sumo dump configure the function with 256MB  and 30s timeout.

Need policies with s3 access for reading and writing, and the lambda
stuff as well, for instance:

<pre>
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "logs:CreateLogGroup",
            "Resource": "arn:aws:logs:us-west-2:<account no>:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:us-west-2:<account no>:log-group:/aws/lambda/myTestFn:*"
            ]
        }
    ]
}

{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1503590552000",
            "Effect": "Allow",
            "Action": [
                "s3:*"
            ],
            "Resource": [
                "arn:aws:s3:::xtds-lambda*"
            ]
        }
    ]
}
</pre>

And for a dead letter queue access:

<pre>
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1503654687000",
            "Effect": "Allow",
            "Action": [
                "sqs:SendMessage",
                "sqs:SendMessageBatch"
            ],
            "Resource": [
                "arn:aws:sqs:us-west-2:930295567417:myTestFnDLQ"
            ]
        }
    ]
}
</pre>


A crawler is used to populate a table that Athena can query:

<pre>
aws glue get-crawler --name call-data-crawler
{
    "Crawler": {
        "CrawlElapsedTime": 195000, 
        "Name": "call-data-crawler", 
        "Schedule": {
            "State": "SCHEDULED", 
            "ScheduleExpression": "cron(0/5 * * * ? *)"
        }, 
        "CreationTime": 1503649256.0, 
        "LastUpdated": 1503651037.0, 
        "Targets": {
            "JdbcTargets": [], 
            "S3Targets": [
                {
                    "Path": "s3://xtds-lambda-sink/", 
                    "Exclusions": []
                }
            ]
        }, 
        "LastCrawl": {
            "Status": "SUCCEEDED", 
            "LogStream": "call-data-crawler", 
            "MessagePrefix": "53e0a4d7-ee6f-4c6f-b6df-80cb8cba4c92", 
            "StartTime": 1503656446.0, 
            "LogGroup": "/aws-glue/crawlers"
        }, 
        "State": "STOPPING", 
        "Version": 2, 
        "Role": "arn:aws:iam::<account no>:role/s3-glue-crawler-role", 
        "DatabaseName": "apicalls", 
        "SchemaChangePolicy": {
            "DeleteBehavior": "DEPRECATE_IN_DATABASE", 
            "UpdateBehavior": "UPDATE_IN_DATABASE"
        }, 
        "Classifiers": []
    }
}
</pre>


Policy in the crawler role: AWSS3FullAccess and AWSGlueServiceRole policies are attached
to the rol


