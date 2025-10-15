package service

import (
	"context"
	"retoBack/internal/dao"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/mapper"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

type administradoresService struct{}

var Administradores = administradoresService{}

// obtener por id (devuelve DTO)
func (a *administradoresService) GetById(ctx context.Context, id int) (*dto.AdministradorDTO, error) {
	adminEntity, err := dao.Administradores().GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToAdministradorDTO(adminEntity), nil
}

// crear (recibe DTO)
func (a *administradoresService) Create(ctx context.Context, adminDTO *dto.AdministradorDTO) error {
	existente, err := dao.Administradores().GetByCorreo(ctx, adminDTO.Correo)
	if err != nil {
		return err
	}
	if existente != nil {
		return gerror.New("El correo ya está registrado")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(adminDTO.Password), bcrypt.DefaultCost)
	adminDTO.Password = string(hashPassword)
	if err != nil {
		return err
	}
	// mapear DTO -> Entity
	adminEntity := mapper.ToAdministradorEntity(adminDTO)
	return dao.Administradores().Create(ctx, adminEntity)
}

// actualizar (recibe DTO)
func (a *administradoresService) Update(ctx context.Context, adminDTO *dto.AdministradorDTO) error {
	existente, err := dao.Administradores().GetByCorreo(ctx, adminDTO.Correo)
	if err != nil {
		return err
	}
	if existente != nil && existente.ID != adminDTO.ID {
		return gerror.New("El correo ya está registrado por otro administrador")
	}

	adminEntity := mapper.ToAdministradorEntity(adminDTO)
	return dao.Administradores().Update(ctx, adminEntity)
}

// eliminar
func (a *administradoresService) Delete(ctx context.Context, id int) error {
	return dao.Administradores().Delete(ctx, id)
}

// obtener todos (devuelve slice de DTOs)
func (a *administradoresService) GetAll(ctx context.Context) ([]*dto.AdministradorDTO, error) {
	adminEntities, err := dao.Administradores().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// mapear lista de entities -> lista de DTOs
	var result []*dto.AdministradorDTO
	for _, e := range adminEntities {
		result = append(result, mapper.ToAdministradorDTO(e))
	}
	return result, nil
}
func (a *administradoresService) Login(ctx context.Context, correo, password string) (*dto.AdministradorDTO, error) {
	// Buscar admin por correo
	adminEntity, err := dao.Administradores().GetByCorreo(ctx, correo)
	if err != nil {
		return nil, err
	}
	if adminEntity == nil {
		return nil, gerror.New("correo no encontrado")
	}

	// Comparar contraseña ingresada con la hasheada en DB
	if err := bcrypt.CompareHashAndPassword([]byte(adminEntity.Password), []byte(password)); err != nil {
		return nil, gerror.New("credenciales inválidas")
	}

	// Mapear Entity -> DTO (sin exponer password)
	adminDTO := mapper.ToAdministradorDTO(adminEntity)
	adminDTO.Password = "" // ⚠️ Nunca devolver password

	return adminDTO, nil
}
