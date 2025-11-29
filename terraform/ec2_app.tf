resource "aws_key_pair" "app_key" {
  key_name   = "whiteboard-app-key"
  public_key = file(var.public_key_path)
}

resource "aws_instance" "app_server" {
  ami             = data.aws_ami.amazon_linux_2.id
  instance_type   = var.instance_type             
  key_name        = aws_key_pair.app_key.key_name
  
  subnet_id       = aws_subnet.public_az_a.id 
  vpc_security_group_ids = [aws_security_group.app_sg.id]
  iam_instance_profile       = aws_iam_instance_profile.ecr_profile.name

  metadata_options {
    http_endpoint = "enabled"
    http_tokens = "required"
    http_put_response_hop_limit = 1
  }

  user_data = <<-EOF
              #!/bin/bash

              sudo yum update -y
              sudo yum install docker -y
              sudo service docker start
              sudo usermod -aG docker ec2-user

              aws ecr get-login-password --region ${var.region} \
                 | sudo docker login --username AWS --password-stdin ${var.aws_account_id}.dkr.ecr.${var.region}.amazonaws.com
              
              sudo docker run -d \
                -p 8080:8080 \
                -e ENV=production \
                -e REDIS_ADDR="${aws_instance.db_server.private_ip}:6379" \
                -e MONGO_URI="mongodb://${aws_instance.db_server.private_ip}:27017" \
                --name server \
                ${var.aws_account_id}.dkr.ecr.${var.region}.amazonaws.com/whiteboard:latest
              EOF

  tags = {
    Name = "Whiteboard-Server"
  }
}

output "app_server_instance_id" {
  description = "The Instance ID of the Application Server (Required for SSM)."
  value       = aws_instance.app_server.id
}
