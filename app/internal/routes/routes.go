package routes

import (
	"fmt"
)

const (
	Home               = "/"
	Login              = "/login"
	Logout             = "/logout"
	Discipline         = "/discipline"
	Material           = "/material"
	MaterialAttachment = Material + "/attachment"
)

type Route[T any] struct {
	Pattern string
	Link    func(T) string
}

var DisciplineID = Route[int64]{
	Pattern: Discipline + "/{id}",
	Link: func(id int64) string {
		return fmt.Sprintf("%s/%d", Discipline, id)
	},
}

var MaterialNew = Route[string]{
	Pattern: Material + "/new/{type}",
	Link: func(s string) string {
		return fmt.Sprintf("%s/new/%s", Material, s)
	},
}

var MaterialID = Route[string]{
	Pattern: Material + "/{id}",
	Link: func(s string) string {
		return fmt.Sprintf("%s/%s", Material, s)
	},
}

var MaterialMove = Route[string]{
	Pattern: Material + "/{id}/move",
	Link: func(s string) string {
		return fmt.Sprintf("%s/%s/move", Material, s)
	},
}

var MaterialAttachmentDelete = Route[string]{
	Pattern: MaterialAttachment + "/{id}/delete",
	Link: func(s string) string {
		return fmt.Sprintf("%s/%s/delete", MaterialAttachment, s)
	},
}

var MaterialEdit = Route[string]{
	Pattern: Material + "/edit/{id}",
	Link: func(s string) string {
		return fmt.Sprintf("%s/edit/%s", Material, s)
	},
}

var MaterialHide = Route[string]{
	Pattern: Material + "/hide/{id}",
	Link: func(s string) string {
		return fmt.Sprintf("%s/hide/%s", Material, s)
	},
}

var MaterialShow = Route[string]{
	Pattern: Material + "/show/{id}",
	Link: func(s string) string {
		return fmt.Sprintf("%s/show/%s", Material, s)
	},
}

var MaterialDelete = Route[string]{
	Pattern: Material + "/delete/{id}",
	Link: func(s string) string {
		return fmt.Sprintf("%s/delete/%s", Material, s)
	},
}
