output "id" {
  value       = "${null_resource.tags.triggers.id}"
  description = "Disambiguated ID for resource which is used for the name tag"
}

output "environment" {
  value       = "${null_resource.tags.triggers.environment}"
  description = "Environment name"
}

output "project" {
  value       = "${null_resource.tags.triggers.project}"
  description = "The name of the project a resource is associated with"
}

output "description" {
  value       = "${null_resource.tags.triggers.description}"
  description = "This is human readable, simple description for the resource"
}

# Merge input tags with required tags
output "tags" {
  value = "${
      merge(
        map(
          "Name", "${null_resource.tags.triggers.id}",
          "Environment", "${null_resource.tags.triggers.environment}",
          "Role", "${null_resource.tags.triggers.role}",
          "Description", "${null_resource.tags.triggers.description}",
          "Created", "${null_resource.tags.triggers.created}",
          "Terraform", "true"
        ), var.tags
      )
    }"

  description = "Map of tags which can be applied to a resource"
}

# Use the same tag map as above, but autoscalers expect a list instead of a map
output "asg_tags" {
  value = [
    "${
        list(
          map("key", "Name", "value", null_resource.tags.triggers.id, "propagate_at_launch", true),
          map("key", "Environment", "value", null_resource.tags.triggers.environment, "propagate_at_launch", true),
          map("key", "Project", "value", null_resource.tags.triggers.project, "propagate_at_launch", true),
          map("key", "Role", "value", null_resource.tags.triggers.role, "propagate_at_launch", true),
          map("key", "Terraform", "value", "true", "propagate_at_launch", true)
        )
      }"]

  description = "Map of tags which instructs autoscaling to propogate tags at launch"
}
