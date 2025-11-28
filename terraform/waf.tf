# Define the rate limiting web ACL (the WAF)
resource "aws_wafv2_web_acl" "rate_limit_acl" {
  name        = "whiteboard-rate-limit-acl"
  scope       = "REGIONAL"
  description = "Rate limiting for DDoS protection"

  default_action {
    allow {}
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "whiteboard-rate-limit-metrics"
    sampled_requests_enabled   = true
  }

  rule {
    name     = "RateLimitRule"
    priority = 1

    action {
      block {} # Action: Block traffic when the limit is exceeded
    }

    statement {
      rate_limit_statement {
        # Limit: 500 requests per IP over a 5-minute window.
        # A starting point for protecting against rapid handshakes/DDoS.
        limit = 500 
        aggregate_key = "IP" # Track limit per originating IP address
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "RateLimitMetrics"
      sampled_requests_enabled   = true
    }
  }
  
  tags = {
    Name = "Whiteboard-WAF"
  }
}

# Associate the WAF with the ALB
resource "aws_wafv2_web_acl_association" "acl_assoc" {
  resource_arn = aws_lb.app_lb.arn
  web_acl_arn  = aws_wafv2_web_acl.rate_limit_acl.arn
}
