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
