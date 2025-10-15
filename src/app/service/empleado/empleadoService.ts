// src/app/service/empleado/empleado.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';
import { Empleado } from '../../models/empleado.model';

@Injectable({
  providedIn: 'root'
})
export class EmpleadoService {
  private apiUrl = 'http://localhost:8000';

  constructor(private http: HttpClient) { }

  getAll(): Observable<Empleado[]> {
    return this.http.get<any>(`${this.apiUrl}/empleados/list`).pipe(
      tap(response => console.log('🔍 Respuesta completa del backend:', response)),
      map(response => {
        // Soportar varias formas que pueda devolver el backend
        const raw = response?.data ?? response?.empleados ?? response ?? [];
        const arr = Array.isArray(raw) ? raw : [];
        return arr.map((e: any) => ({
          id: e.id ?? e.ID,
          nombre: e.nombre ?? e.Nombre,
          apellido: e.apellido ?? e.Apellido,
          documentoTipo: e.documento_tipo ?? e.DocumentoTipo ?? e.documentoTipo,
          documentoNumero: e.documento_numero ?? e.DocumentoNumero ?? e.documentoNumero,
          correoElectronico: e.correo_electronico ?? e.CorreoElectronico ?? e.correoElectronico,
          ciudad: e.ciudad ?? e.Ciudad ?? e.ciudad,
          direccion: e.direccion ?? e.Direccion ?? e.direccion,
          telefono: e.telefono ?? e.Telefono ?? e.telefono,
          bolsaEmpleoID: e.bolsa_empleo_id ?? e.BolsaEmpleoID ?? e.bolsaEmpleoID ?? null,
          bolsaEmpleo: e.bolsaEmpleo ?? e.bolsa_empleo ?? e.BolsaEmpleo ?? null,
          fechaCreacion: e.fecha_creacion ?? e.FechaCreacion ?? e.fechaCreacion,
          fechaModificacion: e.fecha_modificacion ?? e.FechaModificacion ?? e.fechaModificacion
        }));
      })
    );
  }

  getById(id: number): Observable<Empleado> {
    return this.http.get<Empleado>(`${this.apiUrl}/empleados/${id}`);
  }

  create(empleado: any): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/empleados/create`, empleado).pipe(
      tap({
        next: res => console.log('✅ Empleado creado:', res),
        error: err => console.error('❌ Error creando empleado:', err)
      })
    );
  }

  update(empleado: any): Observable<any> {
    return this.http.put<any>(`${this.apiUrl}/empleados/update`, empleado);
  }

  delete(id: number): Observable<any> {
    return this.http.delete<any>(`${this.apiUrl}/empleados/delete`, { params: { id: id.toString() }});
  }
}
