package main

import (
	"github.com/rostamborn/snippetbox/pkg/models"
)

type templateData struct {
    Snippet *models.Snippet
    Snippets []*models.Snippet
}
