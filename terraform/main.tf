terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.46.0"
    }
  }

  required_version = ">= 1.16.1"

  backend "gcs" {
    bucket = "shiron-dev-dotfiles-terraform"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = "shiron-dev"
  region  = "asia-northeast1"
}
