// 📁 models/nomina.model.ts - ACTUALIZADO
export interface Nomina {
  id?: number;
  empleadoId: number;
  fechaPago?: Date;
  salarioBase: number;
  horasExtras: number;
  bonificaciones: number;
  deducciones: number;
  totalPago: number;

  // ✅ AGREGAR estas propiedades que faltaban
  diasFaltantes?: number;
  empleadoNombre?: string;
  empleadoPuesto?: string;

  // Propiedades que podrían venir del backend con nombres diferentes
  empleado_id?: number;
  horas_extras?: number;
  dias_faltantes?: number;
  fecha_pago?: Date;
  salario_base?: number;
  total_pago?: number;
  empleado_nombre?: string;
  empleado_puesto?: string;
}
