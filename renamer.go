package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "strings"
)


func ensureTrailingSlash(pwd string) string {
  return strings.TrimSuffix(pwd, "/") + "/"
}

func splitString(s, pwd string) {
  sub_string := strings.Split(s, ".")
  sections := len(sub_string)
  title, season, episode := "", "", ""
  season_found, episode_found := false, false
  extension := sub_string[sections-1:]
  string_compare := make(map[string]string)

  // Iterate through each dot-seperated slice
  for i := 0; i < sections; i++ {
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
              // All substrings before the season (SXX) are appended to title (and capitalized)
              for k := 0; k < i; k++ {
                title += strings.Title(sub_string[k]) + " "
              }
            }
          }

        // After Season has been extracted, extract the next number as the episode number
        // If there is a dot between S and E, this will fail (S01.E01). This is not common.
        } else if episode_found == false {
          if sub_string[i][j] >= '0' && sub_string[i][j] <= '9' || sub_string[i][j] == '-' {
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
  }

  // Map each input filepath to the new desired file path for item  
  if title != "" && episode != "" && season != "" && extension[0] != "" {
    old_path := ensureTrailingSlash(pwd) + s
    new_path := ensureTrailingSlash(pwd) + title + "S" + season + "E" + episode + "." + extension[0]
    string_compare[old_path] = new_path
  } else {
    string_compare[ensureTrailingSlash(pwd) + s] = "Parse Error. Skipping"
  }

  renameFiles(string_compare)
}


func getFiles(pwd string) {
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

  // call splitString on each file in the directory (to rename it)
  for _, f := range files {
    splitString(f.Name(), abs_dir)
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
      // value for key was "Parse Error. Skipping" from the splitStrings method. File is likely
      // not a video file. Or it's identified a limitation of the program (ex., spaces in file name)
      fmt.Println("Could not parse " + key + yellow + "  [Skipping] " + reset)
    }
  }
}

func main() {

  // if the script is run with the help flag, print the following
  if len(os.Args) >= 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
    fmt.Println(
          "\nRenamer will parse episodes of tv series and rename them with a clear, consistent style. \n\n" +
          "Example: \n\t some.series.title.s01e05.branding{year}.x264.HDTV.mp4  =>  Some Series Title S01E05.mp4 \n" +
          "Usage: \n\t ./renamer /path/to/directory \n" +
          "Limitations: \n\t - won't work if 's' and 'e' in 's01e01' is dot-seperated" +
          "\n\t - won't handle spaces in input name properly \n\t - recursive folder search isn't implemented\n" +
          "Cowardly Disclaimer: \n\t - use at your own risk. Check terminal output to ensure naming was done correctly")

  // if the script is run with a (non-help) argument, then the argument should be the filepath
  } else if len(os.Args) >= 2 {
    user_path := os.Args[1:2][0]
    getFiles(user_path)

  // if no path is provided, simply default to present directory as the path
  } else {
    getFiles("./")
  }
}
