package mapper

import (
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/entity"
)

// Entity -> DTO
func ToAdministradorDTO(admin *entity.Administrador) *dto.AdministradorDTO {
	if admin == nil {
		return nil
	}
	return &dto.AdministradorDTO{
		ID:                admin.ID,
		Nombre:            admin.Nombre,
		Correo:            admin.Correo,
		FechaCreacion:     admin.FechaCreacion,
		FechaModificacion: admin.FechaModificacion,
	}
}

// DTO -> Entity
func ToAdministradorEntity(dto *dto.AdministradorDTO) *entity.Administrador {
	if dto == nil {
		return nil
	}
	return &entity.Administrador{
		ID:                dto.ID,
		Nombre:            dto.Nombre,
		Correo:            dto.Correo,
		Password:          dto.Password,
		FechaCreacion:     dto.FechaCreacion,
		FechaModificacion: dto.FechaModificacion,
	}
}
