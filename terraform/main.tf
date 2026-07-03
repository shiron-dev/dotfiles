terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.39.0"
    }
  }

  required_version = ">= 1.15.7"

  backend "gcs" {
    bucket = "shiron-dev-dotfiles-terraform"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = "shiron-dev"
  region  = "asia-northeast1"
}
