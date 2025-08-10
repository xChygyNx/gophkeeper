package entity

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/xChygyNx/gophkeeper/internal/server/database"
	"github.com/xChygyNx/gophkeeper/internal/server/model"
	custom_errors "github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

type Entity struct {
	db *database.DB
}

func New(db *database.DB) *Entity {
	return &Entity{
		db: db,
	}
}

func (e *Entity) Create(entityRequest *model.CreateEntityRequest) (int64, error) {
	var id int64
	metadata := model.MetadataEntity{Name: entityRequest.Metadata.Name, Description: entityRequest.Metadata.Description, Type: entityRequest.Metadata.Type}
	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return 0, err
	}
	if err = e.db.Pool.QueryRow(
		"INSERT INTO entity (user_id, data, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"+
			" RETURNING entity_id",
		entityRequest.UserID,
		entityRequest.Data,
		jsonMetadata,
		time.Now(),
		time.Now(),
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("error in insert entity for UserID %d in DB: %w", entityRequest.UserID, err)
	}
	return id, nil
}

func (e *Entity) GetList(userID int64, typeEntity string) ([]model.Entity, error) {
	entities := make([]model.Entity, 0)
	rows, err := e.db.Pool.Query("SELECT entity_id, user_id, data, metadata, created_at, updated_at FROM entity "+
		"where user_id = $1 and metadata->>'Type' = $2 and deleted_at IS NULL",
		userID, typeEntity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities, custom_errors.ErrRecordNotFound
		} else {
			return entities, fmt.Errorf("error in get data for UserID %d from DB: %w", userID, err)
		}
	}
	defer func() {
		err = rows.Close()
	}()
	for rows.Next() {
		entity := model.Entity{}

		var jsonEntity string
		err = rows.Scan(&entity.ID, &entity.UserID, &entity.Data, &jsonEntity, &entity.CreatedAt, &entity.UpdatedAt)
		if err != nil {
			return entities, fmt.Errorf("error in scan entity record from DB: %w", err)
		}

		err = json.Unmarshal([]byte(jsonEntity), &entity.Metadata)
		if err != nil {
			return entities, fmt.Errorf("error in unmarshal entity metadata: %w", err)
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

func (e *Entity) Exists(entityRequest *model.CreateEntityRequest) (bool, error) {
	var exists bool
	row := e.db.Pool.QueryRow("SELECT EXISTS(SELECT 1 FROM entity "+
		"where entity.user_id = $1 and entity.metadata->>'Name' = $2 and "+
		"entity.metadata->>'Type' = $3 and entity.deleted_at IS NULL)",
		entityRequest.UserID, entityRequest.Metadata.Name, entityRequest.Metadata.Type)
	if err := row.Scan(&exists); err != nil {
		return exists, fmt.Errorf("error in check entity exists in DB: %w", err)
	}
	return exists, nil
}

func (e *Entity) Delete(userID int64, name string, typeEntity string) (int64, error) {
	var id int64
	if err := e.db.Pool.QueryRow("UPDATE entity SET deleted_at = $1 "+
		"where entity.user_id = $2 and entity.metadata->>'Name' = $3 and "+
		"entity.metadata->>'Type' = $4 RETURNING entity_id",
		time.Now(),
		userID,
		name,
		typeEntity,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("error in mark entity as deleted in DB: %w", err)
	}
	return id, nil
}

func (e *Entity) Update(userID int64, name string, typeEntity string, data []byte) (int64, error) {
	var id int64
	if err := e.db.Pool.QueryRow("UPDATE entity SET data = $1, updated_at = $2 "+
		"where entity.user_id = $3 and entity.metadata->>'Name' = $4 "+
		"and entity.metadata->>'Type' = $5 RETURNING entity_id",
		data,
		time.Now(),
		userID,
		name,
		typeEntity,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("error in mark entity as deleted in DB: %w", err)
	}
	return id, nil
}
