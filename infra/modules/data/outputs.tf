output "incidents_table_name" {
  value = aws_dynamodb_table.incidents.name
}

output "incidents_table_arn" {
  value = aws_dynamodb_table.incidents.arn
}

output "links_table_name" {
  value = aws_dynamodb_table.links.name
}

output "links_table_arn" {
  value = aws_dynamodb_table.links.arn
}
