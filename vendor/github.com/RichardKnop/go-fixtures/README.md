[![Codeship Status for RichardKnop/go-fixtures](https://codeship.com/projects/2972d6c0-37b3-0134-ffd1-02b643534a44/status?branch=master)](https://codeship.com/projects/165845)

# go-fixtures

Django style fixtures for Golang's excellent built-in `database/sql` library. Currently only `YAML` fixtures are supported.

There are two reserved values you can use for `datetime` fields:

* `ON_INSERT_NOW()` will only be used when a row is being inserted
* `ON_UPDATE_NOW()` will only be used when a row is being updated

Example YAML fixture:

```yaml
---

- table: 'some_table'
  pk:
    id: 1
  fields:
    string_field: 'foobar'
    boolean_field: true
    created_at: 'ON_INSERT_NOW()'
    updated_at: 'ON_UPDATE_NOW()'

- table: 'other_table'
  pk:
    id: 2
  fields:
    int_field: 123
    boolean_field: false
    created_at: 'ON_INSERT_NOW()'
    updated_at: 'ON_UPDATE_NOW()'

- table: 'join_table'
  pk:
    some_id: 1
    other_id: 2
```

Example integration for your project:

```go
package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/RichardKnop/go-fixtures"
	"github.com/urfave/cli"
	// Drivers
	_ "github.com/lib/pq"
)

var (
	cliApp *cli.App
)

func init() {
	cliApp = cli.NewApp()
	cliApp.Name = "your-project"
	cliApp.Usage = "Project's usage"
	cliApp.Author = "Your Name"
	cliApp.Email = "your@email"
	cliApp.Version = "0.0.0"
}

func main() {
	db, err := sql.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cliApp.Commands = []cli.Command{
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Action: func(c *cli.Context) error {
				data, err := ioutil.ReadFile(c.Args().First())
				if err != nil {
					return err
				}

				if err := fixtures.Load(data, db, "postgres"); err != nil {
					return err
				}
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) error {
				// Run your web server here
        return nil
			},
		},
	}

	cliApp.Run(os.Args)
}
```
