resource "aws_ssm_parameter" "daily_cap_textract" {
  name        = "/${var.name_prefix}/config/daily_cap_textract"
  description = "daily maximum textract requests to stay within free tier"
  type        = "String"
  value       = "100"
  tier        = "Standard"
}

resource "aws_ssm_parameter" "daily_cap_comprehend" {
  name        = "/${var.name_prefix}/config/daily_cap_comprehend"
  description = "daily maximum comprehend requests to stay within free tier"
  type        = "String"
  value       = "100"
  tier        = "Standard"
}

resource "aws_ssm_parameter" "hotline_overrides" {
  name        = "/${var.name_prefix}/config/hotline_overrides"
  description = "json dictionary overriding regional hotlines"
  type        = "String"
  value       = "{}"
  tier        = "Standard"
}

resource "aws_ssm_parameter" "risk_threshold_high" {
  name        = "/${var.name_prefix}/config/risk_threshold_high"
  description = "score threshold for high risk caregiver alerts"
  type        = "String"
  value       = "70"
  tier        = "Standard"
}

resource "aws_ssm_parameter" "telegram_token" {
  name        = "/sentry/telegram_token"
  description = "telegram bot token for caregiver alerts"
  type        = "SecureString"
  value       = "PLACEHOLDER"
  tier        = "Standard"

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "telegram_chat_id" {
  name        = "/sentry/telegram_chat_id"
  description = "telegram chat id for caregiver alerts"
  type        = "SecureString"
  value       = "PLACEHOLDER"
  tier        = "Standard"

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "telegram_webhook_secret" {
  name        = "/sentry/telegram_webhook_secret"
  description = "telegram webhook verification secret"
  type        = "SecureString"
  value       = "PLACEHOLDER"
  tier        = "Standard"

  lifecycle {
    ignore_changes = [value]
  }
}
