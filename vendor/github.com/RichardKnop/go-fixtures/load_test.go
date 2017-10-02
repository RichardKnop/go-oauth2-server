package fixtures_test

var testSchemaSQLite = `
CREATE TABLE some_table(
  id INT PRIMARY KEY NOT NULL,
  string_field VARCHAR(50) NOT NULL,
  boolean_field BOOL NOT NULL,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE other_table(
  id INT PRIMARY KEY NOT NULL,
  int_field INT NOT NULL,
  boolean_field BOOL NOT NULL,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE join_table(
  some_id INT NOT NULL,
  other_id INT NOT NULL,
  PRIMARY KEY(some_id, other_id)
);

CREATE TABLE string_key_table(
  id VARCHAR(50) PRIMARY KEY NOT NULL,
  created_at DATETIME,
  updated_at DATETIME
);
`

var testSchemaPostgres = `
CREATE TABLE some_table(
  id INT PRIMARY KEY NOT NULL,
  string_field VARCHAR(50) NOT NULL,
  boolean_field BOOL NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE other_table(
  id INT PRIMARY KEY NOT NULL,
  int_field INT NOT NULL,
  boolean_field BOOL NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE join_table(
  some_id INT NOT NULL,
  other_id INT NOT NULL,
  PRIMARY KEY(some_id, other_id)
);

CREATE TABLE string_key_table(
  id VARCHAR(50) PRIMARY KEY NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE
);
`

var testData = `
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
- table: 'string_key_table'
  pk:
    id: 'new_id'
  fields:
    created_at: 'ON_INSERT_NOW()'
    updated_at: 'ON_UPDATE_NOW()'
`

var (
	fixtureFile  = "fixtures/test_fixtures1.yml"
	fixtureFiles = []string{
		"fixtures/test_fixtures1.yml",
		"fixtures/test_fixtures2.yml",
	}
)
