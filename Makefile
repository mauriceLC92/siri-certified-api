HANDLERS_DIR=handlers/authentication
TARGETS=sign-up confirm-verification-code log-in
GO_BUILD_ENV=GOARCH=arm64 GOOS=linux

.PHONY: build-authentication $(TARGETS) cognito destroy-cognito infra destroy

build-authentication: $(TARGETS)

$(TARGETS):
	cd $(HANDLERS_DIR)/$@ && \
	$(GO_BUILD_ENV) go build -o bootstrap main.go && \
	zip $@.zip bootstrap

cognito:
	cd infrastructure/cognito && pulumi up --stack dev --yes

destroy-cognito:
	cd infrastructure/cognito && pulumi destroy --stack dev --yes

infra:
	cd infrastructure && pulumi up --stack dev --yes

destroy:
	cd infrastructure && pulumi destroy --stack dev --yes

build-companies:
	cd handlers/companies && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go && zip companies.zip bootstrap company-data.json

build-users:
	cd handlers/users && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go && zip users.zip bootstrap

build-auth:
	cd handlers/auth && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go && zip auth.zip bootstrap

zip-auth:
	cd handlers/auth && zip auth.zip bootstrap