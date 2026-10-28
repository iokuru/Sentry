resource "aws_iam_role" "scheduler" {
  name = "${var.name_prefix}-scheduler-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "scheduler.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_policy" "scheduler_invoke" {
  name = "${var.name_prefix}-scheduler-invoke-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "lambda:InvokeFunction"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "scheduler" {
  role       = aws_iam_role.scheduler.name
  policy_arn = aws_iam_policy.scheduler_invoke.arn
}

resource "aws_scheduler_schedule" "daily_digest" {
  name       = "${var.name_prefix}-daily-digest"
  group_name = "default"

  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression = "cron(0 9 * * ? *)"

  target {
    arn      = aws_lambda_function.digest.arn
    role_arn = aws_iam_role.scheduler.arn
  }
}
