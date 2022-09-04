#!/bin/sh
aws dynamodb create-table --endpoint-url http://localhost:8000 --cli-input-json file://users_table.json
aws dynamodb batch-write-item --endpoint-url http://localhost:8000 --request-items file://users_table_data.json