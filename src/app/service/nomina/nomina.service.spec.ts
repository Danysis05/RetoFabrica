import { TestBed } from '@angular/core/testing';

import { NominaServiceTs } from './nomina.service.ts';

describe('NominaServiceTs', () => {
  let service: NominaServiceTs;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(NominaServiceTs);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
