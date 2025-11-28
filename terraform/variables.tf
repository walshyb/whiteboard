variable "region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "instance_type" {
  description = "EC2 instance"
  type        = string
  default     = "t4g.micro" 
}

variable "my_ssh_ip_cidr" {
  description = "SSH access IP"
  type        = string
}

variable "aws_account_id" {
  description = "AWS Account ID for ECR authentication."
  type        = string
}
