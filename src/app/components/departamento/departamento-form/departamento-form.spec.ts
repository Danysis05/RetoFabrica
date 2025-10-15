import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DepartamentosFormComponent } from './departamento-form';

describe('DepartamentoForm', () => {
  let component: DepartamentosFormComponent;
  let fixture: ComponentFixture<DepartamentosFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DepartamentosFormComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DepartamentosFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
