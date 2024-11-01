package main

import (
	"fmt"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		description := "Automation for GitHub repository configurations for the Holochain organization."
		selfRepositoryArgs := StandardRepositoryArgs("hc-github-config", &description)
		selfRepositoryArgs.AllowMergeCommit = pulumi.Bool(false)
		selfRepositoryArgs.AllowSquashMerge = pulumi.Bool(false)
		selfRepositoryArgs.AllowRebaseMerge = pulumi.Bool(true)
		self, err := github.NewRepository(ctx, "hc-github-config", &selfRepositoryArgs, pulumi.Import(pulumi.ID("hc-github-config")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "hc-github-config", self); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "hc-github-config", self); err != nil {
			return err
		}

		if _, err := github.NewRepositoryRuleset(ctx, "hc-github-config", &github.RepositoryRulesetArgs{
			Name:        pulumi.String("default"),
			Repository:  self.Name,
			Target:      pulumi.String("branch"),
			Enforcement: pulumi.String("active"),
			Conditions: &github.RepositoryRulesetConditionsArgs{
				RefName: &github.RepositoryRulesetConditionsRefNameArgs{
					Includes: pulumi.StringArray{
						pulumi.String("~DEFAULT_BRANCH"),
					},
					Excludes: pulumi.StringArray{},
				},
			},
			Rules: &github.RepositoryRulesetRulesArgs{
				Creation:              pulumi.Bool(true),
				Update:                pulumi.Bool(false),
				Deletion:              pulumi.Bool(true),
				RequiredLinearHistory: pulumi.Bool(true),
				RequiredSignatures:    pulumi.Bool(false),
				RequiredStatusChecks: &github.RepositoryRulesetRulesRequiredStatusChecksArgs{
					RequiredChecks: github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArray{
						github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
							Context: pulumi.String("check"),
						},
					},
					StrictRequiredStatusChecksPolicy: pulumi.Bool(true),
				},
			},
			BypassActors: github.RepositoryRulesetBypassActorArray{
				&github.RepositoryRulesetBypassActorArgs{
					ActorId:    pulumi.Int(5), // Repository admin
					ActorType:  pulumi.String("RepositoryRole"),
					BypassMode: pulumi.String("always"),
				},
			},
		}); err != nil {
			return err
		}

		// holochain-wasmer
		holochainWasmerRepositoryArgs := StandardRepositoryArgs("holochain-wasmer", nil)
		holochainWasmer, err := github.NewRepository(ctx, "holochain-wasmer", &holochainWasmerRepositoryArgs, pulumi.Import(pulumi.ID("holochain-wasmer")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "holochain-wasmer", holochainWasmer); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "holochain-wasmer", holochainWasmer); err != nil {
			return err
		}

		return nil
	})
}

func StandardRepositoryArgs(name string, description *string) github.RepositoryArgs {
	args := github.RepositoryArgs{
		Name:                pulumi.String(name),
		Description:         nil,
		Visibility:          pulumi.String("public"),
		HasDownloads:        pulumi.Bool(false),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(false),
		VulnerabilityAlerts: pulumi.Bool(true),
		AllowAutoMerge:      pulumi.Bool(true),
		DeleteBranchOnMerge: pulumi.Bool(true),
	}

	if description != nil {
		args.Description = pulumi.String(*description)
	}

	return args
}

func StandardRepositoryAccess(ctx *pulumi.Context, name string, repository *github.Repository) error {
	_, err := github.NewTeamRepository(ctx, fmt.Sprintf("%s-collaborator-core-dev", name), &github.TeamRepositoryArgs{
		Repository: repository.Name,
		Permission: pulumi.String("admin"),
		TeamId:     pulumi.String("core-dev"),
	})
	if err != nil {
		return err
	}
	_, err = github.NewTeamRepository(ctx, fmt.Sprintf("%s-collaborator-holochain-devs", name), &github.TeamRepositoryArgs{
		Repository: repository.Name,
		Permission: pulumi.String("maintain"),
		TeamId:     pulumi.String("holochain-devs"),
	})
	if err != nil {
		return err
	}

	return nil
}

func RequireMainAsDefaultBranch(ctx *pulumi.Context, name string, repository *github.Repository) error {
	_, err := github.NewBranchDefault(ctx, fmt.Sprintf("%s-default-branch", name), &github.BranchDefaultArgs{
		Repository: repository.Name,
		Branch:     pulumi.String("main"),
		Rename:     pulumi.Bool(false),
	})
	return err
}
