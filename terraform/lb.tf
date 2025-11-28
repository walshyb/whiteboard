# Creates the application load balancer (ALB)
resource "aws_lb" "app_lb" {
  name               = "app-whiteboard-lb"
  internal           = false
  load_balancer_type = "application"
  drop_invalid_header_fields = true
  security_groups    = [aws_security_group.lb_sg.id]
  subnets = [
    aws_subnet.public_az_a.id,
    aws_subnet.public_az_b.id,
  ]

  tags = {
    Name = "WhiteboardAppALB"
  }
}

# Creates target group (destination for traffic)
resource "aws_lb_target_group" "app_tg" {
  name     = "whiteboard-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  health_check {
    path                = "/health"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
}

# Create a Listener (tell the ALB where to listen for traffic)
resource "aws_lb_listener" "http_listener" {
  load_balancer_arn = aws_lb.app_lb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app_tg.arn
  }
}

# Attach the App Server to the Target Group
resource "aws_lb_target_group_attachment" "app_server_attach" {
  target_group_arn = aws_lb_target_group.app_tg.arn
  target_id        = aws_instance.app_server.id
  port             = 8080
}
