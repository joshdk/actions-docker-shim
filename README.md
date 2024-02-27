[![License][license-badge]][license-link]
[![Actions][github-actions-badge]][github-actions-link]
[![Releases][github-release-badge]][github-release-link]

# GitHub Actions Docker Shim

🐋 Shim that enables using private ghcr.io images in GitHub Actions.

## Motivations

Currently, there isn't a good story for authoring GitHub Actions backed by private Docker images.
Unlike workflow [services](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idservicesservice_idcredentials) and [jobs](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idcontainercredentials), there is no way to transparently provide registry credentials to image backed Actions.
Some potential workarounds include manually authenticating inside a composite action, deploying your own Action runners, or making your image public.

This repository produces a tool which aims to solve this problem as seamlessly as possible.

## Usage

> [!NOTE]
> These examples assumes that you author an action from a repository named `Example/example-action`, which produces an image named `ghcr.io/example/example-action`.

### Initial Integration

Modify your `action.yml` to use the latest release of `ghcr.io/joshdk/actions-docker-shim` and to add a `token` input with a default value of `${{ github.token }}`.

```diff
name: Example
description: An example action, wow!

+inputs:
+  token:
+    description: GitHub Actions workflow token.
+    default: ${{ github.token }}

runs:
  using: docker
-  image: docker://ghcr.io/example/example-action:v1.2.3
+  image: docker://ghcr.io/joshdk/actions-docker-shim:v0.1.0
```

Modify the caller workflow to grant the `packages: read` permission so that your ghcr.io image can be pulled.

```diff
jobs:
  example:
+    permissions:
+      packages: read

    steps:
      - uses: Example/example-action@v1.2.3
```

#### How this works

When GitHub runs a workflow which references your action, the `ghcr.io/joshdk/actions-docker-shim` image is pulled and run instead of your image.

The shim determines the name of the GitHub repository which hosts this action (`Example/example-action`) and the ref at which this action was references (`v1.2.3`).
Using this information, the name of this image is speculated to be `ghcr.io/example/example-action:v1.2.3`.

The shim performs a login to ghcr.io using the provided `token` input, and then pulls the target image.

Finally, the target image is run on the GitHub runner using the same set of volume mappings, environment variables, & arguments that GitHub **would** have used to run your image.

### Specifying Images

By default, the git tag/branch/sha that is used when referring to your action will be used to determine the image tag to run.
You can see here how the image name is derived.

```yaml
jobs:
  example:
    steps:
      # Will run ghcr.io/example/example-action:v1.2.3
      - uses: Example/example-action@v1.2.3

      # Will run ghcr.io/example/example-action:v1.2
      - uses: Example/example-action@v1.2

      # Will run ghcr.io/example/example-action:v1
      - uses: Example/example-action@v1

      # Will run ghcr.io/example/example-action:master
      - uses: Example/example-action@master
      
      # Will run ghcr.io/example/example-action:a35f...b316
      - uses: Example/example-action@a35f...b316
```

If you want your action to use a specific image tag, then you can set one manually in `action.yml`.

```diff
runs:
  using: docker
  image: docker://ghcr.io/joshdk/actions-docker-shim:v0.1.0
+  args:
+    - --shim-image-tag=snapshot
```

If your image isn't named the same as your action repository, that can be overridden as well.

```diff
runs:
  using: docker
  image: docker://ghcr.io/joshdk/actions-docker-shim:v0.1.0
  args:
    - --shim-image-tag=snapshot
+    - --shim-image=example/some-other-image
```

### Authentication

If you need to provide a token using a custom input name (e.g. to avoid changing the interface for your action) then you can specify the env var name to use.  

```diff
name: Example
description: An example action, wow!

inputs:
-  token:
+  custom-token:
    description: GitHub Actions workflow token.
    default: ${{ github.token }}

runs:
  using: docker
  image: docker://ghcr.io/joshdk/actions-docker-shim:v0.1.0
+  args:
+    - --shim-token-env=INPUT_CUSTOM-TOKEN
```

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[github-actions-badge]:  https://github.com/joshdk/actions-docker-shim/workflows/Build/badge.svg
[github-actions-link]:   https://github.com/joshdk/actions-docker-shim/actions
[github-release-badge]:  https://img.shields.io/github/release/joshdk/actions-docker-shim/all.svg
[github-release-link]:   https://github.com/joshdk/actions-docker-shim/releases
[license-badge]:         https://img.shields.io/badge/license-MIT-green.svg
[license-file]:          https://github.com/joshdk/actions-docker-shim/blob/master/LICENSE.txt
[license-link]:          https://opensource.org/licenses/MIT
