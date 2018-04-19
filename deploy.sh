#!/usr/bin/env bash

set -euo pipefail

timestamp=$(date +"%y%m%d%H%M")
release=$timestamp

cf_bucket=cf-import-${timestamp}
cf_stack=api-stack-${release}

api_bucket=api-mock-data-${release}
api_filename=data.json

cd $( dirname $0 )

run_dir=$(pwd)

aws s3api create-bucket --bucket $cf_bucket --create-bucket-configuration LocationConstraint=eu-west-2

package_yml=$run_dir/package.yml

aws-sam-local package \
       --template-file $run_dir/template.yml \
       --s3-bucket $cf_bucket \
       --output-template-file $package_yml

aws cloudformation deploy \
       --template-file $package_yml \
       --stack-name $cf_stack \
       --capabilities CAPABILITY_IAM \
       --parameter-overrides Release=${release},S3Bucket=${api_bucket},S3Filename=${api_filename}

aws s3 cp $run_dir/mock/${api_filename} s3://${api_bucket}/${api_filename}

aws s3 rm s3://$cf_bucket --recursive

aws s3 rb s3://$cf_bucket



