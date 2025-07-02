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
		selfDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(self, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "hc-github-config-default", &selfDefaultRepositoryRulesetArgs); err != nil {
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
							Context: pulumi.String("ci_pass"),
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
		holochainWasmerDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(holochainWasmer, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "holochain-wasmer-default", &holochainWasmerDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		holochainWasmerReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(holochainWasmer, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "holochain-wasmer-release", &holochainWasmerReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "holochain-wasmer", holochainWasmer); err != nil {
			return err
		}

		//
		// wind tunnel
		//
		description = "Performance testing for Holochain"
		windTunnelRepositoryArgs := StandardRepositoryArgs("wind-tunnel", &description)
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
		windTunnelDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(windTunnel, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "wind-tunnel-default", &windTunnelDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		windTunnelReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(windTunnel, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "wind-tunnel-release", &windTunnelReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		windTunnelConf := config.New(ctx, "wind-tunnel")
		if err = AddNomadAccessTokenSecret(ctx, windTunnelConf, "wind-tunnel"); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "wind-tunnel", windTunnel); err != nil {
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
		jsClientDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(jsClient, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "holochain-client-js-default", &jsClientDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		jsClientReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(jsClient, NewRulesetOptions())
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
		tryoramaDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(tryorama, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "tryorama-default", &tryoramaDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		tryoramaReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(tryorama, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "tryorama-release", &tryoramaReleaseRepositoryRulesetArgs); err != nil {
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
		holonixDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(holonix, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "holonix-default", &holonixDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		holonixReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(holonix, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "holonix-release", &holonixReleaseRepositoryRulesetArgs); err != nil {
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
		sbdDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(sbd, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "sbd-default", &sbdDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		sbdReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(sbd, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "sbd-release", &sbdReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "sbd", sbd); err != nil {
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
		tx5DefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(tx5, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "tx5-default", &tx5DefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		tx5ReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(tx5, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "tx5-release", &tx5ReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "tx5", tx5); err != nil {
			return err
		}

		//
		// Lair Keystore
		//
		description = "secret lair private keystore"
		lairRepositoryArgs := StandardRepositoryArgs("lair", &description)
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
		lairDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(lair, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "lair-default", &lairDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		lairReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(lair, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "lair-release", &lairReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "lair", lair); err != nil {
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
		hcChcServiceDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcChcService, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "hc-chc-service-default", &hcChcServiceDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		hcChcServiceReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(hcChcService, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "hc-chc-service-release", &hcChcServiceReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "hc-chc-service", hcChcService); err != nil {
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
		holochainSerializationDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(holochainSerialization, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "holochain-serialization-default", &holochainSerializationDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		holochainSerializationReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(holochainSerialization, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "holochain-serialization-release", &holochainSerializationReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "holochain-serialization", holochainSerialization); err != nil {
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
		influxiveDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(influxive, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "influxive-default", &influxiveDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		influxiveReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(influxive, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "influxive-release", &influxiveReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "influxive", influxive); err != nil {
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
		nixCacheCheckDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(nixCacheCheck, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "nix-cache-check-default", &nixCacheCheckDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		nixCacheCheckReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(nixCacheCheck, NewRulesetOptions())
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
		kitsune2RepositoryRulesetArgs := DefaultRepositoryRulesetArgs(kitsune2, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "kitsune2-default", &kitsune2RepositoryRulesetArgs); err != nil {
			return err
		}
		kitsune2ReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(kitsune2, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "kitsune2-release", &kitsune2ReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "kitsune2", kitsune2); err != nil {
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
		docsPagesDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(docsPages, NewRulesetOptions().withExtraStatusChecks([]github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
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
		}))
		if _, err = github.NewRepositoryRuleset(ctx, "docs-pages-default", &docsPagesDefaultRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// hc-static-site
		//
		// There's enough different about this repo's config that it makes
		// sense just to start from scratch.
		hcStaticSiteRepositoryArgs := github.RepositoryArgs{
			Name:                pulumi.String("hc-static-site"),
			Description:         pulumi.String("Static website"),
			Visibility:          pulumi.String("private"),
			HasDownloads:        pulumi.Bool(true),
			HasIssues:           pulumi.Bool(true),
			HasProjects:         pulumi.Bool(true),
			HasWiki:             pulumi.Bool(false),
			VulnerabilityAlerts: pulumi.Bool(false),
			AllowAutoMerge:      pulumi.Bool(false),
			DeleteBranchOnMerge: pulumi.Bool(false),
			AllowUpdateBranch:   pulumi.Bool(false),
			AllowSquashMerge:    pulumi.Bool(false),
			AllowRebaseMerge:    pulumi.Bool(false),
			AllowMergeCommit:    pulumi.Bool(false),
			AutoInit:            pulumi.Bool(false),
			VulnerabilityAlerts: pulumi.Bool(false),
		}
		hcStaticSite, err := github.NewRepository(ctx, "hc-static-site", &hcStaticSiteRepositoryArgs, pulumi.Import(pulumi.ID("hc-static-site")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "hc-static-site", hcStaticSite); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "hc-static-site", hcStaticSite); err != nil {
			return err
		}
		hcStaticSiteDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcStaticSite, NewRulesetOptions().withExtraStatusChecks([]github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
			{
				IntegrationId: pulumi.Int(13473), // Netlify
				Context:       pulumi.String("Header rules - holochain-prod"),
			},
			{
				IntegrationId: pulumi.Int(13473), // Netlify
				Context:       pulumi.String("netlify/holochain-prod/deploy-preview"),
			},
			{
				IntegrationId: pulumi.Int(13473), // Netlify
				Context:       pulumi.String("Redirect rules - holochain-prod"),
			},
		}))
		if _, err = github.NewRepositoryRuleset(ctx, "hc-static-site-default", &hcStaticSiteDefaultRepositoryRulesetArgs); err != nil {
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
		scaffoldingDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(scaffolding, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "scaffolding-default", &scaffoldingDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		scaffoldingReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(scaffolding, NewRulesetOptions().noLinearHistoryRequired())
		if _, err = github.NewRepositoryRuleset(ctx, "scaffolding-release", &scaffoldingReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "scaffolding", scaffolding); err != nil {
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
		hcLaunchDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcLaunch, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "hc-launch-default", &hcLaunchDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		hcLaunchReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(hcLaunch, NewRulesetOptions().noLinearHistoryRequired())
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
		hcSpinDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcSpin, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "hc-spin-default", &hcSpinDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		hcSpinReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(hcSpin, NewRulesetOptions().noLinearHistoryRequired())
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
		// Since kangaroo is a Github Template we currently omit mandatory CI checks
		kangarooElectronDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(kangarooElectron, NewRulesetOptions().noStatusChecks())
		if _, err = github.NewRepositoryRuleset(ctx, "kangaroo-electron-default", &kangarooElectronDefaultRepositoryRulesetArgs); err != nil {
			return err
		}

		kangarooElectronReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(kangarooElectron, NewRulesetOptions())
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
		if err = AddAppleAppSigningSecrets(ctx, conf, "kangaroo-electron"); err != nil {
			return err
		}
		if err = AddWindowsCodeSigningCertificates(ctx, conf, "kangaroo-electron"); err != nil {
			return err
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
		dinoAdventureDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(dinoAdventure, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "dino-adventure-default", &dinoAdventureDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		dinoAdventureReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(dinoAdventure, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "dino-adventure-release", &dinoAdventureReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// Dino Adventure - Kangaroo
		//
		dinoAdventureKangarooDescription := "Kangaroo packaging for the dino adventure app"
		dinoAdventureKangarooRepositoryArgs := StandardRepositoryArgs("dino-adventure-kangaroo", &dinoAdventureKangarooDescription)
		dinoAdventureKangaroo, err := github.NewRepository(ctx, "dino-adventure-kangaroo", &dinoAdventureKangarooRepositoryArgs, pulumi.Import(pulumi.ID("dino-adventure-kangaroo")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "dino-adventure-kangaroo", dinoAdventureKangaroo); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "dino-adventure-kangaroo", dinoAdventureKangaroo); err != nil {
			return err
		}
		if err = AddAppleAppSigningSecrets(ctx, conf, "dino-adventure-kangaroo"); err != nil {
			return err
		}
		if err = AddWindowsCodeSigningCertificates(ctx, conf, "dino-adventure-kangaroo"); err != nil {
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
		nomadServerDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(nomadServer, NewRulesetOptions())
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
		hcHttpGwDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(hcHttpGw, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "hc-http-gw-default", &hcHttpGwDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		hcHttpGwReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(hcHttpGw, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "hc-http-gw-release", &hcHttpGwReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "hc-http-gw", hcHttpGw); err != nil {
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
		networkServicesDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(networkServices, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "network-services-default", &networkServicesDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddGithubUserTokenSecret(ctx, conf, "network-services"); err != nil {
			return err
		}
		if err = AddPulumiAccessTokenSecret(ctx, conf, "network-services"); err != nil {
			return err
		}

		//
		// wind-tunnel-runner
		//
		windTunnelRunnerDescription := "The guide and NixOS configuration for setting up a machine to run Wind Tunnel scenarios"
		windTunnelRunnerRepositoryArgs := StandardRepositoryArgs("wind-tunnel-runner", &windTunnelRunnerDescription)
		windTunnelRunner, err := github.NewRepository(ctx, "wind-tunnel-runner", &windTunnelRunnerRepositoryArgs)
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "wind-tunnel-runner", windTunnelRunner); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "wind-tunnel-runner", windTunnelRunner); err != nil {
			return err
		}
		windTunnelRunnerDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(windTunnelRunner, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "wind-tunnel-runner-default", &windTunnelRunnerDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddTailscaleOAuthSecrets(ctx, conf, "wind-tunnel-runner"); err != nil {
			return err
		}

		//
		// must_future
		//
		mustFutureDescription := "A wrapper future marked must_use - mainly to wrap BoxFutures"
		mustFutureRepositoryArgs := StandardRepositoryArgs("must_future", &mustFutureDescription)
		mustFuture, err := github.NewRepository(ctx, "must_future", &mustFutureRepositoryArgs, pulumi.Import(pulumi.ID("must_future")))
		if err != nil {
			return err
		}
		if err = MigrateDefaultBranchToMain(ctx, "must_future", mustFuture); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "must_future", mustFuture); err != nil {
			return err
		}
		mustFutureDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(mustFuture, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "must_future-default", &mustFutureDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		mustFutureReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(mustFuture, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "must_future-release", &mustFutureReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// url2
		//
		url2Description := "ergonomic wrapper around the popular url crate"
		url2RepositoryArgs := StandardRepositoryArgs("url2", &url2Description)
		url2, err := github.NewRepository(ctx, "url2", &url2RepositoryArgs, pulumi.Import(pulumi.ID("url2")))
		if err != nil {
			return err
		}
		if err = MigrateDefaultBranchToMain(ctx, "url2", url2); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "url2", url2); err != nil {
			return err
		}
		url2DefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(url2, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "url2-default", &url2DefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		url2ReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(url2, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "url2-release", &url2ReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// automap-rs
		//
		automapRsDescription := "Simple pattern for expressing Rust maps where the Value type contains the Key"
		automapRsRepositoryArgs := StandardRepositoryArgs("automap-rs", &automapRsDescription)
		automapRs, err := github.NewRepository(ctx, "automap-rs", &automapRsRepositoryArgs, pulumi.Import(pulumi.ID("automap-rs")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "automap-rs", automapRs); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "automap-rs", automapRs); err != nil {
			return err
		}
		automapRsDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(automapRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "automap-rs-default", &automapRsDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		automapRsReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(automapRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "automap-rs-release", &automapRsReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// rand-utf8
		//
		randUtf8Description := "Random utf8 utility"
		randUtf8RepositoryArgs := StandardRepositoryArgs("rand-utf8", &randUtf8Description)
		randUtf8, err := github.NewRepository(ctx, "rand-utf8", &randUtf8RepositoryArgs, pulumi.Import(pulumi.ID("rand-utf8")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "rand-utf8", randUtf8); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "rand-utf8", randUtf8); err != nil {
			return err
		}
		randUtf8DefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(randUtf8, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "rand-utf8-default", &randUtf8DefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		randUtf8ReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(randUtf8, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "rand-utf8-release", &randUtf8ReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "rand-utf8", randUtf8); err != nil {
			return err
		}

		//
		// serde-json
		//
		serdeJsonDescription := "Strongly typed JSON library for Rust"
		serdeJsonRepositoryArgs := StandardRepositoryArgs("serde-json", &serdeJsonDescription)
		serdeJson, err := github.NewRepository(ctx, "serde-json", &serdeJsonRepositoryArgs, pulumi.Import(pulumi.ID("serde-json")))
		if err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "serde-json", serdeJson); err != nil {
			return err
		}
		serdeJsonDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(serdeJson, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "serde-json-default", &serdeJsonDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		serdeJsonReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(serdeJson, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "serde-json-release", &serdeJsonReleaseRepositoryRulesetArgs); err != nil {
			return err
		}
		if err = AddReleaseIntegrationSupport(ctx, conf, "serde-json", serdeJson); err != nil {
			return err
		}

		//
		// isotest-rs
		//
		isoTestRsDescription := "Opinionated way to solve a very particular problem in Rust testing"
		isoTestRsRepositoryArgs := StandardRepositoryArgs("isotest-rs", &isoTestRsDescription)
		isoTestRs, err := github.NewRepository(ctx, "isotest-rs", &isoTestRsRepositoryArgs, pulumi.Import(pulumi.ID("isotest-rs")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "isotest-rs", isoTestRs); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "isotest-rs", isoTestRs); err != nil {
			return err
		}
		isoTestRsDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(isoTestRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "isotest-rs-default", &isoTestRsDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		isoTestRsReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(isoTestRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "isotest-rs-release", &isoTestRsReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// one_err
		//
		oneErrDescription := "OneErr to rule them all"
		oneErrRepositoryArgs := StandardRepositoryArgs("one_err", &oneErrDescription)
		oneErr, err := github.NewRepository(ctx, "one_err", &oneErrRepositoryArgs, pulumi.Import(pulumi.ID("one_err")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "one_err", oneErr); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "one_err", oneErr); err != nil {
			return err
		}
		oneErrDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(oneErr, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "one_err-default", &oneErrDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		oneErrReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(oneErr, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "one_err-release", &oneErrReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// bootstrap
		//
		bootstrapDescription := "Bootstrap nodes onto a network by allowing existing nodes to list themselves under a URL"
		bootstrapRepositoryArgs := StandardRepositoryArgs("bootstrap", &bootstrapDescription)
		bootstrap, err := github.NewRepository(ctx, "bootstrap", &bootstrapRepositoryArgs, pulumi.Import(pulumi.ID("bootstrap")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "bootstrap", bootstrap); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "bootstrap", bootstrap); err != nil {
			return err
		}
		bootstrapDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(bootstrap, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "bootstrap-default", &bootstrapDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		bootstrapReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(bootstrap, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "bootstrap-release", &bootstrapReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// ametrics
		//
		ametricsDescription := "ametrics metric abstraction helpers"
		ametricsRepositoryArgs := StandardRepositoryArgs("ametrics", &ametricsDescription)
		ametrics, err := github.NewRepository(ctx, "ametrics", &ametricsRepositoryArgs, pulumi.Import(pulumi.ID("ametrics")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "ametrics", ametrics); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "ametrics", ametrics); err != nil {
			return err
		}
		ametricsDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(ametrics, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "ametrics-default", &ametricsDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		ametricsReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(ametrics, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "ametrics-release", &ametricsReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// contrafact-rs
		//
		contrafactRsDescription := "Generate test fixtures and check data properties with declarative, modular constraints"
		contrafactRsRepositoryArgs := StandardRepositoryArgs("contrafact-rs", &contrafactRsDescription)
		contrafactRs, err := github.NewRepository(ctx, "contrafact-rs", &contrafactRsRepositoryArgs, pulumi.Import(pulumi.ID("contrafact-rs")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "contrafact-rs", contrafactRs); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "contrafact-rs", contrafactRs); err != nil {
			return err
		}
		contrafactRsDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(contrafactRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "contrafact-rs-default", &contrafactRsDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		contrafactRsReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(contrafactRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "contrafact-rs-release", &contrafactRsReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// task-motel-rs
		//
		taskMotelRsDescription := "An opinionated Tokio task manager"
		taskMotelRsRepositoryArgs := StandardRepositoryArgs("task-motel-rs", &taskMotelRsDescription)
		taskMotelRs, err := github.NewRepository(ctx, "task-motel-rs", &taskMotelRsRepositoryArgs, pulumi.Import(pulumi.ID("task-motel-rs")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "task-motel-rs", taskMotelRs); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "task-motel-rs", taskMotelRs); err != nil {
			return err
		}
		taskMotelRsDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(taskMotelRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "task-motel-rs-default", &taskMotelRsDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		taskMotelRsReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(taskMotelRs, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "task-motel-rs-release", &taskMotelRsReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// devhub-gui
		//
		devHubGuiDescription := "A web-based UI that works with Holochain's collection of DevHub DNAs."
		devHubGuiRepositoryArgs := StandardRepositoryArgs("devhub-gui", &devHubGuiDescription)
		devHubGui, err := github.NewRepository(ctx, "devhub-gui", &devHubGuiRepositoryArgs, pulumi.Import(pulumi.ID("devhub-gui")))
		if err != nil {
			return err
		}
		if err = MigrateDefaultBranchToMain(ctx, "devhub-gui", devHubGui); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "devhub-gui", devHubGui); err != nil {
			return err
		}
		devHubGuiDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(devHubGui, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "devhub-gui-default", &devHubGuiDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		devHubGuiReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(devHubGui, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "devhub-gui-release", &devHubGuiReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// app-store-gui
		//
		appStoreGuiDescription := "A web-based UI that works with Holochain's collection of App Store DNAs."
		appStoreGuiRepositoryArgs := StandardRepositoryArgs("app-store-gui", &appStoreGuiDescription)
		appStoreGui, err := github.NewRepository(ctx, "app-store-gui", &appStoreGuiRepositoryArgs, pulumi.Import(pulumi.ID("app-store-gui")))
		if err != nil {
			return err
		}
		if err = MigrateDefaultBranchToMain(ctx, "app-store-gui", appStoreGui); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "app-store-gui", appStoreGui); err != nil {
			return err
		}
		appStoreGuiDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(appStoreGui, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "app-store-gui-default", &appStoreGuiDefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		appStoreGuiReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(appStoreGui, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "app-store-gui-release", &appStoreGuiReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// bootstrap2
		//
		bootstrap2Description := "Holochain bootstrap peer discovery."
		bootstrap2RepositoryArgs := StandardRepositoryArgs("bootstrap2", &bootstrap2Description)
		bootstrap2, err := github.NewRepository(ctx, "bootstrap2", &bootstrap2RepositoryArgs, pulumi.Import(pulumi.ID("bootstrap2")))
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "bootstrap2", bootstrap2); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "bootstrap2", bootstrap2); err != nil {
			return err
		}
		bootstrap2DefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(bootstrap2, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "bootstrap2-default", &bootstrap2DefaultRepositoryRulesetArgs); err != nil {
			return err
		}
		bootstrap2ReleaseRepositoryRulesetArgs := ReleaseRepositoryRulesetArgs(bootstrap2, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "bootstrap2-release", &bootstrap2ReleaseRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// release-integration
		//
		releaseIntegrationDescription := "Integration of third-party release tools with Holochain repositories"
		releaseIntegrationRepositoryArgs := StandardRepositoryArgs("release-integration", &releaseIntegrationDescription)
		releaseIntegration, err := github.NewRepository(ctx, "release-integration", &releaseIntegrationRepositoryArgs)
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "release-integration", releaseIntegration); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "release-integration", releaseIntegration); err != nil {
			return err
		}
		releaseIntegrationDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(releaseIntegration, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "release-integration-default", &releaseIntegrationDefaultRepositoryRulesetArgs); err != nil {
			return err
		}

		//
		// actions
		//
		actionsDescription := "Actions for common tasks in Holochain repositories"
		actionsRepositoryArgs := StandardRepositoryArgs("actions", &actionsDescription)
		actions, err := github.NewRepository(ctx, "actions", &actionsRepositoryArgs)
		if err != nil {
			return err
		}
		if err = RequireMainAsDefaultBranch(ctx, "actions", actions); err != nil {
			return err
		}
		if err = StandardRepositoryAccess(ctx, "actions", actions); err != nil {
			return err
		}
		actionsDefaultRepositoryRulesetArgs := DefaultRepositoryRulesetArgs(actions, NewRulesetOptions())
		if _, err = github.NewRepositoryRuleset(ctx, "actions-default", &actionsDefaultRepositoryRulesetArgs); err != nil {
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

type RulesetOptions struct {
	extraStatusChecks   []github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs
	withoutStatusChecks bool
	noLinearHistory     bool
}

func NewRulesetOptions() RulesetOptions {
	return RulesetOptions{}
}

func (options RulesetOptions) withExtraStatusChecks(extraStatusChecks []github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs) RulesetOptions {
	if options.withoutStatusChecks {
		panic("withExtraStatusChecks() cannot be called if noStatusChecks() has already been called.")
	}
	options.extraStatusChecks = extraStatusChecks
	return options
}

func (options RulesetOptions) noLinearHistoryRequired() RulesetOptions {
	options.noLinearHistory = true
	return options
}

func (options RulesetOptions) noStatusChecks() RulesetOptions {
	if options.extraStatusChecks != nil {
		panic("noStatusChecks() cannot be called if extraStatusChecks() has already been called.")
	}
	options.withoutStatusChecks = true
	return options
}

func DefaultRepositoryRulesetArgs(repository *github.Repository, options RulesetOptions) github.RepositoryRulesetArgs {
	requiredChecks := github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArray{
		github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
			// Each repository should define a single job that checks all the required checks passed.
			Context: pulumi.String("ci_pass"),
		},
	}
	if options.extraStatusChecks != nil {
		for _, check := range options.extraStatusChecks {
			requiredChecks = append(requiredChecks, &check)
		}
	}
	linearHistory := pulumi.Bool(true)
	if options.noLinearHistory {
		linearHistory = pulumi.Bool(false)
	}
	requiredStatusChecks := &github.RepositoryRulesetRulesRequiredStatusChecksArgs{
		RequiredChecks:                   requiredChecks,
		StrictRequiredStatusChecksPolicy: pulumi.Bool(true),
	}
	if options.withoutStatusChecks {
		requiredStatusChecks = nil
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
			RequiredLinearHistory: linearHistory,
			RequiredSignatures:    pulumi.Bool(false),
			PullRequest: &github.RepositoryRulesetRulesPullRequestArgs{
				DismissStaleReviewsOnPush:      pulumi.Bool(true),
				RequireCodeOwnerReview:         pulumi.Bool(false),
				RequireLastPushApproval:        pulumi.Bool(true),
				RequiredApprovingReviewCount:   pulumi.Int(1),
				RequiredReviewThreadResolution: pulumi.Bool(true),
			},
			RequiredStatusChecks: requiredStatusChecks,
		},
	}
}

func ReleaseRepositoryRulesetArgs(repository *github.Repository, options RulesetOptions) github.RepositoryRulesetArgs {
	requiredChecks := github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArray{
		github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
			// Each repository should define a single job that checks all the required checks passed.
			Context: pulumi.String("ci_pass"),
		},
	}
	if options.extraStatusChecks != nil {
		for _, check := range options.extraStatusChecks {
			requiredChecks = append(requiredChecks, &check)
		}
	}
	if options.withoutStatusChecks {
		requiredChecks = nil
	}
	linearHistory := pulumi.Bool(true)
	if options.noLinearHistory {
		linearHistory = false
	}
	requiredStatusChecks := &github.RepositoryRulesetRulesRequiredStatusChecksArgs{
		RequiredChecks:                   requiredChecks,
		DoNotEnforceOnCreate:             pulumi.Bool(true),
		StrictRequiredStatusChecksPolicy: pulumi.Bool(true),
	}
	if options.withoutStatusChecks {
		requiredStatusChecks = nil
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
					pulumi.String("refs/heads/release-*"),
					pulumi.String("refs/heads/main-*"),
					pulumi.String("refs/heads/develop-*"),
				},
				Excludes: pulumi.StringArray{},
			},
		},
		Rules: &github.RepositoryRulesetRulesArgs{
			Creation:              pulumi.Bool(false),
			Update:                pulumi.Bool(false),
			Deletion:              pulumi.Bool(true),
			RequiredLinearHistory: linearHistory,
			RequiredSignatures:    pulumi.Bool(false),
			PullRequest: &github.RepositoryRulesetRulesPullRequestArgs{
				DismissStaleReviewsOnPush:      pulumi.Bool(true),
				RequireCodeOwnerReview:         pulumi.Bool(false),
				RequireLastPushApproval:        pulumi.Bool(true),
				RequiredApprovingReviewCount:   pulumi.Int(0),
				RequiredReviewThreadResolution: pulumi.Bool(true),
			},
			RequiredStatusChecks: requiredStatusChecks,
		},
		BypassActors: github.RepositoryRulesetBypassActorArray{
			&github.RepositoryRulesetBypassActorArgs{
				ActorId:    pulumi.Int(5),
				ActorType:  pulumi.String("RepositoryRole"),
				BypassMode: pulumi.String("always"),
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

	return err
}

func AddGithubAdminTokenSecret(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	// A GITHUB_TOKEN with more access rights to larger scopes.
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-github-token", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("HRA2_GITHUB_TOKEN"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("hra2GithubAdminToken"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))

	return err
}

func AddCratesIoTokenSecret(ctx *pulumi.Context, cfg *config.Config, name string, repository *github.Repository) error {
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-crates-io-token", name), &github.ActionsSecretArgs{
		Repository: repository.Name,
		SecretName: pulumi.String("HRA2_CRATES_IO_TOKEN"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("hra2CratesIoToken"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))

	return err
}

func AddPulumiAccessTokenSecret(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-pulumi-access-token", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("HRA2_PULUMI_ACCESS_TOKEN"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("hra2PulumiAccessToken"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))

	return err
}

func AddNomadAccessTokenSecret(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-nomad-access-token", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("NOMAD_ACCESS_TOKEN"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("nomadAccessToken"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))

	return err
}

func AddTailscaleOAuthSecrets(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-tailscale-oauth-client-id", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("TS_OAUTH_CLIENT_ID"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("tailscaleOAuthClientId"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-tailscale-oauth-secret", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("TS_OAUTH_SECRET"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("tailscaleOAuthSecret"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))

	return err
}

func AddAppleAppSigningSecrets(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-apple-dev-identity", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("APPLE_DEV_IDENTITY"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("appleDevIdentity"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-apple-id-email", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("APPLE_ID_EMAIL"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("appleIdEmail"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-apple-id-password", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("APPLE_ID_PASSWORD"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("appleIdPassword"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-apple-team-id", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("APPLE_TEAM_ID"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("appleTeamId"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-apple-certificate", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("APPLE_CERTIFICATE"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("appleCertificate"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-apple-certificate-password", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("APPLE_CERTIFICATE_PASSWORD"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("appleCertificatePassword"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	return err
}

func AddWindowsCodeSigningCertificates(ctx *pulumi.Context, cfg *config.Config, repository string) error {
	_, err := github.NewActionsSecret(ctx, fmt.Sprintf("%s-azure-key-vault-uri", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("AZURE_KEY_VAULT_URI"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("azureKeyVaultUri"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-azure-cert-name", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("AZURE_CERT_NAME"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("azureCertName"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-azure-tenant-id", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("AZURE_TENANT_ID"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("azureTenantId"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-azure-client-id", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("AZURE_CLIENT_ID"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("azureClientId"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))
	if err != nil {
		return err
	}

	_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-azure-client-secret", repository), &github.ActionsSecretArgs{
		Repository: pulumi.String(repository),
		SecretName: pulumi.String("AZURE_CLIENT_SECRET"),
		// The GitHub API only accepts encrypted values. This will be encrypted by the provider before being sent.
		PlaintextValue: cfg.RequireSecret("azureClientSecret"),
	}, pulumi.DeleteBeforeReplace(true), pulumi.IgnoreChanges([]string{"encryptedValue"}))

	return err
}

func AddReleaseIntegrationSupport(ctx *pulumi.Context, cfg *config.Config, name string, repository *github.Repository) error {
	if _, err := github.NewIssueLabel(ctx, fmt.Sprintf("%s-hra-release-label", name), &github.IssueLabelArgs{
		Repository: repository.Name,
		// Must match what the holochain_release_integration CLI looks for.
		Name: pulumi.String("hra-release"),
		// Golden Fizz
		Color: pulumi.String("E8F723"),
	}); err != nil {
		return err
	}

	if err := AddGithubUserTokenSecret(ctx, cfg, name); err != nil {
		return err
	}
	if err := AddCratesIoTokenSecret(ctx, cfg, name, repository); err != nil {
		return err
	}

	return nil
}
