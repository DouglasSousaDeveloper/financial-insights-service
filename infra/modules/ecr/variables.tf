variable "environment" {
  type        = string
  description = "Ambiente (dev, hom, prod)"
}

variable "project_name" {
  type        = string
  description = "Nome do projeto"
  default     = "financial-insights"
}
