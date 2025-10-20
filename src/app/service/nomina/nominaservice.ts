// 📁 service/nomina/nominaservice.ts - CORREGIDO
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Nomina } from '../../models/nomina.model';

@Injectable({
  providedIn: 'root'
})
export class NominaService {
  // ✅ CORREGIDO: Puerto 8000 y ruta /nomina (sin 's')
  private apiUrl = 'http://localhost:8000/nomina';

  constructor(private http: HttpClient) {}

  // ✅ CORREGIDO: Rutas que coinciden con el backend
  getAll(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/all`);
  }

  getById(id: number): Observable<any> {
    // ✅ CORREGIDO: Usar formato :id en la URL
    return this.http.get<any>(`${this.apiUrl}/${id}`);
  }

  create(nominaData: any): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/create`, nominaData);
  }

  update(nominaData: any): Observable<any> {
    // ✅ CORREGIDO: Usar el formato :id en la URL
    return this.http.put<any>(`${this.apiUrl}/${nominaData.id}`, nominaData);
  }

  save(nominaData: any): Observable<any> {
    if (nominaData.id) {
      return this.update(nominaData);
    } else {
      return this.create(nominaData);
    }
  }

  delete(id: number): Observable<any> {
    // ✅ CORREGIDO: Usar el formato :id en la URL
    return this.http.delete<any>(`${this.apiUrl}/${id}`);
  }
    getEmpleadoInfo(empleadoId: number): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/empleado-info?empleadoId=${empleadoId}`);
  }
}
