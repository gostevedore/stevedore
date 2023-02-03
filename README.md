
![stevedore-logo](docs/logo/logo_4_stevedore.png "Stevedore logo")


## [Stevedore website](https://gostevedore.github.io/) - [Documentation](https://gostevedore.github.io/documentation/)

---

Stevedore, the docker images' factory

- [Stevedore website - Documentation](#stevedore-website---documentation)
- [What is Stevedore?](#what-is-stevedore)
- [Why stevedore?](#why-stevedore)
- [Getting started](#getting-started)
- [Features](#features)
    - [Automate](#automate)
    - [Standardize](#standardize)
    - [Semver tags](#semver-tags)
    - [Credentials](#credentials)
    - [Promote](#promote)

## What is Stevedore?
Stevedore is a Docker images factory, a tool that allows you to build Docker images at scale. It is not an alternative to Dockerfile or Buildkit, but a way to improve your building and promote Docker images experience.


## Why stevedore?
Stevedore is a helpful tool when you need to build a bunch of Docker images, build it in a standardized way, such as on a microservices architecture. It lets you define how to build your Docker images and their parent-child relationship. It builds automatically the children's images when the parents are ready. In case you are managing your tags using semver, it is possible to generate automatically several tags following the version tree, even you can configure which tags to create.
Stevedore could also store your private registry credentials and log in to them on your behalf during the CI/CD processes.

## Getting started

- Initialize Stevedore
- Add credentials
- Define images
- Start building

## Features

#### Automate
Define the Docker images relationship and automate its building process
> **use case:** You have an image created with a custom setup and that image is configured as the base from many other images, you could understand it as the image `FROM`, on all your microservices. Then, every time you rebuild that image to include new security patches, all those images that depend on it are automatically built

#### Standardize
Standardize and reuse how you create your images
> **use case:** Your all microservices are based on the same skeleton, with a common way to test or built them. Then, you could define one Dockerfile and reuse it to generate the images for all your microservices

#### Semver tags
Generate automatically semver tags when the main tag is semver 2.0.0 compliance
> **use case:** Suppose that you are tagging your images using semantic versioning. In that case, you could also generate automatically new tags based on semver tree

_example:_
_You have the tag `2.1.0`, then tags either `2` and `2.1` are also generated._

#### Credentials
Centralized the Docker registry's credentials store
> **use case:** Suppose that your Docker registry needs to be secured and should only be accessed by continuous delivery procedures. Then, `stevedore` can authenticate on your behalf and push any created image to a Docker registry

#### Promote
Promote images to another Docker registry or to another registry namespace
> **use case:** Your production environment only accepts pulling images from a specific Docker registry, and that docker registry is only used by your production environment. Suppose that you have built and pushed an image to your staging Docker registry and you have passed all your end-to-end tests. Then, you can promote the image from your staging Docker registry to your production one
