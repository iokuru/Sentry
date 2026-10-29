variable "name_prefix" {
  type        = string
  description = "naming prefix"
}

variable "frontend_bucket_id" {
  type        = string
  description = "id of frontend s3 bucket"
}

variable "frontend_bucket_arn" {
  type        = string
  description = "arn of frontend s3 bucket"
}

variable "cognito_user_pool_endpoint" {
  type        = string
  description = "endpoint for cognito jwt authorizer"
  default     = ""
}

variable "cognito_client_id" {
  type        = string
  description = "app client id for jwt audience"
  default     = ""
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

variable "analysis_queue_url" {
  type        = string
  description = "sqs queue url for intake submissions"
  default     = ""
}

variable "analysis_queue_arn" {
  type        = string
  description = "sqs queue arn for intake permissions"
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
