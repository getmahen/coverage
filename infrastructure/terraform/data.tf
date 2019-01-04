data "aws_vpc" "microservice_vpc" {
  tags {
    Role = "microservice_vpc"
  }
}

data "aws_subnet_ids" "application" {
  vpc_id = "${data.aws_vpc.microservice_vpc.id}"

  tags {
    Tier = "Application"
  }
}

data "aws_security_group" "docker" {
  name = "${var.environment}-docker"
}

data "aws_s3_bucket_object" "coverage_pkg" {
  bucket = "credo-${var.environment}-lambdas"
  key    = "${var.lambda_name}.zip"
}
