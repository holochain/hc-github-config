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

The automation user is provided with an access token that can be used in workflows.

To rotate the token, you can run the following command:

```bash
pulumi config set --secret automationUserToken <new-token>
```

This value is encrypted by Pulumi and stored in `Pulumi.github.yaml`.

Then you will need to ask Pulumi to deploy the token to projects that are 
configured to use it. You can do this by running:

```bash
pulumi up
```
