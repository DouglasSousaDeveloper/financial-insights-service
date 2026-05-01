variable "environment" {
  type        = string
  description = "Nome do ambiente (dev, hom, prod)"
}

variable "vpc_cidr" {
  type        = string
  description = "Bloco CIDR para a VPC global"
  default     = "10.0.0.0/16"
}

variable "public_subnet_count" {
  type        = number
  description = "Quantidade de subnets públicas (ideal 2 para multi-az)"
  default     = 2
}

variable "private_subnet_count" {
  type        = number
  description = "Quantidade de subnets privadas (ideal 2 para multi-az)"
  default     = 2
}

variable "enable_nat_gateway" {
  type        = bool
  description = "Exige NAT Gateway para containers acessarem internet (OpenAi)"
  default     = true
}

variable "single_nat_gateway" {
  type        = bool
  description = "Mantém 1 NAT Gateway apenas (ajuda no custo em DEV/HOM)"
  default     = false
}
