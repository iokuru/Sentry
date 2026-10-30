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

variable "enable_cloudfront" {
  type        = bool
  description = "whether to provision cloudfront distribution"
  default     = true
}

variable "ocr_mode" {
  type        = string
  description = "ocr processing mode: auto, textract, or browser"
  default     = "auto"
}
