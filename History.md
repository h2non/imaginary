
1.2.4 / 2020-08-12
==================

  * upgrade: libvips to v8.10.0
  * fix(pipeline): add missing autorate (#326)

v1.2.3 / 2020-08-04
===================

  * feat(#315, #309): autorotate / gracefully fallback failed image type encoding

v1.2.2 / 2020-06-11
===================

  * fix(docs): define mirror as default extend param
  * refactor(params): use mirror as default extend behavior
  * refactor(params): use mirror as default extend behavior

v1.2.1 / 2020-06-08
===================

  * fix(history): release changes
  * Merge branch 'master' of https://github.com/h2non/imaginary
  * feat(version): release v1.1.1
  * [improvement/server-graceful-shutdown] Add support of graceful shutdown (#312)
  * fix

v1.2.0 / 2020-06-07
===================

  * feat(ci): use job stages
  * feat(ci): use job stages
  * feat(ci): use job stages
  * fix(ci): deploy filter
  * feat(ci): only deploy for libvips 8.9.2
  * fix(ci): strip v in semver version value
  * feat(History): add changes
  * New release, minor features, bimg upgrade and several fixes (#311)
  * watermarkImage must be in lowerCamelCase (#255)
  * Pre-release HEIF / HEIC support (#297)
  * [improvement/log-levels] Add support set log levels (#301)
  * chore(license): update year
  * chore(docs): delete not valid contributor
  * Added PlaceholderStatus option (#304)
  * Create FUNDING.yml
  * Delete README.md~
  * Add fly.io (#300)

v1.1.3 / 2020-02-28
===================

  * feat: add history changes
  * Merge branch 'master' of https://github.com/h2non/imaginary
  * refactor
  * add fluentd config example to ingest imaginary logs (#260)

v1.1.2 / 2020-02-08
===================

  * feature: implement interlace parameter (#273)
  * Implement wildcard for paths for the allowed-origins option (#290)
  * Go Modules, Code Refactoring, VIPS 8.8.1, etc. (#269)
  * feature: implement aspect ratio (#275)
  * Fix and add test for 2 buckets example (#281)
  * Fix "--allowed-origins" type (#282)
  * fix(docs): watermarkimage -> watermarkImage
  * refactor(docs): removen image layers badge
  * refactor: version set dev

v1.1.1 / 2019-07-07
===================

  * feat(version): bump patch
  * add validation to allowed-origins to include path (#265)
  * Width and height are required in thumbnail request (#262)
  * Merge pull request #258 from r-antonio/master
  * Cleaned code to check for existence of headers
  * Modified check for empty or undefined headers
  * Merge pull request #259 from h2non/release/next
  * Code Style changes
  * Removing gometalinter in favor of golangci-lint
  * Bumping libvips versions for building
  * Merge branch 'master' into release/next
  * Changed custom headers naming to forward headers
  * Fixed spacing typo and headers checking
  * Changed forwarded headers order, added tests and some fixes
  * Added custom headers forwarding support
  * Merge pull request #254 from nicolasmure/fix/readme
  * apply @Dynom patch to fix cli help
  * Update README.md
  * fix enable-url-source param description in README
  * Merge branch 'NextWithCIBase' into release/next
  * Reverting and reordering
  * ups
  * Moving back to a single file, worst case we need to maintain two again.
  * fixing a var
  * Updating travis config, adding docker build and preparing for automated image building
  * Megacheck has been removed in favor of staticcheck
  * timing the pull separately, by putting it in a before_install
  * Trying with a dev base image
  * Adding .dockerignore and consistently guarding the variables
  * Merging in changes by jbergstroem with some extra changes
  * Merge pull request #229 from Dynom/uniformBuildRefactoring
  * Improving gometalinter config
  * First travis-ci config attempt
  * Making sure vendor is not stale and that our deps are correctly configured
  * gometalinter config
  * Fixing Gopkg.toml
  * Adding a newish Dockerfile

v1.1.0 / 2019-02-21
===================

  * bumping to 1.1.0
  * Merge pull request #243 from Dynom/CheckingIfDefaultValueWasSpecified
  * Updating the documentation
  * Testing if an ImageOption false value, was in fact requested

v1.0.18 / 2019-01-28
====================

  * Bumping version to 1.0.18
  * Isolated the calculation added a test and added Rounding (#242)

v1.0.17 / 2019-01-20
====================

  * Bumping version to 1.0.17
  * Merge pull request #239 from Dynom/RefactoringParameterParsingToImproveFeedback
  * cleanup
  * Bumping Go's version requirement
  * Allow Go 1.9 to fail
  * Refactoring, making things simpler
  * Merge pull request #230 from Dynom/NonFunctionalImprovements
  * minor styling
  * Simplifying some if statements, removing unnecessary parenthesis and added an explicit error ignore
  * Correct casing for multiple words
  * explicitly ignoring errors
  * ErrorReply's return value was never used.
  * Changing the if's to a single map lookup.
  * More literal to constant replacements
  * Correcting comments, casing and constant use
  * Removed unused variable and explicitly ignoring errors
  * Correcting the use of abbreviations
  * Comment fix
  * Removing literals in favor of http constants
  * Unavailable is not actually used, replaced it with _ to better convey intent.
  * Removing unused function
  * Style fixes
  * Merge pull request #227 from Dynom/moarHealthEndpointDetails
  * Merge pull request #228 from Dynom/addingExtraOriginTest
  * Exposing several extra details
  * Including a test case with multiple subdomains
  * docs(watermarkimage): Add docs for the /watermarkimage endpoint (#226)
  * refactor(README): remove header image

v1.0.16 / 2018-12-11
====================

  * Adding LIBVIPS 8.7 and updating libvips URL (#225)

v1.0.15 / 2018-12-10
====================

  * Updating bimg, libvips and Go
  * Updated dockerfile vips repo (#222)
  * fix: correct fit operation with autorotated images by switching width/height in certain cases (#208)
  * Watermark image api (#221)
  * Adding remote url wildcard support (#219)
  * Bump Go versions and use '.x' to always get latest patch versions (#220)
  * Update README.md (#207)
  * Fix typo in documentation (#202)
  * Drop salt as suggested in #194 (#200)
  * Add URL signature feature (#194)
  * fix(docker): remove race detector (#197)
  * feat(version): bump to v1.0.15
  * Changing build steps (#189)

v1.0.14 / 2018-03-05
====================

  * feat(version): bump to v1.0.14
  * Add Docker Compose note to README (#174)
  * Fixes https by installing root CA certificates (#186)

v1.0.13 / 2018-03-01
====================

  * fix(Dockerfile): update version sha2 hash
  * feat(history): update changelog
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

  * fix(type_test): use string for proper formatting
  * feat(version): bump to v1.0.11
  * feat(bimg): update to v1.0.17
  * Merge branch 'realla-add-fit'
  * merge(add-fit): fix conflicts in server_test.go
  * refactor(image): remove else statement
  * fix(test): remove unused variable body
  * Add type=auto using client Accept header to auto negotiate type. (#162)
  * Add /fit action

v1.0.10 / 2017-10-30
====================

  * feat(docs): update CLI usage help
  * feat(#156): support disable endpoints (#160)

v1.0.9 / 2017-10-29
===================

  * feat(version): bump to v1.0.9
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

  * feat(#101): add pipeline endpoint implementation + smart crop (#154)
  * refactor(docs): move sponsor banner
  * feat(docs): add sponsor ad
  * refactor(license): update copyright

v1.0.7 / 2017-09-11
===================

  * feat(version): bump to v1.0.7
  * feat(version): bump to v1.0.6

v1.0.5 / 2017-09-10
===================

  * feat(version): bump to v1.0.5
  * feat(History): update version changes
  * feat(params): add stripmeta params

v1.0.4 / 2017-08-21
===================

  * feat(version): bump to 1.0.4
  * Mapping Blur URL params to the ImageOptions struct fields (#152)

v1.0.3 / 2017-08-20
===================

  * feat(version): bump to v1.0.3
  * Merge branch 'master' of https://github.com/h2non/imaginary
  * fix(docs): CLI spec typo
  * Adding the Gaussian Blur feature plus a few minor formatting with gofmt. (#150)
  * feat(docs): update maintainer note

v1.0.2 / 2017-07-28
===================

  * feat(version): bump to v1.0.2
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

v1.0.0 / 2017-05-27
===================

  * refactor(controller): add height for smart crop form
  * feat(controllers): add smart crop form
  * feat(version): bump to v1.0.0
  * feat(History): update changes
  * Supporting smart crop (#136)

v0.1.31 / 2017-05-18
====================

  * feat(History): update latest changes
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

v0.1.30 / 2017-01-18
====================

  * refacgor(version): add comments
  * feat(version): bump to v0.1.30
  * feat(History): update changes
  * fix(travis): remove libvips 8.5
  * feat(travis): add multi libvips testing environments
  * fix(travis): use proper preinstall.sh URL
  * Update .travis.yml
  * fix(tests): integration with bimg v1.0.7
  * fix(type): bimg v1.0.7 integration
  * fix(type): bimg v1.0.7 integration
  * Update History.md

v0.1.29 / 2016-12-18
====================

  * feat(version): bump to 0.1.29
  * feat(max-allowed-size): add new option max-allowed-size in bytes (#111)
  * Merge pull request #112 from touhonoob/fix-help-allowed-origins
  * fix(usage): correct help message of 'allowed-origins'
  * refactor(docs): remove deprecated sharp benchmark results
  * refactor(docs): update preinstall.sh install URL
  * refactor(docs): use preinstall.sh script from bimg repository
  * fix(docs): Docker image link
  * fix(history)
  * fix(history)

v0.1.28 / 2016-10-02
====================

  * feat(docs): add placeholder docs and several refactors
  * feat(docs): add placeholder docs and several refactors
  * feat(#94): support placeholder image
  * feat(version): bump to v0.1.28
  * feat(version): release v0.1.28
  * feat(core): support bimg@1.0.5, support extend background param
  * chore(history): add Docker Go 1.7.1 support
  * feat(docker): use Go 1.7.1

v0.1.27 / 2016-09-28
====================

  * fix(server): mount route
  * refactor(server): DRYer path prefix
  * Merge pull request #93 from h2non/develop
  * fix(tests): type tests based on libvips runtime support
  * feat(version): bump
  * feat(travis): add Go 1.7
  * feat(docs): add new formats support
  * fix(history): update to bimg@1.0.3
  * fix(controllers): fix binary image processing
  * refactor(controllers)
  * Merge branch 'develop' of github.com:h2non/imaginary into develop
  * feat(core): add additional image formats
  * feat(core): add support for bimg@1.0.2 and new image formats
  * Merge pull request #90 from iosphere/feature/path-prefix
  * Add `path-prefix` flag to bind to an url path
  * Merge pull request #89 from h2non/develop
  * refactor(cli): update flag description
  * feat(docs): improve CLI docs
  * refactor(cli): improve description for -authorization flag

v0.1.26 / 2016-09-06
====================

  * fix(merge): master
  * chore(history): update history changelog
  * feat(docs): update CLI usage and help
  * feat: forward authorization headers support
  * Fix description for URL source, and allowed origins server options (#83)
  * fix(version): ups, editing from iPad
  * fix(version): unresolved conflict
  * merge: fix History conflicts
  * Merge branch 'develop'
  * Fix Expires and Cache-Control headers to be valid (#77)
  * Update README.md (#74)
  * Promote version 0.1.24 (#73)

0.1.25 / 2016-05-27
===================

  * Sync develop (#82)
  * feat(version): bump
  * fix(#79): infer buffer type via magic numbers signature
  * Sync develop (#80)

v0.1.24 / 2016-04-21
====================

  * feat(version): bump
  * Merge branch 'develop' of github.com:h2non/imaginary into develop
  * Sync develop (#72)
  * merge(upstream)
  * refactor(travis)
  * feat(bimg): bump version to v1
  * Sync develop (#71)
  * add background param (#69)
  * Merge pull request #70 from h2non/develop
  * fix(docs): typo
  * fix(docs): minor typos
  * Merge pull request #64 from h2non/develop
  * Merge pull request #63 from h2non/develop

0.1.23 / 2016-04-06
===================

  * feat(docs): add flip flop params
  * feat(version): bump
  * feat(#66): flip/flop support as param
  * feat(timeout): increase read/write timeout to 60 seconds

0.1.22 / 2016-02-20
===================

  * feat(docker): use SHA256 checksum
  * feat: update history
  * feat(version): bump
  * feat(docs)
  * feat(#62): support allowed origins
  * feat(#62): support allowed origins
  * Merge pull request #61 from h2non/master
  * feat(travis): use go 1.6
  * Merge pull request #60 from h2non/develop
  * feat(history): add change log
  * Merge pull request #59 from h2non/develop
  * Merge pull request #58 from h2non/develop

0.1.21 / 2016-02-09
===================

  * feat(version): bump

0.1.20 / 2016-02-06
===================

  * feat(version): bump
  * feat(docs): add PKGCONFIg variable
  * merge(master)
  * Merge pull request #57 from h2non/develop
  * Merge pull request #56 from h2non/develop
  * Merge pull request #55 from h2non/develop
  * Merge pull request #54 from pra85/patch-1
  * Typo fixes
  * fix(docs): typo in scalability
  * Merge pull request #53 from h2non/develop
  * feat(docs): add imaginary badge
  * refactor(docs): improve scalability notes
  * feat(docs): add docker pulls badge
  * feat(docs): add form data spec

0.1.19 / 2016-01-30
===================

  * refactor(form): use previous params
  * feat(docs): add rotate param in endpoints
  * feat(version): bump
  * feat(#49): support custom form field
  * fix(docs): minor typo
  * feat: add more tests, partially document code
  * refactor(controllers): use external struct
  * refactor: follow go idioms
  * refactor(middleware): rename function
  * refactor(middleware): only cache certain requests
  * fix(docs): use proper flag
  * fix(docs): add supported method
  * fix(cli): bad flag description
  * fix
  * refactor(health)
  * refactor
  * feat: ignore imaginary root binary
  * refactor(middleware)
  * feat(docs): add examples
  * feat(docs): update CLI help

0.1.18 / 2015-11-04
===================

  * feat(version): bump
  * fix(badge)
  * fix(badge)
  * merge(upstream)
  * feat(docs): add remote URL support, update badges
  * refactor(cli): change flag
  * feat(#43, #35): support gravity param and health
  * feat(#32): add test coverage
  * feat(#32): initial support for URL processing
  * fix(tests)
  * feat(#32): support flags
  * feat(#32): initial seed implementation
  * Merge pull request #44 from freeformz/master
  * Add Heroku Button Support
  * fix(docs): content typo
  * feat: add glide.yaml for vendording packages
  * feat: add glide.yaml for vendording packages
  * refactor(docs): add performance note
  * refactor(docs)
  * refactor(benchmark): uncomment kill sentence
  * feat(docs): add benchmark notes
  * refactor(image): add default error on panic
  * feat: add panic handler. feat(docs): add error docs

0.1.17 / 2015-10-31
===================

  * feat(version): bump
  * Merge pull request #39 from Dynom/addingHttpCaching
  * Added documentation.
  * More style fixes.
  * Removing redundant construct
  * Fixing coding-style
  * Merge pull request #41 from Dynom/enablingSecureDownloadOfGo
  * Added the CA certs so that the --insecure flag can be removed from the GO installer.
  * Added a sanity check for the value of the -http-cache-ttl flag.
  * Added -http-cache-ttl flag
  * feat(log): add comments
  * refactor(body)
  * refactor(benchmark)
  * refactor(benchmark)
  * refactor: rename function
  * refactor: normalize statements, add minor docs
  * refactor(docs): add link
  * feat(docs): add toc

0.1.16 / 2015-10-06
===================

  * fix(docker): restore to default
  * refactor(docker): uses latest version
  * feat(version): bump
  * fix(#31): use libvips 7.42 docker tag
  * refactor(docs): update descriptiong
  * feat(docs): add libvips version compatibility note
  * merge(upstream)
  * refactor(docs): add root endpoint, fix minor typos
  * refactor(docs): description
  * feat(docs): add sourcegraph badge
  * refactor(docs): minor changes, reorder

0.1.15 / 2015-09-29
===================

  * feat(version): bump
  * merge: upstream
  * feat: expose libvips and bimg version in index route
  * refactor(docs): add docker debug command

0.1.14 / 2015-08-30
===================

  * fix: build
  * refactor(docker): bump Go version
  * feat(version): bump
  * feat: use throttle v2
  * refactor(make): push specific tag

0.1.13 / 2015-08-10
===================

  * feat(version): bump
  * feat(#30)

0.1.12 / 2015-07-29
===================

  * feat(version): bump
  * fix(dependency)
  * refactor: add errors as constants. middleware
  * refactor: add errors as constants. middleware
  *  fix(docs): typo
  * fix(travis): remove go tip build due to install.sh error
  * refactor: server router
  * refactor:
  * fix(docs): add missing params per specific method
  * feat(docs): add image
  * feat(docs): add link

0.1.11 / 2015-07-11
===================

  * feat(version): bump
  * feat(#26): add TLS support
  * feat(#27)
  * feat(#27)
  * refactor(docs)
  * fix(docs): description
  * fix
  * feat: merge
  * refactor(form): dry
  * refactor(docs): http api
  * refactor(main)

0.1.10 / 2015-06-30
===================

  * feat(version): bump
  * refactor(docs)
  * feat(#25): several refactors and test coverage
  * feat(#25): experimental support for local files processing
  * feat: support no profile param
  * feat(http): add bimg version header
  * feat(http): add bimg header
  * refactor(docs): node graph

0.1.9 / 2015-06-12
==================

  * refactor: disable interlace by default (due to performance issues)
  * feat(version): bump
  * feat: add interlace support by default
  * refactor(image): remove debug statement
  * refactor(docs): description
  * fix(form): add proper param for watermark
  * fix(form): add proper param for watermark
  * refactor(docs): description
  * refactor(params): use math function

0.1.8 / 2015-05-24
==================

  * feat(version): bump
  * feat(version): bump
  * fix(form): bad param
  * refactor(docs): scalability
  * refactor(docs): benchmark
  * refactor(docs): description
  * refactor(docs): usage
  * refactor(docs): update sections
  * refactor(bench). feat(docs): add resources and scalability notes
  * refactor(docs)
  * refact(bench)
  * feat(docs): add production note
  * merge
  * refactor(server): isolate throttle to middleware
  * fix(docs): duplicated param

0.1.7 / 2015-04-27
==================

  * fix(extract): bad query param
  * feat(version): bump
  * fix(enlarge): bad params assignment
  * feat(#24): crop by default

0.1.6 / 2015-04-26
==================

  * feat(version): bump (maintenance release)

0.1.5 / 2015-04-25
==================

  * feat(version): bump
  * feat(params): add params for no auto rotate
  * feat(params): add params for no auto rotate
  * refactor(docs): description
  * fix(docs): description
  * refactor(docs): description
  * refactor(docs): add new Heroku steps
  * refactor(buildpack)
  * refactor(buildpack)
  * refactor(buildpack)
  * refactor(buildpack)
  * fix(heroku)
  * refactor: update buildpack

0.1.4 / 2015-04-19
==================

  * feat(version): bump
  * feat: handle HTTP 404
  * feat(heroku): update docs
  * feat: update buildpack
  * feat
  * refactor(heroku)
  * fix(heroku): buildpack order
  * feat(#23)

0.1.3 / 2015-04-19
==================

  * feat(version): bumo
  * fix(port)
  * refactor(docker): remove help flag
  * refactor: heroku
  * refactor
  * refactor
  * refactor
  * refactor
  * feat: add dependencies
  * feat: add dependencies
  * feat: add heroku files
  * feat: add Heroku files
  * refactor(docs): update description
  * refactor(image)
  * fix(docs)

0.1.2 / 2015-04-18
==================

  * refactor(bench)
  * feat(docs): better Heroku docs
  * refactor: split parse params and body read
  * feat(docs): add server clients
  * refactor(docs)
  * fix(form): add query param
  * fix(docs): usage
  * fix(cli): memory release
  * fix(travis)

0.1.1 / 2015-04-15
==================

  * feat(version): bump
  * feat(#20) fix(#2)
  * feat: refactor
  * refactor(bench)
  * refactor(docs)
  * feat(docs): add benchmarks
  * feat(docs): add benchmarks
  * feat(docs): add benchmarks
  * feat(#16): add benchmark
  * refactor(docs)
  * refactor(docs)

0.1.0 / 2015-04-13
==================

  * feat(#18): http docs
  * fix(travis): another attempt
  * fix(travis)
  * fix(docs)
  * fix(travis)
  * fix(travis)
  * refactor(docs): docker
  * refactor(cli): priorize CLI flag
  * refactor(server)
  * fix(travis): pending issue from Coveralls
  * feat: add Makefile
  * fix(package): name
  * feat(docs): add docker badge

0.1.0-rc.0 / 2015-04-12
=======================

  * feat(test): add test coverage
  * refactor
  * feat(#15): docker file settings
  * feat(#19, #10, #13)
  * feat(docs): add image
  * feat(docs): add image
  * feat(docs): add image
  * feat(image): add image
  * feat(#15, #11, #7, #5, #13)
  * feat(#17, #7, #2)
  * refactor: rename
  * refactor: remove file
  * refactor: rename
  * refactor: use bimg
  * refactor(image): options
  * feat: add upload test
  * refactor
  * refactor
  * feat: add Dockerfile
  * refactor: server
  * feat: add test
  * feat: add test
  * feat: add test
  * feat: add sources
  * feat: add files

