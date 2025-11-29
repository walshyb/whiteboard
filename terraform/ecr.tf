resource "aws_ecr_repository" "app_repo" {
  name                 = "whiteboard"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "Whiteboard-App-Repo"
  }
}

output "ecr_repository_url" {
  description = "The URI of the private ECR repository for use in CI/CD pipelines."
  value       = aws_ecr_repository.app_repo.repository_url
}
