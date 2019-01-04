provider "aws" {
  version = "~> 1.37"
  region  = "us-east-2"
  profile = "credo-auth"
  assume_role {
    role_arn = "arn:aws:iam::${lookup(var.env_to_acct_id, var.environment)}:role/${var.role}"
  }
}

variable "env_to_acct_id" {
  type  = "map"
  default = {
    dev     = "674346455231"
    qa      = "772404289823"
    prod    = "465292320167"
  }
}

provider "vault" {
  address = "https://vault.credo${var.environment}.dev:443"
  version = "~> 1.1"
}

provider "consul" {
  scheme = "https"
  address = "consul.credo${var.environment}.dev:443"
  version = "~> 2.1"
}

provider "null" {
  version = "~> 1.0"
}