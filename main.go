package main

import (
	"fmt"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")

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
		selfDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(self, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "hc-github-config-default", &selfDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		selfReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(self, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "hc-github-config-release", &selfReleaseRepositoryRulesetArgs); err != nil {
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
		if err = AddGithubAdminTokenSecret(ctx, conf, "hc-github-config"); err != nil {
			return err
		}
		if err = AddPulumiAccessTokenSecret(ctx, conf, "hc-github-config"); err != nil {
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
		jsClientDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(jsClient, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "holochain-client-js-default", &jsClientDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		jsClientReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(jsClient, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "holochain-client-js-release", &jsClientReleaseRepositoryRulesetArgs); err != nil {
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
		nixCacheCheckDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(nixCacheCheck, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "nix-cache-check-default", &nixCacheCheckDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		nixCacheCheckReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(nixCacheCheck, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "nix-cache-check-release", &nixCacheCheckReleaseRepositoryRulesetArgs); err != nil {
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
		kitsune2RepositoryRulesetArgs := DefaultRepositoryRulesetArgs(kitsune2, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "kitsune2-default", &kitsune2RepositoryRulesetArgs); err != nil {
			return err
		}
		kitsune2ReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(kitsune2, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "kitsune2-release", &kitsune2ReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddGithubUserTokenSecret(ctx, conf, "kitsune2"); err != nil {
			return err
		}

		//
		// ghost_actor
		//
		ghostActorRepositoryArgs := StandardRepositoryArgs("ghost_actor", nil)
		ghostActorRepositoryArgs.Description = pulumi.String("GhostActor makes it simple, ergonomic, and idiomatic to implement async / concurrent code using an Actor model.")
		ghostActor, err := github.NewRepository(ctx, "ghost_actor", &ghostActorRepositoryArgs, pulumi.Import(pulumi.ID("ghost_actor")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "ghost_actor", ghostActor); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "ghost_actor", ghostActor); err != nil {
			return err
		}
		ghostActorDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(ghostActor, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "ghost_actor-default", &ghostActorDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		ghostActorReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(ghostActor, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "ghost_actor-release", &ghostActorReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// docs-pages
		//
		docsPagesRepositoryArgs := StandardRepositoryArgs("docs-pages", nil)
		docsPagesRepositoryArgs.Description = pulumi.String("The hosted static files for the Holochain developer documentation")
		docsPagesRepositoryArgs.HasDiscussions = pulumi.Bool(true)
		docsPagesRepositoryArgs.HomepageUrl = pulumi.String("https://developer.holochain.org")
		docsPages, err := github.NewRepository(ctx, "docs-pages", &docsPagesRepositoryArgs, pulumi.Import(pulumi.ID("docs-pages")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "docs-pages", docsPages); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "docs-pages", docsPages); err != nil {
			return err
		}
		docsPagesDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(docsPages, []github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
			{
				IntegrationId: pulumi.Int(13473), // Netlify
				Context:       pulumi.String("Header rules - developer-portal-production"),
			},
			{
				IntegrationId: pulumi.Int(13473), // Netlify
				Context:       pulumi.String("netlify/developer-portal-production/deploy-preview"),
			},
			{
				IntegrationId: pulumi.Int(13473), // Netlify
				Context:       pulumi.String("Redirect rules - developer-portal-production"),
			},
		})
		if _, err = github.NewRepositoryRuleset(ctx, "docs-pages-default", &docsPagesDefaultRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// scaffolding
		//
		scaffoldingDescription := "Scaffolding tool to quickly generate and modify holochain applications"
		scaffoldingRepositoryArgs := StandardRepositoryArgs("scaffolding", &scaffoldingDescription)
		scaffoldingRepositoryArgs.HomepageUrl = pulumi.String("https://docs.rs/holochain_scaffolding_cli")
		scaffolding, err := github.NewRepository(ctx, "scaffolding", &scaffoldingRepositoryArgs, pulumi.Import(pulumi.ID("scaffolding")))
		if err != nil {
			return err
		}
		if err = MigrateDefaultBranchToMain(ctx, "scaffolding", scaffolding); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "scaffolding", scaffolding); err != nil {
			return err
		}
		if err = AdditionalRepositoryAdmin(ctx, "scaffolding", "c12i", scaffolding); err != nil {
			return err
		}
		scaffoldingDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(scaffolding, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "scaffolding-default", &scaffoldingDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		scaffoldingReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(scaffolding, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "scaffolding-release", &scaffoldingReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// hc-launch
		//
		hcLaunchDescription := "tauri based CLI to run holochain apps in development mode"
		hcLaunchRepositoryArgs := StandardRepositoryArgs("hc-launch", &hcLaunchDescription)

		hcLaunch, err := github.NewRepository(ctx, "hc-launch", &hcLaunchRepositoryArgs, pulumi.Import(pulumi.ID("hc-launch")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "hc-launch", hcLaunch); err != nil {
		  return err
		}
		if err = StandardRepositoryAccess(ctx, "hc-launch", hcLaunch); err != nil {
			return err
		}
		if err = AdditionalRepositoryAdmin(ctx, "hc-launch", "matthme", hcLaunch); err != nil {
			return err
		}
		hcLaunchDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcLaunch, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "hc-launch-default", &hcLaunchDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		hcLaunchReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(hcLaunch, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "hc-launch-release", &hcLaunchReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// hc-spin
		//
		hcSpinDescription := "CLI to run Holochain Apps in Development Mode"
		hcSpinRepositoryArgs := StandardRepositoryArgs("hc-spin", &hcSpinDescription)

		hcSpin, err := github.NewRepository(ctx, "hc-spin", &hcSpinRepositoryArgs, pulumi.Import(pulumi.ID("hc-spin")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "hc-spin", hcSpin); err != nil {
		  return err
		}
		if err = StandardRepositoryAccess(ctx, "hc-spin", hcSpin); err != nil {
			return err
		}
		if err = AdditionalRepositoryAdmin(ctx, "hc-spin", "matthme", hcSpin); err != nil {
			return err
		}
		hcSpinDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcSpin, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "hc-spin-default", &hcSpinDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		hcSpinReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(hcSpin, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "hc-spin-release", &hcSpinReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// kangaroo-electron
		//
		kangarooElectronDescription := "Bundle your holochain app a a standalone electron app with a built-in conductor"
		kangarooElectronRepositoryArgs := StandardRepositoryArgs("kangaroo-electron", &kangarooElectronDescription)
		kangarooElectronRepositoryArgs.IsTemplate = pulumi.Bool(true)

		kangarooElectron, err := github.NewRepository(ctx, "kangaroo-electron", &kangarooElectronRepositoryArgs, pulumi.Import(pulumi.ID("kangaroo-electron")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "kangaroo-electron", kangarooElectron); err != nil {
		  return err
		}
		if err = StandardRepositoryAccess(ctx, "kangaroo-electron", kangarooElectron); err != nil {
			return err
		}
		if err = AdditionalRepositoryAdmin(ctx, "kangaroo-electron", "matthme", kangarooElectron); err != nil {
			return err
		}
		kangarooElectronDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(kangarooElectron, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "kangaroo-electron-default", &kangarooElectronDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		// Since kangaroo is a Github Template we currently omit mandatory CI checks
		kangarooElectronDefaultRepositoryRulesetArgs.Rules = &github.RepositoryRulesetRulesArgs{
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
		}

		kangarooElectronReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(kangarooElectron, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "kangaroo-electron-release", &kangarooElectronReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		// Since kangaroo is a Github Template we currently omit mandatory CI checks
		kangarooElectronReleaseRepositoryRulesetArgs.Rules = &github.RepositoryRulesetRulesArgs{
			Creation:              pulumi.Bool(true),
			Update:                pulumi.Bool(false),
			Deletion:              pulumi.Bool(true),
			RequiredLinearHistory: pulumi.Bool(true),
			RequiredSignatures:    pulumi.Bool(false),
			PullRequest: &github.RepositoryRulesetRulesPullRequestArgs{
				DismissStaleReviewsOnPush:      pulumi.Bool(true),
				RequireCodeOwnerReview:         pulumi.Bool(false),
				RequireLastPushApproval:        pulumi.Bool(true),
				RequiredApprovingReviewCount:   pulumi.Int(0),
				RequiredReviewThreadResolution: pulumi.Bool(true),
			},
		}

		//
		// Dino Adventure
		//
		dinoAdventureDescription := "A dinosaur adventure game for testing Holochain"
		dinoAdventureRepositoryArgs := StandardRepositoryArgs("dino-adventure", &dinoAdventureDescription)
		dinoAdventure, err := github.NewRepository(ctx, "dino-adventure", &dinoAdventureRepositoryArgs)
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "dino-adventure", dinoAdventure); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "dino-adventure", dinoAdventure); err != nil {
			return err
		}
		dinoAdventureDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(dinoAdventure, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "dino-adventure-default", &dinoAdventureDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		dinoAdventureReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(dinoAdventure, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "dino-adventure-release", &dinoAdventureReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// nomad-server
		//
		nomadServerDescription := "A Pulumi definition for deploying a cluster of Nomad servers as DigitalOcean droplets"
		nomadServerRepositoryArgs := StandardRepositoryArgs("nomad-server", &nomadServerDescription)
		nomadServer, err := github.NewRepository(ctx, "nomad-server", &nomadServerRepositoryArgs)
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "nomad-server", nomadServer); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "nomad-server", nomadServer); err != nil {
			return err
		}
		nomadServerDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(nomadServer, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "nomad-server-default", &nomadServerDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddGithubUserTokenSecret(ctx, conf, "nomad-server"); err != nil {
			return err
		}
		if err = AddPulumiAccessTokenSecret(ctx, conf, "nomad-server"); err != nil {
			return err
		}

		//
		// hc-http-gw
		//
		hcHttpGwDescription := "The Holochain HTTP Gateway for providing a way to bridge from the web2 world into Holochain"
		hcHttpGwRepositoryArgs := StandardRepositoryArgs("hc-http-gw", &hcHttpGwDescription)
		hcHttpGw, err := github.NewRepository(ctx, "hc-http-gw", &hcHttpGwRepositoryArgs)
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "hc-http-gw", hcHttpGw); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "hc-http-gw", hcHttpGw); err != nil {
			return err
		}
		hcHttpGwDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcHttpGw, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "hc-http-gw-default", &hcHttpGwDefaultRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// network-services
		//
		networkServicesDescription := "A Pulumi definition for deploying Holochain network services to be used for development"
		networkServicesRepositoryArgs := StandardRepositoryArgs("network-services", &networkServicesDescription)
		networkServices, err := github.NewRepository(ctx, "network-services", &networkServicesRepositoryArgs)
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "network-services", networkServices); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "network-services", networkServices); err != nil {
			return err
		}
		networkServicesDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(networkServices, nil)
		if _, err = github.NewRepositoryRuleset(ctx, "network-services-default", &networkServicesDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddGithubUserTokenSecret(ctx, conf, "network-services"); err != nil {
			return err
		}
		if err = AddPulumiAccessTokenSecret(ctx, conf, "network-services"); err != nil {
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
		AutoInit:            pulumi.Bool(true),
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

func AdditionalRepositoryAdmin(ctx *pulumi.Context, resourceName string, username string, repository *github.Repository) error {
	_, err := github.NewRepositoryCollaborator(ctx, fmt.Sprintf("%s-collaborator-%s", resourceName, username), &github.RepositoryCollaboratorArgs{
		Repository: repository.Name,
		Username:   pulumi.String(username),
		Permission: pulumi.String("admin"),
	})

	return err
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

func DefaultRepositoryRulesetArgs(repository *github.Repository, extraStatusChecks []github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs) github.RepositoryRulesetArgs {
	requiredChecks := github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArray{
		github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
			// Each repository should define a single job that checks all the required checks passed.
			Context: pulumi.String("ci_pass"),
		},
	}
	if extraStatusChecks != nil {
		for _, check := range extraStatusChecks {
			requiredChecks = append(requiredChecks, &check)
		}
	}

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
				RequiredChecks:                   requiredChecks,
				StrictRequiredStatusChecksPolicy: pulumi.Bool(true),
			},
		},
	}
}

func ReleaseRepositoryRulesetArgs(repository *github.Repository, extraStatusChecks []github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs) github.RepositoryRulesetArgs {
	requiredChecks := github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArray{
		github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
			// Each repository should define a single job that checks all the required checks passed.
			Context: pulumi.String("ci_pass"),
		},
	}
	if extraStatusChecks != nil {
		for _, check := range extraStatusChecks {
			requiredChecks = append(requiredChecks, &check)
		}
	}

	return github.RepositoryRulesetArgs{
		Name:        pulumi.String("release"),
		Repository:  repository.Name,
		Target:      pulumi.String("branch"),
		Enforcement: pulumi.String("active"),
		Conditions: &github.RepositoryRulesetConditionsArgs{
			RefName: &github.RepositoryRulesetConditionsRefNameArgs{
				Includes: pulumi.StringArray{
					pulumi.String("refs/heads/release/*"),
					pulumi.String("refs/heads/main-*"),
					pulumi.String("refs/heads/develop-*"),
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
				RequiredApprovingReviewCount:   pulumi.Int(0),
				RequiredReviewThreadResolution: pulumi.Bool(true),
			},
			RequiredStatusChecks: &github.RepositoryRulesetRulesRequiredStatusChecksArgs{
				RequiredChecks:                   requiredChecks,
				StrictRequiredStatusChecksPolicy: pulumi.Bool(true),
			},
		},
	}
}

func AddGithubUserTokenSecret(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	// A GITHUB_TOKEN with standard repository access to be used on most repositories.
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-github-token", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("HRA2_GITHUB_TOKEN"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("hra2GithubUserToken"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	return nil
}

func AddGithubAdminTokenSecret(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	// A GITHUB_TOKEN with more access rights to larger scopes.
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-github-token", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("HRA2_GITHUB_TOKEN"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("hra2GithubAdminToken"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	return nil
}

func AddPulumiAccessTokenSecret(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-pulumi-access-token", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("HRA2_PULUMI_ACCESS_TOKEN"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("hra2PulumiAccessToken"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	return nil
}
