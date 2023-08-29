# 1.8.0 / 2023-08-29

  * Update zh translations (thanks to @YongJie-Xie) (#113)

# 1.7.1 / 2023-08-25

  * Fix: Encode data for download
  * [ci] Add local build target

# 1.7.0 / 2023-08-13

  * [#110] Add interaction buttons for displayed secret (#111)

# 1.6.1 / 2023-08-11

  * Fix: Adjust HTML page title to customized AppTitle (#107)
  * Fix dutch translation for minute (#108)

# 1.6.0 / 2023-08-04

  * Add Polish translation (thanks to @Icikowski) (#106)

# 1.5.1 / 2023-07-07

  * Add missing Catalan translations (thanks to @v0ctor) (#102)
  * Add missing Spanish translations (thanks to @v0ctor) (#103)

# 1.5.0 / 2023-07-06

  * New Features
    * [#97] Add framework for formal language & formal German translation (#98)
    * Add Ukrainian language (thanks to @t0rik) (#99)

# 1.4.0 / 2023-06-27

  * New Features
    * [#85] Allow to customize secret expiry (#93)

# 1.3.0 / 2023-06-24

  * New Features
    * [#91] Add Copy-to-Clipboard button to secret URL
    * [#92] Add detection for write-disabled instances
    * Add Turkish language (thanks to @vehbiyilmaz)

  * Improvements
    * Implement proper tool to manage translations
    * Improve README readability
    * Mitigate possible XSS through `unsafe-inline` script CSP

# 1.2.0 / 2023-06-14

  * Improvements
    * Log API errors in server log

  * Bugfixes
    * [#89] Fix error handling of `fetch` API

# 1.1.0 / 2023-06-14

  * New Features
    * Add QR-code display for secret URL (#61)
    * Implement frontend customizations (#71)

  * Improvements
    * Disable secret creation when secret is empty (#86)
    * Log secret expiry on startup
    * Only mention tool name in footer (#71)
    * Replace redis client, move expiry into creation interface

With this release an old migration was removed and in case you are still using the `REDIS_EXPIRY` environment variable you need to switch to `SECRET_EXPIRY`. Also with the new redis client you might need to adjust the username in your `REDIS_URL` to a proper ACL username (or enable legacy auth in Redis) - see the README for the `REDIS_URL` format.

# 1.0.0 / 2023-04-14

  * Breaking: Replace deprecated / archived crypto library (#80)

# 0.27.0 / 2023-04-10

  * Add pt-BR translation (Thanks to [@imfvieira](https://github.com/imfvieira)!)

# 0.26.0 / 2023-03-29

  * Add Swedish language (Thanks to [@artingu](https://github.com/artingu)!)

# 0.25.0 / 2023-03-17

  * Add Russian language (#79) (Thanks to [@alexovchinnicov](https://github.com/alexovchinnicov)!)

# 0.24.1 / 2023-03-07

  * Update dependencies / fix vulnerabilities
  * CI: Fix release publishing

# 0.24.0 / 2022-11-24

  * Add Traditional Chinese translations (#68) (Thanks to [@DejavuMoe](https://github.com/DejavuMoe)!)
  * Fix: Use full browser provided language tag

# 0.23.0 / 2022-11-21

  * Add Simplified Chinese translations (#67) (Thanks to [@DejavuMoe](https://github.com/DejavuMoe)!)
  * Replace password generation with web-crypto API
  * [typo] comprimise -> compromise (#63)

# 0.22.0 / 2022-04-10

  * Upgrade golang dependencies
  * Upgrade node dependencies
  * Add Catalan translation (#50)
  * Add Spanish translation (#49)
  * Add OpenAPI documentation (#48)
  * Add security HTTP headers (#45)
  * [#46] Remove external font deps, add SRI checks (#47)

# 0.21.0 / 2021-09-18

  * Add Theme-Switcher for Dark-/Light-Mode

# 0.20.1 / 2021-09-08

  * [#44] Fix missing libraries within compiles binary / container

# 0.20.0 / 2021-09-07

  * Switch to structs instead of maps in api (#40)
  * [#35] Encode pipe in secret URL by default
  * Update dependencies, upgrade build utils
  * Add dutch translation (#39)
  * Switch to Go 1.16 embed functionality (#42)
  * Remove duplicate call LastIndex (#41)

Many thanks to [@sorcix](https://github.com/sorcix) for the contributions to this release!

# 0.19.0 / 2021-08-09

  * Change Cache-Control on responses to no-store (#37)

# 0.18.1 / 2020-10-20

  * Fix: Update node dependencies

# 0.18.0 / 2020-08-10

  * Remove unused assets, bundle Latvian translation
  * Latvian translation (#25) (Thanks to [@Stegadons](https://github.com/Stegadons)!)
  * npm update / audit fix / bundle update

# 0.17.3 / 2020-06-02

  * Update node dependencies and rebuild packed frontend

# 0.17.2 / 2020-06-02

  * Add example script to get secret from CLI (#18)
  * npm audit fix & generate bundled js
  * [#14] Document creation of secrets through CLI / script

# 0.17.1 / 2020-01-26

  * [#13] Fix: Secrets in MEM store were instantly expired

# 0.17.0 / 2020-01-24

  * [#12] Add lazy-expiry to mem-storage, unify envvars
  * Update Dockerfiles
  * Switch to Go 1.11+ modules
  * Fix NPM audit alerts
  * Add minimal Dockerfile without alpine base

# 0.16.1 / 2019-07-20

  * Fix: Update assets to include FR translation

# 0.16.0 / 2019-07-20

  * Add french language translation (#10) (Thanks to [@ometra](https://github.com/ometra)!)

# 0.15.0 / 2019-07-14

  * Add explanation of functionality (#9)

# 0.14.0 / 2019-07-14

  * UX: Auto-select secret URL after creation
  * Implement vue/recommended linter
  * Fix eslint missing dependencies
  * Update translation hint in README

# 0.13.4 / 2019-07-14

  * Update fontawesome and vue-i18n
  * Bump lodash from 4.17.11 to 4.17.14 in /src (#7)
  * Fix node package vulnerabilities

# 0.13.3 / 2019-05-13

  * CI: Update build image

# 0.13.2 / 2019-05-13

  * Fix: Encoded hashes were not properly processed (again)
  * Fix eslinter errors

# 0.13.1 / 2019-05-10

  * Fix: Broken version display

# 0.13.0 / 2019-05-10

  * Fix: Cleanup debugging stuff
  * Move frontend to Vue
  * Move translations to frontend
  * Handle json requests to create API
  * Update frontend dependencies
  * Add gzip compression for included assets
  * Update dev-dependencies

# 0.12.0 / 2018-10-22

  * Be more specific about security risks

# 0.11.1 / 2018-10-06

  * Replace uuid library, update vendors

# 0.11.0 / 2018-10-06

  * Port frontend code to pure Javascript

# 0.10.0 / 2018-08-22

  * Auto-resize textareas, use babel to transpile JS
  * Fix: Transmit secret using POST method

# 0.9.0 / 2018-05-05

  * Generate SRI integrity hashes into html

# 0.8.1 / 2018-05-05

  * Update Dockerfile to multi-stage build

# 0.8.0 / 2018-05-05

  * Feat: Internalize previously external libraries and stylesheets
  * Feat: Migrate to bootstrap 4 and fontawesome 5
  * Fix: Fix date and maintainer in LICENSE file
  * Vendor: Switch to dep for vendoring
  * Vendor: Update dependencies

# 0.7.0 / 2018-05-05

  * Introduce data expiry in Redis

# 0.6.0 / 2017-08-19

  * Add view to confirm display and destroy of the secret
  * Add translation information
  * Add version to footer

# 0.5.1 / 2017-08-04

  * Fix: Vendor missing libraries

# 0.5.0 / 2017-08-04

  * Add localization for en-US and de-DE

# 0.4.0 / 2017-08-04

  * Remove option to disable encryption

# 0.3.1 / 2017-08-03

  * Fix: Some messengers mess up the URL

# 0.3.0 / 2017-08-03

  * Add footer

# 0.2.0 / 2017-08-03

  * Follow linter advices

# 0.1.0 / 2017-08-03

  * Initial Version
