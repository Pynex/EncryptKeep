package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"encryptkeep-backend/internal/blockchain"
	"encryptkeep-backend/internal/keymanager"
	"encryptkeep-backend/internal/vault"
	"encryptkeep-backend/internal/vaultmanager"

	"github.com/ethereum/go-ethereum/crypto"
)

// CLI entrypoint
func main() {
	km := keymanager.NewKeyManager(keymanager.KeyManagerConfig{})
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter master password: ")
	masterPassword, err := readLine(reader)
	if err != nil {
		log.Fatalf("read master password: %v", err)
	}
	if len(masterPassword) < 8 {
		log.Fatalf("master password too short (min 8)")
	}

	if km.HasStoredKeys() {
		if err := km.LoadFromStorage(masterPassword); err != nil {
			log.Fatalf("load keys: %v", err)
		}
		fmt.Println("Keys loaded from storage.")
	} else {
		fmt.Print("Enter private key (64 hex chars): ")
		privHex, err := readLine(reader)
		if err != nil {
			log.Fatalf("read private key: %v", err)
		}
		privHex = strings.TrimSpace(privHex)
		if len(privHex) != 64 {
			log.Fatalf("invalid private key length")
		}
		if _, err := hex.DecodeString(privHex); err != nil {
			log.Fatalf("invalid private key hex: %v", err)
		}
		if err := km.InitializeFirstTime(privHex, masterPassword); err != nil {
			log.Fatalf("init keys: %v", err)
		}
		fmt.Println("Keys initialized and stored.")
	}

	privKey, err := km.GetPrivateKey()
	if err != nil {
		log.Fatalf("get private key: %v", err)
	}
	privHex := hex.EncodeToString(crypto.FromECDSA(privKey))

	svc := blockchain.NewBlockchainService(blockchain.GetDefaultConfig())
	if err := svc.Connect(); err != nil {
		log.Fatalf("blockchain connect: %v", err)
	}
	if _, err := svc.StartSession(privHex, masterPassword); err != nil {
		log.Fatalf("start session: %v", err)
	}

	localVault := vault.NewLocalVault()
	if err := svc.SyncVault(localVault); err != nil {
		log.Fatalf("sync vault: %v", err)
	}
	fmt.Printf("Sync complete. Entries: %d, LastSync: %s\n", len(localVault.Entries), localVault.LastSyncTime.Format("2006-01-02 15:04:05"))

	ctx := context.Background()
	vm := vaultmanager.NewVaultManager(svc, masterPassword)

	for {
		fmt.Print("\nCommands: list, get, add, update, delete, sync, exit\n> ")
		cmd, err := readLine(reader)
		if err != nil {
			log.Fatalf("read command: %v", err)
		}
		cmd = strings.ToLower(cmd)

		switch cmd {
		case "list":
			if len(localVault.Entries) == 0 {
				fmt.Println("No entries.")
				continue
			}
			fmt.Println("Entries:")
			for id, e := range localVault.Entries {
				fmt.Printf("- ID: %s | Title: %s | Username: %s | Updated: %s\n",
					id, e.Title, e.Username, e.UpdatedAt.Format("2006-01-02 15:04:05"))
			}
		case "get":
			id := prompt(reader, "Entry ID", false)
			entry, ok := localVault.Entries[id]
			if !ok {
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

			entry := vault.NewPasswordEntry(title, username, password)
			entry.URL = url

			if err := vm.AddEntry(ctx, localVault, entry); err != nil {
				fmt.Printf("add entry error: %v\n", err)
				continue
			}
			fmt.Println("Entry added and synced.")

		case "update":
			id := prompt(reader, "Entry ID", false)
			entry, ok := localVault.Entries[id]
			if !ok {
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

			if err := vm.UpdateEntry(ctx, localVault, entry); err != nil {
				fmt.Printf("update entry error: %v\n", err)
				continue
			}
			fmt.Println("Entry updated and synced.")

		case "delete":
			id := prompt(reader, "Entry ID", false)
			if _, ok := localVault.Entries[id]; !ok {
				fmt.Println("entry not found")
				continue
			}
			if err := vm.DeleteEntry(ctx, localVault, id); err != nil {
				fmt.Printf("delete entry error: %v\n", err)
				continue
			}
			fmt.Println("Entry deleted and synced.")

		case "sync":
			if err := svc.SyncVault(localVault); err != nil {
				fmt.Printf("sync error: %v\n", err)
				continue
			}
			fmt.Printf("Synced. Entries: %d, LastSync: %s\n", len(localVault.Entries), localVault.LastSyncTime.Format("2006-01-02 15:04:05"))

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
