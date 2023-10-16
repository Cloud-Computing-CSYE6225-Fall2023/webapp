packer {
  required_plugins {
    amazon = {
      source = "github.com/hashicorp/amazon"
      version = ">1.0.0"
    }
  }
}

variable "aws_region" {
  type = string
  default = "us-east-1"
}

variable "source_ami" {
  type = string
  default ="ami-06db4d78cb1d3bbf9"
}

variable "ssh_username" {
  type = string
  default = "ubuntu"
}

variable "subnet_id" {
  type = string
  default = "subnet-022e4a40b96c87648"
}

source "amazon-ebs" "my-ami" {
  region = "${var.aws_region}"
  ami_name = "webapp-+${formatdate("YYYY_MM_DD_hh_mm_ss", timestamp())}"
  ami_description = "AMI for webapp",

  ami_regions = [
    "us-east-1",
  ]

  aws_polling {
    delay_seconds = 120
    max_attempts = 50
  }

  instance_type = "t2.micro"
  subnet_id = "${var.subnet_id}"
  source_ami = "${var.source_ami}"
  access_key = "${var.access_key_id}",
  secret_key = "${var.secret_access_key}",
  ssh_username = "${var.ssh_username}"
  ami_users = ["394598842451"]

  launch_block_device_mappings {
    delete_on_termination = true
    device_name = "/dev/xvda"
    volume_size = 8
    volume_type = "gp2"
  }
}

build {
  sources = ["source.amazon-ebs.my-ami"]

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1"
    ]
    inline = [
      "sudo apt-get update",
      "sudo apt-get upgrade -y",
      "sudo apt-get install nginx -y",
      "sudo apt-get clean"
    ]
  }
}
