import { BolsaEmpleo } from "./bolsa.empleo";

// empleado.model.ts
export interface Empleado {
  id?: number;
  nombre: string;
  apellido: string;
  documentoTipo: string;
  documentoNumero: string;
  correoElectronico: string;
  ciudad: string;
  direccion: string;
  telefono?: string;

  // ✅ Asegúrate que coincida con el backend
  bolsaEmpleoID?: number | null;

  // Campos de solo lectura
  departamento?: string;
  bolsaPuesto?: string;
  fechaCreacion?: string | Date;
  fechaModificacion?: string | Date;

  // Relación completa
  bolsaEmpleo?: BolsaEmpleo;
}

export const TIPOS_DOCUMENTO = [
  'Cédula de Ciudadanía',
  'Cédula de Extranjería',
  'Pasaporte',
  'Tarjeta de Identidad'
];
