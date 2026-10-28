resource "aws_sns_topic" "caregiver_alerts" {
  name = "${var.name_prefix}-caregiver-alerts"
}

resource "aws_sns_topic_policy" "caregiver_alerts" {
  arn = aws_sns_topic.caregiver_alerts.arn

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowLambdaPublish"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action   = "sns:Publish"
        Resource = aws_sns_topic.caregiver_alerts.arn
      }
    ]
  })
}
