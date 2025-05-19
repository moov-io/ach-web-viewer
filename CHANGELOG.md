## v0.12.5 (Released 2025-05-19)

BUILD

- build: update moov-io deps

## v0.12.4 (Released 2025-03-21)

IMPROVEMENTS

- filelist: allow missing file headers/control records to be displayed
- webui: input#filename should always have a white background

## v0.12.3 (Released 2025-03-20)

BUILD

- build: cleanup makefile
- fix(deps): update module github.com/moov-io/ach to v1.45.1 (#239)
- refactor: convert project over to go:embed

## v0.12.2 (Released 2024-12-16)

IMPROVEMENTS

- webui: pass ?pattern to pagination links

## v0.12.1 (Released 2024-12-11)

IMPROVEMENTS

- webui: basic exact pattern filtering when listing files

BUILD

- chore(deps): update golang docker tag to v1.23
- fix(deps): update module cloud.google.com/go/storage to v1.48.0 (#234)
- fix(deps): update module github.com/moov-io/ach to v1.44.2 (#235)
- fix(deps): update module github.com/moov-io/base to v0.51.1
- fix(deps): update module github.com/moov-io/cryptfs to v0.7.3
- fix(deps): update module github.com/stretchr/testify to v1.10.0 (#231)

## v0.11.1 (Released 2024-04-08)

IMPROVEMENTS

- filelist: use buffered channel for consistent listings

BUILD

- fix(deps): update module golang.org/x/sync to v0.7.0

## v0.11.0 (Released 2024-04-05)

IMPROVEMENTS

- filelist: support optional path for ODFI files

BUILD

- build(deps): bump github.com/cloudflare/circl from 1.3.3 to 1.3.7
- build(deps): bump github.com/go-jose/go-jose/v3 from 3.0.0 to 3.0.3
- build(deps): bump golang.org/x/net from v0.22.0 to v0.24.0

## v0.10.8 (Released 2024-04-02)

IMPROVEMENTS

- filelist: only scan parts of bucket when listing files

BUILD

- chore(deps): update golang docker tag to v1.22
- fix(deps): update module cloud.google.com/go/storage to v1.40.0
- fix(deps): update module github.com/moov-io/ach to v1.37.2
- fix(deps): update module github.com/moov-io/base to v0.48.5
- fix(deps): update module github.com/stretchr/testify to v1.9.0
- fix(deps): update module gocloud.dev to v0.37.0

## v0.10.7 (Released 2023-08-31)

IMPROVEMENTS

- webui: offer the file creation time localized in the browser

BUILD

- build: update github.com/moov-io/base to v0.46.0
- feat: show file size on file display
- fix(deps): update module github.com/moov-io/ach to v1.32.2
- fix(deps): update module gocloud.dev to v0.34.0

## v0.10.6 (Released 2023-02-22)

BUILD

- chore(deps): update golang docker tag to v1.20
- fix(deps): update module github.com/moov-io/ach to v1.28.1
- fix(deps): update module github.com/moov-io/base to v0.39.0

## v0.10.5 (Released 2023-01-13)

BUILD

- fix(deps): update module github.com/moov-io/ach to v1.28.0
- fix(deps): update module github.com/moov-io/base to v0.38.0
- fix(deps): update module gocloud.dev to v0.28.0

## v0.10.4 (Released 2022-12-08)

BUILD

- build: update github.com/moov-io/ach to v1.26.1

## v0.10.3 (Released 2022-12-08)

BUILD

- fix(deps): update module github.com/moov-io/ach to v1.26.0
- fix(deps): update module github.com/moov-io/base to v0.37.0

## v0.10.2 (Released 2022-11-30)

BUILD

- fix(deps): update module github.com/moov-io/ach to v1.25.0

## v0.10.1 (Released 2022-11-18)

IMPROVEMENTS

- fix: decrypt initial bytes on iterative attempts

## v0.10.0 (Released 2022-11-18)

IMPROVEMENTS

- filelist: allow using multiple decryption keys

## v0.9.1 (Released 2022-11-15)

IMPROVEMENTS

ach-web-viewer supports the additional masking options provided by moov-io/ach v1.23.1

BUILD

- fix(deps): update module github.com/moov-io/ach to v1.23.1
- fix(deps): update module github.com/moov-io/base to v0.36.2
- fix(deps): update module github.com/stretchr/testify to v1.8.1
- fix(deps): update module gocloud.dev to v0.27.0

## v0.8.4 (Released 2022-09-19)

BUILD

- build: require Go 1.19.1 in CI/CD

## v0.8.3 (Released 2022-09-19)

BUILD

- fix(deps): update module github.com/moov-io/ach to v1.20.0
- fix(deps): update module github.com/moov-io/base to v0.35.0
- fix(deps): update module gocloud.dev to v0.26.0

## v0.8.2 (Released 2022-05-18)

BUILD

- fix(deps): update github.com/protonmail/go-crypto digest to 902f79d
- fix(deps): update module github.com/moov-io/ach to v1.15.1
- fix(deps): update module github.com/moov-io/base to v0.29.2

## v0.8.1 (Released 2022-04-04)

BUILD

- fix(deps): update github.com/protonmail/go-crypto digest to 616f957
- fix(deps): update module github.com/moov-io/ach to v1.14.0
- fix(deps): update module gocloud.dev to v0.25.0

## v0.8.0 (Released 2022-03-21)

IMPROVEMENTS

- webui: display the file count for each date

BUILD

- build: update moov-io/base to v0.28.1
- chore(deps): update dependency golang to v1.18
- fix(deps): update github.com/protonmail/go-crypto commit hash to 70ae35b
- fix(deps): update module github.com/moov-io/ach to v1.13.1
- fix(deps): update module github.com/stretchr/testify to v1.7.1

## v0.7.0 (Released 2021-12-21)

ADDITIONS

- filelist: support "achgateway" as a source of files

BUILD

- fix(deps): update github.com/protonmail/go-crypto commit hash to a4f6767

## v0.6.0 (Released 2021-12-14)

IMPROVEMENTS

- web: include helpful links section
- web: include metadata section in response

BUG FIXES

- fix: nil checks found via linting, path traversal

BUILD

- fix(deps): update module github.com/moov-io/base to v0.27.1

## v0.5.0 (Released 2021-12-13)

IMPROVEMENTS

- filelist: support reading ListOpts from http endpoints
- web: print "Previous" and "Next" pagination links

## v0.4.0 (Released 2021-11-08)

BREAKING CHANGES

moov-io/base introduces errors when unexpected configuration attributes are found in the files parsed on startup.

BUILD

- docs: changelog updated for release v0.3.10
- fix(deps): update module github.com/moov-io/ach to v1.12.2
- fix(deps): update module github.com/moov-io/base to v0.27.0
- fix(deps): update module gocloud.dev to v0.24.0
- use ProtonMail crypto openpgp fork

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
