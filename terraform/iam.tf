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
