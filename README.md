
![stevedore-logo](docs/logo/logo_4_stevedore.png "Stevedore logo")


## [Stevedore website](https://gostevedore.github.io/) - [Documentation](https://gostevedore.github.io/documentation/)

---

Stevedore, the docker images factory

## What is Stevedore?

Stevedore is a tool to build Docker images. Is not a Dockerfile's alternative, but to the way you use them to build your images. 
What Stevedore does is accompany you during Docker image's building process when you manage a bunch of images, and it does through an Stevedore's [driver](https://gostevedore.github.io/documentation/getting-started/concepts/#driver).

## Why stevedore?

Stevedore is a useful tool when you need to manage a bunch of Docker images in a standardized way, such on a microservices architecture. It lets you to define how to build your Docker images and their parent-child relationship. It builds automatically the children images when parent ones are done. In case your are managing your tags using semver, it is possible to generate automatically several tags following the version tree, even you can configure which tags to create.
Stevedore could also store you private registry credentials and login to it behalf you during the CI/CD processes.

## Features

#### Automate
Define docker images relationship and automate its building process
> **use case:** You have an image created with a custom setup and that image is configured as the base from many other images, you could understand it as the image `FROM`, on all your microservices. Then, every time you rebuild that image to include new security patches, all those images that depends on it are automatically built

#### Standardize
Standardize and reuse how you create your images
> **use case:** Your all microservices are based on the same skeleton, with a common way to test or built them. Then, you could define one Dockerfile and reuse it to generate the images for all your microservices

#### Semver tags
Generate automatically semver tags when main tag is semver 2.0.0 compliance
> **use case:** Supose that you are tagging your images using semantic versioning. In that case you could also generate automatically new tags based on semver tree

_example:_
_You have the tag `2.1.0`, then tags either `2` and `2.1` are also generated._

#### Credentials
Centralized Docker registry's credentials store
> **use case:** Supose that your Docker registry needs to be secured and should only be accessed by continuous delivery procedures. Then, `stevedore` could be authenticated on your behalf to push any created image to a Docker registry

#### Promote
Promote images to another Docker registry or to another registry namespace
> **use case:** Your production environment only accept to pull images from an specific Docker registry, and that docker registry is only used by your production environment. Supose that you have built and pushed an image to your staging Docker registry and you have passed all your end to end test. Then, you can promote the image from your staging Docker registry to your production one
