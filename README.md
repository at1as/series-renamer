# Series-Renamer

Series-Renamer will rename episodes of TV series on your computer with a clear and consistent style.

## Example

Illegible titles will all be styled to follow the format {title} {S}{season#}{E}{episode#}.{extension}

Ex. some.series.title.s01e05.branding{year}.x264.HDTV.mp4  =>  Some Series Title S01E05.mp4

## Usage

Renamer is written in Go. To use:

- First, Download renamer to your computer
`$ git clone https://github.com/at1as/series-renamer.git`

Renamer can be used any of these two ways:

- Place in the directory whose files you wish to rename

```bash
$ mv ./renamer.go /path/to/directory/renamer.go
$ chmod 775 /path/to/directory/renamer.go
$ go run /path/to/directory/renamer.go
```

- It can also be passed an absolute path parameter:

```bash
$ chmod 775 ./renamer.go
$ go run ./renamer.go /path/to/directory/to/scan
```

See `go run renamer.go --help` for details


The terminal output will show you what has been performed

## Limitations

- Won't work is season and episode are dot-separated (ex., S01.E05)
- Not guarenteed to work on titles that are space-separated
- Does not recurse directories

## Cowardly Disclaimer

* Will rename files (for which there is no Edit/Undo)
  * Ensure video files are in their own folder
  * Currently the script is stateless and there is no rollback for changes

## TODO

* Add '-i' flag to ignore filetypes (like '.srt') or substrings (like '(Dutch).srt')

