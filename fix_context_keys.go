package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	// Files to fix
	files := []string{
		"internal/handler/product_handler.go",
		"internal/handler/category_handler.go",
	}

	replacements := map[string]string{
		`c.Get("companyID")`:    `c.Get("company_id")`,
		`c.Get("branchID")`:     `c.Get("branch_id")`,
		`c.Get("externalRole")`: `c.Get("external_role")`,
		`c.Get("internalRole")`: `c.Get("internal_role")`,
	}

	for _, file := range files {
		fmt.Printf("Processing %s...\n", file)
		
		// Read file
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		contentStr := string(content)
		modified := false

		// Apply replacements
		for old, new := range replacements {
			if strings.Contains(contentStr, old) {
				contentStr = strings.ReplaceAll(contentStr, old, new)
				modified = true
				fmt.Printf("  - Replaced %s with %s\n", old, new)
			}
		}

		// Write back if modified
		if modified {
			err = ioutil.WriteFile(file, []byte(contentStr), 0644)
			if err != nil {
				log.Printf("Error writing %s: %v\n", file, err)
				continue
			}
			fmt.Printf("✅ %s updated successfully\n", file)
		} else {
			fmt.Printf("⏭️  %s - no changes needed\n", file)
		}
	}

	fmt.Println("\n✅ All files processed!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Rebuild: go build -o server.exe cmd/server/main.go")
	fmt.Println("2. Restart server")
	fmt.Println("3. Login ulang untuk dapat token baru")
}
