
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
