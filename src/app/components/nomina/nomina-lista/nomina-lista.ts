import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatButtonModule } from '@angular/material/button';
import { NominaService } from '../../../service/nomina/nominaservice';
import { Nomina } from '../../../models/nomina.model';
import { NominaFormComponent } from '../nomina-form/nomina-form';

@Component({
  selector: 'app-nomina-list',
  templateUrl: './nomina-lista.html',
  styleUrls: ['./nomina-lista.css'],
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    MatButtonModule,
    MatDialogModule
  ]
})
export class NominaListComponent implements OnInit {
  nominas: Nomina[] = [];
  nominasFiltradas: Nomina[] = [];
  filtroForm!: FormGroup;
  cargando = false;

  constructor(
    private fb: FormBuilder,
    private nominaService: NominaService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.filtroForm = this.fb.group({
      empleado: [''],
      fechaDesde: [''],
      fechaHasta: ['']
    });

    this.cargarNominas();

    // 🔄 Aplica filtros automáticamente al cambiar los campos
    this.filtroForm.valueChanges.subscribe(() => this.aplicarFiltro());
  }

  cargarNominas(): void {
    this.cargando = true;
    this.nominaService.getAll().subscribe({
      next: (response: any) => {
        let nominasArray: any[] = [];

        if (response && response.success && Array.isArray(response.data)) {
          nominasArray = response.data;
        } else if (Array.isArray(response)) {
          nominasArray = response;
        } else if (response?.data && Array.isArray(response.data)) {
          nominasArray = response.data;
        }

        this.nominas = nominasArray.map((n) => this.mapearNomina(n));
        this.nominasFiltradas = [...this.nominas];
        this.cargando = false;
      },
      error: (err) => {
        console.error('❌ Error cargando nóminas:', err);
        this.snackBar.open('Error cargando nóminas', 'Cerrar', { duration: 3000 });
        this.cargando = false;
      }
    });
  }

  /** ✅ Corrige y adapta los datos del backend */
  private mapearNomina(nominaBackend: any): Nomina {
    const fechaParseada = this.parsearFecha(
      nominaBackend.fechaPago ||
        nominaBackend.fecha_pago ||
        nominaBackend.FechaPago
    );

    // 🔍 Armar nombre completo del empleado según el formato que venga del backend
    let empleadoNombre = '';
    if (nominaBackend.empleadoNombre || nominaBackend.empleado_nombre) {
      empleadoNombre =
        nominaBackend.empleadoNombre || nominaBackend.empleado_nombre;
    } else if (nominaBackend.empleado) {
      const emp = nominaBackend.empleado;
      empleadoNombre = `${emp.nombre || ''} ${emp.apellido || ''}`.trim();
    }

    //  Obtener puesto (si existe)
    const empleadoPuesto =
      nominaBackend.empleadoPuesto ||
      nominaBackend.empleado_puesto ||
      nominaBackend.empleado?.bolsaEmpleo?.puesto ||
      '';

    return {
      id: nominaBackend.id || nominaBackend.ID || undefined,
      empleadoId:
        nominaBackend.empleadoId ||
        nominaBackend.empleado_id ||
        nominaBackend.EmpleadoId ||
        0,
      fechaPago: fechaParseada,
      salarioBase:
        nominaBackend.salarioBase ||
        nominaBackend.salario_base ||
        nominaBackend.SalarioBase ||
        0,
      horasExtras:
        nominaBackend.horasExtras ||
        nominaBackend.horas_extras ||
        nominaBackend.HorasExtras ||
        0,
      bonificaciones:
        nominaBackend.bonificaciones || nominaBackend.Bonificaciones || 0,
      deducciones:
        nominaBackend.deducciones || nominaBackend.Deducciones || 0,
      totalPago:
        nominaBackend.totalPago ||
        nominaBackend.total_pago ||
        nominaBackend.TotalPago ||
        0,
      empleadoNombre: empleadoNombre,
      empleadoPuesto: empleadoPuesto,
      diasFaltantes:
        nominaBackend.diasFaltantes ||
        nominaBackend.dias_faltantes ||
        nominaBackend.DiasFaltantes ||
        0
    };
  }

  private parsearFecha(fecha: any): Date | undefined {
    if (!fecha) return undefined;
    try {
      const d = new Date(fecha);
      return isNaN(d.getTime()) ? undefined : d;
    } catch {
      return undefined;
    }
  }

  // Convierte cualquier fecha (string o Date) a formato numérico YYYYMMDD
  // Versión mejorada que maneja zona horaria correctamente
private fechaADiaNumero(fecha: string | Date): number | undefined {
  if (!fecha) return undefined;

  try {
    let yyyy: number, mm: number, dd: number;

    if (typeof fecha === 'string') {
      // Extraemos YYYY, MM, DD directamente del string
      const match = fecha.match(/^(\d{4})-(\d{1,2})-(\d{1,2})/);
      if (!match) return undefined;

      yyyy = parseInt(match[1], 10);
      mm = parseInt(match[2], 10);
      dd = parseInt(match[3], 10);
    } else {
      // Para Date objects, creamos una fecha en UTC para evitar desplazamientos
      const utcDate = new Date(Date.UTC(
        fecha.getFullYear(),
        fecha.getMonth(),
        fecha.getDate()
      ));

      yyyy = utcDate.getUTCFullYear();
      mm = utcDate.getUTCMonth() + 1;
      dd = utcDate.getUTCDate();
    }

    return yyyy * 10000 + mm * 100 + dd;
  } catch (error) {
    console.error('Error procesando fecha:', fecha, error);
    return undefined;
  }
}

  aplicarFiltro(): void {
    const { empleado, fechaDesde, fechaHasta } = this.filtroForm.value;

    // DEBUG: Ver qué valores llegan
    console.log('Filtros aplicados:', { empleado, fechaDesde, fechaHasta });

    const desdeNum = fechaDesde ? this.fechaADiaNumero(fechaDesde) : undefined;
    const hastaNum = fechaHasta ? this.fechaADiaNumero(fechaHasta) : undefined;

    // DEBUG: Ver conversiones
    console.log('Conversiones de fecha:', { desdeNum, hastaNum });

    this.nominasFiltradas = this.nominas.filter((n) => {
      // Filtro por empleado
      const coincideEmpleado =
        !empleado ||
        n.empleadoNombre?.toLowerCase().includes(empleado.toLowerCase()) ||
        n.empleadoId?.toString().includes(empleado);

      // Filtro por fecha - usar n.fechaPago directamente (ya es Date)
      const fechaNum = n.fechaPago ? this.fechaADiaNumero(n.fechaPago) : undefined;

      // DEBUG por cada nómina
      console.log('Procesando nómina:', {
        empleado: n.empleadoNombre,
        fechaPago: n.fechaPago,
        fechaNum,
        coincideEmpleado
      });

      let coincideFecha = true;

      if (fechaNum !== undefined) {
        if (desdeNum !== undefined) {
          coincideFecha = fechaNum >= desdeNum;
        }
        if (hastaNum !== undefined) {
          coincideFecha = coincideFecha && (fechaNum <= hastaNum);
        }
      } else {
        // Si no hay fecha en la nómina, solo coincide si no hay filtros de fecha
        coincideFecha = (desdeNum === undefined && hastaNum === undefined);
      }

      const resultado = coincideEmpleado && coincideFecha;
      console.log('Resultado filtro:', resultado);

      return resultado;
    });

    console.log('Nominas filtradas final:', this.nominasFiltradas);
  }

  limpiarFiltros(): void {
    this.filtroForm.reset();
    this.nominasFiltradas = [...this.nominas];
  }

  eliminar(id?: number): void {
    if (!id) return;
    if (!confirm('¿Seguro que deseas eliminar esta nómina?')) return;

    this.nominaService.delete(id).subscribe({
      next: () => {
        this.snackBar.open('Nómina eliminada', 'Cerrar', { duration: 3000 });
        this.cargarNominas();
      },
      error: (err) => {
        const mensajeError = err.error?.error || 'Error eliminando nómina';
        this.snackBar.open(mensajeError, 'Cerrar', { duration: 3000 });
      }
    });
  }

  abrirFormulario(nomina?: Nomina): void {
    const dialogRef = this.dialog.open(NominaFormComponent, {
      width: '500px',
      data: nomina ? { ...nomina } : null
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) this.cargarNominas();
    });
  }

  editar(nomina: Nomina): void {
    this.abrirFormulario(nomina);
  }

  crear(): void {
    this.abrirFormulario();
  }

  getTotalPagado(): number {
    return this.nominasFiltradas.reduce((t, n) => t + (n.totalPago || 0), 0);
  }

  getPromedioSalario(): number {
    return this.nominasFiltradas.length > 0
      ? this.getTotalPagado() / this.nominasFiltradas.length
      : 0;
  }
}
