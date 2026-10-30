variable "name_prefix" {
  type        = string
  description = "prefix for naming storage resources"
}

variable "allowed_origin" {
  type        = string
  description = "allowed origin for cors on uploads bucket"
  default     = ""
}
