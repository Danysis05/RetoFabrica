import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Empleado, TIPOS_DOCUMENTO } from '../../../models/empleado.model';
import { BolsaEmpleo } from '../../../models/bolsa.empleo';
import { EmpleadoService } from '../../../service/empleado/empleadoService';
import { BolsaEmpleoService } from '../../../service/bolsa-empleo/bolsa-empleo.service.ts';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-empleado-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule
  ],
  templateUrl: './empleados-form.html',
  styleUrls: ['./empleados-form.css']
})
export class EmpleadoFormComponent implements OnInit {
  @Input() empleado: Empleado | null = null;
  @Input() empleosDisponibles: BolsaEmpleo[] = [];
  @Output() cerrar = new EventEmitter<boolean>();

  form: FormGroup;
  guardando = false;
  titulo = '';
  tiposDocumento = TIPOS_DOCUMENTO;

  constructor(
    private fb: FormBuilder,
    private empleadoService: EmpleadoService,
    private bolsaEmpleoService: BolsaEmpleoService,
    private snackBar: MatSnackBar
  ) {
    // Inicializamos el formulario con validación mínima; detalles por tipo se aplican en ngOnInit
    this.form = this.fb.group({
      nombre: ['', [Validators.required, Validators.minLength(2), Validators.maxLength(100)]],
      apellido: ['', [Validators.required, Validators.minLength(2), Validators.maxLength(100)]],
      documentoTipo: ['Cédula de Ciudadanía', [Validators.required]],
      documentoNumero: ['', [Validators.required]], // validadores dinámicos en ngOnInit
      correoElectronico: ['', [Validators.required, Validators.email]],
      ciudad: ['', [Validators.required, Validators.minLength(2), Validators.maxLength(100)]],
      direccion: ['', [Validators.required, Validators.minLength(5), Validators.maxLength(255)]],
      telefono: ['', [Validators.maxLength(20)]],
      bolsaEmpleoID: ['', [Validators.required]]
    });
  }

  ngOnInit(): void {
    if (this.empleado && this.empleado.id) {
      this.titulo = 'Editar Empleado';
      this.cargarDatosEmpleado();
    } else {
      this.titulo = 'Nuevo Empleado';
    }

    // Aplicar validaciones dinámicas según el tipo seleccionado
    this.form.get('documentoTipo')?.valueChanges.subscribe(tipo => {
      this.actualizarValidacionesDocumento(tipo);
    });

    // Aplicar validaciones inmediatamente según el valor inicial (si lo hay)
    this.actualizarValidacionesDocumento(this.form.get('documentoTipo')?.value);

    // Debug opcional
    setTimeout(() => this.verificarDatos(), 1000);
  }

  /**
   * Actualiza los validadores del campo documentoNumero según el tipo seleccionado.
   * No modifica nada más en la lógica principal del componente.
   */
  private actualizarValidacionesDocumento(tipo: string | null | undefined): void {
    const control = this.form.get('documentoNumero');
    if (!control) return;

    // Limpiar validadores previos
    control.clearValidators();

    switch (tipo) {
      case 'Cédula de Ciudadanía':
      case 'Cedula de Ciudadania': // por si hay variaciones
      case 'Cédula':
        control.setValidators([
          Validators.required,
          Validators.pattern(/^[0-9]+$/), // solo números
          Validators.minLength(10), // ajuste común para cédula colombiana
          Validators.maxLength(10)
        ]);
        break;

      case 'Cédula de Extranjería':
      case 'Cedula de Extranjería':
      case 'Cedula de Extranjeria':
        control.setValidators([
          Validators.required,
          Validators.pattern(/^[0-9A-Za-z]+$/), // números y letras
          Validators.minLength(6),
          Validators.maxLength(12)
        ]);
        break;

      case 'Pasaporte':
        control.setValidators([
          Validators.required,
          Validators.pattern(/^[A-Za-z0-9]+$/), // alfanumérico
          Validators.minLength(6),
          Validators.maxLength(12)
        ]);
        break;

      case 'Tarjeta de Identidad':
      case 'Tarjeta de identidad':
        control.setValidators([
          Validators.required,
          Validators.pattern(/^[0-9]+$/), // solo números
          Validators.minLength(10),
          Validators.maxLength(10)
        ]);
        break;

      default:
        // Si tipo desconocido, dejar validación básica (requerido)
        control.setValidators([Validators.required]);
        break;
    }

    control.updateValueAndValidity();
  }

  // Para que la vista pueda cambiar el tipo de input (number/text) según el tipo seleccionado
  getTipoInputDocumento(): string {
    const tipo = this.form.get('documentoTipo')?.value;
    // Usamos text para pasaportes y para cualquier que pueda contener letras
    if (!tipo) return 'text';
    if (tipo.toLowerCase().includes('pasaport') || tipo.toLowerCase().includes('extranjer')) return 'text';
    // por defecto usar number (mejor UX en móviles) pero el control seguirá validando según pattern
    return 'number';
  }

  // ===================== MANTENEMOS TODOS TUS MÉTODOS ORIGINALES =====================

  get empleosFiltrados(): BolsaEmpleo[] {
    if (!this.empleosDisponibles || this.empleosDisponibles.length === 0) return [];

    return this.empleosDisponibles.filter(empleo => {
      const estaDisponible = empleo.estado === 'DISPONIBLE' &&
        (empleo.empleado_id === null || empleo.empleado_id === undefined || empleo.empleado_id === 0);
      const esElActual = this.empleado && empleo.id === this.empleado.bolsaEmpleoID;
      return estaDisponible || esElActual;
    });
  }

  verificarDatos(): void {
    console.log('=== DEBUG COMPLETO DE DATOS ===');
    console.log('Empleado actual:', this.empleado);
    console.log('Todos los empleos disponibles:', this.empleosDisponibles);
    console.log('Empleos filtrados:', this.empleosFiltrados);

    if (this.empleosDisponibles && this.empleosDisponibles.length > 0) {
      console.log('=== DETALLE POR EMPLEO ===');
      this.empleosDisponibles.forEach((emp, i) => {
        console.log(`Empleo ${i}:`, {
          id: emp.id,
          puesto: emp.puesto,
          estado: emp.estado,
          empleado_id: emp.empleado_id,
          enFiltrado: this.empleosFiltrados.some(f => f.id === emp.id)
        });
      });
    }
    console.log('=== FIN DEBUG ===');
  }

  cargarDatosEmpleado(): void {
    if (this.empleado) {
      this.form.patchValue({
        nombre: this.empleado.nombre,
        apellido: this.empleado.apellido,
        documentoTipo: this.empleado.documentoTipo,
        documentoNumero: this.empleado.documentoNumero,
        correoElectronico: this.empleado.correoElectronico,
        ciudad: this.empleado.ciudad,
        direccion: this.empleado.direccion,
        telefono: this.empleado.telefono || '',
        bolsaEmpleoID: this.empleado.bolsaEmpleoID || ''
      });
    }
  }

  getPuestoSeleccionado(): BolsaEmpleo | undefined {
    const empleoId = this.form.get('bolsaEmpleoID')?.value;
    return this.empleosDisponibles.find(empleo => empleo.id === empleoId);
  }

  getNombreDepartamento(departamentoId: number): string {
    return `Departamento ID: ${departamentoId}`;
  }

    guardar(): void {
  if (this.form.invalid) {
    this.marcarCamposComoSucios();
    return;
  }

  this.guardando = true;

  const formData = this.form.value;

  // ✅ Ajuste: enviar siempre un número válido, 0 si no hay selección
  const bolsaID = formData.bolsaEmpleoID ? Number(formData.bolsaEmpleoID) : 0;

  const empleoData: any = {
    nombre: formData.nombre?.trim(),
    apellido: formData.apellido?.trim(),
    documentoTipo: formData.documentoTipo,
    documentoNumero: formData.documentoNumero,
    correoElectronico: formData.correoElectronico?.trim(),
    ciudad: formData.ciudad?.trim(),
    direccion: formData.direccion?.trim(),
    telefono: formData.telefono?.trim() || '',
    bolsaEmpleoID: bolsaID
  };

  console.log('🎯 DATOS A ENVIAR AL BACKEND:', empleoData);

  const peticion: Observable<any> =
    this.empleado && this.empleado.id
      ? this.empleadoService.update({ ...empleoData, id: this.empleado.id })
      : this.empleadoService.create(empleoData);

  peticion.subscribe({
    next: (response) => {
      console.log('✅ RESPUESTA COMPLETA DEL BACKEND:', response);

      this.guardando = false;
      this.cerrar.emit(true);

      this.snackBar.open('Empleado guardado exitosamente', 'Cerrar', {
        duration: 3000,
        panelClass: ['snackbar-success']
      });
    },
    error: (err) => {
      console.error('❌ ERROR GUARDANDO EMPLEADO:', err);

      console.error('🔍 DETALLE COMPLETO DEL ERROR:', {
        status: err.status,
        statusText: err.statusText,
        url: err.url,
        error: err.error,
        headers: err.headers
      });

      this.guardando = false;

      let mensajeError = 'Error al guardar el empleado';

      if (err.error) {
        if (err.error.message) mensajeError = err.error.message;
        else if (err.error.error) mensajeError = err.error.error;
        else if (typeof err.error === 'string') mensajeError = err.error;
      } else if (err.status === 500) {
        mensajeError = 'Error interno del servidor.';
      } else if (err.status === 409) {
        mensajeError = 'El correo electrónico o documento ya están en uso.';
      } else if (err.status === 400) {
        mensajeError = 'Datos inválidos. Verifica los campos.';
      }

      this.snackBar.open(mensajeError, 'Cerrar', {
        duration: 5000,
        panelClass: ['snackbar-error']
      });
    }
  });
}


  cancelar(): void {
    this.cerrar.emit(false);
  }

  private marcarCamposComoSucios(): void {
    Object.keys(this.form.controls).forEach(key => {
      const control = this.form.get(key);
      if (control) control.markAsTouched();
    });
  }

  // ✅ Getters para el template (sin cambios)
  get nombre() { return this.form.get('nombre'); }
  get apellido() { return this.form.get('apellido'); }
  get documentoTipo() { return this.form.get('documentoTipo'); }
  get documentoNumero() { return this.form.get('documentoNumero'); }
  get correoElectronico() { return this.form.get('correoElectronico'); }
  get ciudad() { return this.form.get('ciudad'); }
  get direccion() { return this.form.get('direccion'); }
  get telefono() { return this.form.get('telefono'); }
  get bolsaEmpleoID() { return this.form.get('bolsaEmpleoID'); }
}
