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
