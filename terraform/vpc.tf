# Create the VPC (private network container)
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16" 
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "Whiteboard-VPC"
  }
}

# Internet Gateway allows traffic in and out of the VPC
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "Whiteboard-IGW"
  }
}

# Need 2 subnets in different availability zones (AZ's) for ALB
# Public Subnet A 
resource "aws_subnet" "public_az_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24" 
  availability_zone       = data.aws_availability_zones.available.names[0]
  map_public_ip_on_launch = true 

  tags = {
    Name = "Public-Subnet-A"
  }
}

# Public Subnet B
resource "aws_subnet" "public_az_b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.3.0/24"
  availability_zone       = data.aws_availability_zones.available.names[1]
  map_public_ip_on_launch = true 

  tags = {
    Name = "Public-Subnet-B"
  }
}

# Public Route Table (routes traffic to the Internet Gateway)
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }

  tags = {
    Name = "Public-Route-Table"
  }
}

# Associate the Public Subnet with the Public Route Table
resource "aws_route_table_association" "public_a" {
  subnet_id      = aws_subnet.public_az_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_b" {
  subnet_id      = aws_subnet.public_az_b.id
  route_table_id = aws_route_table.public.id
}


# Security Group for the App Server (Box 1)
resource "aws_security_group" "app_sg" {
  vpc_id = aws_vpc.main.id
  name   = "app-server-sg"

  # Allow SSH for administration
  ingress {
    description = "Allow SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.my_ssh_ip_cidr]
  }

  # Allow HTTP/WS ONLY from the Load Balancer's Security Group
  ingress {
    description     = "Allow HTTP/WS from ALB"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.lb_sg.id]
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Security Group for the DB Server (Box 2)
resource "aws_security_group" "db_sg" {
  vpc_id = aws_vpc.main.id
  name   = "db-server-sg"

  # Allow Redis and Mongo only from the App Server's subnet
  ingress {
    description = "Allow DB Access from App Subnet"
    from_port   = 0
    to_port     = 65535 # All ports for simplicity
    protocol    = "tcp"
    cidr_blocks = [aws_subnet.public_az_a.cidr_block] # Only allows traffic from 10.0.1.0/24
  }

  # Rule: Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Security Group for the Application Load Balancer (ALB)
resource "aws_security_group" "lb_sg" {
  name        = "whiteboard-alb-sg"
  description = "Allows HTTP/80 inbound from the Internet to the ALB"
  vpc_id      = aws_vpc.main.id

  # Inbound rule: Allow HTTP/80 from anywhere
  ingress {
    description = "Allow HTTP from internet"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Outbound rule: Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
