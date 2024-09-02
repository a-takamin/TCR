#! /bin/sh
aws configure set aws_access_key_id fake
aws configure set aws_secret_access_key fake

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
          "IndexName": "ManifestDigestIndex",
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
