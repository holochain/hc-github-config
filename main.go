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

		//
		// holochain-wasmer
		//
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

		//
		// wind tunnel
		//
		description = "Performance testing for Holochain"
		windTunnelRepositoryArgs := StandardRepositoryArgs("wind-tunnel", &description)
		windTunnelRepositoryArgs.AllowAutoMerge = pulumi.Bool(false)
		windTunnel, err := github.NewRepository(ctx, "wind-tunnel", &windTunnelRepositoryArgs, pulumi.Import(pulumi.ID("wind-tunnel")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "wind-tunnel", windTunnel); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "wind-tunnel", windTunnel); err != nil {
			return err
		}

		//
		// Holochain JS client
		//
		description = "A JavaScript client for the Holochain Conductor API"
		jsClientRepositoryArgs := StandardRepositoryArgs("holochain-client-js", &description)
		jsClient, err := github.NewRepository(ctx, "holochain-client-js", &jsClientRepositoryArgs, pulumi.Import(pulumi.ID("holochain-client-js")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "holochain-client-js", jsClient); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "holochain-client-js", jsClient); err != nil {
			return err
		}

		//
		// Holochain Rust client
		//
		description = "A Rust client for the Holochain Conductor API"
		rustClientRepositoryArgs := StandardRepositoryArgs("holochain-client-rust", &description)
		rustClient, err := github.NewRepository(ctx, "holochain-client-rust", &rustClientRepositoryArgs, pulumi.Import(pulumi.ID("holochain-client-rust")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "holochain-client-rust", rustClient); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "holochain-client-rust", rustClient); err != nil {
			return err
		}

		//
		// Tryorama
		//
		description = "Toolset to manage Holochain conductors and facilitate test scenarios"
		tryoramaRepositoryArgs := StandardRepositoryArgs("tryorama", &description)
		tryorama, err := github.NewRepository(ctx, "tryorama", &tryoramaRepositoryArgs, pulumi.Import(pulumi.ID("tryorama")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "tryorama", tryorama); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "tryorama", tryorama); err != nil {
			return err
		}

		//
		// Holonix
		//
		description = "Holochain app development environment based on Nix."
		holonixRepositoryArgs := StandardRepositoryArgs("holonix", &description)
		holonix, err := github.NewRepository(ctx, "holonix", &holonixRepositoryArgs, pulumi.Import(pulumi.ID("holonix")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "holonix", holonix); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "holonix", holonix); err != nil {
			return err
		}

		//
		// Binaries
		//
		description = "Holochain binaries for supported platforms"
		binariesRepositoryArgs := StandardRepositoryArgs("binaries", &description)
		binaries, err := github.NewRepository(ctx, "binaries", &binariesRepositoryArgs, pulumi.Import(pulumi.ID("binaries")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "binaries", binaries); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "binaries", binaries); err != nil {
			return err
		}

		//
		// Signal bends decently
		//
		description = "Simple websocket-based message relay servers and clients"
		sbdRepositoryArgs := StandardRepositoryArgs("sbd", &description)
		sbd, err := github.NewRepository(ctx, "sbd", &sbdRepositoryArgs, pulumi.Import(pulumi.ID("sbd")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "sbd", sbd); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "sbd", sbd); err != nil {
			return err
		}

		//
		// Tx5
		//
		description = "Holochain WebRTC P2P Communication Ecosystem"
		tx5RepositoryArgs := StandardRepositoryArgs("tx5", &description)
		tx5RepositoryArgs.SquashMergeCommitTitle = pulumi.String("PR_TITLE")
		tx5, err := github.NewRepository(ctx, "tx5", &tx5RepositoryArgs, pulumi.Import(pulumi.ID("tx5")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "tx5", tx5); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "tx5", tx5); err != nil {
			return err
		}

		//
		// Lair Keystore
		//
		description = "secret lair private keystore"
		lairRepositoryArgs := StandardRepositoryArgs("lair", &description)
		lairRepositoryArgs.AllowRebaseMerge = pulumi.Bool(false)
		lairRepositoryArgs.SquashMergeCommitTitle = pulumi.String("PR_TITLE")
		lair, err := github.NewRepository(ctx, "lair", &lairRepositoryArgs, pulumi.Import(pulumi.ID("lair")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "lair", lair); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "lair", lair); err != nil {
			return err
		}

		//
		// Holochain CHC Service
		//
		description = "A local web server that implements the CHC (Chain Head Coordinator) interface in Rust"
		hcChcServiceRepositoryArgs := StandardRepositoryArgs("hc-chc-service", &description)
		hcChcService, err := github.NewRepository(ctx, "hc-chc-service", &hcChcServiceRepositoryArgs, pulumi.Import(pulumi.ID("hc-chc-service")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "hc-chc-service", hcChcService); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "hc-chc-service", hcChcService); err != nil {
			return err
		}

		//
		// Holochain Serialization
		//
		description = "Abstractions to probably serialize and deserialize things properly without forgetting or doubling"
		holochainSerializationRepositoryArgs := StandardRepositoryArgs("holochain-serialization", &description)
		holochainSerialization, err := github.NewRepository(ctx, "holochain-serialization", &holochainSerializationRepositoryArgs, pulumi.Import(pulumi.ID("holochain-serialization")))
		if err != nil {
			return err
		}
		if err = MigrateDefaultBranchToMain(ctx, "holochain-serialization", holochainSerialization); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "holochain-serialization", holochainSerialization); err != nil {
			return err
		}

		//
		// Influxive
		//
		description = "Opinionated tools for working with InfluxDB from Rust"
		influxiveRepositoryArgs := StandardRepositoryArgs("influxive", &description)
		influxive, err := github.NewRepository(ctx, "influxive", &influxiveRepositoryArgs, pulumi.Import(pulumi.ID("influxive")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "influxive", influxive); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "influxive", influxive); err != nil {
			return err
		}

		//
		// Holochain Python Client
		//
		description = "A Python client for the Holochain Conductor API "
		pythonClientRepositoryArgs := StandardRepositoryArgs("holochain-client-python", &description)
		pythonClientRepositoryArgs.Topics = pulumi.StringArray{
			pulumi.String("python"),
			pulumi.String("python3"),
			pulumi.String("holochain"),
			pulumi.String("conductor-api"),
		}
		pythonClient, err := github.NewRepository(ctx, "holochain-client-python", &pythonClientRepositoryArgs, pulumi.Import(pulumi.ID("holochain-client-python")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "holochain-client-python", pythonClient); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "holochain-client-python", pythonClient); err != nil {
			return err
		}

		//
		// Holochain Python Serialization
		//
		pythonSerializationRepositoryArgs := StandardRepositoryArgs("holochain-serialization-python", nil)
		pythonSerialization, err := github.NewRepository(ctx, "holochain-serialization-python", &pythonSerializationRepositoryArgs, pulumi.Import(pulumi.ID("holochain-serialization-python")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "holochain-serialization-python", pythonSerialization); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "holochain-serialization-python", pythonSerialization); err != nil {
			return err
		}

		//
		// Nix Cache Check
		//
		nixCacheCheckRepositoryArgs := StandardRepositoryArgs("nix-cache-check", nil)
		nixCacheCheck, err := github.NewRepository(ctx, "nix-cache-check", &nixCacheCheckRepositoryArgs, pulumi.Import(pulumi.ID("nix-cache-check")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "nix-cache-check", nixCacheCheck); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "nix-cache-check", nixCacheCheck); err != nil {
			return err
		}

		//
		// Kitsune2
		//
		kitsune2RepositoryArgs := StandardRepositoryArgs("kitsune2", nil)
		kitsune2RepositoryArgs.Description = pulumi.String("p2p / dht communication framework")
		kitsune2, err := github.NewRepository(ctx, "kitsune2", &kitsune2RepositoryArgs, pulumi.Import(pulumi.ID("kitsune2")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "kitsune2", kitsune2); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "kitsune2", kitsune2); err != nil {
			return err
		}
        kitsune2RepositoryRulesetArgs := DefaultRepositoryRulesetArgs(kitsune2)
        if _, err = github.NewRepositoryRuleset(ctx, "kitsune2", &kitsune2RepositoryRulesetArgs); err != nil {
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
		AllowUpdateBranch:   pulumi.Bool(true),
		AllowSquashMerge:    pulumi.Bool(true),
		AllowRebaseMerge:    pulumi.Bool(true),
		AllowMergeCommit:    pulumi.Bool(false),
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

func MigrateDefaultBranchToMain(ctx *pulumi.Context, name string, repository *github.Repository) error {
	_, err := github.NewBranchDefault(ctx, fmt.Sprintf("%s-default-branch-migrate", name), &github.BranchDefaultArgs{
		Repository: repository.Name,
		Branch:     pulumi.String("main"),
		Rename:     pulumi.Bool(true),
	})
	return err
}

func DefaultRepositoryRulesetArgs(repository *github.Repository) github.RepositoryRulesetArgs {
	return github.RepositoryRulesetArgs{
		Name:        pulumi.String("default"),
		Repository:  repository.Name,
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
            PullRequest: &github.RepositoryRulesetRulesPullRequestArgs{
                DismissStaleReviewsOnPush:      pulumi.Bool(true),
                RequireCodeOwnerReview:         pulumi.Bool(false),
                RequireLastPushApproval:        pulumi.Bool(true),
                RequiredApprovingReviewCount:   pulumi.Int(1),
                RequiredReviewThreadResolution: pulumi.Bool(true),
            },
			RequiredStatusChecks: &github.RepositoryRulesetRulesRequiredStatusChecksArgs{
				RequiredChecks: github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArray{
					github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
					    // Each repository should define a single job that checks all the required checks passed.
						Context: pulumi.String("ci_pass"),
					},
				},
				StrictRequiredStatusChecksPolicy: pulumi.Bool(true),
			},
		},
	}
}
