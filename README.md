# get-ipsws

[![Circle CI](https://circleci.com/gh/blacktop/get-ipsws.png?style=shield)](https://circleci.com/gh/blacktop/get-ipsws) [![Github All Releases](https://img.shields.io/github/downloads/blacktop/get-ipsws/total.svg)](https://github.com/https://github.com/blacktop/get-ipsws/releases/releases) [![GitHub release](https://img.shields.io/github/release/blacktop/get-ipsws.svg)](https://github.com/https://github.com/blacktop/get-ipsws/releases/releases) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

> Download ipsw(s) from [ipsw.me](ipsw.me) or [theiphonewiki.com](theiphonewiki.com)

---

## Getting Started

```bash
$ get-ipsws --help

Usage: get-ipsws [OPTIONS] COMMAND [arg...]

IPSW Downloader

Version: 18.09.1, BuildTime: 2018-10-01T02:34:03Z
Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V                    verbose output
  --dec                            decompress the kernelcache
  --device value, -d value         iOS Device [$IOS_DEVICE]
  --ios-version value, --iv value  iOS Version [$IOS_VERSION]
  --build value, -b value          iOS Build [$IOS_BUILD]
  --help, -h                       show help
  --version, -v                    print the version

Commands:
  generate  crawl theiphonewiki.com and create JSON database
  help      Shows a list of commands or help for one command

Run 'get-ipsws COMMAND --help' for more information on a command.
```

### Download an `ipsw` and extract/decompress the kernelcache

```bash
$ get-ipsws --device iPhone11,2 --build 16A366 --dec

   • Downloading IPSW
      • Downloading... 3.6 GB complete
   • Extracting Kernelcache from IPSW
   • Parsing Compressed Kernelcache
      • compressed size: 17842843, uncompressed: 35727352. unknown: 0x3f9543fd, unknown 1: 0x1
   • Decompressing Kernelcache
```

Notice that the `kernelcache` was extracted from the IPSW and decompressed :smiling_imp:

```bash
$ file kernelcache.release.iphone11.decompressed

kernelcache.release.iphone11.decompressed: "Mach-O 64-bit executable arm64"
```

## TODO

- [ ] use https://github.com/gocolly/colly
- [ ] create offline copy of ipsw.me API
- [ ] crawl ipsw.me for non-betas

## Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/blacktop/get-ipsws/issues/new)

## License

MIT Copyright (c) 2018 **blacktop**
