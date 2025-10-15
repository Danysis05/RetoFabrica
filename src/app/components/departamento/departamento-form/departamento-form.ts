import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Departamento } from '../../../models/departamento.model';
import { DepartamentoService } from '../../../service/departamento/departamentoService';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-departamentos-form',
  standalone: true,
  imports: [ReactiveFormsModule, CommonModule],
  templateUrl: './departamento-form.html',
  styleUrls: ['./departamento-form.css']
})
export class DepartamentosFormComponent implements OnInit {
  @Input() departamento: Departamento | null = null;
  @Output() cerrar = new EventEmitter<boolean>();

  form: FormGroup;
  titulo = '';
  guardando = false;
  esEdicion = false;

  constructor(
    private fb: FormBuilder,
    private depService: DepartamentoService
  ) {
    // ✅ FORMULARIO BASE (sin código - para creación)
    this.form = this.fb.group({
      nombre: ['', [Validators.required, Validators.minLength(2)]],
      descripcion: ['', [Validators.required, Validators.minLength(5)]]
    });
  }

  ngOnInit(): void {
    console.log('🎯 Departamento recibido en formulario:', this.departamento);

    if (this.departamento && this.departamento.id) {
      // ✅ MODO EDICIÓN - Agregar campo código dinámicamente
      this.esEdicion = true;
      this.form.addControl('codigo',
        this.fb.control(
          this.departamento.codigo,
          [Validators.required, Validators.minLength(2), Validators.maxLength(10)]
        )
      );

      this.form.patchValue({
        nombre: this.departamento.nombre,
        descripcion: this.departamento.descripcion,
        codigo: this.departamento.codigo // ✅ Cargar código existente
      });

      this.titulo = 'Editar Departamento';
      console.log('📝 MODO EDICIÓN - ID:', this.departamento.id, 'Código:', this.departamento.codigo);
    } else {
      // ✅ MODO CREACIÓN - Sin campo código
      this.esEdicion = false;
      this.titulo = 'Nuevo Departamento';
      console.log('🆕 MODO CREACIÓN - Código se generará automáticamente');
    }
  }

  guardar(): void {
    if (this.form.invalid) {
      this.marcarCamposComoSucios();
      return;
    }

    this.guardando = true;

    // ✅ PREPARAR DATOS SEGÚN MODO
    let depData: any;

    if (this.esEdicion && this.departamento?.id) {
      // ✅ UPDATE: Enviar todos los campos incluyendo código
      depData = {
        nombre: this.form.value.nombre,
        codigo: this.form.value.codigo, // ✅ Enviar código existente
        descripcion: this.form.value.descripcion
      };
      console.log('🔄 UPDATE - Enviando datos:', depData);
      console.log('🎯 ID del departamento:', this.departamento.id);
    } else {
      // ✅ CREATE: Enviar solo nombre y descripción (código se genera automáticamente)
      depData = {
        nombre: this.form.value.nombre,
        descripcion: this.form.value.descripcion
        // ❌ NO enviar código - se generará automáticamente en el backend
      };
      console.log('🆕 CREATE - Enviando datos (código se generará automáticamente):', depData);
    }

    let peticion: Observable<any>;

    if (this.esEdicion && this.departamento?.id) {
      console.log('🔄 EJECUTANDO UPDATE para ID:', this.departamento.id);
      peticion = this.depService.update(this.departamento.id, depData);
    } else {
      console.log('🆕 EJECUTANDO CREATE');
      peticion = this.depService.create(depData);
    }

    peticion.subscribe({
      next: (response) => {
        console.log('✅ Guardado exitoso - Respuesta:', response);
        this.guardando = false;
        this.cerrar.emit(true);
      },
      error: (err) => {
        console.error('❌ ERROR completo al guardar departamento:');
        console.error('Status:', err.status);
        console.error('Status Text:', err.statusText);
        console.error('URL:', err.url);
        console.error('Mensaje:', err.message);
        console.error('Error body:', err.error);

        this.guardando = false;

        // Mensaje de error más específico
        let mensajeError = `Error ${err.status} del servidor.`;
        if (err.error?.error) {
          mensajeError += ` ${err.error.error}`;
        } else if (err.error?.message) {
          mensajeError += ` ${err.error.message}`;
        }
        alert(mensajeError);
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

  // ✅ GETTERS PARA TEMPLATE
  get nombre() { return this.form.get('nombre'); }
  get descripcion() { return this.form.get('descripcion'); }
  get codigo() { return this.form.get('codigo'); } // ✅ Solo disponible en edición
}
