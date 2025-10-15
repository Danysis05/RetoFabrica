import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BolsaEmpleadoListaComponent } from './bolsa-empleado-lista';

describe('BolsaEmpleadoLista', () => {
  let component: BolsaEmpleadoListaComponent;
  let fixture: ComponentFixture<BolsaEmpleadoListaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BolsaEmpleadoListaComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(BolsaEmpleadoListaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
