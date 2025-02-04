# hc-github-config

[![Deploy](https://github.com/holochain/hc-github-config/actions/workflows/deploy.yaml/badge.svg)](https://github.com/holochain/hc-github-config/actions/workflows/deploy.yaml)

Get started with:

```bash
pulumi login
pulumi org set-default holochain
pulumi stack select github
```

Then deploy using:

```bash
pulumi up
```

### Rotating the GitHub access token

The automation user is provided with an access token that can be used in standard workflows.

To rotate the token, you can run the following command:

```bash
pulumi config set --secret hra2GithubUserToken <new-token>
```

This value is encrypted by Pulumi and stored in `Pulumi.github.yaml`.

Then you will need to ask Pulumi to deploy the token to projects that are 
configured to use it. You can do this by running:

```bash
pulumi up
```

### Rotating the GitHub admin access token

There is also a GitHub access token with admin scope to modify other
repositories. This token is likely only meant to be used by this repository
itself.

To rotate the token, you can run the following command:

```bash
pulumi config set --secret hra2GithubAdminToken <new-token>
```

### Rotating the Pulumi access token

The secret `hra2PulumiAccessToken` can be used to give repositories access to
Pulumi itself so that changes can be deployed in the CI.

To rotate the token, you can run the following command:

```bash
pulumi config set --secret hra2PulumiAccessToken <new-token>
```

This value is encrypted by Pulumi and stored in `Pulumi.github.yaml`.

Then you will need to ask Pulumi to deploy the token to projects that are
configured to use it. You can do this by getting a PR merged into `main` and
allowing the CI to deploy it or by manually running:

```bash
pulumi up
```
