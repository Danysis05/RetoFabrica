import { Component, OnInit, Inject } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule, CurrencyPipe } from '@angular/common';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';

import { NominaService } from '../../../service/nomina/nominaservice';
import { EmpleadoService } from '../../../service/empleado/empleadoService';
import { Nomina } from '../../../models/nomina.model';
import { Empleado } from '../../../models/empleado.model';

@Component({
  selector: 'app-nomina-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, CurrencyPipe],
  templateUrl: './nomina-form.html',
  styleUrls: ['./nomina-form.css']
})
export class NominaFormComponent implements OnInit {
  form!: FormGroup;
  nomina?: Nomina;
  empleados: Empleado[] = [];
  guardando = false;

  salarioBaseCalculado = 0;
  bonificacionesCalculadas = 0;
  deduccionesCalculadas = 0;
  totalPagoCalculado = 0;

  // ✅ CORREGIDO: Mejor nombre para la variable
  cargandoSalario = false;
  infoEmpleadoSeleccionado: any = null;

  constructor(
    private fb: FormBuilder,
    private nominaService: NominaService,
    private empleadoService: EmpleadoService,
    private snackBar: MatSnackBar,
    private dialogRef: MatDialogRef<NominaFormComponent>,
    @Inject(MAT_DIALOG_DATA) public data?: Nomina
  ) {
    this.nomina = data;
  }

  ngOnInit(): void {
    this.form = this.fb.group({
      empleadoId: [this.nomina?.empleadoId ?? null, Validators.required],
      horasExtras: [this.nomina?.horasExtras ?? 0, [Validators.min(0)]],
      diasFaltantes: [this.nomina?.diasFaltantes ?? 0, [Validators.min(0)]]
    });

    this.salarioBaseCalculado = this.nomina?.salarioBase ?? 0;
    this.bonificacionesCalculadas = this.nomina?.bonificaciones ?? 0;
    this.deduccionesCalculadas = this.nomina?.deducciones ?? 0;
    this.totalPagoCalculado = this.nomina?.totalPago ?? 0;

    this.cargarEmpleados();

    // ✅ MANTENIDO: Cargar información completa cuando cambia el empleado
    this.form.get('empleadoId')?.valueChanges.subscribe((id) => {
      const empleadoId = Number(id);
      if (empleadoId) {
        this.cargarInformacionEmpleadoCompleta(empleadoId);
      } else {
        this.limpiarInformacionEmpleado();
      }
    });

    this.form.get('horasExtras')?.valueChanges.subscribe(() => this.recalcularTotal());
    this.form.get('diasFaltantes')?.valueChanges.subscribe(() => this.recalcularTotal());
  }

  // ✅ AGREGADO: Método para obtener el puesto en la lista de empleados
  getPuestoEnLista(empleado: Empleado): string {
    // Intentar obtener el puesto de diferentes formas
    if (empleado.bolsaEmpleo && empleado.bolsaEmpleo.puesto) {
      return empleado.bolsaEmpleo.puesto;
    }

    // Usar cualquier propiedad alternativa que pueda tener el puesto
    const anyEmpleado = empleado as any;
    if (anyEmpleado.bolsaPuesto) {
      return anyEmpleado.bolsaPuesto;
    }
    if (anyEmpleado.puesto) {
      return anyEmpleado.puesto;
    }
    if (anyEmpleado.cargo) {
      return anyEmpleado.cargo;
    }

    return 'SIN PUESTO';
  }

  // ✅ CORREGIDO: Método mejorado para cargar información completa del empleado
  cargarInformacionEmpleadoCompleta(empleadoId: number): void {
    this.cargandoSalario = true;
    this.infoEmpleadoSeleccionado = null;

    console.log(`🔍 Cargando información completa para empleado ID: ${empleadoId}`);

    this.nominaService.getEmpleadoInfo(empleadoId).subscribe({
      next: (response) => {
        this.cargandoSalario = false;

        if (response.success && response.data) {
          this.infoEmpleadoSeleccionado = response.data;
          this.salarioBaseCalculado = response.data.salarioBase;

          console.log(`✅ Información del empleado cargada:`, response.data);
          console.log(`💰 Salario: $${this.salarioBaseCalculado}`);
          console.log(`💼 Puesto: ${response.data.puesto}`);
          console.log(`📊 Tiene bolsa activa: ${response.data.tieneBolsaActiva}`);

          // Mostrar advertencia si no tiene puesto
          if (!response.data.tieneBolsaActiva || !response.data.puesto) {
            this.snackBar.open(
              `⚠️ ${response.data.nombre} ${response.data.apellido} no tiene un puesto activo con salario`,
              'Cerrar',
              { duration: 5000 }
            );
          }

          this.recalcularTotal();
        } else {
          console.error('❌ Respuesta inválida del servidor:', response);
          this.usarDatosLocalesComoFallback(empleadoId);
        }
      },
      error: (err) => {
        this.cargandoSalario = false;
        console.error('❌ Error cargando información del empleado:', err);
        this.usarDatosLocalesComoFallback(empleadoId);
      }
    });
  }

  // ✅ NUEVO: Método para usar datos locales como fallback
  usarDatosLocalesComoFallback(empleadoId: number): void {
    const empleado = this.empleados.find(e => e.id === empleadoId);
    if (empleado) {
      // Intentar obtener salario y puesto de diferentes formas
      const salario = empleado.bolsaEmpleo?.salario ||
                     (empleado as any)?.bolsaSalario ||
                     (empleado as any)?.salarioBase || 0;

      const puesto = empleado.bolsaEmpleo?.puesto ||
                    (empleado as any)?.bolsaPuesto ||
                    (empleado as any)?.puesto ||
                    'No asignado';

      this.salarioBaseCalculado = salario;

      // Crear objeto de información local
      this.infoEmpleadoSeleccionado = {
        nombre: empleado.nombre,
        apellido: empleado.apellido,
        puesto: puesto,
        salarioBase: salario,
        tieneBolsaActiva: salario > 0 && puesto !== 'No asignado'
      };

      console.log(`🔄 Usando datos locales - Puesto: ${puesto}, Salario: $${salario}`);

      if (salario === 0) {
        this.snackBar.open('Error cargando información del empleado, usando datos locales', 'Cerrar', { duration: 3000 });
      }

      this.recalcularTotal();
    } else {
      this.salarioBaseCalculado = 0;
      this.recalcularTotal();
    }
  }

  // ✅ NUEVO: Limpiar información del empleado
  limpiarInformacionEmpleado(): void {
    this.salarioBaseCalculado = 0;
    this.infoEmpleadoSeleccionado = null;
    this.bonificacionesCalculadas = 0;
    this.deduccionesCalculadas = 0;
    this.totalPagoCalculado = 0;
  }

  cargarEmpleados(): void {
    this.empleadoService.getAll().subscribe({
      next: (data) => {
        this.empleados = data;
        console.log('👤 Empleados cargados:', data);

        // ✅ DEBUG: Verificar estructura de los empleados
        this.empleados.forEach((emp, index) => {
          console.log(`Empleado ${index + 1}:`, {
            id: emp.id,
            nombre: emp.nombre,
            apellido: emp.apellido,
            bolsaEmpleo: emp.bolsaEmpleo,
            tieneBolsaEmpleo: !!emp.bolsaEmpleo,
            puesto: emp.bolsaEmpleo?.puesto,
            salario: emp.bolsaEmpleo?.salario
          });
        });

        this.form.updateValueAndValidity();

        // Si estamos editando una nómina existente, cargar la información
        if (this.nomina?.empleadoId) {
          this.cargarInformacionEmpleadoCompleta(this.nomina.empleadoId);
        }
      },
      error: (err) => {
        console.error('❌ Error cargando empleados:', err);
        this.snackBar.open('Error cargando empleados', 'Cerrar', { duration: 3000 });
      }
    });
  }

  recalcularTotal(): void {
    const salario = this.salarioBaseCalculado || 0;
    const horasExtras = this.form.get('horasExtras')?.value || 0;
    const diasFaltantes = this.form.get('diasFaltantes')?.value || 0;

    const valorHora = salario / 240;
    this.bonificacionesCalculadas = horasExtras * valorHora * 1.5;
    this.deduccionesCalculadas = diasFaltantes * (salario / 30);
    this.totalPagoCalculado = salario + this.bonificacionesCalculadas - this.deduccionesCalculadas;

    console.log(`🧮 Cálculos actualizados - Salario: $${salario}, Total: $${this.totalPagoCalculado}`);
  }

  guardar(): void {
    if (this.form.invalid) {
      this.snackBar.open('Por favor complete todos los campos requeridos', 'Cerrar', { duration: 3000 });
      return;
    }

    // Validar que el empleado tenga salario
    if (this.salarioBaseCalculado === 0) {
      this.snackBar.open('El empleado seleccionado no tiene un salario asignado', 'Cerrar', { duration: 5000 });
      return;
    }

    this.guardando = true;

    const payload: any = {
      empleado_id: Number(this.form.get('empleadoId')?.value),
      horas_extras: this.form.get('horasExtras')?.value,
      dias_faltantes: this.form.get('diasFaltantes')?.value
    };

    if (this.nomina?.id) {
      payload.id = this.nomina.id;
    }

    console.log('📤 Enviando payload al backend:', payload);

    this.nominaService.save(payload).subscribe({
      next: (response) => {
        this.guardando = false;
        console.log('✅ Respuesta del backend:', response);
        this.snackBar.open(
          this.nomina?.id ? 'Nómina actualizada exitosamente' : 'Nómina creada exitosamente',
          'Cerrar',
          { duration: 3000 }
        );
        this.dialogRef.close(true);
      },
      error: (err) => {
        this.guardando = false;
        console.error('❌ Error guardando nómina:', err);
        this.manejarError(err);
      }
    });
  }

  cancelar(): void {
    this.dialogRef.close(false);
  }

  manejarError(err: any): void {
    let mensajeError = 'Error guardando nómina';
    if (err.error?.error) mensajeError = err.error.error;
    else if (err.error?.message) mensajeError = err.error.message;
    else if (err.error?.details) mensajeError = err.error.details;
    else if (err.message) mensajeError = err.message;

    if (mensajeError.includes('no tiene asignado un puesto')) {
      mensajeError = 'El empleado seleccionado no tiene un puesto con salario.';
    }

    this.snackBar.open(mensajeError, 'Cerrar', { duration: 5000 });
  }

  // ✅ MEJORADO: Método para verificar si el empleado tiene puesto
  empleadoTienePuesto(empleadoId?: number): boolean {
    if (!empleadoId) return false;

    // Primero verificar si tenemos información del backend
    if (this.infoEmpleadoSeleccionado && this.infoEmpleadoSeleccionado.tieneBolsaActiva) {
      return true;
    }

    // Luego verificar en datos locales
    const empleado = this.empleados.find(e => e.id === Number(empleadoId));
    return !!empleado && !!(empleado.bolsaEmpleo?.puesto || (empleado as any).bolsaPuesto);
  }

  // ✅ MEJORADO: Método para obtener mensaje de advertencia
  getMensajeAdvertencia(): string | null {
    const empleadoId = this.form.get('empleadoId')?.value;
    if (!empleadoId) return null;

    // Usar información del backend si está disponible
    if (this.infoEmpleadoSeleccionado) {
      if (!this.infoEmpleadoSeleccionado.tieneBolsaActiva) {
        return `⚠️ ${this.infoEmpleadoSeleccionado.nombre} ${this.infoEmpleadoSeleccionado.apellido} no tiene un puesto activo con salario.`;
      }
      return null;
    }

    // Fallback a datos locales
    const empleado = this.empleados.find(e => e.id === Number(empleadoId));
    if (!empleado) return null;

    if (!(empleado.bolsaEmpleo?.puesto || (empleado as any).bolsaPuesto)) {
      return `⚠️ El empleado ${empleado.nombre} ${empleado.apellido} no tiene puesto asignado.`;
    }

    return null;
  }

  // ✅ NUEVO: Método para obtener el puesto del empleado seleccionado
  getPuestoEmpleado(): string {
    if (this.infoEmpleadoSeleccionado) {
      return this.infoEmpleadoSeleccionado.puesto || 'No asignado';
    }

    const empleadoId = this.form.get('empleadoId')?.value;
    if (empleadoId) {
      const empleado = this.empleados.find(e => e.id === Number(empleadoId));
      if (empleado) {
        return empleado.bolsaEmpleo?.puesto ||
               (empleado as any)?.bolsaPuesto ||
               (empleado as any)?.puesto ||
               'No asignado';
      }
    }

    return 'No seleccionado';
  }

  // ✅ NUEVO: Método para verificar si se está cargando información
  estaCargandoInformacion(): boolean {
    return this.cargandoSalario;
  }
}
