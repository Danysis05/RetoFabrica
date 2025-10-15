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

    this.form.get('empleadoId')?.valueChanges.subscribe((id) => {
      const empleadoId = Number(id);
      const emp = this.empleados.find(e => e.id === empleadoId);
      const salario = emp?.bolsaEmpleo?.salario || (emp as any)?.bolsaSalario || 0;
      this.salarioBaseCalculado = salario;
      this.recalcularTotal();
    });

    this.form.get('horasExtras')?.valueChanges.subscribe(() => this.recalcularTotal());
    this.form.get('diasFaltantes')?.valueChanges.subscribe(() => this.recalcularTotal());
  }

  cargarEmpleados(): void {
    this.empleadoService.getAll().subscribe({
      next: (data) => {
        this.empleados = data;
        console.log('👤 Empleados cargados:', data);
        this.form.updateValueAndValidity();
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
  }

  guardar(): void {
    if (this.form.invalid) {
      this.snackBar.open('Por favor complete todos los campos requeridos', 'Cerrar', { duration: 3000 });
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

  empleadoTienePuesto(empleadoId?: number): boolean {
    if (!empleadoId) return false;
    const empleado = this.empleados.find(e => e.id === Number(empleadoId));
    return !!empleado && !!(empleado.bolsaEmpleo?.puesto || (empleado as any).bolsaPuesto);
  }

  getMensajeAdvertencia(): string | null {
    const empleadoId = this.form.get('empleadoId')?.value;
    if (!empleadoId) return null;
    const empleado = this.empleados.find(e => e.id === Number(empleadoId));
    if (!empleado) return null;

    if (!(empleado.bolsaEmpleo?.puesto || (empleado as any).bolsaPuesto)) {
      return `⚠️ El empleado ${empleado.nombre} ${empleado.apellido} no tiene puesto asignado.`;
    }

    return null;
  }
}
