variable "aws_region" { type = string, default = "us-east-1" }
variable "environment" { type = string, default = "dev" }
variable "project_name" { type = string, default = "financial-insights" }
variable "vpc_cidr" { type = string, default = "10.0.0.0/16" }
variable "public_subnet_count" { type = number, default = 2 }
variable "private_subnet_count" { type = number, default = 2 }
variable "enable_nat_gateway" { type = bool, default = true }
variable "single_nat_gateway" { type = bool, default = true }
variable "app_port" { type = number, default = 8080 }
variable "db_instance_class" { type = string, default = "db.t4g.micro" }
variable "db_allocated_storage" { type = number, default = 20 }
variable "db_name" { type = string, default = "insights_db" }
variable "db_username" { type = string, default = "postgres" }
variable "db_password" { type = string }
variable "fargate_cpu" { type = string, default = "256" }
variable "fargate_memory" { type = string, default = "512" }
variable "app_count" { type = number, default = 1 }
variable "openai_api_key" { type = string }
