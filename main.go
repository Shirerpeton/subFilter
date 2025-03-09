package main

import (
	"errors"
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
	subPath := flag.String("sub", "", "Path to input srt file or directory containing them")
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
	subStat, err := os.Stat(*subPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !subStat.IsDir() {
		var outputPath string
		if *output == "" {
			outputPath = *output
		} else {
			outputPath = getOutputPath(*subPath)
		}
		content, err := os.ReadFile(*subPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		contentStr := string(content)
		sub, err := subParser.GetSubFromSrt(contentStr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sub.Filter(*filter)
		srt := sub.GetSrtFromSub()
		os.WriteFile(outputPath, []byte(srt), 0644)
		fmt.Printf("done - %s\n", outputPath)
	} else {
		outputFolder := *output
		if outputFolder == "" {
			outputFolder = "./output/"
		}
		err := os.Mkdir(outputFolder, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			fmt.Println(err)
			os.Exit(1)
		}
		if !strings.HasSuffix(outputFolder, "/") {
			outputFolder += "/"
		}
		entries, err := os.ReadDir(*subPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, entr := range entries {
			if entr.IsDir() {
				continue
			}
			subPath := *subPath + entr.Name()
			outputPath := outputFolder + getOutputPath(entr.Name())
			content, err := os.ReadFile(subPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			contentStr := string(content)
			sub, err := subParser.GetSubFromSrt(contentStr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			sub.Filter(*filter)
			srt := sub.GetSrtFromSub()
			os.WriteFile(outputPath, []byte(srt), 0644)
			fmt.Printf("done - %s\n", outputPath)
		}
	}
}
