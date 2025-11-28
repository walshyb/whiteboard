# Define a Trust Policy
resource "aws_iam_role" "ec2_ecr_role" {
  name = "EC2-ECR-Pull-Role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}

# grant permission to pull images from ECR.
resource "aws_iam_role_policy_attachment" "ecr_readonly" {
  role       = aws_iam_role.ec2_ecr_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

# Create the Instance Profile, used to attach the role to the EC2 instance
resource "aws_iam_instance_profile" "ecr_profile" {
  name = "ECR-Instance-Profile"
  role = aws_iam_role.ec2_ecr_role.name
}

# IAM Group for CI/CD users
resource "aws_iam_group" "ci_group" {
  name = "ci-cd-docker-push-group"
}

# IAM User for GitHub Actions
resource "aws_iam_user" "ci_user" {
  name = "github-actions-user"
}

resource "aws_iam_group_membership" "ci_membership" {
  name  = "ci-cd-user-membership"
  group = aws_iam_group.ci_group.name
  users = [
    aws_iam_user.ci_user.name,
  ]
}

resource "aws_iam_access_key" "ci_key" {
  user = aws_iam_user.ci_user.name
}

data "aws_caller_identity" "current" {}

resource "aws_iam_policy" "ecr_push_policy" {
  name        = "ECR-Push-Whiteboard-Policy"
  description = "Allows CI/CD user to push to whiteboard ECR repo"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "AllowECRPullPush"
        Effect   = "Allow"
        Action = [
          "ecr:CompleteLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:InitiateLayerUpload",
          "ecr:BatchCheckLayerAvailability",
          "ecr:PutImage",
          "ecr:GetAuthorizationToken",
        ]
        # Restrict this policy to ONLY the whiteboard repository ARN
        Resource = [
          "${aws_ecr_repository.app_repo.arn}"
        ]
      },
      # The GetAuthorizationToken action cannot be scoped to a resource, so we allow it on all resources
      {
        Sid      = "AllowGlobalAuth"
        Effect   = "Allow"
        Action   = "ecr:GetAuthorizationToken"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_group_policy_attachment" "ci_ecr_attach" {
  group      = aws_iam_group.ci_group.name
  policy_arn = aws_iam_policy.ecr_push_policy.arn
}
