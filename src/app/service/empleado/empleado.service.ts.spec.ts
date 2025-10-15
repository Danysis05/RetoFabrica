import { TestBed } from '@angular/core/testing';
import { EmpleadoService } from './empleadoService.js';

describe('EmpleadoService', () => {
  let service: EmpleadoService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EmpleadoService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
