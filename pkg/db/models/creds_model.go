package models

import "encoding/json"

type Creds struct {
	ID      int             `gorm:"primaryKey;autoIncrement"`
	Name    string          `gorm:"not null,size:50"`
	Content json.RawMessage `gorm:"type:json;not null"`
}

func (c *Creds) TableName() string {
	return "creds"
}

func NewCreds(name string, content json.RawMessage) *Creds {
	return &Creds{Name: name, Content: content}
}
