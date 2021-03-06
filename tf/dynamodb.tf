provider "aws" {
  profile             = "e6e-dev"
  region              = "us-west-2"
}

resource "aws_dynamodb_table" "questions" {
  name           = "questions-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "id"

  attribute {
     name = "id"
     type = "S"
  }

  attribute {
     name = "ownerid"
     type = "S"
  }

  global_secondary_index {
    name               = "ownerid-index"
    hash_key           = "ownerid"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }
}

resource "aws_dynamodb_table" "answers" {
  name           = "answers-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "id"
  range_key       = "ownerid"

  attribute {
     name = "id"
     type = "S"
  }
  attribute {
     name = "ownerid"
     type = "S"
  }

  global_secondary_index {
    name               = "id-index"
    hash_key           = "id"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

  global_secondary_index {
    name               = "ownerid-index"
    hash_key           = "ownerid"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}
