resource "aws_lambda_function" "analysis" {
  function_name = "${var.name_prefix}-analysis"
  role          = aws_iam_role.analysis.arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  memory_size   = 256
  timeout       = 60
  filename      = "${path.root}/../bin/analysis.zip"

  environment {
    variables = {
      UPLOADS_BUCKET             = var.uploads_bucket
      CAREGIVER_ALERTS_TOPIC_ARN = aws_sns_topic.caregiver_alerts.arn
      INCIDENTS_TABLE            = var.incidents_table
      LINKS_TABLE                = var.links_table
      ANALYSIS_DLQ_URL           = aws_sqs_queue.analysis_dlq.url
    }
  }

  depends_on = [aws_cloudwatch_log_group.lambdas]
}

resource "aws_lambda_event_source_mapping" "analysis" {
  event_source_arn = aws_sqs_queue.analysis.arn
  function_name    = aws_lambda_function.analysis.arn
  batch_size       = 1
}

resource "aws_lambda_function" "digest" {
  function_name = "${var.name_prefix}-digest"
  role          = aws_iam_role.digest.arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  memory_size   = 128
  timeout       = 30
  filename      = "${path.root}/../bin/digest.zip"

  environment {
    variables = {
      CAREGIVER_ALERTS_TOPIC_ARN = aws_sns_topic.caregiver_alerts.arn
      INCIDENTS_TABLE            = var.incidents_table
      LINKS_TABLE                = var.links_table
    }
  }

  depends_on = [aws_cloudwatch_log_group.lambdas]
}
