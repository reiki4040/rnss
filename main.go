package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	EnvAWSRgion = "AWS_REGION"
)

var (
	version  string
	revision string

	optUsage         bool
	optFlushCache    bool
	optRegion        string
	optShowCommand   bool
	optListFromStdin bool
)

func init() {
	flag.BoolVar(&optUsage, "h", false, "show usage")
	flag.BoolVar(&optUsage, "help", false, "show usage")
	flag.BoolVar(&optFlushCache, "f", false, "flush ec2 list cache. set this if you change state EC2.")
	flag.StringVar(&optRegion, "region", "", "target aws region.")
	flag.BoolVar(&optShowCommand, "show-command", false, "show aws ssm start-session command, NOT run it.")
	flag.BoolVar(&optListFromStdin, "stdin", false, "instance list from stdin. required line starts with `instance-id<tab>`")
	flag.Parse()
}

func showHelp() {
	usage := `rnss is instance selection helper for ssm start session.
you can show EC2 instances and select in CUI then start ssm sesion.

  rnss is simple wrapper that call below.
    aws ssm start-session --target <you selected instance-id>

[Usage]

  rnss [Options] [filter phrase]

[Optoins]
`
	usageLast := `
[filter phrase]

  command args are filter ec2 instance info in line.
  ex) if you call with 'web', then filtered instances that incluide web in Name tag.
`
	fmt.Printf("rnss %s[%s]\n", version, revision)
	fmt.Println(usage)
	flag.PrintDefaults()
	fmt.Println(usageLast)
}

func main() {
	if optUsage {
		showHelp()
		return
	}

	var ec2list []string
	var err error
	if optListFromStdin {
		if isStdinEmpty() {
			fmt.Println("set -stdin however stdin is empty.")
			os.Exit(1)
		}
		in, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("failed read list from stdin: %v", err)
			os.Exit(1)
		}

		for _, l := range strings.Split(string(in), "\n") {
			if len(l) == 0 {
				continue
			}
			ec2list = append(ec2list, l)
		}
	} else {
		region := ""
		if optRegion != "" {
			region = optRegion
		}

		ctx := context.Background()
		if optFlushCache {
			ec2list, err = getEC2InfoAndStoreCache(ctx, region)
			if err != nil {
				fmt.Printf("failed get ec2 info: %v\n", err)
				os.Exit(1)
			}
		} else {
			// read from cache. call aws if failed cache reading.
			var cacheList []string
			c, err := NewEC2Cache(region)
			if err != nil {
				fmt.Printf("failed get ec2 info from cache: %v\n", err)
				cacheList, err = getEC2InfoAndStoreCache(ctx, region)
				if err != nil {
					fmt.Printf("failed get ec2 info: %v\n", err)
					os.Exit(1)
				}
			} else {
				cacheList, err = c.Get()
				if err != nil || len(cacheList) == 0 {
					if err != nil && !os.IsNotExist(err) {
						fmt.Printf("failed get ec2 info from cache: %v\n", err)
					}
					cacheList, err = getEC2InfoAndStoreCache(ctx, region)
					if err != nil {
						fmt.Printf("failed get ec2 info: %v\n", err)
						os.Exit(1)
					}
				}
			}
			ec2list = cacheList
		}
	}

	if len(ec2list) == 0 {
		fmt.Println("there is no running instance.")
		return
	}

	// args are filter filterPhrase
	filterPhrase := []rune{}
	if len(flag.Args()) > 0 {
		phraseStr := strings.Join(flag.Args(), " ")
		filterPhrase = []rune(phraseStr)
	}

	cui, err := NewSelectionCUI(ec2list, filterPhrase)
	if err != nil {
		fmt.Printf("failed initialise selection CUI: %v", err)
		os.Exit(1)
	}
	cui.list.Title = "select start-session instance"

	// start selection CUI by Bubble Tea
	p := tea.NewProgram(cui,
		tea.WithAltScreen(),
	)
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	selectedInstanceId := ""
	retm, ok := m.(model)
	if ok {
		if retm.selected {
			selectedInstanceId = retm.selectedItem.InstanceId
		} else {
			// canceled
			fmt.Println("none selected instance ID.")
			return
		}
	} else {
		fmt.Printf("Why return other struct??: %v\n", retm)
		os.Exit(1)
	}

	awsCmdArgs := []string{
		"ssm",
		"start-session",
		"--target",
		selectedInstanceId,
	}
	if optShowCommand {
		fmt.Println("aws", strings.Join(awsCmdArgs, " "))
		return
	}

	// run aws ssm start-session --target <instance-id>
	cmd := exec.CommandContext(
		context.Background(),
		"aws",
		awsCmdArgs...,
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	signal.Ignore(os.Interrupt)
	if err := cmd.Run(); err != nil {
		fmt.Printf("failed aws ssm start-session: %v", err)
	}
}

func isStdinEmpty() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode()&os.ModeCharDevice) != 0 || stat.Size() <= 0
}

func getEC2InfoAndStoreCache(ctx context.Context, region string) ([]string, error) {
	ec2list, err := GetEC2ListForStartSession(ctx, region)
	if err != nil {
		fmt.Printf("failed get ec2 info: %v\n", err)
		os.Exit(1)
	}

	// store cache
	c, err := NewEC2Cache(region)
	if err != nil {
		fmt.Printf("failed store ec2 info: %v, however still selection...\n", err)
	} else {
		err = c.Store(ec2list)
		if err != nil {
			fmt.Printf("failed store ec2 info: %v, however still selection...\n", err)
		}
	}
	return ec2list, nil
}
