
v1.0.15 / 2018-03-15
====================

  * fix(#188, #185): Install SSL CA certificates. 
  * chore(Dockerfile): upgrade libvips to v8.6.3


v1.0.14 / 2018-03-05
====================

  * Add Docker Compose note to README (#174)
  * Fixes https by installing root CA certificates (#186)

v1.0.13 / 2018-03-01
====================

  * feat(version): bump to v1.0.13
  * feat(Docker): upgrade libvips to v8.6.2 (#184)
  * feat(vendor): upgrade bimg to v1.0.18
  * fix(debug): implement custom debug function
  * feat: add docker-compose.yml
  * Merge branch 'master' of https://github.com/h2non/imaginary
  * refactor(vendor): remove go-debug package from vendor
  * refactor(docs): remove codesponsor :(
  * fix testdata image links (#173)
  * Log hours in 24 hour clock (#165)
  * refactor(docs): update CLI usage and minimum requirements

v1.0.11 / 2017-11-14
====================

  * feat(bimg): update to v1.0.17
  * Merge branch 'realla-add-fit'
  * merge(add-fit): fix conflicts in server_test.go
  * refactor(image): remove else statement
  * fix(test): remove unused variable body
  * Add type=auto using client Accept header to auto negotiate type. (#162)
  * Add /fit action

v1.0.10 / 2017-10-30
====================

  * feat(bimg): upgrade to v1.0.16
  * fix(docs): remove no-form docs
  * feat(#156): support disable endpoints

v1.0.9 / 2017-10-29
===================

  * fix(#157): disable gzip compression support
  * debug(travis)
  * debug(travis)
  * debug(travis)
  * debug(travis)
  * debug(travis)
  * refactor(Dockerfile): use local source copy
  * refactor(requirements): add Go 1.6+ as minimum requirement
  * feat(vendor): support dependencies vendoring
  * refactor(Gopkg): define version
  * feat(vendor): add vendor dependencies
  * feat(travis): add Go 1.9 support
  * refactor(docs): specify POST payloads in description
  * feat(docs): add imagelayer badge
  * feat(docs): add imagelayer badge
  * feat(docs): add imagelayer badge
  * feat(docs): add imagelayer badge

v1.0.8 / 2017-10-06
===================

  * feat(docker): upgrade Go to v1.9.1
  * feat(docker): upgrade bimg to v1.0.15
  * feat(api): add smartcrop endpoint
  * feat(#101): add pipeline endpoint implementation + smart crop
  * fix(api): properly parse and use `embed`
  * feat(docs): add note about pipeline max operations
  * fix(tests): refactor Options -> Params
  * refactor(docs): update author notes
  * refactor(docs): update internal docs links
  * refactor(docs): move sponsor banner
  * feat(docs): add sponsor ad
  * refactor(license): update copyright

v1.0.7 / 2017-09-11
===================

  * feat(version): bump to v1.0.6

v1.0.6 / 2017-09-11
===================

  * feat(bimg): upgrade to v1.0.13
  * feat(version): bump to v1.0.6


v1.0.5 / 2017-09-10
===================

  * feat(params): add stripmeta params
  * feat(bimg): use bimg v1.0.12
  * feat(Docker): upgrade Go version to 1.9 in Docker image

v1.0.4 / 2017-08-21
===================

  * Mapping Blur URL params to the ImageOptions struct fields (#152)

v1.0.3 / 2017-08-20
===================

  * Merge branch 'master' of https://github.com/h2non/imaginary
  * fix(docs): CLI spec typo
  * Adding the Gaussian Blur feature plus a few minor formatting with gofmt. (#150)
  * feat(docs): update maintainer note

v1.0.2 / 2017-07-28
===================

  * fix(#146): handle proper response code range for max allowed size
  * Typos and minor language in help text (#144)
  * Update README.md (#143)
  * feat(History): add missing Docker changes
  * fix(server_test): assert content type header is present
  * fix(Docker): use proper SHA256 hash
  * feat(Docker): upgrade Go to v1.8.3 and libvips to v8.5.6
  * feat(changelog): update v1.0.1 changes
  * feat(version): bump to v1.0.1
  * feat(#140): expose Content-Length header

v1.0.1 / 2017-06-26
===================

  * feat(version): bump to v1.0.1
  * feat(#140): expose Content-Length header
  * feat(bimg): upgrade to bimg v1.0.10
  * feat(Docker): use Go v1.8.3.
  * feat(Docker): use libvips v8.5.6.

## v1.0.0 / 2017-05-27

  * Supporting smart crop (#136).
  * Deprecate Go < 1.5 support.
  * Uses `bimg` v1.0.9.
  * Uses `libvips` v8.5.5 in Docker image.

## v0.1.31 / 2017-05-18

  * feat(version): bump to 0.1.31
  * feat(Dockerfile): use libvips v8.5.5, Go v1.8.1 and bimg v1.0.8
  * Correcting the documentation, caching headers are always sent, regardless of being fetched from mount or by URL. (#133)
  * fix(docs): move toc top level sections
  * feat(docs): add new maintainer notice (thanks to @kirillDanshin)
  * feat(travis): use Go 1.8
  * refactor(docs): update support badges
  * feat(docs): add maintainers section
  * fix(.godir): add project name
  * fix(#124): fast workaround to unblock Heroku deployment until the buildpack can be updated
  * Deploy on Cloud Foundry PaaS (#122)
  * Add backers & sponsors from open collective (#119)
  * 1. remove the .godir as Heroku and Cloud Foundry remove the support. (#117)

## v0.1.30 / 2017-01-18

- fix(resizer): calculate proper crop width/height if only one axis is provided
- feat(travis): add multiple libvips testing environments
- fix(travis): use proper preinstall.sh URL
- fix(tests): integration with bimg v1.0.7
- fix(type): bimg v1.0.7 integration

## v0.1.29 / 2016-12-18

- feat(max-allowed-size): add new option max-allowed-size in bytes
- fix(max-allowed-size): HEAD response handling
- fix(usage): correct help message of 'allowed-origins'

## v0.1.28 / 01-10-2016

- feat(#95): use `libvips@8.4.1`.
- fix(#75): use `bimg@1.0.5`, which provides extract area fix.
- feat(api): supports `extend` and `embed` query params. See HTTP API params docs for more details.
- feat(#94): add placeholder image support in case of error.
- refactor(heroku): remove defaults flags in `Procfile` (user most specify them via Heroku app init settings).

## v0.1.27 / 27-09-2016

- feat(#90): adds `path-prefix` flag to bind to an url path.
- feat(core): adds support for `bimg@1.0.3` and `libvips@8.3+`.
- feat(core): adds support for GIF, SVG, TIFF and PDF formats.
- fix(controllers): fix application/octet-stream image processing issue.
- feat(docker): use Go 1.7.1 in Docker image.

## v0.1.26 / 05-09-2016

- feat: add support for authorization headers forwarding to remote image server origins.

## v0.1.25 / 26-05-2016

- fix(#79): make payload MIME type inference more versatile checking the file magic numbers signature.
- fix(#77): fix cache HTTP header expression.

## v0.1.24 / 21-04-2016

- feat(bimg): uses bimg `1.0.0`. No breaking changes introduced.

## v0.1.23 / 06-04-2016

- feat(api): support flip/flip query param arguments.

## v0.1.22 / 20-02-2016

- feat(#62): restrict remote URL source origins.
- feat(docker): container now uses Go 1.6.

## v0.1.21 / 09-02.2016

- feat(bimg): uses bimg `0.1.24`.

## v0.1.20 / 06-02-2016

- feat(bimg): uses bimg `0.1.23`.
