package dao

import (
	"context"
	"retoBack/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type administradoresDao struct{}

var administradores = administradoresDao{}

// Obtener instancia del DAO
func Administradores() *administradoresDao {
	return &administradores
}

// Nombre de la tabla
func (a *administradoresDao) table() string {
	return "administradores"
}

// Obtener por ID
func (a *administradoresDao) GetById(ctx context.Context, id int) (*entity.Administrador, error) {
	var administrador entity.Administrador
	err := g.Model(a.table()).Where("id", id).Scan(&administrador)
	if err != nil {
		return nil, err
	}
	if administrador.ID == 0 {
		return nil, nil
	}
	return &administrador, nil
}

// Crear administrador
func (a *administradoresDao) Create(ctx context.Context, administrador *entity.Administrador) error {
	_, err := g.Model(a.table()).Data(administrador).Insert()
	return err
}

// Actualizar administrador
func (a *administradoresDao) Update(ctx context.Context, administrador *entity.Administrador) error {
	_, err := g.Model(a.table()).Data(administrador).Where("id", administrador.ID).Update()
	return err
}

// Eliminar administrador
func (a *administradoresDao) Delete(ctx context.Context, id int) error {
	_, err := g.Model(a.table()).Where("id", id).Delete()
	return err
}

// Obtener todos los administradores
func (a *administradoresDao) GetAll(ctx context.Context) ([]*entity.Administrador, error) {
	var administradores []*entity.Administrador
	err := g.Model(a.table()).Scan(&administradores)
	return administradores, err
}

// Obtener administrador por correo
func (a *administradoresDao) GetByCorreo(ctx context.Context, correo string) (*entity.Administrador, error) {
	var administrador entity.Administrador
	err := g.Model(a.table()).Where("correo", correo).Scan(&administrador)
	if err != nil {
		// Si el error es "no rows", retornamos nil sin error
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &administrador, nil
}
