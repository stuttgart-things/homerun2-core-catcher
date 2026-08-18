## [0.13.4](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.13.3...v0.13.4) (2026-08-18)


### Bug Fixes

* **deps:** update module charm.land/bubbletea/v2 to v2.0.8 ([3130a5e](https://github.com/stuttgart-things/homerun2-core-catcher/commit/3130a5efd864a2a7593562c4a7e6ee34061f6271))
* **deps:** update module github.com/redis/go-redis/v9 to v9.22.0 ([56b0c09](https://github.com/stuttgart-things/homerun2-core-catcher/commit/56b0c096821e8b645e444c4da986f1449bc7eb83))

## [0.13.3](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.13.2...v0.13.3) (2026-08-18)


### Bug Fixes

* **deps:** update module charm.land/lipgloss/v2 to v2.0.6 ([fd5e6c9](https://github.com/stuttgart-things/homerun2-core-catcher/commit/fd5e6c963373e5c3ac1eca103730d56c71a15364))

## [0.13.2](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.13.1...v0.13.2) (2026-05-28)


### Bug Fixes

* **kcl:** stop emitting Namespace from kustomize OCI ([0ae946f](https://github.com/stuttgart-things/homerun2-core-catcher/commit/0ae946f4284ae3a6e32572093812b501b556bcb3)), closes [#70](https://github.com/stuttgart-things/homerun2-core-catcher/issues/70)

## [0.13.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.13.0...v0.13.1) (2026-05-27)


### Bug Fixes

* **kcl:** default catcherMode to "web" in schema ([7a32143](https://github.com/stuttgart-things/homerun2-core-catcher/commit/7a3214364dc8120aa49851b8d29dff2edd5b6fd6))

# [0.13.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.12.2...v0.13.0) (2026-05-16)


### Features

* ship HTTPRoute in the kustomize OCI base (Option B for argocd[#116](https://github.com/stuttgart-things/homerun2-core-catcher/issues/116)) ([32a3bcd](https://github.com/stuttgart-things/homerun2-core-catcher/commit/32a3bcd79a49ae3d9a0844b7e6e4f7b53966f4c1)), closes [#59](https://github.com/stuttgart-things/homerun2-core-catcher/issues/59)

## [0.12.2](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.12.1...v0.12.2) (2026-05-16)


### Bug Fixes

* **deps:** update module charm.land/lipgloss/v2 to v2.0.3 ([9f2b891](https://github.com/stuttgart-things/homerun2-core-catcher/commit/9f2b8913e06e2b0ce83432404a0f89617897628c))
* **deps:** update module github.com/stuttgart-things/homerun-library/v3 to v3.1.0 ([4bed7f1](https://github.com/stuttgart-things/homerun2-core-catcher/commit/4bed7f1810cc1658618b4f1107f0faa1f1c3a149))

## [0.12.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.12.0...v0.12.1) (2026-05-15)


### Bug Fixes

* **deps:** update module charm.land/bubbletea/v2 to v2.0.6 ([5d6f841](https://github.com/stuttgart-things/homerun2-core-catcher/commit/5d6f8413e29ac06ebac5164808f38b435e5ab970))
* **deps:** update module github.com/redis/go-redis/v9 to v9.19.0 ([0ba5d3e](https://github.com/stuttgart-things/homerun2-core-catcher/commit/0ba5d3e05dbbfba415723bf2bc2e5a931b30fe82))

# [0.12.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.11.0...v0.12.0) (2026-05-14)


### Features

* **ci:** PR-preview workflows (build-pr + kustomize + cleanup + bot) ([82be9f6](https://github.com/stuttgart-things/homerun2-core-catcher/commit/82be9f684b9b5314632024082cae5314d86b376e)), closes [#110](https://github.com/stuttgart-things/homerun2-core-catcher/issues/110) [stuttgart-things/homerun2-omni-pitcher#116](https://github.com/stuttgart-things/homerun2-omni-pitcher/issues/116) [stuttgart-things/homerun2-omni-pitcher#116](https://github.com/stuttgart-things/homerun2-omni-pitcher/issues/116)

# [0.11.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.10.1...v0.11.0) (2026-04-15)


### Features

* multi-stream subscription (list of Redis streams) ([#53](https://github.com/stuttgart-things/homerun2-core-catcher/issues/53)) ([fe91f40](https://github.com/stuttgart-things/homerun2-core-catcher/commit/fe91f40633145409023b9e6aa9d95beb998b78a6)), closes [#52](https://github.com/stuttgart-things/homerun2-core-catcher/issues/52) [stuttgart-things/homerun-library#83](https://github.com/stuttgart-things/homerun-library/issues/83)

## [0.10.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.10.0...v0.10.1) (2026-04-08)


### Bug Fixes

* hydrate in-memory store from RedisJSON on startup ([163be77](https://github.com/stuttgart-things/homerun2-core-catcher/commit/163be77fc19f507783d34ed6a0c61b4ac8d82275)), closes [#47](https://github.com/stuttgart-things/homerun2-core-catcher/issues/47)

# [0.10.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.9.2...v0.10.0) (2026-04-08)


### Features

* clickable banner to reset filters and local favicon ([2923974](https://github.com/stuttgart-things/homerun2-core-catcher/commit/29239745beffebd2df10fbf8e601f290e56828a5)), closes [#44](https://github.com/stuttgart-things/homerun2-core-catcher/issues/44) [#45](https://github.com/stuttgart-things/homerun2-core-catcher/issues/45)

## [0.9.2](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.9.1...v0.9.2) (2026-04-07)


### Bug Fixes

* use case-insensitive comparison for severity, system, and author filters ([8786ef9](https://github.com/stuttgart-things/homerun2-core-catcher/commit/8786ef945836a6428b8815e2f75112f6b996b82e)), closes [#42](https://github.com/stuttgart-things/homerun2-core-catcher/issues/42)

## [0.9.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.9.0...v0.9.1) (2026-04-07)


### Bug Fixes

* **deps:** update module github.com/stuttgart-things/homerun-library/v3 to v3.0.5 ([4e7a687](https://github.com/stuttgart-things/homerun2-core-catcher/commit/4e7a68750a8a6c600fb5015d84668fc6c9842ff0))

# [0.9.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.8.0...v0.9.0) (2026-04-03)


### Features

* add favicon to web dashboard ([78d5461](https://github.com/stuttgart-things/homerun2-core-catcher/commit/78d5461126880296e6f89a8a483d64f7b1797714)), closes [#38](https://github.com/stuttgart-things/homerun2-core-catcher/issues/38)

# [0.8.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.7.4...v0.8.0) (2026-03-15)


### Features

* web UI tweaks, build date fix, and local dev tasks ([e235e3a](https://github.com/stuttgart-things/homerun2-core-catcher/commit/e235e3af968060f93117a406889e73ff617300fe))

## [0.7.4](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.7.3...v0.7.4) (2026-03-15)


### Bug Fixes

* **deps:** update module github.com/stuttgart-things/homerun-library/v2 to v3 ([24d2f27](https://github.com/stuttgart-things/homerun2-core-catcher/commit/24d2f27382aa2ac278c2151b0b31aaefab42ac32)), closes [#23](https://github.com/stuttgart-things/homerun2-core-catcher/issues/23)

## [0.7.3](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.7.2...v0.7.3) (2026-03-15)


### Bug Fixes

* **deps:** update module charm.land/lipgloss/v2 to v2.0.2 ([1cff7c9](https://github.com/stuttgart-things/homerun2-core-catcher/commit/1cff7c96f055f6da785d8cbf7362aa7bd4b3872e)), closes [#19](https://github.com/stuttgart-things/homerun2-core-catcher/issues/19)

## [0.7.2](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.7.1...v0.7.2) (2026-03-15)


### Bug Fixes

* **deps:** update module charm.land/bubbletea/v2 to v2.0.2 ([18f1d1e](https://github.com/stuttgart-things/homerun2-core-catcher/commit/18f1d1e6a32a257c640ce7fbf926479eea99f011))
* **deps:** update module github.com/redis/go-redis/v9 to v9.18.0 ([5ae1996](https://github.com/stuttgart-things/homerun2-core-catcher/commit/5ae199601fc21ce60bd8b41c7585649065323a51))

## [0.7.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.7.0...v0.7.1) (2026-03-15)


### Bug Fixes

* resolve errcheck and staticcheck lint issues ([6d6d408](https://github.com/stuttgart-things/homerun2-core-catcher/commit/6d6d4088faf8a81b152feb272f3a0942098022cc))

# [0.7.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.6.1...v0.7.0) (2026-03-15)


### Features

* add version/build info to UI and adopt clusterbook color scheme ([cc6debb](https://github.com/stuttgart-things/homerun2-core-catcher/commit/cc6debbea9bf301a7f83ebfab3e9760c4234f94e)), closes [#4f46e5](https://github.com/stuttgart-things/homerun2-core-catcher/issues/4f46e5) [#818cf8](https://github.com/stuttgart-things/homerun2-core-catcher/issues/818cf8) [#f97316](https://github.com/stuttgart-things/homerun2-core-catcher/issues/f97316)

## [0.6.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.6.0...v0.6.1) (2026-03-13)


### Bug Fixes

* add catcherMode web to KCL deploy profile ([d293c14](https://github.com/stuttgart-things/homerun2-core-catcher/commit/d293c14861eb5a7f99ac8b7334d879e64bf94bdf))

# [0.6.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.5.1...v0.6.0) (2026-03-10)


### Features

* add filter dropdowns for System, Severity, Author and time range ([#16](https://github.com/stuttgart-things/homerun2-core-catcher/issues/16)) ([a934cb3](https://github.com/stuttgart-things/homerun2-core-catcher/commit/a934cb3623f11065a7ae3c5c92eab525447c8a40)), closes [#15](https://github.com/stuttgart-things/homerun2-core-catcher/issues/15)

## [0.5.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.5.0...v0.5.1) (2026-03-09)


### Bug Fixes

* correct gateway parentRef in movie-scripts profile ([b9d6211](https://github.com/stuttgart-things/homerun2-core-catcher/commit/b9d6211c4a93b64dcfa6335c23690926d90c0abd))

# [0.5.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.4.0...v0.5.0) (2026-03-09)


### Features

* add Service and HTTPRoute KCL manifests for web mode ([e07c1b6](https://github.com/stuttgart-things/homerun2-core-catcher/commit/e07c1b6be5ffff462c6bc6b6a300927dec80270b))

# [0.4.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.3.0...v0.4.0) (2026-03-08)


### Features

* add web mode with HTMX-based dashboard UI ([d4aec83](https://github.com/stuttgart-things/homerun2-core-catcher/commit/d4aec834084551cd62efb898c7800d90c8cbc6bd)), closes [#6](https://github.com/stuttgart-things/homerun2-core-catcher/issues/6)

# [0.3.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.2.0...v0.3.0) (2026-03-08)


### Features

* add CLI mode with interactive Bubble Tea TUI ([fcbd713](https://github.com/stuttgart-things/homerun2-core-catcher/commit/fcbd7131c2392b88d8a6dbc47b5bb572f26b3830)), closes [#5](https://github.com/stuttgart-things/homerun2-core-catcher/issues/5)

# [0.2.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.1.1...v0.2.0) (2026-03-08)


### Features

* add MessageStore, log mode, FileCatcher, and mode selection ([063bef1](https://github.com/stuttgart-things/homerun2-core-catcher/commit/063bef1d4fa3e92f105e854aa5a2328ffc9afa98)), closes [#3](https://github.com/stuttgart-things/homerun2-core-catcher/issues/3) [#4](https://github.com/stuttgart-things/homerun2-core-catcher/issues/4) [#7](https://github.com/stuttgart-things/homerun2-core-catcher/issues/7) [#8](https://github.com/stuttgart-things/homerun2-core-catcher/issues/8)

## [0.1.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.1.0...v0.1.1) (2026-03-08)


### Bug Fixes

* resolve full message payload from Redis JSON using messageID ([42ce156](https://github.com/stuttgart-things/homerun2-core-catcher/commit/42ce156e4936d9e4de604e44e4bdfe2cc69dda4a)), closes [#1](https://github.com/stuttgart-things/homerun2-core-catcher/issues/1)

# [0.1.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v0.0.0...v0.1.0) (2026-03-08)


### Bug Fixes

* add goreleaser config to exclude windows builds ([bec2f88](https://github.com/stuttgart-things/homerun2-core-catcher/commit/bec2f88c37655581eb376a2e8c357be885d8c2d5))
* set BufferSize and Concurrency for redis consumer ([30398c5](https://github.com/stuttgart-things/homerun2-core-catcher/commit/30398c5e27b9e25b683d822333ff1029f8dc439b))


### Features

* add end-to-end integration test with pitcher via Dagger ([44c0047](https://github.com/stuttgart-things/homerun2-core-catcher/commit/44c00473681f544f43fa7ca7ed594dbb8b4f801d))
* add KCL deployment module and deploy tasks ([c0006b9](https://github.com/stuttgart-things/homerun2-core-catcher/commit/c0006b9041b4e2eaaafc9454ca19b7b7fe0a1796))
* add mkdocs, GitHub Pages workflow, and enable kustomize push ([ac1219c](https://github.com/stuttgart-things/homerun2-core-catcher/commit/ac1219c565e1c972aab2528489c641e2d1a2bde3))

# 1.0.0 (2026-03-08)


### Bug Fixes

* add goreleaser config to exclude windows builds ([bec2f88](https://github.com/stuttgart-things/homerun2-core-catcher/commit/bec2f88c37655581eb376a2e8c357be885d8c2d5))
* set BufferSize and Concurrency for redis consumer ([30398c5](https://github.com/stuttgart-things/homerun2-core-catcher/commit/30398c5e27b9e25b683d822333ff1029f8dc439b))


### Features

* add end-to-end integration test with pitcher via Dagger ([44c0047](https://github.com/stuttgart-things/homerun2-core-catcher/commit/44c00473681f544f43fa7ca7ed594dbb8b4f801d))
* add KCL deployment module and deploy tasks ([c0006b9](https://github.com/stuttgart-things/homerun2-core-catcher/commit/c0006b9041b4e2eaaafc9454ca19b7b7fe0a1796))
* add mkdocs, GitHub Pages workflow, and enable kustomize push ([ac1219c](https://github.com/stuttgart-things/homerun2-core-catcher/commit/ac1219c565e1c972aab2528489c641e2d1a2bde3))
* initial homerun2-core-catcher MVP ([3b57c91](https://github.com/stuttgart-things/homerun2-core-catcher/commit/3b57c91b4e8f9bb395acd7ddac92f842ef1a1880))

## [1.2.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v1.2.0...v1.2.1) (2026-03-08)


### Bug Fixes

* add goreleaser config to exclude windows builds ([bec2f88](https://github.com/stuttgart-things/homerun2-core-catcher/commit/bec2f88c37655581eb376a2e8c357be885d8c2d5))

# [1.2.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v1.1.0...v1.2.0) (2026-03-08)


### Features

* add KCL deployment module and deploy tasks ([c0006b9](https://github.com/stuttgart-things/homerun2-core-catcher/commit/c0006b9041b4e2eaaafc9454ca19b7b7fe0a1796))

# [1.1.0](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v1.0.1...v1.1.0) (2026-03-08)


### Features

* add end-to-end integration test with pitcher via Dagger ([44c0047](https://github.com/stuttgart-things/homerun2-core-catcher/commit/44c00473681f544f43fa7ca7ed594dbb8b4f801d))

## [1.0.1](https://github.com/stuttgart-things/homerun2-core-catcher/compare/v1.0.0...v1.0.1) (2026-03-08)


### Bug Fixes

* set BufferSize and Concurrency for redis consumer ([30398c5](https://github.com/stuttgart-things/homerun2-core-catcher/commit/30398c5e27b9e25b683d822333ff1029f8dc439b))

# 1.0.0 (2026-03-08)


### Features

* initial homerun2-core-catcher MVP ([3b57c91](https://github.com/stuttgart-things/homerun2-core-catcher/commit/3b57c91b4e8f9bb395acd7ddac92f842ef1a1880))
