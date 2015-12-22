# Series-Renamer

Series-Renamer will rename episodes of TV series on your computer with a defined clear and consistent style.

## Example

Illegible titles will can be styled to follow either of the following formats: 

* **{series title} {S}{season#}{E}{episode#}.{extension}**
* **{series title} {S}{season#}{E}{episode#} {episode title}.{extension}**

Examples: 

* *some.series.title.s01e05.branding{1990}.x264.HDTV.mp4*  =>  *Some Series Title S01E05.mp4*
* *some series title s01e05 episode title branding 1990 x264 HDTV.mp4*  =>  *Some Series Title S01E05 Episode Title.mp4*

## Usage

```bash
$ git clone https://github.com/at1as/series-renamer.git
$ go run /path/to/directory/renamer.go {flags}
```
#### Flags

**-p, --path**
*	**Usage**: path to folder containing files to be renamed
*	**Default**: "./"

**-u, --unescape**
* **Usage**: whether to html unescape filename (ex. "Hello%20World" -> "Hello World")
*	**Default**: false

**-k, --keep**
* **Usage**: whether to retain individual episode titles (see examples above)
*	**Default**: false

**-c, --cut**
* **Usage**: text to cut from end of filename (branding, encoding details, etc)
*	**Default**: ""

See `go run renamer.go --help` for details

The terminal output will show you what has been performed

#### Examples

```
$ go run renamer.go -p ~/Movies/ -u -k -c "720p"
$ go run renamer.go
$ go run renamer.go --path /Users/johndoe/vids/
```

## Limitations

- Won't work is season and episode are dot-separated (ex., 'S01.E05')
- Does not recurse directories

## Cowardly Disclaimer

* Will rename files (for which there is no Edit/Undo)
  * Ensure video files are in their own folder
  * Currently the script is stateless and there is no rollback for changes

## TODO

* Add '-i' flag to ignore filetypes (like '.srt') or substrings (like '(Dutch).srt')
