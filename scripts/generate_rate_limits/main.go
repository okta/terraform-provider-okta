package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	rateLimitMappingsTxt string
	rootCmd              = &cobra.Command{
		Use:   "go run main.go -h",
		Short: "generate rate limits code",
		Run: func(cmd *cobra.Command, args []string) {
			mappingsPath := cmd.Flag("mappings").Value
			mappingsFile, err := os.Open(mappingsPath.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read %q: %v\n", mappingsPath, err)
				return
			}

			lines := []string{}
			scanner := bufio.NewScanner(mappingsFile)

			reID := regexp.MustCompile(`{[\w]+}`)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "/api/v1/internal") ||
					!(strings.HasPrefix(line, "/.well-known") ||
						strings.HasPrefix(line, "/api/v1") ||
						strings.HasPrefix(line, "/oauth2")) {
					continue
				}
				// 0 path, 1 method, 3 bucket, 4 type
				values := strings.Split(line, " ")
				if len(values) < 5 {
					fmt.Fprintf(os.Stderr, "unknown format of mapping line: %s\n", line)
					return
				}
				if values[4] != "URL" {
					continue
				}
				path := reID.ReplaceAllString(values[0], "ID")
				lines = append(lines, fmt.Sprintf("%s %s %s", path, values[1], values[3]))
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "reading lines from %q failed: %v\n", mappingsPath, err)
				return
			}
			sort.Strings(lines)

			tmplPath := "rate_limit_lines.tmpl"
			tmpl, err := os.ReadFile(tmplPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read %q: %v\n", tmplPath, err)
			}

			t := template.Must(template.New("tmpl").Parse(string(tmpl)))
			rlGoPath := "../../okta/internal/apimutex/rate_limit_lines.go"
			rlGoFile, err := os.Create(rlGoPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open %q: %v\n", rlGoPath, err)
				return
			}
			t.Execute(rlGoFile, lines)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&rateLimitMappingsTxt, "mappings", "m", "", "path to Okta rate limit mappings fixture for ITs - any of the files in monolith source components/tests/api/webapp/src/test/resources/ratelimits")
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: '%s'", err)
		rootCmd.Help()
		os.Exit(1)
	}
}
