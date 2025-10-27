import { TestBed } from '@angular/core/testing';

import { Shortener } from './shortener';

describe('Shortener', () => {
  let service: Shortener;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Shortener);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
