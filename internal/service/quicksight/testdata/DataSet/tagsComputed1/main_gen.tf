# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "null" {}

resource "aws_quicksight_data_set" "test" {
  data_set_id = var.rName
  name        = var.rName
  import_mode = "SPICE"

  physical_table_map {
    physical_table_map_id = var.rName
    s3_source {
      data_source_arn = aws_quicksight_data_source.test.arn
      input_columns {
        name = "Column1"
        type = "STRING"
      }
      upload_settings {
        format = "JSON"
      }
    }
  }

  tags = {
    (var.unknownTagKey) = null_resource.test.id
  }
}

# testAccDataSetConfig_base

resource "aws_quicksight_data_source" "test" {
  data_source_id = var.rName
  name           = var.rName

  parameters {
    s3 {
      manifest_file_location {
        bucket = aws_s3_bucket.test.bucket
        key    = aws_s3_object.test.key
      }
    }
  }

  type = "S3"
}

# testAccDataSourceConfig_base

data "aws_partition" "current" {}

resource "aws_s3_bucket" "test" {
  bucket        = var.rName
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_ownership_controls" "test" {
  bucket = aws_s3_bucket.test.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "test" {
  depends_on = [
    aws_s3_bucket_public_access_block.test,
    aws_s3_bucket_ownership_controls.test,
  ]

  bucket = aws_s3_bucket.test.id
  acl    = "public-read"
}

resource "aws_s3_object" "test_data" {
  depends_on = [aws_s3_bucket_acl.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = "${var.rName}-test-data"
  content = <<EOF
[
	{
		"Column1": "aaa",
		"Column2": 1
	},
	{
		"Column1": "bbb",
		"Column2": 1
	}
]
EOF
  acl     = "public-read"
}

resource "aws_s3_object" "test" {
  depends_on = [aws_s3_bucket_acl.test]

  bucket  = aws_s3_bucket.test.bucket
  key     = var.rName
  content = <<EOF
{
  "fileLocations": [
      {
          "URIs": [
              "https://${aws_s3_bucket.test.bucket}.s3.${data.aws_partition.current.dns_suffix}/${var.rName}-test-data"
          ]
      }
  ],
  "globalUploadSettings": {
      "format": "JSON"
  }
}
EOF
  acl     = "public-read"
}

resource "null_resource" "test" {}

variable "rName" {
  description = "Name for resource"
  type        = string
  nullable    = false
}

variable "unknownTagKey" {
  type     = string
  nullable = false
}
