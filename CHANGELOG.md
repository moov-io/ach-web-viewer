## v0.3.10 (Released 2021-10-26)

IMPROVEMENTS

- gpg: use ProtonMail crypto openpgp fork in place of deprecated golang.org/x/crypto

## v0.3.9 (Released 2021-09-14)

IMPROVEMENTS

- webui: Better styling for long filepaths

## v0.3.8 (Released 2021-09-14)

IMPROVEMENTS

- filelist: show as much of partial files as possible

## v0.3.7 (Released 2021-09-14)

BUG FIXES

- web: guard for nil with ACH file validation

## v0.3.6 (Released 2021-09-13)

BUG FIXES

- filelist: skip missing directories

## v0.3.5 (Released 2021-08-18)

IMPROVEMENTS

- web: display a bit of the parent location for a file

BUILD

- fix(deps): update module github.com/moov-io/ach to v1.11.0
- fix(deps): update module github.com/moov-io/base to v0.23.0

## v0.3.4 (Released 2021-07-26)

BUILD

- fix(deps): update module github.com/moov-io/ach to v1.10.1

## v0.3.3 (Released 2021-07-16)

IMPROVEMENTS

- docs: include full config and screenshots of new design

BUILD

- fix(deps): update golang.org/x/crypto commit hash to a769d52
- fix: Dockerfile to reduce vulnerabilities

## v0.3.2 (Released 2021-07-07)

IMPROVEMENTS

- web: fix ordering of date headers and display

## v0.3.1 (Released 2021-07-07)

IMPROVEMENTS

- web: improve display of date groups

## v0.3.0 (Released 2021-07-07)

IMPROVEMENTS

- web: improve style with fancy stylesheet
- web: pass through optional validation error when viewing file

## v0.2.1 (Released 2021-06-28)

This release contains MacOS and Windows binaries.

## v0.2.0 (Released 2021-06-28)

ADDITIONS

- web: include "Back" link when viewing a specific file

BUILD

- fix(deps): update golang.org/x/crypto commit hash to 5ff15b2
- fix(deps): update module github.com/moov-io/ach to v1.9.3

## v0.1.5 (Released 2021-06-09)

IMPROVEMENTS

- filelist: skip checking filepath stat info on service startup

## v0.1.4 (Released 2021-05-21)

IMPROVEMENTS

- service: setup HTTP server to run off sub-path
- web: include basePath in hyperlinks

## v0.1.0 (Released 2021-05-18)

Initial Release

- Basic website [listing all Source files](http://localhost:8585/)
   - Supporting filesystem reads and S3-compatiable (e.g. GCS) blob storage
