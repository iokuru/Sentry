variable "name_prefix" {
  type        = string
  description = "naming prefix"
}

variable "budget_alert_email" {
  type        = string
  description = "email recipient for budget warnings"
  default     = "caregiver@example.com"
}
