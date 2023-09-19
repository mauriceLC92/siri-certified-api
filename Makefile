infra:
	cd infrastructure && pulumi up --stack dev --yes

destroy:
	cd infrastructure && pulumi destroy --stack dev --yes

build:
	cd handler && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go

zip:
	cd handler && zip myFunction.zip bootstrap