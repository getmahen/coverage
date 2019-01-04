resource "aws_iam_role" "log_forwarder_execution_role" {
  name = "log_forwarder_${var.forwarder_name}_execution_role"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": [
          "lambda.amazonaws.com",
          "logs.us-east-2.amazonaws.com"
        ]
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}

data "aws_s3_bucket_object" "log_forwarderzip" {
  bucket = "credo-${var.environment}-lambdas"
  key = "sumologic_forwarder.js.zip"
}

resource "aws_lambda_function" "log_forwarder" {
  function_name = "sumologic_forwarder_${var.forwarder_name}"
  s3_bucket         = "${data.aws_s3_bucket_object.log_forwarderzip.bucket}"
  s3_key            = "${data.aws_s3_bucket_object.log_forwarderzip.key}"
  s3_object_version = "${data.aws_s3_bucket_object.log_forwarderzip.version_id}"
  handler = "sumologic_forwarder.handler"
  runtime = "nodejs4.3"
  role = "${aws_iam_role.log_forwarder_execution_role.arn}"
  environment {
    variables {
      SOURCE_CATEGORY_OVERRIDE = "${var.source_category}"
      SUMO_ENDPOINT = "https://endpoint2.collection.us2.sumologic.com/receiver/v1/http/ZaVnC4dhaV2wI7oylSyJpbwJyLP0GN0kTd57IxzL9Fq5G1X-9L6nrHAsoahsLz9CRgSI_jm8xejmwEbbSQJLHTsTkiFMD8e7sRoef0hZ0gsGRgiYfixIrw=="
    }
  }
}

resource "aws_lambda_permission" "log_forwarder_cwl_permission" {
  statement_id = "log_forwarder_${var.forwarder_name}_AllowExecutionFromCloudwatchLogsSubscription"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.log_forwarder.arn}"
  principal = "logs.us-east-2.amazonaws.com"
}

resource "aws_cloudwatch_log_subscription_filter" "cloudwatch_log_forwarder" {
  name = "log_forwarder_${var.forwarder_name}"
  destination_arn = "${aws_lambda_function.log_forwarder.arn}"
  filter_pattern = ""
  log_group_name = "${var.log_group_name}"
}
