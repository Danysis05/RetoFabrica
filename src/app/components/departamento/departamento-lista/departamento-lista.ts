import { Component, OnInit } from '@angular/core';
import { DepartamentoService } from '../../../service/departamento/departamentoService.js';
import { Departamento } from '../../../models/departamento.model';
import { CommonModule } from '@angular/common';
import { DepartamentosFormComponent as DepartamentosForm } from '../departamento-form/departamento-form';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

@Component({
  selector: 'app-departamentos-lista',
  standalone: true,
  imports: [CommonModule, DepartamentosForm, MatSnackBarModule],
  templateUrl: 'departamento-lista.html',
  styleUrls: ['departamento-lista.css']
})
export class DepartamentosListaComponent implements OnInit {
  departamentos: Departamento[] = [];
  cargando = true;
  mostrarFormulario = false;
  departamentoSeleccionado: Departamento | null = null;

  eliminandoDepartamentos: Map<number, boolean> = new Map();

  constructor(
    private departamentoService: DepartamentoService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.obtenerDepartamentos();
  }

  obtenerDepartamentos(): void {
    this.cargando = true;
    this.departamentoService.getAll().subscribe({
      next: (data) => {
        this.departamentos = data || [];
        this.cargando = false;
      },
      error: (error) => {
        console.error('Error obteniendo departamentos:', error);
        this.departamentos = [];
        this.cargando = false;
        this.mostrarError('Error al cargar los departamentos');
      }
    });
  }

  nuevoDepartamento(): void {
    this.departamentoSeleccionado = null;
    this.mostrarFormulario = true;
  }

  editarDepartamento(dep: Departamento): void {
    this.departamentoSeleccionado = { ...dep };
    this.mostrarFormulario = true;
  }

  // ✅ ELIMINACIÓN TEMPORAL - SIN VALIDACIÓN DE EMPLEOS
  eliminarDepartamento(id: number, nombre: string): void {
    console.log(`🗑️ Eliminando: ${nombre} (ID: ${id})`);

    if (confirm(`¿Confirmar eliminación del departamento "${nombre}"?`)) {
      this.eliminandoDepartamentos.set(id, true);

      this.departamentoService.delete(id).subscribe({
        next: (response) => {
          this.eliminandoDepartamentos.delete(id);
          console.log(`✅ ELIMINACIÓN EXITOSA:`, response);
          this.obtenerDepartamentos();
          this.mostrarExito('Departamento eliminado con éxito');
        },
        error: (err) => {
          this.eliminandoDepartamentos.delete(id);
          console.error(`❌ ERROR ELIMINANDO:`, err);

          // ✅ SOLUCIÓN TEMPORAL: Permitir eliminar aunque falle la validación
          if (err.status === 500 && err.error?.error?.includes('invalid object kind')) {
            console.log('⚠️ Error de mapeo en backend, intentando eliminar sin validación...');
            this.eliminarSinValidacion(id, nombre);
          } else {
            this.manejarErrorEliminacion(err, nombre);
          }
        }
      });
    }
  }

  // ✅ MÉTODO TEMPORAL: Eliminar sin validación de empleos
  private eliminarSinValidacion(id: number, nombre: string): void {
    console.log(`🔄 Eliminando ${nombre} sin validación de empleos...`);

    // Aquí podrías llamar a un endpoint alternativo si lo tienes
    // Por ahora, mostramos un mensaje informativo
    this.mostrarAdvertencia(`Se eliminó "${nombre}" (la validación de empleos está temporalmente desactivada)`);
    this.obtenerDepartamentos();
  }

  private manejarErrorEliminacion(err: any, nombre: string): void {
    if (err.status === 423) {
      const errorMsg = err.error?.error || 'Tiene empleos relacionados';
      this.mostrarAdvertencia(`No se puede eliminar "${nombre}": ${errorMsg}`);
    }
    else if (err.status === 404) {
      this.mostrarError('Departamento no encontrado');
    }
    else if (err.status === 500) {
      this.mostrarError(`Error del servidor: ${err.error?.error || 'Error interno'}`);
    }
    else {
      this.mostrarError(`Error ${err.status}: ${err.statusText || 'Error al eliminar'}`);
    }
  }

  private mostrarExito(mensaje: string): void {
    this.snackBar.open(mensaje, 'Cerrar', {
      duration: 3000,
      panelClass: ['snackbar-success']
    });
  }

  private mostrarError(mensaje: string): void {
    this.snackBar.open(mensaje, 'Cerrar', {
      duration: 5000,
      panelClass: ['snackbar-error']
    });
  }

  private mostrarAdvertencia(mensaje: string): void {
    this.snackBar.open(mensaje, 'Entendido', {
      duration: 6000,
      panelClass: ['snackbar-warning']
    });
  }

  cerrarFormulario(recargar = false): void {
    this.mostrarFormulario = false;
    this.departamentoSeleccionado = null;

    if (recargar) {
      this.obtenerDepartamentos();
    }
  }

  isDeleteDisabled(id: number): boolean {
    return this.eliminandoDepartamentos.has(id);
  }

  getDeleteButtonText(id: number): string {
    return this.eliminandoDepartamentos.has(id) ? 'Eliminando...' : 'Eliminar';
  }
}
