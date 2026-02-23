package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cmd := flag.String("cmd", "up", "migration command: up or down")
	dir := flag.String("dir", "./migrations", "migration directory")
	flag.Parse()

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	files, err := collectSQLFiles(*dir, *cmd)
	if err != nil {
		log.Fatalf("failed to collect migration files: %v", err)
	}

	if len(files) == 0 {
		log.Printf("no migration files found for cmd=%s", *cmd)
		return
	}

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", f, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("migration failed at %s: %v", f, err)
		}
		fmt.Printf("applied %s\n", f)
	}

	log.Printf("migration completed (%s)", *cmd)
}

func collectSQLFiles(dir, cmd string) ([]string, error) {
	pattern := filepath.Join(dir, "*."+cmd+".sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	if cmd == "down" {
		sort.Sort(sort.Reverse(sort.StringSlice(files)))
	}
	return files, nil
}
