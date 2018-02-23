import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {HttpClient} from '@angular/common/http';
import {MatStepper} from '@angular/material';

@Component({
  selector: 'app-createstickerset',
  templateUrl: './createstickerset.component.html',
  styleUrls: ['./createstickerset.component.css']
})
export class CreatestickersetComponent implements OnInit {
  firstFormGroup: FormGroup;
  secondFormGroup: FormGroup;
  fileToUpload: FileList = null;
  finalStickerName = '';
  loading=true;

  stickerSetExistWarning = false;
  constructor(private _formBuilder: FormBuilder,
              private http: HttpClient) { }

  ngOnInit() {
    this.firstFormGroup = this._formBuilder.group({
      userId: ['', Validators.required],
      stickerName: ['', Validators.required],
      stickerTitle: ['', Validators.required]
    });
    this.secondFormGroup = this._formBuilder.group({
      files: ['', Validators.required]
    });
  }

  checkStickerSetExist(event: any) {
    this.http.get<any>('/api/stickerset/' + this.firstFormGroup.controls['stickerName'].value)
      .subscribe(resp => {
          if (resp.result) {
            this.stickerSetExistWarning = true;
          }
          console.log(resp);
        },
        err => {
          console.log(err);
        }
      );
  }

  handleFileInput(files: FileList) {
    this.fileToUpload = files;
    console.log(files);

  }

  upload() {
    var files;
    files = this.fileToUpload
        const endpoint = '/api/stickerset';
    const formData: FormData = new FormData();
    formData.append('userId', this.firstFormGroup.controls['userId'].value);
    formData.append('stickerName', this.firstFormGroup.controls['stickerName'].value);
    formData.append('stickerTitle', this.firstFormGroup.controls['stickerTitle'].value);
    for (let i = 0; i < files.length; i++) {
      let file;
      file = files.item(i);
      file = files[i];
      formData.append('files', file, file.name);
    }
    console.log('Start upload')
    this.http
      .post<any>(endpoint, formData)
      .subscribe(resp => {
          console.log(resp);
          if(resp.result.stickerName){
            this.finalStickerName = resp.result.stickerName;
            this.loading = false;
          }
        },
        err => {
        console.log(err);
          alert(err.error.error);
        });
  }

}
