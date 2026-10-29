resource "aws_apigatewayv2_api" "http" {
  name          = "${var.name_prefix}-http-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PATCH", "OPTIONS"]
    allow_headers = ["Content-Type", "Authorization"]
    max_age       = 3600
  }
}

resource "aws_apigatewayv2_authorizer" "jwt" {
  api_id           = aws_apigatewayv2_api.http.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "${var.name_prefix}-caregiver-jwt"

  jwt_configuration {
    audience = [var.cognito_client_id]
    issuer   = "https://${var.cognito_user_pool_endpoint}"
  }
}

resource "aws_apigatewayv2_route" "presign" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "POST /api/presign"
  target    = "integrations/${aws_apigatewayv2_integration.presign.id}"
}

resource "aws_apigatewayv2_route" "incidents" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "GET /api/incidents"
  target    = "integrations/${aws_apigatewayv2_integration.caregiver_list.id}"
}

resource "aws_apigatewayv2_route" "caregiver_incidents" {
  api_id             = aws_apigatewayv2_api.http.id
  route_key          = "GET /api/caregiver/incidents"
  target             = "integrations/${aws_apigatewayv2_integration.caregiver_list.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.jwt.id
}

resource "aws_apigatewayv2_route" "caregiver_status" {
  api_id             = aws_apigatewayv2_api.http.id
  route_key          = "POST /api/caregiver/status"
  target             = "integrations/${aws_apigatewayv2_integration.caregiver_status.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.jwt.id
}

resource "aws_apigatewayv2_route" "telegram_hook" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "POST /telegram/hook"
  target    = "integrations/${aws_apigatewayv2_integration.telegram_hook.id}"
}

resource "aws_apigatewayv2_route" "caregiver_invites" {
  api_id             = aws_apigatewayv2_api.http.id
  route_key          = "POST /api/caregiver/invites"
  target             = "integrations/${aws_apigatewayv2_integration.caregiver_invite.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.jwt.id
}

resource "aws_apigatewayv2_route" "caregiver_links" {
  api_id             = aws_apigatewayv2_api.http.id
  route_key          = "GET /api/caregiver/links"
  target             = "integrations/${aws_apigatewayv2_integration.caregiver_invite.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.jwt.id
}

resource "aws_apigatewayv2_route" "caregiver_links_delete" {
  api_id             = aws_apigatewayv2_api.http.id
  route_key          = "DELETE /api/caregiver/links/{senior_id}"
  target             = "integrations/${aws_apigatewayv2_integration.caregiver_invite.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.jwt.id
}

resource "aws_apigatewayv2_route" "senior_pair" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "POST /api/senior/pair"
  target    = "integrations/${aws_apigatewayv2_integration.senior_pair.id}"
}
