import { TestBed } from '@angular/core/testing';

import { BolsaEmpleoService } from './bolsa-empleo.service.ts';

describe('BolsaEmpleoService', () => {
  let service: BolsaEmpleoService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BolsaEmpleoService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
