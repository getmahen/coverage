output "forwarder_lambda.arn" {
  value = "${aws_lambda_function.log_forwarder.arn}"
}