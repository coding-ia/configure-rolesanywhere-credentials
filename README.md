# Configure AWS Roles Anywhere for GitHub Actions

This action allows you to use Roles Anywhere for getting credentials for an AWS account.

Before using this action, you will need to [setup](https://docs.aws.amazon.com/rolesanywhere/latest/userguide/getting-started.html) IAM Roles Anywhere, and setup the proper certificates and roles.

# Using the action

This is the minimal setup required to use the action:

```
- name: AWS Roles Anywhere action
  uses: coding-ia/configure-rolesanywhere-credentials@main
  with:
    aws-region: us-east-1
    private-key: ${{ secrets.PRIVATE_KEY }}
    certificate: ${{ secrets.CERTIFICATE }}
    role-arn: arn:aws:iam::111122223333:role/TestRole
    profile-arn: arn:aws:rolesanywhere:us-east-1:111122223333:profile/eef5f646-6218-4871-bf41-db7e13f7acad
    trust-anchor-arn: arn:aws:rolesanywhere:us-east-1:111122223333:trust-anchor/6dd648ad-92a7-4a44-a59c-2b1d56dfbbee
```

The PEM formatted private key and certificate are stored in GitHub secrets.