import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NominaFormComponent } from './nomina-form';

describe('NominaForm', () => {
  let component: NominaFormComponent;
  let fixture: ComponentFixture<NominaFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NominaFormComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(NominaFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
