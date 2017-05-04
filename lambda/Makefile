package:
	rm -f lambda
	rm -f lambda.zip
	go get github.com/bitly/go-simplejson
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o lambda
	zip lambda.zip index.js lambda

clean:
	rm -f lambda
	rm -f lambda.zip
