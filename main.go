package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/shirerpeton/subFilter/internal/subParser"
)

func getOutputPath(input string) string {
	parts := strings.Split(input, ".")
	if len(parts) == 1 {
		return input + "_filter"
	}
	parts[len(parts) - 2] = parts[len(parts) - 2] + "_filter"
	return strings.Join(parts, ".")
}

func main() {
	subPath := flag.String("sub", "", "Path to input subtitle file or directory containing them")
	filter := flag.String("filter", "", "All dialog lines containing this string will be cut out")
	output := flag.String("out", "", "Path to output subtitle file, defaults to input filename with _filter suffix, for diretory processing must be a directory name as well")
	flag.Parse()

	if *subPath == "" {
		fmt.Println("provide input subtitle file path")
		os.Exit(1)
	}
	if *filter == "" {
		fmt.Println("provide filter value different which is not empty")
		os.Exit(1)
	}
	content, err := os.ReadFile(*subPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	contentStr := string(content)
	outputPath := ""
	if *output != "" {
		outputPath = *output
	} else {
		outputPath = getOutputPath(*subPath)
	}
	sub, err := subParser.GetSubFromSrt(contentStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sub.Filter(*filter)
	srt := sub.GetSrtFromSub()
	os.WriteFile(outputPath, []byte(srt), 0644)
}
