import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatePeriodComponent } from './create-period.component';

describe('CreatePeriodComponent', () => {
  let component: CreatePeriodComponent;
  let fixture: ComponentFixture<CreatePeriodComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [CreatePeriodComponent]
    });
    fixture = TestBed.createComponent(CreatePeriodComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
