output "user_pool_id" {
  value = aws_cognito_user_pool.caregivers.id
}

output "user_pool_arn" {
  value = aws_cognito_user_pool.caregivers.arn
}

output "user_pool_client_id" {
  value = aws_cognito_user_pool_client.portal.id
}

output "user_pool_endpoint" {
  value = aws_cognito_user_pool.caregivers.endpoint
}
