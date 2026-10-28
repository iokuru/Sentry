variable "name_prefix" {
  type        = string
  description = "naming prefix"
}

variable "budget_alert_email" {
  type        = string
  description = "email recipient for budget warnings"
  default     = "caregiver@example.com"
}

variable "uploads_bucket" {
  type        = string
  description = "s3 uploads bucket name"
  default     = ""
}

variable "incidents_table" {
  type        = string
  description = "dynamodb incidents table name"
  default     = ""
}

variable "links_table" {
  type        = string
  description = "dynamodb links table name"
  default     = ""
}
