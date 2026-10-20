resource "aws_dynamodb_table" "links" {
  name         = "${var.name_prefix}-links"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "caregiver_id"
  range_key    = "senior_id"

  attribute {
    name = "caregiver_id"
    type = "S"
  }

  attribute {
    name = "senior_id"
    type = "S"
  }

  global_secondary_index {
    name            = "senior_id-index"
    hash_key        = "senior_id"
    range_key       = "caregiver_id"
    projection_type = "ALL"
  }

  server_side_encryption {
    enabled = true
  }
}
