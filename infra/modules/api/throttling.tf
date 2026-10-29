resource "aws_apigatewayv2_route" "intake" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "POST /api/intake"
  target    = "integrations/${aws_apigatewayv2_integration.intake.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http.id
  name        = "$default"
  auto_deploy = true

  route_settings {
    route_key              = "POST /api/intake"
    throttling_burst_limit = 5
    throttling_rate_limit  = 2.0
  }

  route_settings {
    route_key              = "POST /api/presign"
    throttling_burst_limit = 5
    throttling_rate_limit  = 2.0
  }

  route_settings {
    route_key              = "POST /api/senior/pair"
    throttling_burst_limit = 5
    throttling_rate_limit  = 2.0
  }
}
