terraform {
  required_providers {
    ctfg = {
      source = "ctfg/ctfg"
    }
  }
}

provider "ctfg" {
}

variable "test" {
  type    = string
  default = "1"
}

resource "ctfg_parameter" "test" {
  id    = "test"
  value = var.test
  type  = "number"
}

output "test" {
  value = ctfg_parameter.test.value
}