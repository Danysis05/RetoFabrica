import { Routes } from '@angular/router';
import { EmpleadoListaComponent } from './components/empleados/empleados-lista/empleados-lista';
import { DepartamentosListaComponent } from './components/departamento/departamento-lista/departamento-lista';

import { NominaListComponent } from './components/nomina/nomina-lista/nomina-lista';

export const routes: Routes = [
  { path: '', redirectTo: 'departamentos', pathMatch: 'full' },
  { path: 'empleados', component: EmpleadoListaComponent },
  { path: 'departamentos', component: DepartamentosListaComponent },
  { path: 'bolsa-empleados', loadComponent: () => import('./components/bolsa-empleado/bolsa-empleado-lista/bolsa-empleado-lista').then(m => m.BolsaEmpleadoListaComponent) },
  { path: 'nomina', component: NominaListComponent}
];
