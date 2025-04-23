terraform {
  required_providers {
    multipass = {
      source  = "larstobi/multipass"
      version = "~> 1.4.2"
    }
  }
}

resource "multipass_instance" "lilo-finance-manager" {
  name    = "lilo-finance-manager"
  cpus    = 6
  disk    = "15GiB"
  memory  = "7GiB"
  image   = "jammy"
}
