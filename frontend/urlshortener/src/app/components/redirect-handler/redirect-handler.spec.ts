import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RedirectHandler } from './redirect-handler';

describe('RedirectHandler', () => {
  let component: RedirectHandler;
  let fixture: ComponentFixture<RedirectHandler>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RedirectHandler]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RedirectHandler);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
