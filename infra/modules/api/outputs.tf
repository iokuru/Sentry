output "api_endpoint" {
  value = aws_apigatewayv2_api.http.api_endpoint
}

output "api_id" {
  value = aws_apigatewayv2_api.http.id
}

output "cloudfront_distribution_id" {
  value = var.enable_cloudfront && length(aws_cloudfront_distribution.frontend) > 0 ? aws_cloudfront_distribution.frontend[0].id : ""
}

output "cloudfront_domain_name" {
  value = var.enable_cloudfront && length(aws_cloudfront_distribution.frontend) > 0 ? aws_cloudfront_distribution.frontend[0].domain_name : ""
}
