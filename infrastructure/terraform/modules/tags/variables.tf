variable "environment" {
  description = "Environment name, e.g. 'dev', 'qa', 'prod'"
}
variable "project" {
  description = "The name of the project a resource is associated with"
}

variable "role" {
  description = "This is a machine parsable short name for the resource (spaces are replaced with the delimeter)."
}

variable "description" {
  description = "This is human readable, simple description for the resource"
}

variable "delimiter" {
  type        = "string"
  default     = "_"
  description = "(Default: _) A delimiter to be used between strings"
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "A map of additional tags add to this object (e.g. `map('TagKey`,`TagValue`)"
}
