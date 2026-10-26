# drone-cli

Command line client for the Drone continuous integration server.

Documentation: https://docs.drone.io/cli

## Build requirements

Go 1.26.10 or newer is required to build the CLI.

The CI Go image and the amd64/arm64 Linux runtime images come from the authenticated HHI registry at `harness0.harness.io`. Local Docker builds require an authorized `docker login harness0.harness.io`; CI runners must have equivalent image-pull credentials configured before starting the pipeline. The ARMv7 runtime uses public Alpine, but its CI Go image still requires HHI access.

Each Linux Dockerfile defaults to its named architecture. Building its package-install stages requires a compatible worker or container emulation.

## Release procedure

Run the changelog generator.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u harness -p drone-cli -t <secret github token>
```

You can generate a token by logging into your GitHub account and going to Settings -> Personal access tokens.

Next we tag the PR's with the fixes or enhancements labels. If the PR does not fufil the requirements, do not add a label.

Run the changelog generator again with the future version according to semver.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u harness -p drone-cli -t <secret token> --future-release v1.0.0
```

Create your pull request for the release. Get it merged then tag the release.

## Community and Support
[Harness Community Slack](https://join-community-slack.harness.io/) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report or Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 
