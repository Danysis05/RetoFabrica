export interface Nomina {
  id?: number;

  // ✅ CAMBIAR a camelCase (consistente con el resto de la app)
  empleadoId: number;
  fechaPago?: Date;
  salarioBase: number;
  horasExtras: number;
  bonificaciones: number;
  deducciones: number;
  totalPago: number;
  diasFaltantes?: number;

  // Campos para mostrar información
  empleadoNombre?: string;
  empleadoPuesto?: string;

  // ✅ MANTENER campos del backend para mapeo
  empleado_id?: number;
  horas_extras?: number;
  dias_faltantes?: number;
  fecha_pago?: Date;
  salario_base?: number;
  total_pago?: number;
  empleado_nombre?: string;
  empleado_puesto?: string;
}

// ✅ INTERFAZ PARA PAYLOAD (usa snake_case para el backend)
export interface NominaPayload {
  empleado_id: number;
  horas_extras: number;
  dias_faltantes: number;
}
