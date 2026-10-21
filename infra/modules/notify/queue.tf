resource "aws_sqs_queue" "analysis_dlq" {
  name                      = "${var.name_prefix}-analysis-dlq"
  message_retention_seconds = 1209600
}

resource "aws_sqs_queue" "analysis" {
  name                       = "${var.name_prefix}-analysis"
  visibility_timeout_seconds = 120
  message_retention_seconds  = 86400

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.analysis_dlq.arn
    maxReceiveCount     = 3
  })
}

resource "aws_cloudwatch_metric_alarm" "dlq_depth" {
  alarm_name          = "${var.name_prefix}-analysis-dlq-depth"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = 60
  statistic           = "Sum"
  threshold           = 0
  alarm_description   = "triggers when analysis dead letter queue has failed jobs"

  dimensions = {
    QueueName = aws_sqs_queue.analysis_dlq.name
  }
}
