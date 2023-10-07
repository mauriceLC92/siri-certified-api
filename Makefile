infra:
	cd infrastructure && pulumi up --stack dev --yes

destroy:
	cd infrastructure && pulumi destroy --stack dev --yes

build-companies:
	cd handlers/companies && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go

zip-companies:
	cd handlers/companies && zip companies.zip bootstrap company-data.json

build-users:
	cd handlers/users && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go

zip-users:
	cd handlers/users && zip users.zip bootstrap

build-auth:
	cd handlers/auth && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go && zip auth.zip bootstrap

zip-auth:
	cd handlers/auth && zip auth.zip bootstrap