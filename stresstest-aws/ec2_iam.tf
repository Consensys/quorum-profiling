resource "aws_iam_role_policy" "ec2_cloudwatch" {
  name_prefix = local.network_name
  policy      = data.aws_iam_policy.CloudWatchAgentServerPolicy.policy
  role        = aws_iam_role.ec2.id
}

data "aws_iam_policy_document" "ec2" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "aws_iam_policy" "CloudWatchAgentServerPolicy" {
  arn = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
}

resource "aws_iam_role" "ec2" {
  name_prefix = local.network_name

  assume_role_policy = data.aws_iam_policy_document.ec2.json

  tags = {
    By   = "quorum"
    Name = local.network_name
  }
}

resource "aws_iam_instance_profile" "node" {
  name_prefix = local.network_name
  role        = aws_iam_role.ec2.name
}