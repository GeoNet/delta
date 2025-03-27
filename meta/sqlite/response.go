package sqlite

const responseCreate = `
DROP TABLE IF EXISTS response;
CREATE TABLE IF NOT EXISTS response (
  response_id INTEGER PRIMARY KEY NOT NULL,
  response TEXT NOT NULL,
  xml TEXT NOT NULL,
  UNIQUE (response)
);`

var response = Table{
	Create: responseCreate,
	Select: func() string {
		return "SELECT response_id FROM response WHERE response = ?"
	},
	Insert: func() string {
		return "INSERT INTO response (response, xml) VALUES (?, ?);"
	},
	Fields: []string{"Response", "XML"},
}
