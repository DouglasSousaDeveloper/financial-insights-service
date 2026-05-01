provider "aws" {
  region = var.aws_region
}

module "network" {
  source               = "../../modules/network"
  environment          = var.environment
  vpc_cidr             = var.vpc_cidr
  public_subnet_count  = var.public_subnet_count
  private_subnet_count = var.private_subnet_count
  enable_nat_gateway   = var.enable_nat_gateway
  single_nat_gateway   = var.single_nat_gateway
}

module "ecr" {
  source       = "../../modules/ecr"
  environment  = var.environment
  project_name = var.project_name
}

module "alb" {
  source            = "../../modules/alb"
  environment       = var.environment
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids
  app_port          = var.app_port
}

module "rds" {
  source               = "../../modules/rds"
  environment          = var.environment
  vpc_id               = module.network.vpc_id
  private_subnet_ids   = module.network.private_subnet_ids
  vpc_cidr_block       = module.network.vpc_cidr_block
  db_instance_class    = var.db_instance_class
  db_allocated_storage = var.db_allocated_storage
  db_name              = var.db_name
  db_username          = var.db_username
  db_password          = var.db_password
}

module "ecs" {
  source                = "../../modules/ecs"
  environment           = var.environment
  vpc_id                = module.network.vpc_id
  private_subnet_ids    = module.network.private_subnet_ids
  aws_region            = var.aws_region
  alb_security_group_id = module.alb.alb_security_group_id
  alb_target_group_arn  = module.alb.target_group_arn
  alb_listener_arn      = module.alb.alb_listener_arn
  app_port              = var.app_port
  fargate_cpu           = var.fargate_cpu
  fargate_memory        = var.fargate_memory
  app_image             = "${module.ecr.repository_url}:latest"
  app_count             = var.app_count
  database_url          = "postgres://${var.db_username}:${var.db_password}@${module.rds.db_endpoint}/${var.db_name}?sslmode=disable"
  openai_api_key        = var.openai_api_key
}
