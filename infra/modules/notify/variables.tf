variable "name_prefix" {
  type        = string
  description = "naming prefix"
}

variable "budget_alert_email" {
  type        = string
  description = "email recipient for budget warnings"
  default     = "caregiver@example.com"
}

variable "use_comprehend" {
  type        = string
  description = "comprehend mode: auto, on, or off"
  default     = "auto"
}

variable "ocr_mode" {
  type        = string
  description = "ocr processing mode: auto, textract, or browser"
  default     = "auto"
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

variable "uploads_bucket_arn" {
  type        = string
  description = "arn of uploads s3 bucket"
  default     = ""
}

variable "incidents_table_arn" {
  type        = string
  description = "arn of dynamodb incidents table"
  default     = ""
}

variable "links_table_arn" {
  type        = string
  description = "arn of dynamodb links table"
  default     = ""
}
