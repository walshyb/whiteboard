# dynamically looks up the latest Amazon Linux 2 AMI for Graviton
data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-arm64-gp2"] 
  }

  filter {
    name   = "architecture"
    values = ["arm64"] # required for t4g.micro
  }
}
