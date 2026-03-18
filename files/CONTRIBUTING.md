# Contributing to {{REPO_NAME}}

Thank you for your interest in contributing to {{REPO_NAME}}!

## Getting Started

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Submit a pull request

## Commits

The commits are used to generate the changelog upon a release, therefore keep
them clean. To help with keeping them clean, follow these principles:

- Use [conventional commits](https://www.conventionalcommits.org/) for commit
  messages. You can also use markdown, especially in the body
- Keep commits atomic, meaning that each commit should build and the unit tests
  should pass. If a new feature is added or a bug is fixed then fix the tests
  in the same commit
- In case of changes on the origin branch, always rebase, **NEVER** merge
- Avoid referencing the issue or PR number in the commit message as the
  automated changelog will add this reference for you

## Pull Requests

- All changes require a pull request and at least one approving review
- Add a short, clear, and hand-written description of the change proposed by
  the PR. AI tooling may add an in-depth summary below your description so keep
  it high-level but make sure to add one
- PRs must pass CI checks before merging
- Keep PRs focused — one logical change per PR
- Add comments to your own PR before requesting review to explain the approach
  or to add questions for the reviewers
- Ensure all review threads are resolved before merging
- Use `fixup!` commits to address review comments and make additional changes
  once the PR is open, this allows the reviewers to see what has changed
  between reviews
- When addressing a review comment, link to the commit that addressed the
  comment, even if it is a `fixup!` commit that will be squashed later
- After an approval, rebase on the base branch. If this is done via the GitHub
  UI then no re-approval will be needed. However, if this is done locally then
  a re-approval is required
- Squash all `fixup!` commits and clean-up the commit history before merging
  into the base branch. You can do this at the same time as rebasing with
  `git rebase --autosquash <base_branch>`

## Reporting Issues

- Use GitHub Issues to report bugs or request features
- Search existing issues before creating a new one
- Include reproduction steps for bug reports

## Development Setup

Refer to the repository's README for specific setup instructions.

## Code of Conduct

We are committed to providing a welcoming and inclusive experience for everyone.
Please be respectful and constructive in all interactions.

## License

By contributing, you agree that your contributions will be licensed under
the same license as the project.
