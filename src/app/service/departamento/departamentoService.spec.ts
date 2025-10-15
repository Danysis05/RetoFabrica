import { TestBed } from '@angular/core/testing';

import { DepartamentoService } from './departamentoService.js';

describe('DepartamentoServiceTs', () => {
  let service: DepartamentoService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(DepartamentoService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
