# SIRI Certified API

This repository contains application and infrastructure code for a serverless API built with Go and AWS Lambda.

This is the API powering the SIRI Certified Web app.

## Getting Started

If you want to create this API within your own AWS account, you can follow the following steps using the Makefile commands provided. You will need a [Pulumi](https://www.pulumi.com/) account and have the [Pulumi CLI](https://www.pulumi.com/docs/cli/) installed. The subsequent steps will not work without it.

1. Activate your AWS profile

```bash
export AWS_PROFILE=your-profile-name-here
```

2. Create the Coginto User Pool and associated Client App. Take note of the user pool ID and client app ID exported at the end. You will need those next.

```bash
make cognito
```

3. Create a `.env` file inside the `infrastructure` repository

```bash
touch infrastructure/.env
```

4. Add the following environment variables to your `.env` file. Use the values exported from step 2 above.

```
COGNITO_CLIENT_APP_ID=
COGNITO_ENDPOINT_URL=
```

5. Build the code for the authentication, companies and users Lambdas

```bash
make build-authentication && make build-companies && make build-users
```

6. Create the rest of the infrastructure. Take note of the API endpoint URL exported. This is the endpoint used for your API gateway

```bash
make infra
```

If you ever want to tear down the infrastructure, run the following:

```bash
make destroy-cognito && make destroy
```
