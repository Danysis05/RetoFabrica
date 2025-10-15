// bolsa.empleo.ts - CORREGIDO
export interface BolsaEmpleo {
  id?: number;
  puesto: string;
  descripcion: string;
  salario: number;
  estado: string;
  departamento_id: number;
  empleado_id: number | null;
  fecha_publicacion?: string;
  fecha_ocupacion?: string | null;
}

// Modelo para departamentos
export interface Departamento {
  id: number;
  nombre: string;
  descripcion: string;
}
