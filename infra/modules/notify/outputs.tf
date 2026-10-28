output "analysis_queue_url" {
  value = aws_sqs_queue.analysis.url
}

output "analysis_queue_arn" {
  value = aws_sqs_queue.analysis.arn
}

output "alerts_topic_arn" {
  value = aws_sns_topic.caregiver_alerts.arn
}

output "analysis_role_arn" {
  value = aws_iam_role.analysis.arn
}

output "analysis_role_name" {
  value = aws_iam_role.analysis.name
}

output "scheduler_role_arn" {
  value = aws_iam_role.scheduler.arn
}
