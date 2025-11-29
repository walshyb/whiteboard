resource "aws_instance" "db_server" {
  ami             = data.aws_ami.amazon_linux_2.id
  instance_type   = var.instance_type              
  key_name        = aws_key_pair.app_key.key_name
  
  subnet_id       = aws_subnet.public_az_a.id 
  vpc_security_group_ids = [aws_security_group.db_sg.id]
  iam_instance_profile = aws_iam_instance_profile.ecr_profile.name

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

              sudo docker run -d \
                --restart unless-stopped \
                -p 6379:6379 \
                --name redis-c \
                redis:6.2-alpine
              
              sudo docker run -d \
                --restart unless-stopped \
                -v /data/mongodb:/data/db \
                -p 27017:27017 \
                --name mongo-c \
                mongo:6.0
              EOF

  tags = {
    Name = "DB-Server"
  }
}
