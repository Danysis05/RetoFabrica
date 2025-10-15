import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BolsaEmpleadoFormComponent } from './bolsa-empleado-form';

describe('BolsaEmpleadoForm', () => {
  let component:  BolsaEmpleadoFormComponent;
  let fixture: ComponentFixture< BolsaEmpleadoFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BolsaEmpleadoFormComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent( BolsaEmpleadoFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
