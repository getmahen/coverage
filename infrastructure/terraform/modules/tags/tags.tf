resource "null_resource" "tags" {
  triggers = {
    id          = "${lower(join(var.delimiter, compact(concat(list(var.environment, var.project, var.role)))))}"
    created     = "${timestamp()}"
    description = "${var.description}"
    environment = "${lower(format("%v", var.environment))}"
    project     = "${lower(format("%v", var.project))}"
    # Enforce 150 length
    role        = "${substr(lower(format("%v", var.role)),0,  length(var.role) > 150 ? 150 : -1    )}"
  }

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [
      "triggers.created"
    ]
  }
}
