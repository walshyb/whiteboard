# Import the manually created hosted zone
data "aws_route53_zone" "app" {
  name = "app.bwal.sh"
}

# Import the manually created certificate
data "aws_acm_certificate" "app_cert" {
  domain   = "app.bwal.sh"
  statuses = ["ISSUED"]
}

# Create the A record pointing to ALB
resource "aws_route53_record" "app" {
  zone_id = data.aws_route53_zone.app.zone_id
  name    = "app.bwal.sh"
  type    = "A"

  alias {
    name                   = aws_lb.app_lb.dns_name
    zone_id                = aws_lb.app_lb.zone_id
    evaluate_target_health = true
  }
}
