variable "aws_region" {
  type        = string
  description = "aws region for sentry resources"
  default     = "us-east-1"
}

variable "environment" {
  type        = string
  description = "deployment environment"
  default     = "prod"
}

variable "project_name" {
  type        = string
  description = "project identifier"
  default     = "sentry"
}

variable "budget_alert_email" {
  type        = string
  description = "email recipient for zero dollar budget warnings"
  default     = "caregiver@example.com"
}

variable "use_comprehend" {
  type        = string
  description = "comprehend mode: auto, on, or off"
  default     = "auto"
}
