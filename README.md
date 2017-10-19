# randos
TL;DR: Randos generates a number of files, of various sizes, each filled with _cryptographically random_ binary data.

## What do?
This application creates a set of files, of various sizes, as well as their SHA-512 checksums for verification. The
data inside the file is cryptographically random.

## Why?
So many times, I've needed files of various sizes, full of whatever, to test endpoints, measure performance, and more.
This tool fulfills my needs. It's important to understand that the nature of the data in these files makes it so they
cannot be compressed. However, this makes them ideal for "challenging" work on their data, so if you need to benchmark
encryption or data conversion, they're ideal. In fact, if you can compress these files, file for a patent!

## Build It
You're going to want gb, found [here](https://getgb.io/).

```bash
gb build all
```

That's it.
