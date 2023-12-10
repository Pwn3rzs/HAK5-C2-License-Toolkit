# Hak5 C2 Licensing Toolkit

### Usage

- Get the patched binaries here:
  - [Telegram]()
  - [CyberArsenal]()
- `git clone https://github.com/Pwn3rzs/HAK5-C2-License-Toolkit`
- go build
  - Remember to specify `GOOS` and `GOARCH` for different OS and ARCH
- move `./HAK5-C2-License-Toolkit` inside the same path of `c2-<x.x.x>_<arch>_<os>(.exe)` binary
- run it and choose options
  - `generate` to generate a test License struct hex string
  - `decode` to decode a License / Status struct hex string
  - `read` to read the values inside `Setup[License]` or `Status[Status]` buckets struct hex string from DB
  - `crack` to start process of inserting license values inside DB
    - **MAKE SURE YOU HAVE REPLACED THE BINARY OR LICENSE WILL RESET**


### How To

Read the public "blog" / "article" / "tutorial" here:
- [CyberArsenal Post]()
- [Telegram Channel]()