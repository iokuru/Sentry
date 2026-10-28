output "cloudfront_url" {
  description = "cloudfront domain url serving the sentry frontend"
  value       = "https://${module.api.cloudfront_domain_name}"
}

output "api_url" {
  description = "http api gateway endpoint"
  value       = module.api.api_endpoint
}

output "uploads_bucket" {
  description = "private s3 uploads bucket for screenshots"
  value       = module.storage.uploads_bucket_id
}

output "user_pool_id" {
  description = "cognito user pool id for caregivers"
  value       = module.auth.user_pool_id
}

output "user_pool_client_id" {
  description = "cognito app client id for caregiver portal"
  value       = module.auth.user_pool_client_id
}

output "region" {
  description = "aws deployment region"
  value       = var.aws_region
}

output "frontend_bucket" {
  description = "s3 bucket id for frontend static hosting"
  value       = module.storage.frontend_bucket_id
}

output "cloudfront_distribution_id" {
  description = "id of cloudfront distribution for cache invalidation"
  value       = module.api.cloudfront_distribution_id
}
