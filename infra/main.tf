locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

module "storage" {
  source         = "./modules/storage"
  name_prefix    = local.name_prefix
  allowed_origin = var.allowed_origin
}

module "data" {
  source      = "./modules/data"
  name_prefix = local.name_prefix
}

module "auth" {
  source      = "./modules/auth"
  name_prefix = local.name_prefix
}

module "notify" {
  source              = "./modules/notify"
  name_prefix         = local.name_prefix
  budget_alert_email  = var.budget_alert_email
  use_comprehend      = var.use_comprehend
  ocr_mode            = var.ocr_mode
  uploads_bucket      = module.storage.uploads_bucket_id
  uploads_bucket_arn  = module.storage.uploads_bucket_arn
  incidents_table     = module.data.incidents_table_name
  incidents_table_arn = module.data.incidents_table_arn
  links_table         = module.data.links_table_name
  links_table_arn     = module.data.links_table_arn
}

module "api" {
  source                     = "./modules/api"
  name_prefix                = local.name_prefix
  enable_cloudfront          = var.enable_cloudfront && var.frontend_host == "cloudfront"
  allowed_origin             = var.allowed_origin
  ocr_mode                   = var.ocr_mode
  frontend_bucket_id         = module.storage.frontend_bucket_id
  frontend_bucket_arn        = module.storage.frontend_bucket_arn
  cognito_user_pool_endpoint = module.auth.user_pool_endpoint
  cognito_client_id          = module.auth.user_pool_client_id
  uploads_bucket             = module.storage.uploads_bucket_id
  uploads_bucket_arn         = module.storage.uploads_bucket_arn
  incidents_table            = module.data.incidents_table_name
  incidents_table_arn        = module.data.incidents_table_arn
  links_table                = module.data.links_table_name
  links_table_arn            = module.data.links_table_arn
  analysis_queue_url         = module.notify.analysis_queue_url
  analysis_queue_arn         = module.notify.analysis_queue_arn
}
