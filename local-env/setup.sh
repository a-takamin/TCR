#! /bin/sh
aws configure set aws_access_key_id fake
aws configure set aws_secret_access_key fakefake

aws dynamodb create-table \
  --region \
      ap-northeast-1 \
  --endpoint-url \
      http://dynamodb-local:8000 \
  --table-name \
      dynamodb-local-table \
  --attribute-definitions \
      AttributeName=Digest,AttributeType=S \
      AttributeName=Tag,AttributeType=S \
  --key-schema \
      AttributeName=Digest,KeyType=HASH \
  --billing-mode \
      PAY_PER_REQUEST \
  --global-secondary-indexes \
      '[
        {
          "IndexName": "ManifestTagIndex",
          "KeySchema": [
            {
              "AttributeName":"Tag","KeyType":"HASH"
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
      blob-upload-progress \
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
      blob-local