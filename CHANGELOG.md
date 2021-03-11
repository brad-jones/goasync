## [2.1.2](https://github.com/brad-jones/goasync/compare/v2.1.1...v2.1.2) (2021-03-11)


### Bug Fixes

* **deps:** bump github.com/brad-jones/goerr/v2 from 2.1.1 to 2.1.3 ([#48](https://github.com/brad-jones/goasync/issues/48)) ([0713234](https://github.com/brad-jones/goasync/commit/0713234fd095e6baee33121414eb73f005821667))
* **deps:** bump github.com/stretchr/testify from 1.6.1 to 1.7.0 ([#25](https://github.com/brad-jones/goasync/issues/25)) ([e12a02f](https://github.com/brad-jones/goasync/commit/e12a02fa1305309841d3792c47cbf684b91144fd))

## [2.1.1](https://github.com/brad-jones/goasync/compare/v2.1.0...v2.1.1) (2020-09-18)


### Bug Fixes

* **deps:** bump github.com/brad-jones/goerr/v2 from 2.1.0 to 2.1.1 ([#1](https://github.com/brad-jones/goasync/issues/1)) ([53f968b](https://github.com/brad-jones/goasync/commit/53f968be88fcdd89af620557b2cb469144d74811))

# [2.1.0](https://github.com/brad-jones/goasync/compare/v2.0.0...v2.1.0) (2020-09-18)


### Bug Fixes

* **await:** ensure we don't have any goroutine leaks ([d0b80ef](https://github.com/brad-jones/goasync/commit/d0b80effe3bea6a80eea4aa7a5c6059c04194ccd))


### Features

* **await:** added a Stream() function ([c9a6703](https://github.com/brad-jones/goasync/commit/c9a6703fdb4bd796c3e264bc9b5672b81d58c2eb))
* **await:** added async versions of all methods ([1de1b5c](https://github.com/brad-jones/goasync/commit/1de1b5cbed1c115a967e7c5a784769c17b2af835))
* removed our awaitable and stoppable interfaces ([2594a17](https://github.com/brad-jones/goasync/commit/2594a17a0a54adb2a453aa4d31f55c5cca179f9b))
* **task:** accept func() or func(t *Internal) ([ac21599](https://github.com/brad-jones/goasync/commit/ac215998cbcdf75ec7984fc58440c1b2e4a13f27))
* **task:** added IsCompleted method ([c7c9bd2](https://github.com/brad-jones/goasync/commit/c7c9bd2b20a2b6539e26336d88e2625b403286ec))
* **task:** added Then() method ([2fa58e5](https://github.com/brad-jones/goasync/commit/2fa58e51601f127b379ec95859e6c0e089d26f4b))
* **task:** added Wait() method ([67e30d7](https://github.com/brad-jones/goasync/commit/67e30d7b10700e28678fbeb81f9c0f6faa532f5e))
* **task:** made the done channel public ([3254147](https://github.com/brad-jones/goasync/commit/3254147abb444bcb2584dcc3b331aa7c8fd0059e))

# [2.0.0](https://github.com/brad-jones/goasync/compare/v1.0.0...v2.0.0) (2020-09-11)


### Features

* initial version that uses v2 of goerr ([876f172](https://github.com/brad-jones/goasync/commit/876f17218c7b454ef43d273d1c2a8bb8daaae018))


### BREAKING CHANGES

* uses a completely re-written goerr package
