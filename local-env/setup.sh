#! /bin/sh
aws configure set aws_access_key_id fake
aws configure set aws_secret_access_key fakefake

aws dynamodb create-table \
  --region \
      ap-northeast-1 \
  --endpoint-url \
      http://dynamodb-local:8000 \
  --table-name \
      tcr-manifest-local \
  --attribute-definitions \
      AttributeName=Name,AttributeType=S \
      AttributeName=Digest,AttributeType=S \
      AttributeName=Tag,AttributeType=S \
  --key-schema \
      AttributeName=Name,KeyType=HASH \
      AttributeName=Digest,KeyType=RANGE \
  --billing-mode \
      PAY_PER_REQUEST \
  --local-secondary-indexes \
      '[
        {
          "IndexName": "ManifestTagIndex",
          "KeySchema": [
            {
              "AttributeName":"Name","KeyType":"HASH"
            },
            {
              "AttributeName":"Tag","KeyType":"RANGE"
            }
          ],
          "Projection": {
            "ProjectionType": "INCLUDE",
            "NonKeyAttributes": ["Manifest"]
          }
        }
      ]'

aws dynamodb create-table \
  --region \
      ap-northeast-1 \
  --endpoint-url \
      http://dynamodb-local:8000 \
  --table-name \
      tcr-blob-upload-progress-local \
  --attribute-definitions \
      AttributeName=Uuid,AttributeType=S \
  --key-schema \
      AttributeName=Uuid,KeyType=HASH \
  --billing-mode \
      PAY_PER_REQUEST 

aws dynamodb create-table \
  --region \
      ap-northeast-1 \
  --endpoint-url \
      http://dynamodb-local:8000 \
  --table-name \
      tcr-repository-local \
  --attribute-definitions \
      AttributeName=Name,AttributeType=S \
  --key-schema \
      AttributeName=Name,KeyType=HASH \
  --billing-mode \
      PAY_PER_REQUEST 

aws s3api create-bucket \
  --region \
      ap-northeast-1 \
  --endpoint-url \
      http://s3-local:9000 \
  --bucket \
      tcr-blob-local
