import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DepartamentosListaComponent } from './departamento-lista';

describe('DepartamentoLista', () => {
  let component: DepartamentosListaComponent;
  let fixture: ComponentFixture<DepartamentosListaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DepartamentosListaComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DepartamentosListaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
