resource "aws_instance" "db_server" {
  ami             = data.aws_ami.amazon_linux_2.id
  instance_type   = var.instance_type              
  key_name        = aws_key_pair.app_key.key_name
  
  subnet_id       = aws_subnet.private.id 
  security_groups = [aws_security_group.db_sg.id]
  iam_instance_profile = aws_iam_instance_profile.ecr_profile.name

  user_data = <<-EOF
              #!/bin/bash
              sudo yum update -y
              sudo yum install docker -y
              sudo service docker start
              sudo usermod -aG docker ec2-user

              sudo docker run -d \
                -p 6379:6379 \
                --name redis-c \
                redis:latest
              
              sudo docker run -d \
                -p 27017:27017 \
                --name mongo-c \
                mongo:latest
              EOF

  tags = {
    Name = "DB-Server"
  }
}
