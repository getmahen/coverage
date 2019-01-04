resource "vault_policy" "coverage" {
  name = "${var.lambda_name}"

  policy = <<EOF
path "sys/*" {
  policy = "deny"
}

path "secret/services/coverage/private_key" {
  policy = "read"
}

path "auth/token/lookup-self" {
  policy = "read"
}
EOF
}

resource "vault_aws_auth_backend_role" "coverage" {
  role                    = "${var.lambda_name}"
  policies                = ["${vault_policy.coverage.name}"]
  auth_type               = "iam"
  max_ttl                 = 5
  ttl                     = 5
  bound_iam_principal_arn = "${aws_iam_role.role.arn}"
}
