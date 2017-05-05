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