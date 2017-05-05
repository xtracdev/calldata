# calldata

This project provides a type to capture data related to call timings,
and a data pipeline to scrape call timings from logs dropped into an
S3 bucket and run them through a simple pipeline to store them in Amazon
Redshift.

The pipeline is triggered by S3 ObjectCreated put notification events
on an S3 bucket, which invokes a lambda function that scrapes the 
events from the file. Call records and service call data is written
to a Firehose stream (one for call records, one for service call timings),
which has Redshift as the sink.

The cloudformation directory contains the cloud formation templates, and
an iPython notebook that uses the templates to build the pipeline
stages. The schema is installed in Redshift using flyway - the 
flyway directory contains the config, migrations, etc.

## Contributing

To contribute, you must certify you agree with the [Developer Certificate of Origin](http://developercertificate.org/)
by signing your commits via `git -s`. To create a signature, configure your user name and email address in git.
Sign with your real name, do not use pseudonyms or submit anonymous commits.


In terms of workflow:

0. For significant changes or improvement, create an issue before commencing work.
1. Fork the respository, and create a branch for your edits.
2. Add tests that cover your changes, unit tests for smaller changes, acceptance test
for more significant functionality.
3. Run gofmt on each file you change before committing your changes.
4. Run golint on each file you change before committing your changes.
5. Make sure all the tests pass before committing your changes.
6. Commit your changes and issue a pull request.

## License

(c) 2017 Fidelity Investments
Licensed under the Apache License, Version 2.0