import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NominaListComponent } from './nomina-lista';

describe('NominaLista', () => {
  let component: NominaListComponent;
  let fixture: ComponentFixture<NominaListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NominaListComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(NominaListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
