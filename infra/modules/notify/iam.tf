resource "aws_iam_role" "analysis" {
  name = "${var.name_prefix}-analysis-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

locals {
  comprehend_statement = var.use_comprehend != "off" ? [
    {
      Sid    = "ComprehendAnalysis"
      Effect = "Allow"
      Action = [
        "comprehend:DetectSentiment",
        "comprehend:DetectKeyPhrases",
        "comprehend:DetectPiiEntities",
        "comprehend:DetectDominantLanguage"
      ]
      Resource = "*"
    }
  ] : []

  textract_statement = var.ocr_mode != "browser" ? [
    {
      Sid    = "TextractDetectText"
      Effect = "Allow"
      Action = [
        "textract:DetectDocumentText"
      ]
      Resource = "*"
    }
  ] : []

  analysis_statements = concat([
    {
      Sid    = "CloudWatchLogs"
      Effect = "Allow"
      Action = [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
      Resource = "arn:aws:logs:*:*:*"
    },
    {
      Sid    = "SqsAnalysisQueue"
      Effect = "Allow"
      Action = [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:SendMessage"
      ]
      Resource = [
        aws_sqs_queue.analysis.arn,
        aws_sqs_queue.analysis_dlq.arn
      ]
    },
    {
      Sid    = "S3UploadsRead"
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:DeleteObject"
      ]
      Resource = "${var.uploads_bucket_arn}/*"
    }
  ],
  local.textract_statement,
  local.comprehend_statement,
  [
    {
      Sid    = "DynamoStore"
      Effect = "Allow"
      Action = [
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:UpdateItem",
        "dynamodb:Query"
      ]
      Resource = [
        var.incidents_table_arn,
        "${var.incidents_table_arn}/index/*",
        var.links_table_arn,
        "${var.links_table_arn}/index/*"
      ]
    },
    {
      Sid    = "SnsAlerts"
      Effect = "Allow"
      Action = [
        "sns:Publish"
      ]
      Resource = aws_sns_topic.caregiver_alerts.arn
    },
    {
      Sid    = "SsmConfig"
      Effect = "Allow"
      Action = [
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:GetParametersByPath"
      ]
      Resource = "arn:aws:ssm:*:*:parameter/${var.name_prefix}/*"
    }
  ])
}

resource "aws_iam_policy" "analysis" {
  name = "${var.name_prefix}-analysis-policy"

  policy = jsonencode({
    Version   = "2012-10-17"
    Statement = local.analysis_statements
  })
}

resource "aws_iam_role_policy_attachment" "analysis" {
  role       = aws_iam_role.analysis.name
  policy_arn = aws_iam_policy.analysis.arn
}

resource "aws_iam_role" "digest" {
  name = "${var.name_prefix}-digest-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_policy" "digest" {
  name = "${var.name_prefix}-digest-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "CloudWatchLogs"
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Sid    = "DynamoStore"
        Effect = "Allow"
        Action = [
          "dynamodb:Scan",
          "dynamodb:Query",
          "dynamodb:GetItem"
        ]
        Resource = [
          var.incidents_table_arn,
          "${var.incidents_table_arn}/index/*"
        ]
      },
      {
        Sid    = "SnsAlerts"
        Effect = "Allow"
        Action = [
          "sns:Publish"
        ]
        Resource = aws_sns_topic.caregiver_alerts.arn
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "digest" {
  role       = aws_iam_role.digest.name
  policy_arn = aws_iam_policy.digest.arn
}
