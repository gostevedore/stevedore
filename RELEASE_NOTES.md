# RELEASE NOTES

## [undefined]

## Fixed

- Fixed a bug in Console's ReadPassword function that always returned an error in the defer function. Now it is using the package golang.org/x/term v0.16.0 to read the password from the console.

## Added

- Included a message to notify the user when a credential is successfully created
- Included test to read password. It being used the package github.com/kr/pty v1.1.8 to create a pty and read the password from it.

## Bumped

- Bump github.com/apenella/go-docker-builder from v0.8.1 to v0.8.3.
- Bump github.com/apenella/go-ansible v1.2.2
- Bump github.com/aws/aws-sdk-go-v2 v1.24.1
- Bump github.com/aws/aws-sdk-go-v2/config v1.26.3
- Bump github.com/aws/aws-sdk-go-v2/credentials v1.16.14
- Bump github.com/aws/aws-sdk-go-v2/service/ecr v1.24.7
- Bump github.com/aws/aws-sdk-go-v2/service/sts v1.26.7
- Bump github.com/fatih/color v1.16.0
- Bump github.com/gruntwork-io/terratest v0.46.9
- Bump github.com/spf13/afero v1.11.0
- Bump github.com/spf13/viper v1.18.2
- Bump golang.org/x/term v0.16.0
