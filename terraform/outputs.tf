output "ci_access_key_id" {
  description = "The Access Key ID for the CI/CD user (MUST be saved to GitHub Secrets)."
  value       = aws_iam_access_key.ci_key.id
}
