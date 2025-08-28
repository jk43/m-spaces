package forms

import (
	"cms/models"

	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"github.com/rs/zerolog"
)

type FormElems = service.FormElems
type Forms = service.Forms

type CMSFormElems struct {
	DB     models.DatabaseRepo
	Logger *zerolog.Logger
	Casbin *utils.CasbinManager
}

func (a *CMSFormElems) FormFuncMap(fMap map[utils.Host]service.FormElemsFunc) {
	fMap["base"] = a.Base
	fMap["local.hotdev.com"] = a.LocalHotdevCom
}
