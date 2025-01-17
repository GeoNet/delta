package sqlite

const placenameCreate = `
DROP TABLE IF EXISTS placename;
CREATE TABLE IF NOT EXISTS placename (
  placename_id INTEGER PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  level INTEGER NOT NULL,
  UNIQUE(name)
);
`

var placename = Table{
	Create: placenameCreate,
	Insert: func() string {
		return "INSERT INTO placename (name, latitude, longitude, level) VALUES (?, ?, ?, ?);"
	},
	Fields: []string{"Name", "Latitude", "Longitude", "Level"},
}

const citationCreate = `
DROP TABLE IF EXISTS citation;
CREATE TABLE IF NOT EXISTS citation (
  citation_id INTEGER PRIMARY KEY NOT NULL,
  key TEXT NOT NULL,
  author TEXT NOT NULL,
  year REAL NOT NULL,
  title TEXT NOT NULL,
  published TEXT NULL,
  volume TEXT NULL,
  pages TEXT NULL,
  doi TEXT NULL,
  link TEXT NULL,
  retrieved TEXT NULL,
  UNIQUE(key)
);
`

var citation = Table{
	Create: citationCreate,
	Select: func() string {
		return "SELECT citation_id FROM citation WHERE key = ?"
	},
	Insert: func() string {
		return "INSERT INTO citation (key, author, year, title, published, volume, pages, doi, link, retrieved) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	},
	Fields: []string{"Key", "Author", "Year", "Title", "Published", "Volume", "Pages", "DOI", "Link", "Retrieved"},
	Nulls:  []string{"Published", "Volume", "Pages", "DOI", "Link", "Retrieved"},
}
