# Changelog

## [v0.7.2](https://github.com/yuuki/lstf/compare/v0.7.1...v0.7.2) (2021-03-14)

* fix ENOENT error because of TOCTOU [#31](https://github.com/yuuki/lstf/pull/31) ([yuuki](https://github.com/yuuki))

## [v0.7.1](https://github.com/yuuki/lstf/compare/v0.7.0...v0.7.1) (2021-03-14)

* avoid too many open files error [#30](https://github.com/yuuki/lstf/pull/30) ([yuuki](https://github.com/yuuki))
* fix typo [#29](https://github.com/yuuki/lstf/pull/29) ([dmitris](https://github.com/dmitris))
* bump version Go 1.15 [#28](https://github.com/yuuki/lstf/pull/28) ([yuuki](https://github.com/yuuki))

## [v0.7.0](https://github.com/yuuki/lstf/compare/v0.6.0...v0.7.0) (2021-01-11)

* Migrate CI to GitHub Action [#27](https://github.com/yuuki/lstf/pull/27) ([yuuki](https://github.com/yuuki))
* Watch mode [#24](https://github.com/yuuki/lstf/pull/24) ([yuuki](https://github.com/yuuki))
* fix an error of parsing comm including white spaces [#26](https://github.com/yuuki/lstf/pull/26) ([yuuki](https://github.com/yuuki))
* Fix an error of parsing spesific comm pattern [#25](https://github.com/yuuki/lstf/pull/25) ([yuuki](https://github.com/yuuki))
* improving resource utilization [#23](https://github.com/yuuki/lstf/pull/23) ([yuuki](https://github.com/yuuki))
* Go 1.13 and updating some modules [#21](https://github.com/yuuki/lstf/pull/21) ([yuuki](https://github.com/yuuki))
* Display both name and ipaddr [#20](https://github.com/yuuki/lstf/pull/20) ([yuuki](https://github.com/yuuki))

## [v0.6.0](https://github.com/yuuki/lstf/compare/v0.5.4...v0.6.0) (2019-08-20)

* --filter option [#19](https://github.com/yuuki/lstf/pull/19) ([yuuki](https://github.com/yuuki))

## [v0.5.4](https://github.com/yuuki/lstf/compare/v0.5.3...v0.5.4) (2019-06-14)

* JSONize process field [#18](https://github.com/yuuki/lstf/pull/18) ([yuuki](https://github.com/yuuki))
* Continue if it occurs permission denied error when specified --processes [#17](https://github.com/yuuki/lstf/pull/17) ([yuuki](https://github.com/yuuki))

## [v0.5.3](https://github.com/yuuki/lstf/compare/v0.5.2...v0.5.3) (2019-06-12)


## [v0.5.2](https://github.com/yuuki/lstf/compare/v0.5.1...v0.5.2) (2019-06-12)


## [v0.5.1](https://github.com/yuuki/lstf/compare/v0.5.0...v0.5.1) (2019-06-12)

* Ignore 'readlink: no such file or directory' [#16](https://github.com/yuuki/lstf/pull/16) ([yuuki](https://github.com/yuuki))
* Debug option [#15](https://github.com/yuuki/lstf/pull/15) ([yuuki](https://github.com/yuuki))
* Improve displaying process [#14](https://github.com/yuuki/lstf/pull/14) ([yuuki](https://github.com/yuuki))

## [v0.5.0](https://github.com/yuuki/lstf/compare/v0.4.3...v0.5.0) (2019-06-08)

* Aggregating process group [#13](https://github.com/yuuki/lstf/pull/13) ([yuuki](https://github.com/yuuki))
* Add --processes option [#11](https://github.com/yuuki/lstf/pull/11) ([yuuki](https://github.com/yuuki))
* Avoid to fallback to prasing procfs except netlink error [#12](https://github.com/yuuki/lstf/pull/12) ([yuuki](https://github.com/yuuki))
* Migrate to Go modules [#10](https://github.com/yuuki/lstf/pull/10) ([yuuki](https://github.com/yuuki))

## [v0.4.3](https://github.com/yuuki/lstf/compare/v0.4.2...v0.4.3) (2018-06-19)

* Improve performance by removing the inode and pid process [#9](https://github.com/yuuki/lstf/pull/9) ([yuuki](https://github.com/yuuki))

## [v0.4.2](https://github.com/yuuki/lstf/compare/v0.4.1...v0.4.2) (2018-06-19)

* Add fallback to procfs after try netlink [#8](https://github.com/yuuki/lstf/pull/8) ([yuuki](https://github.com/yuuki))

## [v0.4.1](https://github.com/yuuki/lstf/compare/v0.4.0...v0.4.1) (2018-06-17)

* Embed oss credits to go binary [#7](https://github.com/yuuki/lstf/pull/7) ([yuuki](https://github.com/yuuki))

## [v0.4.0](https://github.com/yuuki/lstf/compare/v0.3.0...v0.4.0) (2018-06-17)

* Improve performance for getting connection stats in linux by netlink [#6](https://github.com/yuuki/lstf/pull/6) ([yuuki](https://github.com/yuuki))
* Reduce the number of code to read /proc/net/tcp [#5](https://github.com/yuuki/lstf/pull/5) ([yuuki](https://github.com/yuuki))

## [v0.3.0](https://github.com/yuuki/lstf/compare/v0.2.1...v0.3.0) (2018-03-23)

* Use local ip addr instead of localhost [#4](https://github.com/yuuki/lstf/pull/4) ([yuuki](https://github.com/yuuki))

## [v0.2.1](https://github.com/yuuki/lstf/compare/v0.2.0...v0.2.1) (2018-03-16)

* Fix misdetection to listening ports [#3](https://github.com/yuuki/lstf/pull/3) ([yuuki](https://github.com/yuuki))

## [v0.2.0](https://github.com/yuuki/lstf/compare/v0.1.0...v0.2.0) (2018-03-04)

* Support json [#2](https://github.com/yuuki/lstf/pull/2) ([yuuki](https://github.com/yuuki))
* Numeric options [#1](https://github.com/yuuki/lstf/pull/1) ([yuuki](https://github.com/yuuki))

## [v0.1.0](https://github.com/yuuki/lstf/compare/...v0.1.0) (2018-03-04)

- Initial release.
