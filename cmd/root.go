/*
Copyright © 2024 Joy Bordhen <jbordhen.jb@gmail.com>
*/
package cmd

import (
	"downloader/download"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "downloader",
	Short: "A cli downloader application",
	Long:  `A cli application to download files`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println(args)

		if len(args) == 0 {
			fmt.Println("Welcome to Downloader. Please enter the url you want to download in `downloader link` format. Type `downloader --help` to learn more.")
			return
		}

		var link, _ = cmd.Flags().GetString("link")
		var path, _ = cmd.Flags().GetString("path")
		var concurrency, _ = cmd.Flags().GetInt("concurrency")

		// fmt.Println(link, path, concurrency)

		if link == "" {
			fmt.Println("Link can not be empty, pass --link=`file link`")
			return
		}

		var dl = download.Download{
			Link:        link,
			Path:        path,
			Concurrency: concurrency,
			Size:        0,
			Name:        "",
			StartTime:   time.Now(),
		}

		contentLength, acceptRange, err := dl.Connect()

		dl.Size = contentLength

		if acceptRange {
			fmt.Println("Server supports download by file range.")
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("file size: %v\n", dl.Size)
		fmt.Printf("file name: %v\n", dl.Name)

		downloadErr := dl.Download()

		if downloadErr != nil {
			fmt.Printf("There was an error!\n%v", err)
		} else {
			fmt.Printf("%v downloaded successfully.\n", dl.Name)
			downloadTime := time.Since(dl.StartTime)
			fmt.Printf("Time to complete download: %v\n", downloadTime)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.downloader.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("link", "l", "", "Link to download a file")
	rootCmd.Flags().StringP("path", "p", ".", "Path to save the file")
	rootCmd.Flags().IntP("cocurrency", "c", 10, "Number of concurrent connections while downloading the file")
}
