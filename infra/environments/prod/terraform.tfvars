aws_region           = "us-east-1"
environment          = "prod"
project_name         = "financial-insights"
vpc_cidr             = "10.2.0.0/16"
public_subnet_count  = 3
private_subnet_count = 3
enable_nat_gateway   = true
single_nat_gateway   = false # Cada subnet em prod terá seu NAT próprio de AZ para High Availability tolerante e falhas completas da AWS.
app_port             = 8080
db_instance_class    = "db.t4g.medium"
db_allocated_storage = 50
db_name              = "insights_prod"
fargate_cpu          = "1024"
fargate_memory       = "2048"
app_count            = 3

# db_password      = "..."
# openai_api_key   = "..."
