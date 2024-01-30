package htmlRender

import (
	"fmt"
	"os"
	"text/template"
)

// RenderHTMLFromTemplate renders the html template with the given data and writes the output to the given file
func RenderHTMLFromTemplate(templatePath string, data interface{}, outputFilePath string)error {
	htmlTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}
	
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer outputFile.Close()

	// Execute the template and write the output to the file
	err = htmlTemplate.Execute(outputFile, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	return nil
}
