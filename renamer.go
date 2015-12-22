package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/url"
  "os"
  "path/filepath"
  "strings"

  "github.com/codegangsta/cli"
)


func ensureTrailingSlash(pwd string) string {
  return strings.TrimSuffix(pwd, "/") + "/"
}

func splitSpaceOrDot(raw_string string) []string {
  return strings.FieldsFunc(raw_string, func (r rune) bool {
    return r == ' ' || r == '.'
  })
}

func htmlUnescape(filename string, modify bool) string {
  // Remove URL escaping (ex. "%20" -> " ")
  if modify == false {
    return filename
  }

  unencoded_str, err := url.QueryUnescape(filename)
  if err != nil {
    return filename
  }
  return unencoded_str
}

func cutFromTitle(title string, text string, modify bool) string {
  // Cut text from end of filename:
  //  ex. cutFromTitle("filenameABC.mkv", "BC") -> "filenameA.mkv"
  if modify == true && text != "" {
    dot_separated := strings.Split(title, ".")

    extension := dot_separated[len(dot_separated) - 1]
    name := strings.Join(dot_separated[0:(len(dot_separated) - 1)], " ")
    name = strings.TrimRight(name, text)

    return name + "." + extension
  }
  return title
}

func generateFileName(s, pwd string, unescape bool, keep bool, cut string) {
  // Create mapping of filename to newly formatted filename
  filename := htmlUnescape(s, unescape)
  filename = cutFromTitle(filename, cut, keep)
  sub_string := splitSpaceOrDot(filename)
  sections := len(sub_string)

  title, season, episode := "", "", ""
  trailing_text := ""
  season_found, episode_found := false, false

  extension := sub_string[sections-1:]
  string_compare := make(map[string]string)

  // Iterate through each dot-seperated slice (up to extension)
  for i := 0; i < sections - 1; i++ {

    if episode != "" {
      episode_found = true
    }
    // Sections beginning with "s" or "S" *may* denote the season
    if strings.HasPrefix(sub_string[i], "s") || strings.HasPrefix(sub_string[i], "S")   {
      section := len(sub_string[i])

      for j := 1; j < section; j++ {

        if season_found == false {
          // If the character directly after an 's'/'S' is a number, append it to season variable
          if sub_string[i][j] >= '0' && sub_string[i][j] <= '9' {
            season = season + string(sub_string[i][j])
          } else {
            // If at least one int was found after the "S", the season information has been extracted
            if season != "" {
              season_found = true
              // All substrings before the season (SXX) are appended to series title (and capitalized)
              for k := 0; k < i; k++ {
                title += strings.Title(sub_string[k]) + " "
              }
            }
          }

        // After Season has been extracted, extract the next number as the episode number
        // If there is a dot between S and E, this will fail (S01.E01). This is not common.
        } else if episode_found == false {
          if sub_string[i][j] >= '0' && sub_string[i][j] <= '9' || sub_string[i][j] == '-' || sub_string[i][j] == 'E' {
            episode = episode + string(sub_string[i][j])
          } else {
            // If at least one int was found after the E, the episode information has been extracted
            // "-" included, to handle less common cases such as S01E01-02 => E01-02
            if episode != "" {
              episode_found = true
            }
          }
        }
      }
    }
    // Include episode title after S01E01 if flag was set
    if keep && season_found && episode_found {
      trailing_text += " " + strings.Title(sub_string[i])
    }
  }

  // Map each input filepath to the new desired file path for item  
  if title != "" && episode != "" && season != "" && extension[0] != "" {
    old_path := ensureTrailingSlash(pwd) + s
    new_path := ensureTrailingSlash(pwd) + title + "S" + season + "E" + episode + trailing_text + "." + extension[0]
    string_compare[old_path] = new_path
  } else {
    string_compare[ensureTrailingSlash(pwd) + s] = "Parse Error. Skipping"
  }

  renameFiles(string_compare)
}


func getFiles(pwd string, unescape bool, keep bool, cut string) {
  // Since abs_dir must be pre-declared, so must err. If abs_dir and error are defined in the
  // if statement they have no scope outside of it. Pre-declaration is the best mitigation I can find...
  abs_dir := ""
  var err error = nil

  // if path is relative to script, obtain absolute path. Otherwise use pwd argument as absolute path
  if pwd == "./" {
    abs_dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
      log.Fatal(err)
    }
  } else {
    abs_dir = pwd
  }

  // read all files in specified directory
  files, _ := ioutil.ReadDir(abs_dir)

  // call generateFileName on each file in the directory (to rename it)
  for _, f := range files {
    generateFileName(f.Name(), abs_dir, unescape, keep, cut)
  }
}

func renameFiles(files map[string]string) {
  // Colored terminal outputs. This works in unix/linux/mac os only.
  reset, red, green, yellow := "\x1b[0m", "\x1b[31m", "\x1b[32m", "\x1b[33m"

  for key, val := range files {
    if val != "Parse Error. Skipping" {
      fmt.Println("\nRenaming:  \t" + key + green + "  [...]" + reset)
      err := os.Rename(key, val)
      if err != nil {
        // if error was returned from the attempt to rename file, notify user of failure
        fmt.Println("Unknown error renaming this file" + red + "  [Error]" + reset)
      } else {
        fmt.Println("Updated: \t" + val + green + "  [Success]" + reset + "\n")
      }
    } else {
      // value for key was "Parse Error. Skipping" from the generateFileNames method. File is likely
      // not a video file. Or it's identified a limitation of the program (ex., spaces in file name)
      fmt.Println("Could not parse " + key + yellow + "  [Skipping] " + reset)
    }
  }
}

func main() {
  var unescape bool
  var keep bool
  var cut string
  var path string

  app := cli.NewApp()
  app.Name = "Series Renamer"
  app.Flags = []cli.Flag {
    cli.BoolFlag{
      Name:        "u, unescape",
      Usage:       "whether to html unescape filename",
      Destination: &unescape,
    },
    cli.BoolFlag {
      Name:        "k, keep",
      Usage:       "whether to keep episode names",
      Destination: &keep,
    },
    cli.StringFlag{
      Name:        "c, cut",
      Value:       "",
      Usage:       "text to cut from end of filename (branding, encoding details, etc)",
      Destination: &cut,
    },
    cli.StringFlag{
      Name:        "p, path",
      Value:       "./",
      Usage:       "path to folder containing files to be renamed",
      Destination: &path,
    },
  }

  app.Action = func(c *cli.Context) {
    getFiles(path, unescape, keep, cut)
  }

  app.Run(os.Args)
}

