locals {
  lambda_functions = [
    "presign",
    "intake",
    "analysis",
    "caregiver_list",
    "caregiver_status",
    "digest",
    "telegram_hook"
  ]
}

resource "aws_cloudwatch_log_group" "lambdas" {
  for_each          = toset(local.lambda_functions)
  name              = "/aws/lambda/${var.name_prefix}-${each.key}"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "api_gateway" {
  name              = "/aws/apigateway/${var.name_prefix}-http-api"
  retention_in_days = 7
}
