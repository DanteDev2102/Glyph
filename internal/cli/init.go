package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/DanteDev2102/Glyph/config"
	"github.com/DanteDev2102/Glyph/internal/gitutils"
	"github.com/DanteDev2102/Glyph/internal/parser"
	"github.com/spf13/cobra"
)

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Security check: Skip symbolic links at the source to prevent information disclosure
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		// Security check: Reject symbolic links at the destination to prevent arbitrary file overwrites
		if info, err := os.Lstat(target); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("destination %s is a symbolic link", target)
			}
		}

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(target)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

func getAuthor(tmplAuthor string, globalAuthor string) string {
	if author != "" {
		return author
	}
	if tmplAuthor != "" {
		return tmplAuthor
	}
	if globalAuthor != "" {
		return globalAuthor
	}
	out, err := exec.Command("git", "config", "user.name").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return "Unknown Author"
}

func replaceInFile(filePath string, replacements map[string]string) {
	// Security check: Don't follow symlinks to prevent path traversal
	info, err := os.Lstat(filePath)
	if err != nil {
		return
	}
	if info.Mode()&os.ModeSymlink != 0 {
		fmt.Printf("Warning: Skipping replacement in %s as it is a symlink\n", filePath)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	s := string(content)
	for k, v := range replacements {
		s = strings.ReplaceAll(s, "{{."+k+"}}", v)
	}

	os.WriteFile(filePath, []byte(s), 0644)
}

func chargeTemplates(cli *Base, initCmd *cobra.Command, commands *[]parser.Command) {
	for i := range *commands {
		command := (*commands)[i]

		initCmd.AddCommand(&cobra.Command{
			Use:   fmt.Sprintf("%s [path]", command.Key),
			Short: command.Short,
			Long:  command.Long,
			Run: func(_ *cobra.Command, args []string) {
				if len(args) == 0 {
					fmt.Println("Directory path is required")
					return
				}

				if len(branch) > 0 && len(tag) > 0 {
					fmt.Println("use only branch or only tag for init project")
					return
				}

				dstPath := filepath.Clean(args[0])
				if dstPath == "/" || dstPath == "." || dstPath == ".." || len(strings.TrimSpace(dstPath)) <= 1 {
					fmt.Println("Dangerous or invalid path provided")
					return
				}

				// Security check: Don't allow initialization into a symbolic link
				if info, err := os.Lstat(dstPath); err == nil {
					if info.Mode()&os.ModeSymlink != 0 {
						fmt.Printf("Destination path %s is a symbolic link, which is not allowed for security reasons.\n", dstPath)
						return
					}
				}

				var err error
				if command.LocalPath != "" {
					err = copyDir(command.LocalPath, dstPath)
				} else {
					var b, t string
					if tag != "" {
						t = tag
					} else if branch != "" {
						b = branch
					} else if command.Tag != "" {
						t = command.Tag
					} else if command.Branch != "" {
						b = command.Branch
					}
					err = gitutils.CloneRepo(command.Repo, dstPath, b, t)
				}

				if err != nil {
					fmt.Println(err)
					return
				}

				// Remove existing .git if it's a clone
				os.RemoveAll(filepath.Join(dstPath, ".git"))

				// Re-initialize git
				exec.Command("git", "init", "--", dstPath).Run()

				// Variables replacement
				projectAuthor := getAuthor(command.Author, cli.Conf.Config.Author)
				projectName := filepath.Base(dstPath)
				replacements := map[string]string{
					"ProjectName": projectName,
					"Author":      projectAuthor,
					"Year":        fmt.Sprintf("%d", time.Now().Year()),
				}

				// Replace in README.md
				readmePath := filepath.Join(dstPath, "README.md")
				replaceInFile(readmePath, replacements)

				// License injection
				selectedLicense := license
				if command.License != "" && license == "MIT" { // only override default MIT if template specifies another
					selectedLicense = command.License
				}

				// Check for --no-license (we'll assume if license is set to "none" via some mechanism)
				// Since we don't have a boolean flag yet, let's check if the user passed it somehow
				// or if we should add it to main.go. I'll check the flags I added.
				// I added: Cli.Root.PersistentFlags().StringVarP(&license, "license", "L", "MIT", ...)
				// User said: "el usuario debe colocar un flag --no-licence en caso de que no quiera usar ninguna licencia"

				// I'll check if a "no-license" flag was passed. I'll need to add it to main.go too.

				if !noLicense && license != "none" {
					licenseName := strings.ToLower(selectedLicense)
					if strings.Contains(licenseName, "..") || strings.Contains(licenseName, "/") || strings.Contains(licenseName, "\\") {
						fmt.Println("Invalid license name, skipping license injection.")
						return
					}
					home, _ := os.UserHomeDir()
					// We should probably look in the installation directory for licenses,
					// but for now let's look in a predictable place or embedded?
					// Let's assume they are in ~/.config/Glyph/licenses/
					licensePath := filepath.Join(home, ".config", "Glyph", "licenses", "LICENSE."+licenseName)

					// If not found there, try the local assets (if running from source)
					if _, err := os.Stat(licensePath); os.IsNotExist(err) {
						licensePath = filepath.Join("internal", "assets", "licenses", "LICENSE."+licenseName)
					}

					if _, err := os.Stat(licensePath); err == nil {
						dstLicense := filepath.Join(dstPath, "LICENSE")
						copyFile(licensePath, dstLicense)
						replaceInFile(dstLicense, replacements)
					}
				}
			},
		})
	}
}

func copyFile(src, dst string) error {
	// Security check: Skip symbolic links at the source to prevent information disclosure
	if info, err := os.Lstat(src); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("source %s is a symbolic link", src)
		}
	}

	// Security check: Reject symbolic links at the destination to prevent arbitrary file overwrites
	if info, err := os.Lstat(dst); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("destination %s is a symbolic link", dst)
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// InitCmd initializes the CLI with the "init" command.
func (cli *Base) InitCmd() {
	initCmd := &cobra.Command{
		Use:   config.InitUse,
		Short: config.InitSummary,
		Long:  config.InitDescription,
	}

	chargeTemplates(cli, initCmd, &cli.Conf.Commmands)

	cli.Root.AddCommand(initCmd)
}
