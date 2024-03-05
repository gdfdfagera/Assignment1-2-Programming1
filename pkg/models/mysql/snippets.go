package mysql

import (
	"database/sql"
	"errors"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.02 (Project Setup and Enabling
	// Modules) so that the import statement looks like this:
	// "{your-module-path}/pkg/models".
	"snippetbox.disa.net/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, author, year, genre, annotation string) (int, error) {
	stmt := `INSERT INTO book (title, author, year, genre, annotation, created)
	VALUES(?, ?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, author, year, genre, annotation)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT book_id, title, author, year, genre, annotation, created FROM book
WHERE book_id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Author, &s.Year, &s.Genre, &s.Annotation, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT book_id, title, author, year, genre, annotation, created FROM book
	ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Author, &s.Year, &s.Genre, &s.Annotation, &s.Created)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) Update(id int, title, author, year, genre, annotation string) error {
	stmt := `UPDATE book SET title=?, author=?, year=?, genre=?, annotation=?, created=UTC_TIMESTAMP() WHERE book_id=?`

	result, err := m.DB.Exec(stmt, title, author, year, genre, annotation, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNoRecord // Indicates that no record was updated (ID not found)
	}

	return nil
}

func (m *SnippetModel) Delete(id int) error {
	stmt := `DELETE FROM book WHERE book_id = ?`

	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return nil
}
