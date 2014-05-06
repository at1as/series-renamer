# Renamer

Renamer will rename episodes of TV series on your computer with a clear and consistent style.

## Example

Illegible titles will all be styled to follow the format {title} {S}{season#}{E}{episode#}.{extension}

Ex. a.series.title.s01e05.branding{year}.x264.HDTV.mp4  =>  A Series Title S01E05.mp4

## Usage

Renamer is written in Go, and can be used two ways:

- Download renamer to your computer
- Place in the directory whose files you wish to rename
- $cdmod 775 /path/to/directory/renamer
- $/path/to/directory/renamer

It can also be passed an absolute path parameter

- Download renamer to your computer
- $cd download/location
- $chmod 775 ./renamer
- ./renamer /path/to/directory/to/scan

See ./renamer --help for details


The terminal output will show you what has been performed

## Limitations

- Won't work is season and episode are dot-separated (ex., S01.E05)
- Probably won't work on titles that are space-separated
- Recursive folder search isn't implemented
- Built and tested on unix/linux/mac environments. Unlikely to work on Windows

## Cowardly Disclaimer

I haven't encountered any issues, however this script is renaming files on your system. Pass it the correct path, watch the terminal output to confirm the result is good, and use at your own risk.

## Why Go?

Perhaps it's not the most convenient language for most people for such a tool, and it certainly doesn't need the efficiency Go offers. However, it's mostly an opportunity to play around and have some fun with what looks like a very interesting language.
