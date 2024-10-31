package main

import (
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		self, err := github.NewRepository(ctx, "literate-umbrella", &github.RepositoryArgs{
			Name:                pulumi.String("literate-umbrella"),
			Description:         nil,
			Visibility:          pulumi.String("public"),
			HasDownloads:        pulumi.Bool(false),
			HasIssues:           pulumi.Bool(true),
			HasProjects:         pulumi.Bool(true),
			HasWiki:             pulumi.Bool(false),
			VulnerabilityAlerts: pulumi.Bool(true),
			AllowAutoMerge:      pulumi.Bool(true),
			DeleteBranchOnMerge: pulumi.Bool(true),
		}, pulumi.Import(pulumi.ID("literate-umbrella")))
		if err != nil {
			return err
		}

		if _, err := github.NewRepositoryRuleset(ctx, "literate-umbrella", &github.RepositoryRulesetArgs{
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
				Update:                pulumi.Bool(true),
				Deletion:              pulumi.Bool(true),
				RequiredLinearHistory: pulumi.Bool(true),
				RequiredSignatures:    pulumi.Bool(false),
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

		return nil
	})
}
