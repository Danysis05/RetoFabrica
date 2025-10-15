import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { BolsaEmpleo } from '../../../models/bolsa.empleo';
import { Departamento } from '../../../models/departamento.model';
import { BolsaEmpleoService } from '../../../service/bolsa-empleo/bolsa-empleo.service.ts';
import { DepartamentoService } from '../../../service/departamento/departamentoService';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-bolsa-empleado-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule
  ],
  templateUrl: './bolsa-empleado-form.html',
  styleUrls: ['./bolsa-empleado-form.css']
})
export class BolsaEmpleadoFormComponent implements OnInit {
  @Input() empleo: BolsaEmpleo | null = null;
  @Output() cerrar = new EventEmitter<boolean>();

  form: FormGroup;
  titulo = '';
  guardando = false;
  estados = ['DISPONIBLE', 'OCUPADO', 'CERRADO'];
  departamentos: Departamento[] = [];
  cargandoDepartamentos = true;
  mostrarCampoEstado = false; // Nueva propiedad para controlar visibilidad del estado

  constructor(
    private fb: FormBuilder,
    private bolsaEmpleoService: BolsaEmpleoService,
    private departamentoService: DepartamentoService
  ) {
    this.form = this.fb.group({
      puesto: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(100)]],
      descripcion: ['', [Validators.required, Validators.minLength(10), Validators.maxLength(255)]],
      salario: [0, [Validators.required, Validators.min(0)]],
      estado: ['DISPONIBLE', Validators.required],
      departamento_id: ['', Validators.required],
      empleado_id: [null],
      fecha_cierre: [''],
      fecha_ocupacion: ['']
    });
  }

  ngOnInit(): void {
    this.cargarDepartamentos();

    if (this.empleo && this.empleo.id) {
      this.titulo = 'Editar Empleo';
      this.mostrarCampoEstado = true; // Mostrar estado solo al editar
    } else {
      this.titulo = 'Nuevo Empleo';
      this.mostrarCampoEstado = false; // Ocultar estado al crear

      // Forzar estado DISPONIBLE en nuevos empleos
      this.form.patchValue({
        estado: 'DISPONIBLE'
      });
    }
  }

  cargarDepartamentos(): void {
    this.cargandoDepartamentos = true;
    this.departamentoService.getAll().subscribe({
      next: (departamentos) => {
        this.departamentos = departamentos;
        this.cargandoDepartamentos = false;

        if (this.empleo && this.empleo.id) {
          // Al editar, cargar todos los datos incluyendo el estado actual
          this.form.patchValue({
            puesto: this.empleo.puesto,
            descripcion: this.empleo.descripcion,
            salario: this.empleo.salario,
            estado: this.empleo.estado,
            departamento_id: this.empleo.departamento_id,
            empleado_id: this.empleo.empleado_id || null,
            fecha_ocupacion: this.empleo.fecha_ocupacion || '',

          });
        } else {
          // Al crear, asegurar que siempre sea DISPONIBLE
          this.form.patchValue({
            estado: 'DISPONIBLE'
          });
        }
      },
      error: (error) => {
        console.error('Error cargando departamentos:', error);
        this.cargandoDepartamentos = false;
      }
    });
  }

  guardar(): void {
    if (this.form.invalid) {
      this.marcarCamposComoSucios();
      return;
    }

    this.guardando = true;

    const formData = this.form.value;

    // Crear objeto con estado forzado a DISPONIBLE si es nuevo empleo
    const empleoData: any = {
      puesto: formData.puesto,
      descripcion: formData.descripcion,
      salario: Number(formData.salario),
      estado: this.empleo && this.empleo.id ? formData.estado : 'DISPONIBLE', // Forzar DISPONIBLE en nuevos
      departamento_id: Number(formData.departamento_id),
      empleado_id: formData.empleado_id ? Number(formData.empleado_id) : null,
      fecha_ocupacion: formData.fecha_ocupacion || null,
      fecha_cierre: formData.fecha_cierre || null
    };

    console.log('🎯 Enviando al backend:', empleoData);

    let peticion: Observable<any>;

    if (this.empleo && this.empleo.id) {
      empleoData.id = this.empleo.id;
      peticion = this.bolsaEmpleoService.update(this.empleo.id, empleoData);
    } else {
      peticion = this.bolsaEmpleoService.create(empleoData);
    }

    peticion.subscribe({
      next: (response) => {
        console.log('✅ Respuesta del backend:', response);
        this.guardando = false;
        this.cerrar.emit(true);
      },
      error: (err) => {
        console.error('❌ Error del backend:', err);
        this.guardando = false;
        alert('Error del backend: ' + err.message);
      }
    });
  }

  cancelar(): void {
    this.cerrar.emit(false);
  }

  private marcarCamposComoSucios(): void {
    Object.keys(this.form.controls).forEach(key => {
      const control = this.form.get(key);
      if (control) {
        control.markAsTouched();
      }
    });
  }

  get puesto() { return this.form.get('puesto'); }
  get descripcion() { return this.form.get('descripcion'); }
  get salario() { return this.form.get('salario'); }
  get estado() { return this.form.get('estado'); }
  get departamento_id() { return this.form.get('departamento_id'); }
}
