aws_region           = "us-east-1"
environment          = "dev"
project_name         = "financial-insights"
vpc_cidr             = "10.0.0.0/16"
public_subnet_count  = 2
private_subnet_count = 2
enable_nat_gateway   = true
single_nat_gateway   = true
app_port             = 8080
db_instance_class    = "db.t4g.micro"
db_allocated_storage = 20
db_name              = "insights_dev"
fargate_cpu          = "256"
fargate_memory       = "512"
app_count            = 1

# Secrets that must be injected dynamically via CLI or Environment Variables:
# db_password      = "..."
# openai_api_key   = "..."
