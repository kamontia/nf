# nf Backend

This directory contains the AWS serverless backend for the `nf` notification service.

## Architecture

The backend is simple and designed to be cost-effective, running entirely within the AWS Free Tier for typical usage.

`nf CLI` -> `API Gateway (HTTP API)` -> `AWS Lambda (Go)` -> `AWS SNS` -> `Mobile App`

1.  The `nf` command-line tool makes an HTTP POST request to an API Gateway endpoint.
2.  API Gateway triggers a Go-based Lambda function.
3.  The Lambda function publishes the notification details to an SNS Topic.
4.  The user's mobile app subscribes to this SNS Topic to receive push notifications.

## Deployment Instructions

These instructions guide you through deploying the backend manually using the AWS CLI.

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or later)
- [AWS CLI](https://aws.amazon.com/cli/) configured with your credentials.

### Step 1: Build the Lambda Function

1.  Navigate to the `nf` project root directory.
2.  Build the Lambda binary. It must be compiled for Linux.

    ```sh
    # From the project root directory
    GOOS=linux GOARCH=amd64 go build -o backend/bootstrap ./backend/lambda
    ```

3.  Create a zip file containing the compiled binary.

    ```sh
    # From the project root directory
    zip backend/function.zip backend/bootstrap
    ```

### Step 2: Create the IAM Role

The Lambda function needs permission to write to SNS and CloudWatch Logs.

1.  Create a trust policy file named `trust-policy.json`:
    ```json
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "Service": "lambda.amazonaws.com"
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }
    ```
2.  Create the IAM role:
    ```sh
    aws iam create-role --role-name nf-lambda-role --assume-role-policy-document file://trust-policy.json
    ```
3.  Attach the necessary AWS managed policies:
    ```sh
    # For writing to SNS
    aws iam attach-role-policy --role-name nf-lambda-role --policy-arn arn:aws:iam::aws:policy/AmazonSNSFullAccess
    # For writing logs
    aws iam attach-role-policy --role-name nf-lambda-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
    ```
    *Note: For better security, you could create a custom policy that only allows the `sns:Publish` action on the specific SNS topic. The `nf setup-app` command also requires `sns:ListTopics` permission to find the topic ARN.*

### Step 3: Create the SNS Topic

Create an SNS topic that the Lambda will publish to.

```sh
aws sns create-topic --name nf-notifications
```
Take note of the `TopicArn` returned by this command. You will need it later.

### Step 4: Create the Lambda Function

1.  Get the ARN of the IAM role you created. Replace `<YOUR_ACCOUNT_ID>` with your AWS Account ID.
    ```sh
    ROLE_ARN="arn:aws:iam::<YOUR_ACCOUNT_ID>:role/nf-lambda-role"
    ```
2.  Get the ARN of the SNS topic you created.
    ```sh
    TOPIC_ARN="<YOUR_SNS_TOPIC_ARN>"
    ```
3.  Create the Lambda function:
    ```sh
    aws lambda create-function --function-name nf-notify-handler \
      --runtime provided.al2 --handler bootstrap \
      --role $ROLE_ARN \
      --zip-file fileb://backend/function.zip \
      --environment "Variables={SNS_TOPIC_ARN=$TOPIC_ARN}"
    ```

### Step 5: Create the API Gateway

1.  Create an HTTP API Gateway.
    ```sh
    aws apigatewayv2 create-api --name "nf-api" --protocol-type HTTP --target "arn:aws:apigateway:<YOUR_REGION>:lambda:path/2015-03-31/functions/arn:aws:lambda:<YOUR_REGION>:<YOUR_ACCOUNT_ID>:function:nf-notify-handler/invocations"
    ```
    Replace `<YOUR_REGION>` and `<YOUR_ACCOUNT_ID>`.
2.  The command will return an `ApiEndpoint`. This is the URL you will use in your `nf` CLI configuration (`api_url`).

### Step 6: Add API Gateway Permissions

Allow the API Gateway to invoke your Lambda function.

```sh
aws lambda add-permission \
  --function-name nf-notify-handler \
  --statement-id "apigateway-invoke-permission" \
  --action "lambda:InvokeFunction" \
  --principal "apigateway.amazonaws.com" \
  --source-arn "arn:aws:execute-api:<YOUR_REGION>:<YOUR_ACCOUNT_ID>:<YOUR_API_ID>/*/*"
```
Replace `<YOUR_REGION>`, `<YOUR_ACCOUNT_ID>`, and `<YOUR_API_ID>` (from the `create-api` output).

## Cleanup

To avoid ongoing charges (if any), remember to delete the created resources when you are done.
-   Delete the API Gateway.
-   Delete the Lambda function.
-   Delete the SNS topic.
-   Detach policies from and delete the IAM role.
-   Delete the `trust-policy.json` file.
-   Delete the `backend/function.zip` and `backend/bootstrap` files.
