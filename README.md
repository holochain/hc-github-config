# hc-github-config

[![Deploy](https://github.com/holochain/hc-github-config/actions/workflows/deploy.yaml/badge.svg)](https://github.com/holochain/hc-github-config/actions/workflows/deploy.yaml)

Get started with:

```bash
pulumi login
pulumi org set-default holochain
pulumi stack select github
```

You may also need to authenticate with github in order for Pulumi to be able to properly access repository settings. Using the [GitHub CLI](https://cli.github.com/), you can run:
```
gh auth login
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

### Importing a repository

Importing a repository is a little different to creating a new one. Pulumi requires that you describe the current state
of the repository in order to import it. Once you have done that, you can make changes.

So to get started, find an existing imported repository in `main.go` and copy it. We'll use `holochain-serialization`
as an example.

```go
description = "Abstractions to probably serialize and deserialize things properly without forgetting or doubling"
holochainSerializationRepositoryArgs := StandardRepositoryArgs("holochain-serialization", &description)
holochainSerialization, err := github.NewRepository(ctx, "holochain-serialization", &holochainSerializationRepositoryArgs, pulumi.Import(pulumi.ID("holochain-serialization")))
if err != nil {
    return err
}
```

Note the `pulumi.Import`, which is telling Pulumi that this repository already exists. Rename all the fields and
variables to describe the repository you are importing and adapt variable assignement where necessary for the code to compile.

Let's call the repository `example`:

```diff
- description = "Abstractions to probably serialize and deserialize things properly without forgetting or doubling"
+ description = "My repo description"
- holochainSerializationRepositoryArgs := StandardRepositoryArgs("holochain-serialization", &description)
+ exampleRepositoryArgs := StandardRepositoryArgs("example", &description)
- holochainSerialization, err := github.NewRepository(ctx, "holochain-serialization", &holochainSerializationRepositoryArgs, pulumi.Import(pulumi.ID("holochain-serialization")))
- if err != nil {
-     return err
- }
+ if _, err = github.NewRepository(ctx, "example", &exampleRepositoryArgs, pulumi.Import(pulumi.ID("example"))); err != nil {
+     return err
+ }

```

Now you can check which repository settings don't match. Try running `pulumi preview` and seeing what fields are reported as
changed. You may need to override the repository settings in Pulumi to match the existing settings for the repository. This might look something like:

```diff
description = "My repo description"
exampleRepositoryArgs := StandardRepositoryArgs("example", &description)
+ exampleRepositoryArgs.AllowRebaseMerge = pulumi.Bool(false)
+ exampleRepositoryArgs.SquashMergeCommitTitle = pulumi.String("PR_TITLE")
if _, err = github.NewRepository(ctx, "example", &exampleRepositoryArgs, pulumi.Import(pulumi.ID("example"))); err != nil {
    return err
}
```

Once there are no differences, you will be able to do the import by running `pulumi up`.

Next, you can configure the repository with the standard settings. Remove any overrides that you had to include for the
import that aren't intended to be kept. Then you need to:

- Either require or migrate the default branch to be `main`
- Set the access rules, which is the groups that are given roles against the repository
- Add rulesets which control how changes are made to the default branch and release branches

```diff
description = "My repo description"
exampleRepositoryArgs := StandardRepositoryArgs("example", &description)
- exampleRepositoryArgs.AllowRebaseMerge = pulumi.Bool(false)
- exampleRepositoryArgs.SquashMergeCommitTitle = pulumi.String("PR_TITLE")
- if _, err = github.NewRepository(ctx, "example", &exampleRepositoryArgs, pulumi.Import(pulumi.ID("example"))); err != nil {
-     return err
- }
+ example, err := github.NewRepository(ctx, "example", &exampleRepositoryArgs, pulumi.Import(pulumi.ID("example")))
+ if err != nil {
+     return err
+ }
+ if err = RequireMainAsDefaultBranch(ctx, "example", example); err != nil {
+     return err
+ }
+ if err = StandardRepositoryAccess(ctx, "example", example); err != nil {
+     return err
+ }
+ exampleDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(example, nil)
+ if _, err = github.NewRepositoryRuleset(ctx, "example-default", &exampleDefaultRepositoryRulesetArgs); err != nil {
+     return err
+ }
+ exampleReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(example, nil)
+ if _, err = github.NewRepositoryRuleset(ctx, "example-release", &exampleReleaseRepositoryRulesetArgs); err != nil {
+     return err
+ }
```

Finally, apply these changes with `pulumi up`.
