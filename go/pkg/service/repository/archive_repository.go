package repository

import (
	"l-semi-chat/pkg/domain"
)

type archiveRepository struct {
	SQLHandler SQLHandler
}

// ArchiveRepository archive repository
type ArchiveRepository interface {
	FindArchiveByThreadID(threadID string) (domain.Archive, error)
	StoreArchive(archiveID, path, password, threadID string, isPublic int) error
	UpdateArchive(threadID, password string, isPublic int) error
	DeleteArchive(threadID string) error

	FindUserByID(userID string) (domain.User, error)

	FindThreadByThreadID(threadID string) (domain.Thread, error)
}

// NewArchiveRepository create archive repository
func NewArchiveRepository(sh SQLHandler) ArchiveRepository {
	return &archiveRepository{
		SQLHandler: sh,
	}
}

func (ar *archiveRepository) FindArchiveByThreadID(threadID string) (archive domain.Archive, err error) {
	row := ar.SQLHandler.QueryRow(
		`SELECT archives.id, archives.path, archives.password, archives.is_public, archives.thread_id, threads.name, threads.description, threads.limit_users, threads.is_public, threads.created_at, threads.updated_at, threads.user_id, users.name, users.mail, users.image, users.profile
		FROM archives
		INNER JOIN threads
		ON archives.thread_id = threads.id
		INNER JOIN users
		ON users.id = threads.user_id
		WHERE thread_id=?`,
		threadID)
	if err = row.Scan(&archive.ID, &archive.Path, &archive.Password, &archive.IsPublic, &archive.Thread.ID, &archive.Thread.Name, &archive.Thread.Description, &archive.Thread.LimitUsers, &archive.Thread.IsPublic, &archive.Thread.CreatedAt, &archive.Thread.UpdatedAt, &archive.Thread.Admin.ID, &archive.Thread.Admin.Name, &archive.Thread.Admin.Mail, &archive.Thread.Admin.Image, &archive.Thread.Admin.Profile); err != nil {
		return archive, domain.InternalServerError(err)
	}
	return
}

func (ar *archiveRepository) StoreArchive(archiveID, path, password, threadID string, isPublic int) error {
	_, err := ar.SQLHandler.Execute(
		"INSERT INTO archives(id, path, password, is_public, thread_id) VALUES (?,?,?,?,?)",
		archiveID,
		path,
		password,
		isPublic,
		threadID,
	)
	return domain.InternalServerError(err)
}

func (ar *archiveRepository) UpdateArchive(threadID, password string, isPublic int) error {
	query := "UPDATE archives SET"
	var values []interface{}
	if password != "" {
		query += " password=?,"
		values = append(values, password)
	}
	query += " is_public=? WHERE thread_id=?"
	values = append(values, isPublic, threadID)

	_, err := ar.SQLHandler.Execute(query, values...)
	return domain.InternalServerError(err)
}

func (ar *archiveRepository) DeleteArchive(threadID string) error {
	_, err := ar.SQLHandler.Execute("DELETE FROM archives WHERE thread_id=?", threadID)
	return domain.InternalServerError(err)
}

func (ar *archiveRepository) FindThreadByThreadID(threadID string) (thread domain.Thread, err error) {
	row := ar.SQLHandler.QueryRow(
		`SELECT threads.id, threads.name, threads.description, threads.limit_users, threads.user_id, threads.is_public, threads.created_at, threads.updated_at, users.user_id, users.name, users.image, users.profile
		FROM threads
		INNER JOIN users
		ON users.id = threads.user_id
		WHERE threads.id = ?`,
		threadID,
	)
	if err = row.Scan(&thread.ID, &thread.Name, &thread.Description, &thread.LimitUsers, &thread.Admin.ID, &thread.IsPublic, &thread.CreatedAt, &thread.UpdatedAt, &thread.Admin.UserID, &thread.Admin.Name, &thread.Admin.Image, &thread.Admin.Profile); err != nil {
		return thread, domain.InternalServerError(err)
	}
	return
}

func (ar *archiveRepository) FindUserByID(userID string) (user domain.User, err error) {
	row := ar.SQLHandler.QueryRow("SELECT id, user_id, name, mail, image, profile FROM users WHERE id=?", userID)
	if err = row.Scan(&user.ID, &user.UserID, &user.Name, &user.Mail, &user.Image, &user.Profile); err != nil {
		return user, domain.InternalServerError(err)
	}
	return user, nil
}
