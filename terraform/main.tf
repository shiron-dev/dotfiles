terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.30.0"
    }
  }

  required_version = ">= 1.14.9"

  backend "gcs" {
    bucket = "shiron-dev-dotfiles-terraform"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = "shiron-dev"
  region  = "asia-northeast1"
}
