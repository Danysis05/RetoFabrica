// models/departamento.model.ts - ACTUALIZAR
import { BolsaEmpleo } from './bolsa.empleo';

// models/departamento.model.ts - TEMPORAL para debug
export interface Departamento {
  id?: number;
  nombre: string;
  codigo: string;
  descripcion: string;
  bolsaEmpleo?: any[]; // Cambiar temporalmente a any[]
  // Agregar otras posibles propiedades
  BolsaEmpleos?: any[];
  bolsaEmpleos?: any[];
}

// Interface para la respuesta del API
export interface ApiResponse {
  success: boolean;
  data: Departamento[];
  count: number;
  message?: string;
}
