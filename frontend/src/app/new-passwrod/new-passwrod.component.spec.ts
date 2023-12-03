import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NewPasswrodComponent } from './new-passwrod.component';

describe('NewPasswrodComponent', () => {
  let component: NewPasswrodComponent;
  let fixture: ComponentFixture<NewPasswrodComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [NewPasswrodComponent]
    });
    fixture = TestBed.createComponent(NewPasswrodComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
