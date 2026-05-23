package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"encryptkeep-backend/internal/appsvc"
)

// CLI entrypoint
func main() {
	app := appsvc.NewService()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter master password: ")
	masterPassword, err := readLine(reader)
	if err != nil {
		log.Fatalf("read master password: %v", err)
	}
	if len(masterPassword) < 8 {
		log.Fatalf("master password too short (min 8)")
	}
	ctx := context.Background()

	if app.HasStoredKeys() {
		if err := app.Unlock(ctx, masterPassword); err != nil {
			log.Fatalf("unlock: %v", err)
		}
		fmt.Println("Keys loaded from storage.")
	} else {
		fmt.Print("Enter private key (64 hex chars): ")
		privHex, err := readLine(reader)
		if err != nil {
			log.Fatalf("read private key: %v", err)
		}
		if err := app.InitializeNewKeys(ctx, privHex, masterPassword); err != nil {
			log.Fatalf("init keys: %v", err)
		}
		fmt.Println("Keys initialized and stored.")
	}

	n, lastSync, err := app.VaultStats()
	if err != nil {
		log.Fatalf("vault stats: %v", err)
	}
	fmt.Printf("Sync complete. Entries: %d, LastSync: %s\n", n, lastSync.Format("2006-01-02 15:04:05"))

	for {
		fmt.Print("\nCommands: list, get, add, update, delete, sync, exit\n> ")
		cmd, err := readLine(reader)
		if err != nil {
			log.Fatalf("read command: %v", err)
		}
		cmd = strings.ToLower(cmd)

		switch cmd {
		case "list":
			entries, err := app.ListEntries()
			if err != nil {
				fmt.Printf("list error: %v\n", err)
				continue
			}
			if len(entries) == 0 {
				fmt.Println("No entries.")
				continue
			}
			fmt.Println("Entries:")
			for _, e := range entries {
				fmt.Printf("- ID: %s | Title: %s | Username: %s | Updated: %s\n",
					e.ID, e.Title, e.Username, e.UpdatedAt.Format("2006-01-02 15:04:05"))
			}
		case "get":
			id := prompt(reader, "Entry ID", false)
			entry, err := app.GetEntry(id)
			if err != nil {
				fmt.Println("entry not found")
				continue
			}

			fmt.Printf("ID: %s\nTitle: %s\nUsername: %s\nPassword: %s\nURL: %s\nUpdated: %s\n",
				entry.ID, entry.Title, entry.Username, entry.Password, entry.URL, entry.UpdatedAt.Format("2006-01-02 15:04:05"))
		case "add":
			title := prompt(reader, "Title", false)
			username := prompt(reader, "Username", false)
			password := prompt(reader, "Password", false)
			url := prompt(reader, "URL (optional)", true)

			if err := app.AddEntry(ctx, title, username, password, url); err != nil {
				fmt.Printf("add entry error: %v\n", err)
				continue
			}
			fmt.Println("Entry added and synced.")

		case "update":
			id := prompt(reader, "Entry ID", false)
			entry, err := app.GetEntry(id)
			if err != nil {
				fmt.Println("entry not found")
				continue
			}

			title := prompt(reader, fmt.Sprintf("Title [%s]", entry.Title), true)
			username := prompt(reader, fmt.Sprintf("Username [%s]", entry.Username), true)
			password := prompt(reader, "Password [leave empty to keep]", true)
			url := prompt(reader, fmt.Sprintf("URL [%s]", entry.URL), true)

			if title != "" {
				entry.Title = title
			}
			if username != "" {
				entry.Username = username
			}
			if password != "" {
				entry.Password = password
			}
			if url != "" {
				entry.URL = url
			}
			entry.UpdatedAt = time.Now()

			if err := app.UpdateEntry(ctx, entry); err != nil {
				fmt.Printf("update entry error: %v\n", err)
				continue
			}
			fmt.Println("Entry updated and synced.")

		case "delete":
			id := prompt(reader, "Entry ID", false)
			if _, err := app.GetEntry(id); err != nil {
				fmt.Println("entry not found")
				continue
			}
			if err := app.DeleteEntry(ctx, id); err != nil {
				fmt.Printf("delete entry error: %v\n", err)
				continue
			}
			fmt.Println("Entry deleted and synced.")

		case "sync":
			if err := app.Sync(ctx); err != nil {
				fmt.Printf("sync error: %v\n", err)
				continue
			}
			n, lastSync, err := app.VaultStats()
			if err != nil {
				fmt.Printf("stats error: %v\n", err)
				continue
			}
			fmt.Printf("Synced. Entries: %d, LastSync: %s\n", n, lastSync.Format("2006-01-02 15:04:05"))

		case "exit", "quit":
			fmt.Println("Bye.")
			return

		default:
			fmt.Println("Unknown command.")
		}
	}
}

func readLine(r *bufio.Reader) (string, error) {
	text, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func prompt(r *bufio.Reader, label string, allowEmpty bool) string {
	for {
		fmt.Printf("%s: ", label)
		txt, err := readLine(r)
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			continue
		}
		if txt == "" && !allowEmpty {
			fmt.Println("value cannot be empty")
			continue
		}
		return txt
	}
}
