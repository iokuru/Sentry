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
