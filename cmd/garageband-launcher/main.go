package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("GarageBand Launcher")

	// Create template selector
	templates := getTemplates()
	templateSelect := widget.NewSelect(templates, nil)
	if len(templates) > 0 {
		templateSelect.Selected = templates[0]
	}

	// Create project name input
	projectNameInput := widget.NewEntry()
	projectNameInput.SetPlaceHolder("Project name (optional)")

	// Create status label
	status := widget.NewLabel("Ready")
	status.Alignment = fyne.TextAlignCenter

	// Create button
	btn := widget.NewButton("Create New Project", nil)
	if len(templates) == 0 {
		btn.Disable()
		status.SetText("No templates found")
	}

	// Define button action
	btn.OnTapped = func() {
		if templateSelect.Selected == "" {
			status.SetText("Please select a template")
			return
		}

		// Disable button while processing
		btn.Disable()
		status.SetText("Creating project...")

		// Run the script in a goroutine
		go func() {
			err := createProject(templateSelect.Selected, projectNameInput.Text)

			// Update UI
			if err != nil {
				status.SetText("Error: " + err.Error())
				// Reset status after 3 seconds
				time.Sleep(3 * time.Second)
				status.SetText("Ready")
			} else {
				status.SetText("Project created successfully!")
				// Reset status after 2 seconds
				time.Sleep(2 * time.Second)
				status.SetText("Ready")
				// Clear the project name input
				projectNameInput.SetText("")
			}
			btn.Enable()
		}()
	}

	// Create layout
	content := container.NewVBox(
		widget.NewLabel(""), // Spacer
		widget.NewLabelWithStyle(
			"GarageBand Launcher",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true},
		),
		widget.NewLabel(""), // Spacer
		widget.NewLabel("Select Template:"),
		templateSelect,
		widget.NewLabel(""), // Spacer
		widget.NewLabel("Project Name:"),
		projectNameInput,
		widget.NewLabel(""), // Spacer
		btn,
		status,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(300, 350))
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.ShowAndRun()
}

func getTemplates() []string {
	templatesDir := filepath.Join(os.Getenv("HOME"), "Music/GarageBand/Templates")
	templates := []string{}

	// Create Templates directory if it doesn't exist
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		fmt.Printf("Error creating Templates directory: %v\n", err)
		return templates
	}

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return templates
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".band") {
			name := strings.TrimSuffix(entry.Name(), ".band")
			templates = append(templates, name)
		}
	}

	return templates
}

func createProject(templateName, projectName string) error {
	templatesDir := filepath.Join(os.Getenv("HOME"), "Music/GarageBand/Templates")
	outputDir := filepath.Join(os.Getenv("HOME"), "Music/GarageBand")

	// Ensure template exists
	templatePath := filepath.Join(templatesDir, templateName+".band")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template not found at %s", templatePath)
	}

	// Generate project name if not provided
	if projectName == "" {
		projectName = fmt.Sprintf("project-%s", time.Now().Format("20060102-150405"))
	}

	// Add .band extension if needed
	if !strings.HasSuffix(projectName, ".band") {
		projectName += ".band"
	}

	// Create full path for new project
	newProjectPath := filepath.Join(outputDir, projectName)

	// Check if project already exists
	if _, err := os.Stat(newProjectPath); !os.IsNotExist(err) {
		return fmt.Errorf("project '%s' already exists", projectName)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Copy template to new project
	if err := copyDir(templatePath, newProjectPath); err != nil {
		return fmt.Errorf("failed to create project: %v", err)
	}

	// Open the project
	cmd := exec.Command("open", newProjectPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open project: %v", err)
	}

	return nil
}

// Add this helper function to copy directories recursively
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}
