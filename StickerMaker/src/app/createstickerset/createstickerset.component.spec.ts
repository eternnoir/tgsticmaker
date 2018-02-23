import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatestickersetComponent } from './createstickerset.component';

describe('CreatestickersetComponent', () => {
  let component: CreatestickersetComponent;
  let fixture: ComponentFixture<CreatestickersetComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreatestickersetComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreatestickersetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
