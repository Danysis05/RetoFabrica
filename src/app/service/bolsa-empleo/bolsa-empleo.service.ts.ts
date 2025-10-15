import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { BolsaEmpleo } from '../../models/bolsa.empleo';

@Injectable({
  providedIn: 'root'
})
export class BolsaEmpleoService {
  private apiUrl = 'http://localhost:8000/bolsa';

  constructor(private http: HttpClient) { }

  getAll(): Observable<BolsaEmpleo[]> {
    return this.http.get<BolsaEmpleo[]>(`${this.apiUrl}/`);
  }

  getByDepartamento(departamentoId: number): Observable<BolsaEmpleo[]> {
    return this.http.get<BolsaEmpleo[]>(`${this.apiUrl}/departamento/${departamentoId}`);
  }

  create(empleo: BolsaEmpleo): Observable<BolsaEmpleo> {
    console.log('🚀 CREANDO EMPLEO - Datos recibidos:', empleo);
    console.log('📦 Enviando al backend:', empleo);

    return this.http.post<BolsaEmpleo>(`${this.apiUrl}/`, empleo);
  }

  update(id: number, empleo: BolsaEmpleo): Observable<BolsaEmpleo> {
    return this.http.put<BolsaEmpleo>(`${this.apiUrl}/${id}`, empleo);
  }

  delete(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }

  getDisponibles(): Observable<BolsaEmpleo[]> {
    return this.http.get<BolsaEmpleo[]>(`${this.apiUrl}/?estado=DISPONIBLE`);
  }

  getDisponiblesByDepartamento(departamentoId: number): Observable<BolsaEmpleo[]> {
    return this.http.get<BolsaEmpleo[]>(`${this.apiUrl}/departamento/${departamentoId}?estado=DISPONIBLE`);
  }
}
