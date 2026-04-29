// Orphan columns: PostgreSQL columns that exist on migrated tables but are not defined
// on any GORM model field that participates in migration (same rules as AutoMigrate).
//
// Usage (from backend/):
//
//	go run ./cmd/orphan_columns
//	go run ./cmd/orphan_columns -json   # JSON report (orphan_columns, missing_tables)
//	go run ./cmd/orphan_columns -fail   # exit 1 if any orphan columns
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type orphan struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

type report struct {
	OrphanColumns  []orphan `json:"orphan_columns"`
	MissingTables  []string `json:"missing_tables,omitempty"`
	OrphanCount    int      `json:"orphan_count"`
	MissingCount   int      `json:"missing_table_count"`
}

func main() {
	jsonOut := flag.Bool("json", false, "print results as JSON")
	fail := flag.Bool("fail", false, "exit with status 1 if any orphan columns are found")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("info: no .env in cwd: %v", err)
	}
	bootstrap.LoadEnv()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := bootstrap.InitDB(dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	db = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})

	namer := db.Config.NamingStrategy
	if namer == nil {
		namer = schema.NamingStrategy{}
	}

	var schemaCache sync.Map
	// table -> set of expected column names (unioned if multiple models shared a table)
	expectedByTable := make(map[string]map[string]struct{})
	// one model pointer per table for HasTable / ColumnTypes
	repByTable := make(map[string]interface{})

	for _, model := range models.AllModels() {
		s, err := schema.Parse(model, &schemaCache, namer)
		if err != nil {
			log.Fatalf("schema %T: %v", model, err)
		}
		tbl := s.Table
		if expectedByTable[tbl] == nil {
			expectedByTable[tbl] = make(map[string]struct{})
			repByTable[tbl] = model
		}
		for _, dbName := range s.DBNames {
			f := s.FieldsByDBName[dbName]
			if f != nil && !f.IgnoreMigration {
				expectedByTable[tbl][dbName] = struct{}{}
			}
		}
	}

	var orphans []orphan
	var missingTables []string

	tables := make([]string, 0, len(expectedByTable))
	for t := range expectedByTable {
		tables = append(tables, t)
	}
	sort.Strings(tables)

	for _, tbl := range tables {
		rep := repByTable[tbl]
		if !db.Migrator().HasTable(rep) {
			missingTables = append(missingTables, tbl)
			continue
		}
		colTypes, err := db.Migrator().ColumnTypes(rep)
		if err != nil {
			log.Fatalf("column types %q: %v", tbl, err)
		}
		want := expectedByTable[tbl]
		for _, ct := range colTypes {
			name := ct.Name()
			if _, ok := want[name]; !ok {
				orphans = append(orphans, orphan{Table: tbl, Column: name})
			}
		}
	}

	sort.Strings(missingTables)

	rep := report{
		OrphanColumns: orphans,
		MissingTables: missingTables,
		OrphanCount:   len(orphans),
		MissingCount:  len(missingTables),
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			log.Fatal(err)
		}
	} else {
		if len(missingTables) > 0 {
			fmt.Println("Tables expected by models but not present in the database:")
			for _, t := range missingTables {
				fmt.Printf("  %s\n", t)
			}
			fmt.Println()
		}
		if len(orphans) == 0 {
			fmt.Println("No orphan columns (no extra columns beyond model definitions).")
		} else {
			fmt.Println("Orphan columns (present in DB, not defined on migrated model fields):")
			for _, o := range orphans {
				fmt.Printf("  %s.%s\n", o.Table, o.Column)
			}
			fmt.Printf("\nTotal: %d orphan column(s) across %d table(s).\n", len(orphans), countDistinctTables(orphans))
		}
	}

	if *fail && len(orphans) > 0 {
		os.Exit(1)
	}
}

func countDistinctTables(orphans []orphan) int {
	seen := make(map[string]struct{})
	for _, o := range orphans {
		seen[o.Table] = struct{}{}
	}
	return len(seen)
}
