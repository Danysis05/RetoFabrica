import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { ApiResponse, Departamento } from '../../models/departamento.model';

@Injectable({
  providedIn: 'root'
})
export class DepartamentoService {
  private baseUrl = 'http://localhost:8000/departamentos';

  constructor(private http: HttpClient) {}

  // ✅ MÉTODOS EXISTENTES (NO MODIFICAR)
  getAll(): Observable<Departamento[]> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/list`)
      .pipe(
        map(res => res.data)
      );
  }

  getById(id: number): Observable<Departamento> {
    return this.http.get<Departamento>(`${this.baseUrl}/${id}`);
  }

  create(dep: Departamento): Observable<any> {
    console.log('📤 CREATE - Enviando:', dep);
    return this.http.post(`${this.baseUrl}/create`, dep);
  }

  update(id: number, dep: Departamento): Observable<any> {
    const url = `${this.baseUrl}/update/${id}`;
    console.log('📤 UPDATE - URL:', url);
    console.log('📦 UPDATE - Datos:', dep);
    return this.http.put(url, dep);
  }

  delete(id: number): Observable<void> {
    const url = `${this.baseUrl}/delete/${id}`;
    console.log('🗑️ DELETE - URL:', url);
    return this.http.delete<void>(url);
  }

  // ✅ NUEVOS MÉTODOS PARA LA VALIDACIÓN PROFESIONAL
  getAllWithEmpleos(): Observable<Departamento[]> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/with-empleos`)
      .pipe(
        map(res => res.data || [])
      );
  }

  canDelete(id: number): Observable<{canDelete: boolean; reason: string}> {
    const url = `${this.baseUrl}/can-delete/${id}`;
    console.log('🔍 CAN-DELETE - URL:', url);

    return this.http.get<{success: boolean; canDelete: boolean; reason: string}>(url)
      .pipe(
        map(response => ({
          canDelete: response.canDelete,
          reason: response.reason || 'Error desconocido'
        }))
      );
  }

  // ✅ CORREGIDO: Este endpoint probablemente no existe, así que lo comentamos
  // getWithEmpleos(id: number): Observable<Departamento> {
  //   const url = `${this.baseUrl}/${id}/with-empleos`;
  //   console.log('📋 GET WITH EMPLEOS - URL:', url);

  //   return this.http.get<ApiResponse>(url)
  //     .pipe(
  //       map(res => res.data)
  //     );
  // }
}
