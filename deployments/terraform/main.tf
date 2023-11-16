terraform {
  required_version = ">=1.5.0"
  required_providers {
    google = {
      source = "hashicorp/google"
      version = ">=4.76.0"
    }
  }

  backend "gcs" {
    bucket = "ops-operation-center"
    prefix = "api_gateway"
  }
}

module "data_project" {
  source = "./data-project"
  project_id = var.data_project_id
  region = var.region
  zone = var.zone
}

module "devops_project" {
  source = "./devops-project"
  project_id = var.devops_project_id
  region = var.region
  zone = var.zone
}