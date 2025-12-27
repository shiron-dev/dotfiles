resource "google_project_service" "cloudkms_api" {
  project = "shiron-dev"
  service = "cloudkms.googleapis.com"

  disable_dependent_services = false
}

resource "google_kms_key_ring" "dotfiles_sops" {
  name     = "dotfiles-sops"
  location = "global"

  depends_on = [google_project_service.cloudkms_api]
}

resource "google_kms_crypto_key" "dotfiles_sops_key" {
  name     = "dotfiles-sops-key"
  key_ring = google_kms_key_ring.dotfiles_sops.id
  purpose  = "ENCRYPT_DECRYPT"

  rotation_period = "7776000s"
}

output "kms_keyring_name" {
  description = "Name of the KMS key ring for SOPS encryption"
  value       = google_kms_key_ring.dotfiles_sops.name
}

output "kms_key_name" {
  description = "Name of the KMS crypto key for SOPS encryption"
  value       = google_kms_crypto_key.dotfiles_sops_key.name
}

output "kms_key_id" {
  description = "ID of the KMS crypto key for SOPS encryption"
  value       = google_kms_crypto_key.dotfiles_sops_key.id
}
