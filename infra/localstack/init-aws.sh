#!/usr/bin/env bash

set -euo pipefail

awslocal s3 mb s3://my-bucket

awslocal s3 cp /tmp/random_transactions.csv s3://my-bucket/random_transactions.csv
