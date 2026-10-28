# IAM role for api lambdas
resource "aws_iam_role" "api_lambda" {
  name = "${var.name_prefix}-api-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_policy" "api_lambda" {
  name = "${var.name_prefix}-api-lambda-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "CloudWatchLogs"
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Sid    = "S3Uploads"
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject"
        ]
        Resource = "*"
      },
      {
        Sid    = "DynamoStore"
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = "*"
      },
      {
        Sid    = "SqsSend"
        Effect = "Allow"
        Action = [
          "sqs:SendMessage"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "api_lambda" {
  role       = aws_iam_role.api_lambda.name
  policy_arn = aws_iam_policy.api_lambda.arn
}

# presign lambda
resource "aws_lambda_function" "presign" {
  function_name = "${var.name_prefix}-presign"
  role          = aws_iam_role.api_lambda.arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  memory_size   = 128
  timeout       = 10
  filename      = "${path.root}/../bin/presign.zip"

  environment {
    variables = {
      UPLOADS_BUCKET = var.uploads_bucket
    }
  }
}

resource "aws_apigatewayv2_integration" "presign" {
  api_id                 = aws_apigatewayv2_api.http.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.presign.arn
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "presign" {
  statement_id  = "AllowAPIGatewayPresign"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.presign.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http.execution_arn}/*/*"
}

# intake lambda
resource "aws_lambda_function" "intake" {
  function_name = "${var.name_prefix}-intake"
  role          = aws_iam_role.api_lambda.arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  memory_size   = 128
  timeout       = 10
  filename      = "${path.root}/../bin/intake.zip"

  environment {
    variables = {
      ANALYSIS_QUEUE_URL = var.analysis_queue_url
      INCIDENTS_TABLE    = var.incidents_table
      LINKS_TABLE        = var.links_table
    }
  }
}

resource "aws_apigatewayv2_integration" "intake" {
  api_id                 = aws_apigatewayv2_api.http.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.intake.arn
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "intake" {
  statement_id  = "AllowAPIGatewayIntake"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.intake.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http.execution_arn}/*/*"
}

# caregiver list and incident lookup lambda
resource "aws_lambda_function" "caregiver_list" {
  function_name = "${var.name_prefix}-caregiver_list"
  role          = aws_iam_role.api_lambda.arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  memory_size   = 128
  timeout       = 10
  filename      = "${path.root}/../bin/caregiver_list.zip"

  environment {
    variables = {
      INCIDENTS_TABLE = var.incidents_table
      LINKS_TABLE     = var.links_table
    }
  }
}

resource "aws_apigatewayv2_integration" "caregiver_list" {
  api_id                 = aws_apigatewayv2_api.http.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.caregiver_list.arn
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "caregiver_list" {
  statement_id  = "AllowAPIGatewayCaregiverList"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.caregiver_list.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http.execution_arn}/*/*"
}

# caregiver status lambda
resource "aws_lambda_function" "caregiver_status" {
  function_name = "${var.name_prefix}-caregiver_status"
  role          = aws_iam_role.api_lambda.arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  memory_size   = 128
  timeout       = 10
  filename      = "${path.root}/../bin/caregiver_status.zip"

  environment {
    variables = {
      INCIDENTS_TABLE = var.incidents_table
      LINKS_TABLE     = var.links_table
    }
  }
}

resource "aws_apigatewayv2_integration" "caregiver_status" {
  api_id                 = aws_apigatewayv2_api.http.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.caregiver_status.arn
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "caregiver_status" {
  statement_id  = "AllowAPIGatewayCaregiverStatus"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.caregiver_status.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http.execution_arn}/*/*"
}
