import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { EmpleadoService } from '../../../service/empleado/empleadoService.js';
import { BolsaEmpleoService } from '../../../service/bolsa-empleo/bolsa-empleo.service.ts';
import { DepartamentoService } from '../../../service/departamento/departamentoService.js';
import { Empleado } from '../../../models/empleado.model';
import { EmpleadoFormComponent } from '../empleados-form/empleados-form';
import { BolsaEmpleo } from '../../../models/bolsa.empleo';

@Component({
  selector: 'app-empleados-lista',
  standalone: true,
  imports: [CommonModule, FormsModule, EmpleadoFormComponent],
  templateUrl: './empleados-lista.html',
  styleUrls: ['./empleados-lista.css']
})
export class EmpleadoListaComponent implements OnInit {
  empleados: Empleado[] = [];
  empleadosFiltrados: Empleado[] = [];
  empleosDisponibles: BolsaEmpleo[] = [];
  departamentos: any[] = [];

  filtroNombre: string = '';
  filtroConEmpleo: string = '';
  filtroDepartamento: string = '';

  mostrarModal: boolean = false;
  empleadoSeleccionado: Empleado | null = null;

  cargando: boolean = true;
  cargandoEmpleos: boolean = true;

  constructor(
    private empleadoService: EmpleadoService,
    private bolsaEmpleoService: BolsaEmpleoService,
    private departamentoService: DepartamentoService
  ) {}

  ngOnInit(): void {
    // Primero cargamos los empleos para poder asignar puestos correctamente
    this.cargarEmpleosDisponibles();
    this.cargarDepartamentos();
    this.cargarEmpleados();
  }

  cargarEmpleados(): void {
    this.cargando = true;
    this.empleadoService.getAll().subscribe({
      next: (data: any) => {
        let empleadosRaw: Empleado[] = [];

        if (Array.isArray(data)) empleadosRaw = data;
        else if (data?.data && Array.isArray(data.data)) empleadosRaw = data.data;
        else if (data?.empleados && Array.isArray(data.empleados)) empleadosRaw = data.empleados;
        else empleadosRaw = [];

        this.empleados = empleadosRaw;
        this.actualizarPuestosEmpleados();
        this.cargando = false;
      },
      error: (err) => {
        console.error('❌ Error cargando empleados:', err);
        this.empleados = [];
        this.empleadosFiltrados = [];
        this.cargando = false;
      }
    });
  }

  cargarEmpleosDisponibles(): void {
    this.cargandoEmpleos = true;
    this.bolsaEmpleoService.getDisponibles().subscribe({
      next: (data: any) => {
        if (Array.isArray(data)) this.empleosDisponibles = data;
        else if (data?.data && Array.isArray(data.data)) this.empleosDisponibles = data.data;
        else this.empleosDisponibles = [];

        // Actualizamos los puestos de empleados después de cargar los empleos
        this.actualizarPuestosEmpleados();
        this.cargandoEmpleos = false;
      },
      error: (err) => {
        console.error('Error cargando empleos disponibles:', err);
        this.empleosDisponibles = [];
        this.cargandoEmpleos = false;
      }
    });
  }

  cargarDepartamentos(): void {
    this.departamentoService.getAll().subscribe({
      next: (data: any) => {
        if (Array.isArray(data)) this.departamentos = data;
        else if (data?.data && Array.isArray(data.data)) this.departamentos = data.data;
        else this.departamentos = [];
      },
      error: (err) => console.error('Error cargando departamentos:', err)
    });
  }

  // 🔄 Actualiza los puestos de los empleados según la bolsa de empleo
  actualizarPuestosEmpleados(): void {
    this.empleados = this.empleados.map(emp => {
      if (!emp.bolsaPuesto && emp.bolsaEmpleoID) {
        const empleo = this.empleosDisponibles.find(e => e.id === emp.bolsaEmpleoID);
        emp.bolsaPuesto = empleo?.puesto ?? 'No asignado';
      }
      return emp;
    });

    this.aplicarFiltros();
  }

  aplicarFiltros(): void {
  if (!Array.isArray(this.empleados)) {
    this.empleadosFiltrados = [];
    return;
  }

  let empleadosFiltrados = [...this.empleados];

  // 🔹 Filtro por nombre/apellido
  if (this.filtroNombre) {
    const filtro = this.filtroNombre.toLowerCase();
    empleadosFiltrados = empleadosFiltrados.filter(empleado =>
      (empleado.nombre?.toLowerCase().includes(filtro) || false) ||
      (empleado.apellido?.toLowerCase().includes(filtro) || false)
    );
  }

  // 🔹 Filtro por estado (con/sin empleo)
  if (this.filtroConEmpleo) {
    const tieneEmpleo = this.filtroConEmpleo === 'true';
    empleadosFiltrados = empleadosFiltrados.filter(empleado =>
      tieneEmpleo ? empleado.bolsaEmpleoID : !empleado.bolsaEmpleoID
    );
  }

  // 🔹 Filtro por departamento
  if (this.filtroDepartamento) {
    const deptoId = Number(this.filtroDepartamento);

    empleadosFiltrados = empleadosFiltrados.filter(emp => {
      // Intentamos obtener el departamento directamente desde bolsaEmpleo
      let empDeptoId = emp.bolsaEmpleo?.departamento_id;

      // Si no existe, buscamos en empleosDisponibles usando bolsaEmpleoID
      if (!empDeptoId && emp.bolsaEmpleoID) {
        const empleo = this.empleosDisponibles.find(e => e.id === emp.bolsaEmpleoID);
        empDeptoId = empleo?.departamento_id;
      }

      return empDeptoId === deptoId;
    });
  }

  this.empleadosFiltrados = empleadosFiltrados;
}

  limpiarFiltros(): void {
    this.filtroNombre = '';
    this.filtroConEmpleo = '';
    this.filtroDepartamento = '';
    this.aplicarFiltros();
  }

  getNombreDepartamento(id: string | number | undefined): string {
    if (!id) return 'No asignado';
    const idNumber = typeof id === 'string' ? parseInt(id) : id;
    const depto = this.departamentos.find(d => d.id === idNumber);
    return depto ? depto.nombre : `Departamento ID: ${idNumber}`;
  }

  abrirModalAgregar(): void {
    this.empleadoSeleccionado = null;
    this.mostrarModal = true;
  }

  abrirModalEditar(empleado: Empleado): void {
    this.empleadoSeleccionado = { ...empleado };
    this.mostrarModal = true;
  }

  cerrarModal(recargar: boolean = false): void {
    this.mostrarModal = false;
    this.empleadoSeleccionado = null;
    if (recargar) this.cargarEmpleados();
  }

  eliminarEmpleado(id: number | undefined): void {
    if (!id) return console.error('ID de empleado no válido');
    if (confirm('¿Seguro que deseas eliminar este empleado?')) {
      this.empleadoService.delete(id).subscribe({
        next: () => this.cargarEmpleados(),
        error: err => console.error('Error al eliminar empleado:', err)
      });
    }
  }

  trackByEmpleado(index: number, empleado: Empleado): number {
    return empleado.id || index;
  }

  getEstadisticas(): any {
    const total = this.empleados.length;
    const asignados = this.empleados.filter(e => e.bolsaEmpleoID).length;
    const sinPuesto = this.empleados.filter(e => !e.bolsaEmpleoID).length;
    return { total, asignados, sinPuesto };
  }

  get hayFiltrosActivos(): boolean {
    return !!this.filtroNombre || !!this.filtroConEmpleo || !!this.filtroDepartamento;
  }

  get hayEmpleadosFiltrados(): boolean {
    return Array.isArray(this.empleadosFiltrados) && this.empleadosFiltrados.length > 0;
  }

  get empleadosFiltradosSeguro(): Empleado[] {
    return Array.isArray(this.empleadosFiltrados) ? this.empleadosFiltrados : [];
  }
}
