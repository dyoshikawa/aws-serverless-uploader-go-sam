package:
	sam package --template-file template.yaml --s3-bucket aws-serverless-uploader-go --output-template-file packaged.yaml

deploy:
	make build && \
	make package && \
	sam deploy --template-file packaged.yaml --stack-name aws-serverless-uploader-go --capabilities CAPABILITY_IAM

deps:
	make deps-images && \
	make deps-batch-destroy

build:
	make build-images

test:
	make test-images

deps-images:
	cd ./src/images && \
	make deps

build-images:
	cd ./src/images && \
	make build

test-images:
	cd ./src/images && \
	make test
	