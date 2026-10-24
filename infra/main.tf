locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

module "storage" {
  source      = "./modules/storage"
  name_prefix = local.name_prefix
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
  source             = "./modules/notify"
  name_prefix        = local.name_prefix
  budget_alert_email = var.budget_alert_email
}

module "api" {
  source                     = "./modules/api"
  name_prefix                = local.name_prefix
  frontend_bucket_id         = module.storage.frontend_bucket_id
  frontend_bucket_arn        = module.storage.frontend_bucket_arn
  cognito_user_pool_endpoint = module.auth.user_pool_endpoint
  cognito_client_id          = module.auth.user_pool_client_id
}
