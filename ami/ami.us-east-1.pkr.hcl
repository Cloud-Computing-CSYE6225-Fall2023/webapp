packer {
  required_plugins {
    amazon = {
      source  = "github.com/hashicorp/amazon"
      version = ">1.0.0"
    }
  }
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "source_ami" {
  type    = string
  default = "ami-06db4d78cb1d3bbf9"
}

variable "ami_users" {
  type    = string
  default = ""
}

variable "ssh_username" {
  type    = string
  default = "admin"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

variable "device_name" {
  type    = string
  default = "/dev/xvda"
}

variable "volume_size" {
  type    = number
  default = 8
}

variable "volume_type" {
  type    = string
  default = "gp2"
}

source "amazon-ebs" "my-ami" {
  source_ami_filter {
    filters = {
      virtualization-type = "hvm"
      architecture        = "x86_64"
      root-device-type    = "ebs"
    }
    owners      = ["136693071363"]
    most_recent = true
  }

  ami_name        = "webapp_${formatdate("YYYY_MM_DD_hh_mm_ss", timestamp())}"
  ami_description = "AMI for webapp"
  region          = "${var.aws_region}"

  ami_regions = [
    "${var.aws_region}",
  ]

  aws_polling {
    delay_seconds = 120
    max_attempts  = 50
  }

  instance_type = "${var.instance_type}"
  ssh_username  = "${var.ssh_username}"
  ami_users     = ["${var.ami_users}"]

  launch_block_device_mappings {
    delete_on_termination = true
    device_name           = "${var.device_name}"
    volume_size           = var.volume_size
    volume_type           = "${var.volume_type}"
  }
}

build {
  name    = "Build Webapp"
  sources = ["source.amazon-ebs.my-ami"]

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1"
    ]
    inline = [
      "sudo apt-get update",
      "sudo apt-get upgrade -y",
      "sudo apt-get clean",
      "sudo apt install zip",
      "sudo apt install unzip",
      "mkdir -p github.com/shivasaicharanruthala",
    ]
  }

  provisioner "file" {
    destination = "/home/admin/github.com/shivasaicharanruthala"
    source      = "../../webapp/webapp.zip"
  }

  provisioner "shell" {
    inline = [
      "cd github.com/shivasaicharanruthala",
      "sudo unzip -q webapp.zip",
      "cd webapp",
      "ls",
      "sudo chmod +x ./startup-scripts/setup-go.sh",
      "sudo chmod +x ./startup-scripts/setup-postgres.sh",
      "sudo ./startup-scripts/setup-go.sh",
      "export PATH=$PATH:/usr/local/go/bin",
      "export GOPATH=/home/admin/github.com/shivasaicharanruthala",
      "go version",
      "sudo ./startup-scripts/setup-postgres.sh -u $PKR_VAR_DB_USER -p $PKR_VAR_DB_PASSWORD -d $PKR_VAR_DB_NAME",
      "go get -v ./..."
    ]
  }
}
