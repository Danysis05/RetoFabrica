import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { BolsaEmpleoService } from '../../../service/bolsa-empleo/bolsa-empleo.service.ts';
import { DepartamentoService } from '../../../service/departamento/departamentoService.js';
import { BolsaEmpleo } from '../../../models/bolsa.empleo';
import { Departamento } from '../../../models/departamento.model';
import { BolsaEmpleadoFormComponent } from '../bolsa-empleado-form/bolsa-empleado-form.js';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

@Component({
  selector: 'app-bolsa-empleado-lista',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    BolsaEmpleadoFormComponent,
    MatSnackBarModule
  ],
  templateUrl: './bolsa-empleado-lista.html',
  styleUrls: ['./bolsa-empleado-lista.css']
})
export class BolsaEmpleadoListaComponent implements OnInit {
  empleos: BolsaEmpleo[] = [];
  empleosFiltrados: BolsaEmpleo[] = [];
  departamentos: Departamento[] = [];
  cargando = true;
  cargandoDepartamentos = true;
  mostrarFormulario = false;
  empleoSeleccionado: BolsaEmpleo | null = null;

  // Filtros
  mostrarSoloDisponibles = false;
  departamentoFiltro: number | null = null;

  constructor(
    private bolsaEmpleoService: BolsaEmpleoService,
    private departamentoService: DepartamentoService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.obtenerEmpleos();
    this.obtenerDepartamentos();
  }

  obtenerEmpleos(): void {
    this.cargando = true;
    this.bolsaEmpleoService.getAll().subscribe({
      next: (data: BolsaEmpleo[]) => {
        this.empleos = data || [];
        this.aplicarFiltros();
        this.cargando = false;
      },
      error: (error: any) => {
        console.error('Error obteniendo empleos:', error);
        this.empleos = [];
        this.empleosFiltrados = [];
        this.cargando = false;
        this.snackBar.open('Error al cargar los empleos', 'Cerrar', {
          duration: 3000,
          panelClass: ['snackbar-error']
        });
      }
    });
  }

  obtenerDepartamentos(): void {
    this.cargandoDepartamentos = true;
    this.departamentoService.getAll().subscribe({
      next: (deps: Departamento[]) => {
        this.departamentos = deps || [];
        this.cargandoDepartamentos = false;
      },
      error: (error: any) => {
        console.error('Error cargando departamentos:', error);
        this.departamentos = [];
        this.cargandoDepartamentos = false;
      }
    });
  }

  aplicarFiltros(): void {
    let empleosFiltrados = [...this.empleos];

    if (this.mostrarSoloDisponibles) {
      empleosFiltrados = empleosFiltrados.filter(empleo =>
        empleo.estado === 'DISPONIBLE'
      );
    }

    if (this.departamentoFiltro) {
      empleosFiltrados = empleosFiltrados.filter(empleo =>
        empleo.departamento_id === this.departamentoFiltro
      );
    }

    this.empleosFiltrados = empleosFiltrados;
  }

  toggleFiltroDisponibles(): void {
    this.mostrarSoloDisponibles = !this.mostrarSoloDisponibles;
    this.aplicarFiltros();
  }

  cambiarDepartamentoFiltro(departamentoId: string): void {
    const id = departamentoId ? parseInt(departamentoId) : null;
    this.departamentoFiltro = isNaN(id as number) ? null : id;
    this.aplicarFiltros();
  }

  limpiarFiltros(): void {
    this.mostrarSoloDisponibles = false;
    this.departamentoFiltro = null;
    this.aplicarFiltros();
  }

  getNombreDepartamento(departamentoId: number): string {
    if (!departamentoId || !this.departamentos || this.departamentos.length === 0) {
      return `ID: ${departamentoId || 'N/A'}`;
    }
    const dep = this.departamentos.find(d => d.id === departamentoId);
    return dep ? dep.nombre : `ID: ${departamentoId}`;
  }

  nuevoEmpleo(): void {
    this.empleoSeleccionado = null;
    this.mostrarFormulario = true;
  }

  editarEmpleo(empleo: BolsaEmpleo): void {
    this.empleoSeleccionado = { ...empleo };
    this.mostrarFormulario = true;
  }
  ultimoError: string | null = null;

  eliminarEmpleo(id: number): void {
    if (!id) {
      this.snackBar.open('Error: ID de empleo no válido', 'Cerrar', {
        duration: 3000,
        panelClass: ['snackbar-error']
      });
      return;
    }

    if (confirm('¿Seguro que deseas eliminar este empleo de la bolsa?')) {
      this.bolsaEmpleoService.delete(id).subscribe({
        next: (res: any) => {
          // ✅ Si el backend responde con success=false y mensaje
          if (res?.success === false) {
            const msg = res?.error || 'No se pudo eliminar el empleo';
            this.snackBar.open(msg, 'Cerrar', {
              duration: 4000,
              panelClass: ['snackbar-error']
            });
            return;
          }

          // ✅ Eliminación exitosa
          this.obtenerEmpleos();
          this.snackBar.open('Empleo eliminado con éxito', 'Cerrar', {
            duration: 3000,
            panelClass: ['snackbar-success']
          });
        },
        error: (err: any) => {
          console.error('Error eliminando empleo:', err);

          // ✅ Si el backend envía un mensaje específico (como el del empleado asignado)
          const msg = err?.error?.error || 'Error al eliminar el empleo';
          this.snackBar.open(msg, 'Cerrar', {
            duration: 4000,
            panelClass: ['snackbar-error']
          });
        }
      });
    }
  }

  cerrarFormulario(recargar = false): void {
    this.mostrarFormulario = false;
    this.empleoSeleccionado = null;

    if (recargar) {
      this.obtenerEmpleos();

      if (this.empleoSeleccionado) {
        this.snackBar.open('Empleo actualizado con éxito', 'Cerrar', {
          duration: 3000,
          panelClass: ['snackbar-success']
        });
      } else {
        this.snackBar.open('Empleo creado con éxito', 'Cerrar', {
          duration: 3000,
          panelClass: ['snackbar-success']
        });
      }
    }
  }

  getEstadoClass(estado: string): string {
    switch (estado) {
      case 'DISPONIBLE': return 'estado-disponible';
      case 'OCUPADO': return 'estado-ocupado';
      case 'CERRADO': return 'estado-cerrado';
      default: return 'estado-default';
    }
  }

  trackByEmpleoId(index: number, empleo: BolsaEmpleo): number {
    return empleo.id || index;
  }

  getEstadisticas(): any {
    const total = this.empleos.length;
    const disponibles = this.empleos.filter(e => e.estado === 'DISPONIBLE').length;
    const ocupados = this.empleos.filter(e => e.estado === 'OCUPADO').length;
    const cerrados = this.empleos.filter(e => e.estado === 'CERRADO').length;

    return { total, disponibles, ocupados, cerrados };
  }

  get hayFiltrosActivos(): boolean {
    return this.mostrarSoloDisponibles || this.departamentoFiltro !== null;
  }

  getEmpleosDisponiblesPorDepartamento(departamentoId: number): number {
    if (!departamentoId) return 0;
    return this.empleos.filter(empleo =>
      empleo.departamento_id === departamentoId &&
      empleo.estado === 'DISPONIBLE'
    ).length;
  }
}
