aws_region           = "us-east-1"
environment          = "hom"
project_name         = "financial-insights"
vpc_cidr             = "10.1.0.0/16"
public_subnet_count  = 2
private_subnet_count = 2
enable_nat_gateway   = true
single_nat_gateway   = true 
app_port             = 8080
db_instance_class    = "db.t4g.micro"
db_allocated_storage = 20
db_name              = "insights_hom"
fargate_cpu          = "512"
fargate_memory       = "1024"
app_count            = 2

# db_password      = "..."
# openai_api_key   = "..."
