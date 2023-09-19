infra:
	cd infrastructure && pulumi up --stack dev --yes

destroy:
	cd infrastructure && pulumi destroy --stack dev --yes

build-companies:
	cd handlers/companies && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go

zip-companies:
	cd handlers/companies && zip companies.zip bootstrap