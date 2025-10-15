package mapper

import (
	"fmt"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/entity"
	"strings"
	"time"
)

// ----------------------
// DTO ↔ Entity
// ----------------------

// ToEmpleadoDTO convierte una entidad Empleados a DTO
func ToEmpleadoDTO(e *entity.Empleados) *dto.EmpleadoDTO {
	if e == nil {
		return nil
	}

	bolsaPuesto := ""
	if e.BolsaEmpleo != nil {
		bolsaPuesto = e.BolsaEmpleo.Puesto
	}

	return &dto.EmpleadoDTO{
		ID:                e.ID,
		Nombre:            e.Nombre,
		Apellido:          e.Apellido,
		DocumentoTipo:     e.DocumentoTipo,
		DocumentoNumero:   e.DocumentoNumero,
		CorreoElectronico: e.CorreoElectronico,
		Ciudad:            e.Ciudad,
		Direccion:         e.Direccion,
		Telefono:          e.Telefono,
		BolsaEmpleoID:     e.BolsaEmpleoID,
		BolsaPuesto:       bolsaPuesto,
		FechaCreacion:     e.FechaCreacion,
		FechaModificacion: e.FechaModificacion,
	}
}

// ToEmpleadoDTOList convierte una lista de entidades a lista de DTOs
func ToEmpleadoDTOList(list []*entity.Empleados) []*dto.EmpleadoDTO {
	if list == nil {
		return nil
	}
	dtos := make([]*dto.EmpleadoDTO, len(list))
	for i, e := range list {
		dtos[i] = ToEmpleadoDTO(e)
	}
	return dtos
}

// ToEmpleadoEntityFromDTO convierte un DTO a entidad (para creación o actualización)
func ToEmpleadoEntityFromDTO(d *dto.EmpleadoDTO) *entity.Empleados {
	if d == nil {
		return nil
	}
	return &entity.Empleados{
		ID:                d.ID,
		Nombre:            strings.TrimSpace(d.Nombre),
		Apellido:          strings.TrimSpace(d.Apellido),
		DocumentoTipo:     strings.TrimSpace(d.DocumentoTipo),
		DocumentoNumero:   strings.TrimSpace(d.DocumentoNumero),
		CorreoElectronico: strings.ToLower(strings.TrimSpace(d.CorreoElectronico)),
		Ciudad:            strings.TrimSpace(d.Ciudad),
		Direccion:         strings.TrimSpace(d.Direccion),
		Telefono:          strings.TrimSpace(d.Telefono),
		BolsaEmpleoID:     maxZero(d.BolsaEmpleoID),
		FechaCreacion:     time.Now(),
		FechaModificacion: time.Now(),
		BolsaEmpleo:       nil, // Evitar problemas con GORM
	}
}

// ToEmpleadoEntityFromData crea una entidad desde datos simples
func ToEmpleadoEntityFromData(nombre, apellido, docTipo, docNumero, correo, ciudad, direccion, telefono string, bolsaID int) *entity.Empleados {
	return &entity.Empleados{
		Nombre:            strings.TrimSpace(nombre),
		Apellido:          strings.TrimSpace(apellido),
		DocumentoTipo:     strings.TrimSpace(docTipo),
		DocumentoNumero:   strings.TrimSpace(docNumero),
		CorreoElectronico: strings.ToLower(strings.TrimSpace(correo)),
		Ciudad:            strings.TrimSpace(ciudad),
		Direccion:         strings.TrimSpace(direccion),
		Telefono:          strings.TrimSpace(telefono),
		BolsaEmpleoID:     maxZero(bolsaID),
		FechaCreacion:     time.Now(),
		BolsaEmpleo:       nil,
	}
}

// ----------------------
// Sanitización y helpers
// ----------------------

// SanitizeEmpleadoDTO limpia y normaliza los campos de un DTO
func SanitizeEmpleadoDTO(d *dto.EmpleadoDTO) *dto.EmpleadoDTO {
	if d == nil {
		return nil
	}
	s := *d
	s.Nombre = strings.TrimSpace(d.Nombre)
	s.Apellido = strings.TrimSpace(d.Apellido)
	s.DocumentoTipo = strings.TrimSpace(d.DocumentoTipo)
	s.DocumentoNumero = strings.TrimSpace(d.DocumentoNumero)
	s.CorreoElectronico = strings.ToLower(strings.TrimSpace(d.CorreoElectronico))
	s.Ciudad = strings.TrimSpace(d.Ciudad)
	s.Direccion = strings.TrimSpace(d.Direccion)
	s.Telefono = strings.TrimSpace(d.Telefono)
	s.BolsaEmpleoID = maxZero(d.BolsaEmpleoID)
	return &s
}

// GetEmpleadoBasicInfo devuelve info resumida de un empleado
func GetEmpleadoBasicInfo(e *entity.Empleados) string {
	if e == nil {
		return "empleado nil"
	}
	bolsaInfo := "sin bolsa"
	if e.BolsaEmpleoID > 0 {
		bolsaInfo = fmt.Sprintf("bolsa:%d", e.BolsaEmpleoID)
	}
	return fmt.Sprintf("%s %s (ID:%d, Doc:%s, %s)", e.Nombre, e.Apellido, e.ID, e.DocumentoNumero, bolsaInfo)
}

// maxZero asegura que el int nunca sea negativo
func maxZero(i int) int {
	if i < 0 {
		return 0
	}
	return i
}
