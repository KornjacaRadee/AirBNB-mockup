import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddPicutresAccommComponent } from './add-picutres-accomm.component';

describe('AddPicutresAccommComponent', () => {
  let component: AddPicutresAccommComponent;
  let fixture: ComponentFixture<AddPicutresAccommComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AddPicutresAccommComponent]
    });
    fixture = TestBed.createComponent(AddPicutresAccommComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
